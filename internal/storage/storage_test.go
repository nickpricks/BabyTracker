package storage

import (
	"os"
	"testing"
	"time"

	"babytracker/internal/models"
)

func setupTestStorage(t *testing.T) *StorageManager {
	t.Helper()
	dir := t.TempDir()
	return &StorageManager{dataDir: dir}
}

func TestNextID(t *testing.T) {
	if got := nextID([]int{}); got != 1 {
		t.Errorf("nextID(empty) = %d, want 1", got)
	}
	if got := nextID([]int{1, 2, 3}); got != 4 {
		t.Errorf("nextID([1,2,3]) = %d, want 4", got)
	}
	if got := nextID([]int{5, 2, 8}); got != 9 {
		t.Errorf("nextID([5,2,8]) = %d, want 9", got)
	}
}

func TestFeedRoundTrip(t *testing.T) {
	sm := setupTestStorage(t)
	// Override global for this test
	origGlobal := globalStorage
	globalStorage = sm
	defer func() { globalStorage = origGlobal }()

	feed := &models.FeedEntry{
		Date:     "2025-06-22",
		Time:     time.Now(),
		Type:     models.FeedTypeBottle,
		Quantity: 120.0,
		Notes:    "test feed",
	}

	if err := SaveFeed(feed); err != nil {
		t.Fatalf("SaveFeed failed: %v", err)
	}
	if feed.ID != 1 {
		t.Errorf("expected ID 1, got %d", feed.ID)
	}

	feeds, err := LoadFeeds()
	if err != nil {
		t.Fatalf("LoadFeeds failed: %v", err)
	}
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}
	if feeds[0].Type != models.FeedTypeBottle {
		t.Errorf("expected type %s, got %s", models.FeedTypeBottle, feeds[0].Type)
	}
	if feeds[0].Quantity != 120.0 {
		t.Errorf("expected quantity 120, got %f", feeds[0].Quantity)
	}
}

func TestSleepRoundTrip(t *testing.T) {
	sm := setupTestStorage(t)
	origGlobal := globalStorage
	globalStorage = sm
	defer func() { globalStorage = origGlobal }()

	entry := &models.SleepEntry{
		Date:    "2025-06-22",
		Type:    models.SleepTypeNap,
		Quality: models.SleepQualityGood,
		Notes:   "test nap",
	}

	if err := SaveSleep(entry); err != nil {
		t.Fatalf("SaveSleep failed: %v", err)
	}
	if entry.ID != 1 {
		t.Errorf("expected ID 1, got %d", entry.ID)
	}

	entries, err := LoadSleep()
	if err != nil {
		t.Fatalf("LoadSleep failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Type != models.SleepTypeNap {
		t.Errorf("expected type Nap, got %s", entries[0].Type)
	}
}

func TestGrowthRoundTrip(t *testing.T) {
	sm := setupTestStorage(t)
	origGlobal := globalStorage
	globalStorage = sm
	defer func() { globalStorage = origGlobal }()

	entry := &models.GrowthEntry{
		Date:              "2025-06-22",
		Weight:            4.5,
		Height:            55.0,
		HeadCircumference: 36.0,
	}

	if err := SaveGrowth(entry); err != nil {
		t.Fatalf("SaveGrowth failed: %v", err)
	}
	if entry.ID != 1 {
		t.Errorf("expected ID 1, got %d", entry.ID)
	}

	entries, err := LoadGrowth()
	if err != nil {
		t.Fatalf("LoadGrowth failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Weight != 4.5 {
		t.Errorf("expected weight 4.5, got %f", entries[0].Weight)
	}
}

func TestDiaperRoundTrip(t *testing.T) {
	sm := setupTestStorage(t)
	origGlobal := globalStorage
	globalStorage = sm
	defer func() { globalStorage = origGlobal }()

	entry := &models.DiaperEntry{
		Date:  "2025-06-22",
		Time:  time.Now(),
		Type:  models.DiaperTypeWet,
		Notes: "test change",
	}

	if err := SaveDiaper(entry); err != nil {
		t.Fatalf("SaveDiaper failed: %v", err)
	}
	if entry.ID != 1 {
		t.Errorf("expected ID 1, got %d", entry.ID)
	}

	entries, err := LoadDiapers()
	if err != nil {
		t.Fatalf("LoadDiapers failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Type != models.DiaperTypeWet {
		t.Errorf("expected type Wet, got %s", entries[0].Type)
	}
}

func TestLoadFromEmptyDir(t *testing.T) {
	sm := setupTestStorage(t)
	origGlobal := globalStorage
	globalStorage = sm
	defer func() { globalStorage = origGlobal }()

	feeds, err := LoadFeeds()
	if err != nil {
		t.Fatalf("LoadFeeds from empty dir failed: %v", err)
	}
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds, got %d", len(feeds))
	}
}

func TestMultipleSavesIncrementID(t *testing.T) {
	sm := setupTestStorage(t)
	origGlobal := globalStorage
	globalStorage = sm
	defer func() { globalStorage = origGlobal }()

	for i := 1; i <= 3; i++ {
		feed := &models.FeedEntry{Date: "2025-06-22", Type: models.FeedTypeBottle}
		if err := SaveFeed(feed); err != nil {
			t.Fatalf("SaveFeed %d failed: %v", i, err)
		}
		if feed.ID != i {
			t.Errorf("feed %d: expected ID %d, got %d", i, i, feed.ID)
		}
	}

	feeds, err := LoadFeeds()
	if err != nil {
		t.Fatalf("LoadFeeds failed: %v", err)
	}
	if len(feeds) != 3 {
		t.Errorf("expected 3 feeds, got %d", len(feeds))
	}
}

func TestGetDataDirectory(t *testing.T) {
	sm := setupTestStorage(t)
	origGlobal := globalStorage
	globalStorage = sm
	defer func() { globalStorage = origGlobal }()

	dir, err := GetDataDirectory()
	if err != nil {
		t.Fatalf("GetDataDirectory failed: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("data directory does not exist: %s", dir)
	}
}
