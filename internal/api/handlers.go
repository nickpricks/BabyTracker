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

// jsonResponse writes a JSON response with the given status code.
func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// handleListFeeds returns all feed entries.
func handleListFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := storage.LoadFeeds()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, feeds)
}

// handleLogFeed logs a new feed entry.
func handleLogFeed(w http.ResponseWriter, r *http.Request) {
	var feed models.FeedEntry
	if err := json.NewDecoder(r.Body).Decode(&feed); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if feed.Type == "" || feed.Date == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
		return
	}
	log.Printf("Log Feed: %+v\n", feed)
	if err := storage.SaveFeed(&feed); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusCreated, feed)
}

// handleGetFeed returns a single feed entry by ID.
func handleGetFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid feed ID"})
		return
	}
	feeds, err := storage.LoadFeeds()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	for _, feed := range feeds {
		if feed.ID == id {
			jsonResponse(w, http.StatusOK, feed)
			return
		}
	}
	jsonResponse(w, http.StatusNotFound, map[string]string{"error": "feed not found"})
}
