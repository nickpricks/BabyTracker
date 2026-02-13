package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRouter sets up the mux router, CORS middleware, and all API endpoints.
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

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
