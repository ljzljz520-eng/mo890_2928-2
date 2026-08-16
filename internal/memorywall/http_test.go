package memorywall

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPImportAndPendingReview(t *testing.T) {
	server := httptest.NewServer(NewHTTPHandler(NewStore()))
	defer server.Close()
	payload := []byte(`[{"externalId":"api-1","elderId":"elder-001","kind":"photo","title":"旧照片","content":"/photos/old.jpg","visibility":"public"}]`)
	response, err := http.Post(server.URL+"/api/import-batches", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", response.StatusCode)
	}
	response.Body.Close()
	pending, err := http.Get(server.URL + "/api/submissions/pending")
	if err != nil {
		t.Fatal(err)
	}
	if pending.StatusCode != http.StatusOK {
		t.Fatalf("pending status = %d", pending.StatusCode)
	}
	pending.Body.Close()
}
