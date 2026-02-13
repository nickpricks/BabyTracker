package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"babytracker/internal/models"
	"babytracker/internal/storage"
)

func setupTestRouter(t *testing.T) {
	t.Helper()
	// Override storage to use temp dir
	sm, err := storage.NewStorageManager()
	if err != nil {
		t.Fatalf("failed to create storage manager: %v", err)
	}
	_ = sm // storage uses global; we rely on temp dir via env or default
}

func TestHandleListFeeds_Empty(t *testing.T) {
	router := SetupRouter()
	req := httptest.NewRequest("GET", "/api/feeds", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestHandleLogFeed_InvalidJSON(t *testing.T) {
	router := SetupRouter()
	req := httptest.NewRequest("POST", "/api/feeds", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleLogFeed_MissingFields(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(models.FeedEntry{})
	req := httptest.NewRequest("POST", "/api/feeds", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleLogSleep_MissingFields(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(models.SleepEntry{})
	req := httptest.NewRequest("POST", "/api/sleep", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleLogGrowth_MissingDate(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(models.GrowthEntry{})
	req := httptest.NewRequest("POST", "/api/growth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleLogDiaper_MissingFields(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(models.DiaperEntry{})
	req := httptest.NewRequest("POST", "/api/diapers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleGetFeed_NotFound(t *testing.T) {
	router := SetupRouter()
	req := httptest.NewRequest("GET", "/api/feeds/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestCORSHeaders(t *testing.T) {
	router := SetupRouter()
	// Test CORS on a regular GET request (OPTIONS routing depends on mux config)
	req := httptest.NewRequest("GET", "/api/feeds", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header Access-Control-Allow-Origin: *")
	}
}

func TestListEndpoints_ReturnJSON(t *testing.T) {
	router := SetupRouter()

	endpoints := []string{"/api/feeds", "/api/sleep", "/api/growth", "/api/diapers"}
	for _, ep := range endpoints {
		req := httptest.NewRequest("GET", ep, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d", ep, w.Code)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("GET %s: expected application/json, got %s", ep, ct)
		}
	}
}
