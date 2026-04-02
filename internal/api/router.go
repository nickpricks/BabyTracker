package api

import (
	"net/http"
	"strings"

	"babytracker/internal/config"

	"github.com/gorilla/mux"
)

// SetupRouter sets up the mux router, CORS, auth, and all API endpoints.
// Returns an http.Handler (not *mux.Router) because CORS wraps the router
// to intercept OPTIONS preflight before mux's method matching rejects it.
func SetupRouter(cfg *config.Config) http.Handler {
	r := mux.NewRouter()

	// Request body size limit — 1MB max (FINDING-08)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.Body = http.MaxBytesReader(w, req.Body, 1<<20)
			next.ServeHTTP(w, req)
		})
	})

	// API key auth middleware — skip if no key configured, skip OPTIONS
	if cfg.APIKey != "" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method == "OPTIONS" {
					next.ServeHTTP(w, req)
					return
				}
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

	// CORS wraps the entire router so OPTIONS preflight is handled before
	// mux rejects it with 405 (routes only register GET/POST).
	// If configured origin is localhost, accept any localhost port for dev.
	return corsHandler(cfg.CORSOrigin, r)
}

func corsHandler(corsOrigin string, next http.Handler) http.Handler {
	isLocalhost := strings.HasPrefix(corsOrigin, "http://localhost")
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		allowedOrigin := corsOrigin
		if isLocalhost && strings.HasPrefix(origin, "http://localhost") {
			allowedOrigin = origin
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, req)
	})
}
