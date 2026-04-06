// Bench data generator for BabyTracker.
//
// WARNING: This replaces all data in ~/.babytracker/ with 10,000 entries per module.
// Existing data is backed up to ~/.babytracker/.backup/ before overwrite.
// Run `make bench-restore` to restore the backup.
//
// Usage:
//
//	make bench          # generate 10k entries (backs up existing data first)
//	make bench-restore  # restore backed-up data
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"babytracker/internal/models"
)

const count = 10_000

var dataDir = filepath.Join(os.Getenv("HOME"), ".babytracker")
var backupDir = filepath.Join(dataDir, ".backup")

func main() {
	fmt.Println("⚠️  BENCH DATA GENERATOR")
	fmt.Println("   This will REPLACE all data in", dataDir)
	fmt.Println("   Existing data backed up to", backupDir)
	fmt.Printf("   Generating %d entries per module (%d total)\n\n", count, count*4)

	if err := backup(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Backup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Backup complete")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	start := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	feeds := genFeeds(r, start)
	sleeps := genSleep(r, start)
	growths := genGrowth(r, start)
	diapers := genDiapers(r, start)

	writeJSON("feeds.json", feeds)
	writeJSON("sleep.json", sleeps)
	writeJSON("growth.json", growths)
	writeJSON("diapers.json", diapers)

	fmt.Printf("\n✅ Generated %d feeds, %d sleep, %d growth, %d diapers\n", len(feeds), len(sleeps), len(growths), len(diapers))
	fmt.Println("   Run `make bench-restore` to restore original data")
}

func backup() error {
	os.MkdirAll(backupDir, 0700)
	files := []string{"feeds.json", "sleep.json", "growth.json", "diapers.json"}
	for _, f := range files {
		src := filepath.Join(dataDir, f)
		dst := filepath.Join(backupDir, f)
		data, err := os.ReadFile(src)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if err := os.WriteFile(dst, data, 0600); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(name string, v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Marshal %s: %v\n", name, err)
		os.Exit(1)
	}
	path := filepath.Join(dataDir, name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Write %s: %v\n", name, err)
		os.Exit(1)
	}
	fmt.Printf("   wrote %s (%d bytes)\n", name, len(data))
}

// ── Generators ──────────────────────────────────────────────────

var feedTypes = []string{"Bottle", "Breast (Left)", "Breast (Right)", "Breast (Both)", "Solid Food"}
var feedNotes = []string{"", "", "", "Good latch", "Fussy at start", "Fell asleep during feed", "Very hungry", "Refused at first", ""}

func genFeeds(r *rand.Rand, start time.Time) []models.FeedEntry {
	entries := make([]models.FeedEntry, count)
	t := start
	for i := range entries {
		t = t.Add(time.Duration(90+r.Intn(180)) * time.Minute)
		ft := feedTypes[r.Intn(len(feedTypes))]
		var qty float64
		if ft == "Bottle" {
			qty = float64(60+r.Intn(140)) // 60-200ml
		}
		var dur int
		if ft != "Bottle" && ft != "Solid Food" {
			dur = 5 + r.Intn(25) // 5-30 min
		}
		entries[i] = models.FeedEntry{
			ID:       i + 1,
			Date:     t.Format("2006-01-02"),
			Time:     models.FlexTime{Time: t},
			Type:     ft,
			Quantity: qty,
			Notes:    feedNotes[r.Intn(len(feedNotes))],
			Duration: dur,
		}
	}
	return entries
}

var sleepTypes = []string{"Nap", "Night"}
var qualities = []string{"Good", "Fair", "Poor"}
var sleepNotes = []string{"", "", "", "Woke once", "Slept through", "Restless", "Needed rocking", "White noise helped", ""}

func genSleep(r *rand.Rand, start time.Time) []models.SleepEntry {
	entries := make([]models.SleepEntry, count)
	t := start
	for i := range entries {
		t = t.Add(time.Duration(3+r.Intn(6)) * time.Hour)
		st := sleepTypes[r.Intn(len(sleepTypes))]
		dur := 30 + r.Intn(360) // 30-390 min
		if st == "Nap" {
			dur = 20 + r.Intn(120) // 20-140 min
		}
		entries[i] = models.SleepEntry{
			ID:        i + 1,
			Date:      t.Format("2006-01-02"),
			StartTime: models.FlexTime{Time: t},
			EndTime:   models.FlexTime{Time: t.Add(time.Duration(dur) * time.Minute)},
			Duration:  dur,
			Type:      st,
			Quality:   qualities[r.Intn(len(qualities))],
			Notes:     sleepNotes[r.Intn(len(sleepNotes))],
		}
	}
	return entries
}

func genGrowth(r *rand.Rand, start time.Time) []models.GrowthEntry {
	entries := make([]models.GrowthEntry, count)
	weight := 3.2 + r.Float64()*0.5
	height := 49.0 + r.Float64()*3.0
	head := 34.0 + r.Float64()*1.5
	t := start
	for i := range entries {
		t = t.Add(time.Duration(4+r.Intn(10)) * time.Hour)
		weight += r.Float64() * 0.05
		height += r.Float64() * 0.03
		head += r.Float64() * 0.015
		entries[i] = models.GrowthEntry{
			ID:                i + 1,
			Date:              t.Format("2006-01-02"),
			Weight:            float64(int(weight*100)) / 100,
			Height:            float64(int(height*10)) / 10,
			HeadCircumference: float64(int(head*10)) / 10,
			Notes:             "",
		}
	}
	return entries
}

var diaperTypes = []string{"Wet", "Dirty", "Mixed"}
var diaperNotes = []string{"", "", "", "", "Blowout", "Rash looking better", "Applied cream", ""}

func genDiapers(r *rand.Rand, start time.Time) []models.DiaperEntry {
	entries := make([]models.DiaperEntry, count)
	t := start
	for i := range entries {
		t = t.Add(time.Duration(60+r.Intn(180)) * time.Minute)
		entries[i] = models.DiaperEntry{
			ID:    i + 1,
			Date:  t.Format("2006-01-02"),
			Time:  models.FlexTime{Time: t},
			Type:  diaperTypes[r.Intn(len(diaperTypes))],
			Notes: diaperNotes[r.Intn(len(diaperNotes))],
		}
	}
	return entries
}
