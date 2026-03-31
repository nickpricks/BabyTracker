package api

import (
	"net/http"
	"strings"

	"babytracker/internal/config"

	"github.com/gorilla/mux"
)

// SetupRouter sets up the mux router, CORS middleware, and all API endpoints.
func SetupRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()

	r.Use(mux.CORSMethodMiddleware(r))

	// CORS middleware — use configured origin instead of wildcard (FINDING-03)
	corsOrigin := cfg.CORSOrigin
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	// Request body size limit — 1MB max (FINDING-08)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.Body = http.MaxBytesReader(w, req.Body, 1<<20)
			next.ServeHTTP(w, req)
		})
	})

	// API key auth middleware — skip if no key configured
	if cfg.APIKey != "" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
				if token != cfg.APIKey {
					http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, req)
			})
		})
	}

	// Feed endpoints
	r.HandleFunc("/api/feeds", handleListFeeds).Methods("GET")
	r.HandleFunc("/api/feeds", handleLogFeed).Methods("POST")
	r.HandleFunc("/api/feeds/{id:[0-9]+}", handleGetFeed).Methods("GET")

	// Sleep endpoints
	r.HandleFunc("/api/sleep", handleListSleep).Methods("GET")
	r.HandleFunc("/api/sleep", handleLogSleep).Methods("POST")
	r.HandleFunc("/api/sleep/{id:[0-9]+}", handleGetSleep).Methods("GET")

	// Growth endpoints
	r.HandleFunc("/api/growth", handleListGrowth).Methods("GET")
	r.HandleFunc("/api/growth", handleLogGrowth).Methods("POST")
	r.HandleFunc("/api/growth/{id:[0-9]+}", handleGetGrowth).Methods("GET")

	// Diaper endpoints
	r.HandleFunc("/api/diapers", handleListDiapers).Methods("GET")
	r.HandleFunc("/api/diapers", handleLogDiaper).Methods("POST")
	r.HandleFunc("/api/diapers/{id:[0-9]+}", handleGetDiaper).Methods("GET")

	return r
}
