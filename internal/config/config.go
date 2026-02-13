// Package config provides centralized configuration for the Baby Tracker application.
// Values are read from environment variables with sensible defaults.
// Set values via a .env file (loaded by the caller) or export them directly.
package config

import (
	"os"
	"path/filepath"
)

// Config holds all application configuration.
type Config struct {
	APIPort  string // HTTP port for the API server
	DataDir  string // Directory for JSON data files
	AppTitle string // Desktop window title
}

// Default values
const (
	DefaultAPIPort  = "8080"
	DefaultDataDir  = ".babytracker"
	DefaultAppTitle = "Baby Tracker"
)

// Load reads configuration from environment variables, falling back to defaults.
//
// Supported environment variables:
//
//	PORT           - API server port (default: 8080)
//	DATA_DIR       - Absolute path for data storage (default: ~/.babytracker)
//	APP_TITLE      - Desktop window title (default: Baby Tracker)
func Load() (*Config, error) {
	cfg := &Config{
		APIPort:  envOr("PORT", DefaultAPIPort),
		AppTitle: envOr("APP_TITLE", DefaultAppTitle),
	}

	// Data directory: use DATA_DIR if set, otherwise ~/.babytracker
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		cfg.DataDir = dir
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cfg.DataDir = filepath.Join(homeDir, DefaultDataDir)
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
