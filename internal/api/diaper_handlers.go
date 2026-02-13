package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"babytracker/internal/models"
	"babytracker/internal/storage"
)

func handleListDiapers(w http.ResponseWriter, r *http.Request) {
	entries, err := storage.LoadDiapers()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, entries)
}

func handleLogDiaper(w http.ResponseWriter, r *http.Request) {
	var entry models.DiaperEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if entry.Date == "" || entry.Type == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required fields (date, type)"})
		return
	}
	log.Printf("Log Diaper: %+v\n", entry)
	if err := storage.SaveDiaper(&entry); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusCreated, entry)
}

func handleGetDiaper(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	entries, err := storage.LoadDiapers()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	for _, e := range entries {
		if e.ID == id {
			jsonResponse(w, http.StatusOK, e)
			return
		}
	}
	jsonResponse(w, http.StatusNotFound, map[string]string{"error": "diaper entry not found"})
}
