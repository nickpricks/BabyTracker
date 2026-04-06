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

const defaultLimit = 10

// PaginatedResponse wraps a list response with pagination metadata.
type PaginatedResponse struct {
	Items  interface{} `json:"items"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

// jsonResponse writes a JSON response with the given status code.
func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// parsePagination extracts limit and offset from query params.
// Returns limit (default 10, 0 = all), offset (default 0).
func parsePagination(r *http.Request) (limit, offset int) {
	limit = defaultLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return
}

// paginateReverse reverses a slice in-place and applies offset/limit.
// Returns the page slice, total count, and clamped limit/offset.
func paginateReverse[T any](items []T, limit, offset int) ([]T, int) {
	total := len(items)
	// Reverse to newest-first
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	// Apply offset
	if offset >= total {
		return []T{}, total
	}
	items = items[offset:]
	// Apply limit (0 = all)
	if limit > 0 && limit < len(items) {
		items = items[:limit]
	}
	return items, total
}

// handleListFeeds returns feed entries (newest-first, paginated).
func handleListFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := storage.LoadFeeds()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	limit, offset := parsePagination(r)
	page, total := paginateReverse(feeds, limit, offset)
	jsonResponse(w, http.StatusOK, PaginatedResponse{Items: page, Total: total, Limit: limit, Offset: offset})
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

func handleUpdateFeed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid feed ID"})
		return
	}
	var feed models.FeedEntry
	if err := json.NewDecoder(r.Body).Decode(&feed); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if feed.Type == "" || feed.Date == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
		return
	}
	log.Printf("Update Feed ID %d: %+v\n", id, feed)
	if err := storage.UpdateFeed(id, &feed); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	feed.ID = id
	jsonResponse(w, http.StatusOK, feed)
}

func handleDeleteFeed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid feed ID"})
		return
	}
	log.Printf("Delete Feed ID %d\n", id)
	if err := storage.DeleteFeed(id); err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}
