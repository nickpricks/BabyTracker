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
	limit, offset := parsePagination(r)
	page, total := paginateReverse(entries, limit, offset)
	jsonResponse(w, http.StatusOK, PaginatedResponse{Items: page, Total: total, Limit: limit, Offset: offset})
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

func handleUpdateDiaper(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	var entry models.DiaperEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if entry.Date == "" || entry.Type == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required fields (date, type)"})
		return
	}
	log.Printf("Update Diaper ID %d: %+v\n", id, entry)
	if err := storage.UpdateDiaper(id, &entry); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	entry.ID = id
	jsonResponse(w, http.StatusOK, entry)
}

func handleDeleteDiaper(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	log.Printf("Delete Diaper ID %d\n", id)
	if err := storage.DeleteDiaper(id); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}
