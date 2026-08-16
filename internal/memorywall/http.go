package memorywall

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(store *Store) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)
	router.HandleFunc("/api/elders", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusOK, store.ListPublicElders()) }).Methods(http.MethodGet)
	router.HandleFunc("/api/elders/{id}", func(w http.ResponseWriter, r *http.Request) {
		elder, ok := store.GetElder(mux.Vars(r)["id"])
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "elder not found"})
			return
		}
		writeJSON(w, http.StatusOK, elder)
	}).Methods(http.MethodGet)
	router.HandleFunc("/api/import-batches", func(w http.ResponseWriter, r *http.Request) {
		var items []ImportItem
		if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		writeJSON(w, http.StatusOK, store.ImportBatch(items))
	}).Methods(http.MethodPost)
	router.HandleFunc("/api/submissions/pending", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusOK, store.PendingSubmissions()) }).Methods(http.MethodGet)
	router.HandleFunc("/api/submissions/{id}/review", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Approve bool `json:"approve"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		submission, err := store.ReviewSubmission(mux.Vars(r)["id"], request.Approve)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, submission)
	}).Methods(http.MethodPost)
	router.HandleFunc("/api/elders/{id}/visibility", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Visibility Visibility `json:"visibility"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		elder, err := store.SetVisibility(mux.Vars(r)["id"], request.Visibility)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, elder)
	}).Methods(http.MethodPut)
	router.HandleFunc("/api/export", func(w http.ResponseWriter, _ *http.Request) {
		data, err := store.Export()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=memory-wall-export.json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}).Methods(http.MethodGet)
	return router
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
