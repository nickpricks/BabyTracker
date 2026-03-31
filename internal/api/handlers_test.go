package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"babytracker/internal/config"
	"babytracker/internal/models"
	"babytracker/internal/storage"

	"github.com/gorilla/mux"
)

func testConfig() *config.Config {
	return &config.Config{APIPort: "8080", CORSOrigin: "http://localhost:3000"}
}

func testRouter(t *testing.T) *mux.Router {
	t.Helper()
	if err := storage.Init(t.TempDir()); err != nil {
		t.Fatalf("failed to init test storage: %v", err)
	}
	return SetupRouter(testConfig())
}

func TestHandleListFeeds_Empty(t *testing.T) {
	router := testRouter(t)
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
	router := testRouter(t)
	req := httptest.NewRequest("POST", "/api/feeds", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleLogFeed_MissingFields(t *testing.T) {
	router := testRouter(t)
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
	router := testRouter(t)
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
	router := testRouter(t)
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
	router := testRouter(t)
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
	router := testRouter(t)
	req := httptest.NewRequest("GET", "/api/feeds/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestCORSHeaders(t *testing.T) {
	router := testRouter(t)
	// Test CORS on a regular GET request (OPTIONS routing depends on mux config)
	req := httptest.NewRequest("GET", "/api/feeds", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("expected CORS header http://localhost:3000, got %s", got)
	}
}

func TestListEndpoints_ReturnJSON(t *testing.T) {
	router := testRouter(t)

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
