package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear env vars that could affect the test
	os.Unsetenv("PORT")
	os.Unsetenv("DATA_DIR")
	os.Unsetenv("APP_TITLE")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.APIPort != DefaultAPIPort {
		t.Errorf("APIPort = %q, want %q", cfg.APIPort, DefaultAPIPort)
	}
	if cfg.AppTitle != DefaultAppTitle {
		t.Errorf("AppTitle = %q, want %q", cfg.AppTitle, DefaultAppTitle)
	}
	if cfg.DataDir == "" {
		t.Error("DataDir should not be empty")
	}
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATA_DIR", "/tmp/test-babytracker")
	t.Setenv("APP_TITLE", "Test Tracker")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.APIPort != "9090" {
		t.Errorf("APIPort = %q, want %q", cfg.APIPort, "9090")
	}
	if cfg.DataDir != "/tmp/test-babytracker" {
		t.Errorf("DataDir = %q, want %q", cfg.DataDir, "/tmp/test-babytracker")
	}
	if cfg.AppTitle != "Test Tracker" {
		t.Errorf("AppTitle = %q, want %q", cfg.AppTitle, "Test Tracker")
	}
}
