package memorywall

import (
	"encoding/json"
	"testing"
)

func TestPublicWallFixture(t *testing.T) {
	store := NewStore()
	elders := store.ListPublicElders()
	if len(elders) != 1 || elders[0].ID != "elder-001" {
		t.Fatalf("unexpected public elders: %#v", elders)
	}
	if len(elders[0].Stories) != 1 || len(elders[0].ImportantYears) != 2 {
		t.Fatalf("fixture content missing: %#v", elders[0])
	}
}

func TestReviewSubmissionPublishesStory(t *testing.T) {
	store := NewStore()
	result := store.ImportBatch([]ImportItem{{ExternalID: "family-1", ElderID: "elder-001", Kind: ContentStory, Title: "一封信", Content: "记得常回家看看", Author: "林晓"}})
	if len(result.Batch.SubmissionIDs) != 1 {
		t.Fatalf("submission was not accepted: %#v", result)
	}
	if _, err := store.ReviewSubmission(result.Batch.SubmissionIDs[0], true); err != nil {
		t.Fatal(err)
	}
	elder, _ := store.GetElder("elder-001")
	if len(elder.Stories) != 2 || elder.Stories[1].Title != "一封信" {
		t.Fatalf("approved story missing: %#v", elder.Stories)
	}
}

func TestExportIsDeterministicJSON(t *testing.T) {
	store := NewStore()
	first, err := store.Export()
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Export()
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatalf("exports differ: %q and %q", first, second)
	}
	var decoded ExportData
	if err := json.Unmarshal(first, &decoded); err != nil || len(decoded.Elders) != 2 {
		t.Fatalf("invalid export: %s", first)
	}
}

func TestImportBatchReportsMalformedItem(t *testing.T) {
	store := NewStore()
	result := store.ImportBatch([]ImportItem{
		{ExternalID: "family-2", ElderID: "elder-001", Kind: ContentMessage, Author: "林晓", Content: "周末见"},
		{ExternalID: "broken-1", ElderID: "elder-001", Kind: ContentYear, Title: "没有年份", Content: "二零零一年"},
	})
	if result.Batch.Status != BatchFailed {
		t.Fatalf("batch status = %q", result.Batch.Status)
	}
	if len(result.Batch.Errors) != 1 || result.Batch.Errors[0].ExternalID != "broken-1" {
		t.Fatalf("batch errors = %#v", result.Batch.Errors)
	}
}
