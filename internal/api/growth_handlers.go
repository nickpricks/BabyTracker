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

func handleListGrowth(w http.ResponseWriter, r *http.Request) {
	entries, err := storage.LoadGrowth()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	limit, offset := parsePagination(r)
	page, total := paginateReverse(entries, limit, offset)
	jsonResponse(w, http.StatusOK, PaginatedResponse{Items: page, Total: total, Limit: limit, Offset: offset})
}

func handleLogGrowth(w http.ResponseWriter, r *http.Request) {
	var entry models.GrowthEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if entry.Date == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required field (date)"})
		return
	}
	log.Printf("Log Growth: %+v\n", entry)
	if err := storage.SaveGrowth(&entry); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusCreated, entry)
}

func handleGetGrowth(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	entries, err := storage.LoadGrowth()
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
	jsonResponse(w, http.StatusNotFound, map[string]string{"error": "growth entry not found"})
}

func handleUpdateGrowth(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	var entry models.GrowthEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if entry.Date == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required field (date)"})
		return
	}
	log.Printf("Update Growth ID %d: %+v\n", id, entry)
	if err := storage.UpdateGrowth(id, &entry); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	entry.ID = id
	jsonResponse(w, http.StatusOK, entry)
}

func handleDeleteGrowth(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	log.Printf("Delete Growth ID %d\n", id)
	if err := storage.DeleteGrowth(id); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}
