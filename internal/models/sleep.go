package models

import "time"

// SleepEntry represents a single sleep session record.
type SleepEntry struct {
	ID        int       `json:"id"`
	Date      string    `json:"date"`       // YYYY-MM-DD
	StartTime time.Time `json:"start_time"` // When sleep began
	EndTime   time.Time `json:"end_time"`   // When sleep ended
	Duration  int       `json:"duration"`   // Duration in minutes
	Type      string    `json:"type"`       // nap, night
	Quality   string    `json:"quality"`    // good, fair, poor
	Notes     string    `json:"notes"`
}

// Sleep type constants
const (
	SleepTypeNap   = "Nap"
	SleepTypeNight = "Night"
)

// Sleep quality constants
const (
	SleepQualityGood = "Good"
	SleepQualityFair = "Fair"
	SleepQualityPoor = "Poor"
)

// IsNap checks if this is a nap entry.
func (s *SleepEntry) IsNap() bool {
	return s.Type == SleepTypeNap
}

// IsNightSleep checks if this is a night sleep entry.
func (s *SleepEntry) IsNightSleep() bool {
	return s.Type == SleepTypeNight
}
