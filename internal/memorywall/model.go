package memorywall

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityFamily  Visibility = "family"
	VisibilityPrivate Visibility = "private"
)

type ContentKind string

const (
	ContentStory   ContentKind = "story"
	ContentPhoto   ContentKind = "photo"
	ContentYear    ContentKind = "year"
	ContentMessage ContentKind = "message"
)

type Elder struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Profile        string          `json:"profile"`
	Visibility     Visibility      `json:"visibility"`
	Stories        []Story         `json:"stories"`
	Photos         []Photo         `json:"photos"`
	ImportantYears []ImportantYear `json:"importantYears"`
	FamilyMessages []FamilyMessage `json:"familyMessages"`
}

type Story struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Source string `json:"source"`
}

type Photo struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
}

type ImportantYear struct {
	ID    string `json:"id"`
	Year  int    `json:"year"`
	Label string `json:"label"`
}

type FamilyMessage struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type ImportItem struct {
	ExternalID string      `json:"externalId"`
	ElderID    string      `json:"elderId"`
	Kind       ContentKind `json:"kind"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Author     string      `json:"author"`
	Year       int         `json:"year"`
	Visibility Visibility  `json:"visibility"`
}

type Submission struct {
	ID         string      `json:"id"`
	BatchID    string      `json:"batchId"`
	ExternalID string      `json:"externalId"`
	ElderID    string      `json:"elderId"`
	Kind       ContentKind `json:"kind"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Author     string      `json:"author"`
	Year       int         `json:"year"`
	Visibility Visibility  `json:"visibility"`
	State      string      `json:"state"`
}

type Batch struct {
	ID            string       `json:"id"`
	Status        string       `json:"status"`
	SubmissionIDs []string     `json:"submissionIds"`
	Errors        []BatchError `json:"errors"`
}

type BatchError struct {
	ExternalID string `json:"externalId"`
	Message    string `json:"message"`
}

type ImportResult struct {
	Batch Batch `json:"batch"`
}

type ExportData struct {
	Elders      []Elder      `json:"elders"`
	Submissions []Submission `json:"submissions"`
	Batches     []Batch      `json:"batches"`
}
