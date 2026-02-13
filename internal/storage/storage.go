// Package storage handles all data persistence operations for the Baby Tracker application.
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"babytracker/internal/models"
)

// StorageManager handles all data persistence operations.
type StorageManager struct {
	dataDir string
}

// NewStorageManager creates a storage manager with the default data directory (~/.babytracker).
func NewStorageManager() (*StorageManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	return NewStorageManagerWithDir(filepath.Join(homeDir, ".babytracker"))
}

// NewStorageManagerWithDir creates a storage manager using the given directory.
func NewStorageManagerWithDir(dataDir string) (*StorageManager, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	return &StorageManager{dataDir: dataDir}, nil
}

var globalStorage *StorageManager

// Init initializes the global storage with a specific data directory.
// Call this early in main() before any Save/Load calls.
// If never called, the default (~/.babytracker) is used.
func Init(dataDir string) error {
	sm, err := NewStorageManagerWithDir(dataDir)
	if err != nil {
		return err
	}
	globalStorage = sm
	return nil
}

func getStorage() (*StorageManager, error) {
	if globalStorage == nil {
		var err error
		globalStorage, err = NewStorageManager()
		if err != nil {
			return nil, err
		}
	}
	return globalStorage, nil
}

// --- Generic JSON helpers ---

func loadJSON[T any](sm *StorageManager, filename string) ([]T, error) {
	filePath := filepath.Join(sm.dataDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []T{}, nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
	}
	return items, nil
}

func saveJSON[T any](sm *StorageManager, filename string, items []T) error {
	filePath := filepath.Join(sm.dataDir, filename)
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", filename, err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	return nil
}

// nextID returns max(existing IDs) + 1 to avoid collisions.
func nextID(ids []int) int {
	max := 0
	for _, id := range ids {
		if id > max {
			max = id
		}
	}
	return max + 1
}

// --- Feeds ---

func SaveFeed(feed *models.FeedEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	feeds, err := loadJSON[models.FeedEntry](sm, "feeds.json")
	if err != nil {
		feeds = []models.FeedEntry{}
	}
	ids := make([]int, len(feeds))
	for i, f := range feeds {
		ids[i] = f.ID
	}
	feed.ID = nextID(ids)
	feeds = append(feeds, *feed)
	return saveJSON(sm, "feeds.json", feeds)
}

func LoadFeeds() ([]models.FeedEntry, error) {
	sm, err := getStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	return loadJSON[models.FeedEntry](sm, "feeds.json")
}

// --- Sleep ---

func SaveSleep(entry *models.SleepEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	entries, err := loadJSON[models.SleepEntry](sm, "sleep.json")
	if err != nil {
		entries = []models.SleepEntry{}
	}
	ids := make([]int, len(entries))
	for i, e := range entries {
		ids[i] = e.ID
	}
	entry.ID = nextID(ids)
	entries = append(entries, *entry)
	return saveJSON(sm, "sleep.json", entries)
}

func LoadSleep() ([]models.SleepEntry, error) {
	sm, err := getStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	return loadJSON[models.SleepEntry](sm, "sleep.json")
}

// --- Growth ---

func SaveGrowth(entry *models.GrowthEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	entries, err := loadJSON[models.GrowthEntry](sm, "growth.json")
	if err != nil {
		entries = []models.GrowthEntry{}
	}
	ids := make([]int, len(entries))
	for i, e := range entries {
		ids[i] = e.ID
	}
	entry.ID = nextID(ids)
	entries = append(entries, *entry)
	return saveJSON(sm, "growth.json", entries)
}

func LoadGrowth() ([]models.GrowthEntry, error) {
	sm, err := getStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	return loadJSON[models.GrowthEntry](sm, "growth.json")
}

// --- Diapers ---

func SaveDiaper(entry *models.DiaperEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	entries, err := loadJSON[models.DiaperEntry](sm, "diapers.json")
	if err != nil {
		entries = []models.DiaperEntry{}
	}
	ids := make([]int, len(entries))
	for i, e := range entries {
		ids[i] = e.ID
	}
	entry.ID = nextID(ids)
	entries = append(entries, *entry)
	return saveJSON(sm, "diapers.json", entries)
}

func LoadDiapers() ([]models.DiaperEntry, error) {
	sm, err := getStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	return loadJSON[models.DiaperEntry](sm, "diapers.json")
}

// GetDataDirectory returns the directory where data files are stored.
func GetDataDirectory() (string, error) {
	sm, err := getStorage()
	if err != nil {
		return "", err
	}
	return sm.dataDir, nil
}
