// Package storage handles all data persistence operations for the Baby Tracker application.
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"babytracker/internal/models"
)

// StorageManager handles all data persistence operations.
type StorageManager struct {
	dataDir string
	mu      sync.Mutex
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
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	return &StorageManager{dataDir: dataDir}, nil
}

var (
	globalStorage *StorageManager
	initOnce      sync.Once
	initErr       error
)

// Init initializes the global storage with a specific data directory.
// Call this early in main() before any Save/Load calls.
// If never called, the default (~/.babytracker) is used.
func Init(dataDir string) error {
	sm, err := NewStorageManagerWithDir(dataDir)
	if err != nil {
		return err
	}
	globalStorage = sm
	initOnce.Do(func() {}) // mark as done so getStorage won't re-init
	return nil
}

func getStorage() (*StorageManager, error) {
	if globalStorage != nil {
		return globalStorage, nil
	}
	initOnce.Do(func() {
		globalStorage, initErr = NewStorageManager()
	})
	if initErr != nil {
		return nil, initErr
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
	// Atomic write: temp file + rename prevents corruption on crash (FINDING-25)
	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("failed to rename %s: %w", filename, err)
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
	sm.mu.Lock()
	defer sm.mu.Unlock()
	feeds, err := loadJSON[models.FeedEntry](sm, "feeds.json")
	if err != nil {
		return fmt.Errorf("refusing to save over unreadable data file: %w", err)
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

func UpdateFeed(id int, updated *models.FeedEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	feeds, err := loadJSON[models.FeedEntry](sm, "feeds.json")
	if err != nil {
		return fmt.Errorf("refusing to update over unreadable data file: %w", err)
	}
	for i, f := range feeds {
		if f.ID == id {
			updated.ID = id
			feeds[i] = *updated
			return saveJSON(sm, "feeds.json", feeds)
		}
	}
	return fmt.Errorf("feed with ID %d not found", id)
}

func DeleteFeed(id int) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	feeds, err := loadJSON[models.FeedEntry](sm, "feeds.json")
	if err != nil {
		return fmt.Errorf("refusing to delete from unreadable data file: %w", err)
	}
	for i, f := range feeds {
		if f.ID == id {
			feeds = append(feeds[:i], feeds[i+1:]...)
			return saveJSON(sm, "feeds.json", feeds)
		}
	}
	return fmt.Errorf("feed with ID %d not found", id)
}

// --- Sleep ---

func SaveSleep(entry *models.SleepEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.SleepEntry](sm, "sleep.json")
	if err != nil {
		return fmt.Errorf("refusing to save over unreadable data file: %w", err)
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

func UpdateSleep(id int, updated *models.SleepEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.SleepEntry](sm, "sleep.json")
	if err != nil {
		return fmt.Errorf("refusing to update over unreadable data file: %w", err)
	}
	for i, e := range entries {
		if e.ID == id {
			updated.ID = id
			entries[i] = *updated
			return saveJSON(sm, "sleep.json", entries)
		}
	}
	return fmt.Errorf("sleep entry with ID %d not found", id)
}

func DeleteSleep(id int) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.SleepEntry](sm, "sleep.json")
	if err != nil {
		return fmt.Errorf("refusing to delete from unreadable data file: %w", err)
	}
	for i, e := range entries {
		if e.ID == id {
			entries = append(entries[:i], entries[i+1:]...)
			return saveJSON(sm, "sleep.json", entries)
		}
	}
	return fmt.Errorf("sleep entry with ID %d not found", id)
}

// --- Growth ---

func SaveGrowth(entry *models.GrowthEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.GrowthEntry](sm, "growth.json")
	if err != nil {
		return fmt.Errorf("refusing to save over unreadable data file: %w", err)
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

func UpdateGrowth(id int, updated *models.GrowthEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.GrowthEntry](sm, "growth.json")
	if err != nil {
		return fmt.Errorf("refusing to update over unreadable data file: %w", err)
	}
	for i, e := range entries {
		if e.ID == id {
			updated.ID = id
			entries[i] = *updated
			return saveJSON(sm, "growth.json", entries)
		}
	}
	return fmt.Errorf("growth entry with ID %d not found", id)
}

func DeleteGrowth(id int) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.GrowthEntry](sm, "growth.json")
	if err != nil {
		return fmt.Errorf("refusing to delete from unreadable data file: %w", err)
	}
	for i, e := range entries {
		if e.ID == id {
			entries = append(entries[:i], entries[i+1:]...)
			return saveJSON(sm, "growth.json", entries)
		}
	}
	return fmt.Errorf("growth entry with ID %d not found", id)
}

// --- Diapers ---

func SaveDiaper(entry *models.DiaperEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.DiaperEntry](sm, "diapers.json")
	if err != nil {
		return fmt.Errorf("refusing to save over unreadable data file: %w", err)
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

func UpdateDiaper(id int, updated *models.DiaperEntry) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.DiaperEntry](sm, "diapers.json")
	if err != nil {
		return fmt.Errorf("refusing to update over unreadable data file: %w", err)
	}
	for i, e := range entries {
		if e.ID == id {
			updated.ID = id
			entries[i] = *updated
			return saveJSON(sm, "diapers.json", entries)
		}
	}
	return fmt.Errorf("diaper entry with ID %d not found", id)
}

func DeleteDiaper(id int) error {
	sm, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	entries, err := loadJSON[models.DiaperEntry](sm, "diapers.json")
	if err != nil {
		return fmt.Errorf("refusing to delete from unreadable data file: %w", err)
	}
	for i, e := range entries {
		if e.ID == id {
			entries = append(entries[:i], entries[i+1:]...)
			return saveJSON(sm, "diapers.json", entries)
		}
	}
	return fmt.Errorf("diaper entry with ID %d not found", id)
}

// GetDataDirectory returns the directory where data files are stored.
func GetDataDirectory() (string, error) {
	sm, err := getStorage()
	if err != nil {
		return "", err
	}
	return sm.dataDir, nil
}
