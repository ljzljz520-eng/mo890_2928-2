package memorywall

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	BatchCompleted     = "completed"
	BatchFailed        = "failed"
	SubmissionPending  = "pending"
	SubmissionApproved = "approved"
	SubmissionRejected = "rejected"
)

type Store struct {
	mu             sync.RWMutex
	elders         map[string]*Elder
	submissions    map[string]*Submission
	batches        map[string]*Batch
	nextBatch      int
	nextSubmission int
}

func NewStore() *Store {
	return &Store{
		elders: map[string]*Elder{
			"elder-001": {
				ID:             "elder-001",
				Name:           "林秀梅",
				Profile:        "从纺织厂退休的社区志愿者，喜欢讲述老街的变迁。",
				Visibility:     VisibilityPublic,
				Stories:        []Story{{ID: "story-001", Title: "河畔的夏夜", Body: "年轻时和邻里在河边乘凉，收音机里总有熟悉的歌。", Source: "社区整理"}},
				Photos:         []Photo{{ID: "photo-001", Caption: "社区合唱队合影", URL: "/fixtures/choir.jpg"}},
				ImportantYears: []ImportantYear{{ID: "year-001", Year: 1968, Label: "来到这座城市"}, {ID: "year-002", Year: 1995, Label: "加入社区志愿队"}},
				FamilyMessages: []FamilyMessage{{ID: "message-001", Author: "女儿 林晓", Message: "谢谢您把温柔和勇气留给我们。"}},
			},
			"elder-002": {
				ID:         "elder-002",
				Name:       "周国安",
				Profile:    "老木匠，退休后在社区开设手作课堂。",
				Visibility: VisibilityFamily,
			},
		},
		submissions: make(map[string]*Submission),
		batches:     make(map[string]*Batch),
	}
}

func (s *Store) ListPublicElders() []Elder {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Elder, 0, len(s.elders))
	for _, elder := range s.elders {
		if elder.Visibility == VisibilityPublic {
			items = append(items, cloneElder(*elder))
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (s *Store) GetElder(id string) (Elder, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	elder, ok := s.elders[id]
	if !ok {
		return Elder{}, false
	}
	return cloneElder(*elder), true
}

func (s *Store) SetVisibility(id string, visibility Visibility) (Elder, error) {
	if !validVisibility(visibility) {
		return Elder{}, errors.New("visibility must be public, family, or private")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	elder, ok := s.elders[id]
	if !ok {
		return Elder{}, errors.New("elder not found")
	}
	elder.Visibility = visibility
	return cloneElder(*elder), nil
}

func (s *Store) ImportBatch(items []ImportItem) ImportResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	batchID := fmt.Sprintf("batch-%04d", s.nextBatch+1)
	s.nextBatch++
	record := &Batch{ID: batchID, Status: BatchCompleted, SubmissionIDs: []string{}, Errors: []BatchError{}}
	for _, item := range items {
		submission, err := s.processItemLocked(batchID, item)
		if err != nil {
			continue
		}
		record.SubmissionIDs = append(record.SubmissionIDs, submission.ID)
	}
	s.batches[batchID] = record
	return ImportResult{Batch: cloneBatch(*record)}
}

func (s *Store) processItemLocked(batchID string, item ImportItem) (*Submission, error) {
	if strings.TrimSpace(item.ExternalID) == "" {
		return nil, errors.New("externalId is required")
	}
	if _, ok := s.elders[item.ElderID]; !ok {
		return nil, errors.New("elder not found")
	}
	if !validVisibility(item.Visibility) {
		item.Visibility = VisibilityFamily
	}
	if item.Visibility == "" {
		item.Visibility = VisibilityFamily
	}
	submission := &Submission{
		ID:         fmt.Sprintf("submission-%04d", s.nextSubmission+1),
		BatchID:    batchID,
		ExternalID: item.ExternalID,
		ElderID:    item.ElderID,
		Kind:       item.Kind,
		Title:      strings.TrimSpace(item.Title),
		Content:    strings.TrimSpace(item.Content),
		Author:     strings.TrimSpace(item.Author),
		Year:       item.Year,
		Visibility: item.Visibility,
		State:      SubmissionPending,
	}
	switch item.Kind {
	case ContentStory:
		if submission.Title == "" || submission.Content == "" {
			return nil, errors.New("story title and content are required")
		}
	case ContentPhoto:
		if submission.Content == "" {
			return nil, errors.New("photo url is required")
		}
	case ContentYear:
		if submission.Year < 1 {
			parsed, err := strconv.Atoi(submission.Content)
			if err != nil || parsed < 1 {
				return nil, errors.New("year must be a positive number")
			}
			submission.Year = parsed
		}
		if submission.Title == "" {
			return nil, errors.New("year label is required")
		}
	case ContentMessage:
		if submission.Content == "" || submission.Author == "" {
			return nil, errors.New("message and author are required")
		}
	default:
		return nil, errors.New("unsupported content kind")
	}
	s.nextSubmission++
	s.submissions[submission.ID] = submission
	return submission, nil
}

func (s *Store) PendingSubmissions() []Submission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Submission, 0)
	for _, submission := range s.submissions {
		if submission.State == SubmissionPending {
			result = append(result, cloneSubmission(*submission))
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (s *Store) ReviewSubmission(id string, approve bool) (Submission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	submission, ok := s.submissions[id]
	if !ok {
		return Submission{}, errors.New("submission not found")
	}
	if submission.State != SubmissionPending {
		return Submission{}, errors.New("submission already reviewed")
	}
	if !approve {
		submission.State = SubmissionRejected
		return cloneSubmission(*submission), nil
	}
	elder, ok := s.elders[submission.ElderID]
	if !ok {
		return Submission{}, errors.New("elder not found")
	}
	switch submission.Kind {
	case ContentStory:
		elder.Stories = append(elder.Stories, Story{ID: submission.ID, Title: submission.Title, Body: submission.Content, Source: submission.Author})
	case ContentPhoto:
		elder.Photos = append(elder.Photos, Photo{ID: submission.ID, Caption: submission.Title, URL: submission.Content})
	case ContentYear:
		elder.ImportantYears = append(elder.ImportantYears, ImportantYear{ID: submission.ID, Year: submission.Year, Label: submission.Title})
	case ContentMessage:
		elder.FamilyMessages = append(elder.FamilyMessages, FamilyMessage{ID: submission.ID, Author: submission.Author, Message: submission.Content})
	}
	submission.State = SubmissionApproved
	return cloneSubmission(*submission), nil
}

func (s *Store) Export() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	elders := make([]Elder, 0, len(s.elders))
	for _, elder := range s.elders {
		elders = append(elders, cloneElder(*elder))
	}
	sort.Slice(elders, func(i, j int) bool { return elders[i].ID < elders[j].ID })
	submissions := make([]Submission, 0, len(s.submissions))
	for _, submission := range s.submissions {
		submissions = append(submissions, cloneSubmission(*submission))
	}
	sort.Slice(submissions, func(i, j int) bool { return submissions[i].ID < submissions[j].ID })
	batches := make([]Batch, 0, len(s.batches))
	for _, batch := range s.batches {
		batches = append(batches, cloneBatch(*batch))
	}
	sort.Slice(batches, func(i, j int) bool { return batches[i].ID < batches[j].ID })
	return json.MarshalIndent(ExportData{Elders: elders, Submissions: submissions, Batches: batches}, "", "  ")
}

func validVisibility(value Visibility) bool {
	return value == VisibilityPublic || value == VisibilityFamily || value == VisibilityPrivate
}

func cloneElder(value Elder) Elder {
	value.Stories = append([]Story(nil), value.Stories...)
	value.Photos = append([]Photo(nil), value.Photos...)
	value.ImportantYears = append([]ImportantYear(nil), value.ImportantYears...)
	value.FamilyMessages = append([]FamilyMessage(nil), value.FamilyMessages...)
	return value
}

func cloneSubmission(value Submission) Submission { return value }

func cloneBatch(value Batch) Batch {
	value.SubmissionIDs = append([]string(nil), value.SubmissionIDs...)
	value.Errors = append([]BatchError(nil), value.Errors...)
	return value
}
