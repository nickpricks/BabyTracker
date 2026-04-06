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

func handleListSleep(w http.ResponseWriter, r *http.Request) {
	entries, err := storage.LoadSleep()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	limit, offset := parsePagination(r)
	page, total := paginateReverse(entries, limit, offset)
	jsonResponse(w, http.StatusOK, PaginatedResponse{Items: page, Total: total, Limit: limit, Offset: offset})
}

func handleLogSleep(w http.ResponseWriter, r *http.Request) {
	var entry models.SleepEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if entry.Date == "" || entry.Type == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required fields (date, type)"})
		return
	}
	log.Printf("Log Sleep: %+v\n", entry)
	if err := storage.SaveSleep(&entry); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusCreated, entry)
}

func handleGetSleep(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	entries, err := storage.LoadSleep()
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
	jsonResponse(w, http.StatusNotFound, map[string]string{"error": "sleep entry not found"})
}

func handleUpdateSleep(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	var entry models.SleepEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if entry.Date == "" || entry.Type == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required fields (date, type)"})
		return
	}
	log.Printf("Update Sleep ID %d: %+v\n", id, entry)
	if err := storage.UpdateSleep(id, &entry); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	entry.ID = id
	jsonResponse(w, http.StatusOK, entry)
}

func handleDeleteSleep(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	log.Printf("Delete Sleep ID %d\n", id)
	if err := storage.DeleteSleep(id); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}
