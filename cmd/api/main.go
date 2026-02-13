package main

import (
	"log"
	"net/http"

	"babytracker/internal/api"
	"babytracker/internal/config"
	"babytracker/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := storage.Init(cfg.DataDir); err != nil {
		log.Fatalf("Failed to initialize storage at %s: %v", cfg.DataDir, err)
	}

	r := api.SetupRouter()
	log.Printf("Baby Tracker API server running on http://localhost:%s", cfg.APIPort)
	log.Printf("Data directory: %s", cfg.DataDir)
	log.Fatal(http.ListenAndServe(":"+cfg.APIPort, r))
}
