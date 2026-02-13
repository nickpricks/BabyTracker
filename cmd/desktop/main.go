package main

import (
	"log"

	"babytracker/internal/config"
	"babytracker/internal/desktop"
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

	app := desktop.NewApp()
	if app == nil {
		log.Fatal("Failed to initialize Baby Tracker application")
	}

	log.Printf("Starting Baby Tracker (data: %s)", cfg.DataDir)
	app.Run()
	log.Println("Baby Tracker closed.")
}
