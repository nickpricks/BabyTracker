package models

import "time"

// DiaperEntry represents a single diaper change record.
type DiaperEntry struct {
	ID    int       `json:"id"`
	Date  string    `json:"date"` // YYYY-MM-DD
	Time  time.Time `json:"time"` // When the change occurred
	Type  string    `json:"type"` // wet, dirty, mixed
	Notes string    `json:"notes"`
}

// Diaper type constants
const (
	DiaperTypeWet   = "Wet"
	DiaperTypeDirty = "Dirty"
	DiaperTypeMixed = "Mixed"
)

// IsWet checks if this was a wet diaper.
func (d *DiaperEntry) IsWet() bool {
	return d.Type == DiaperTypeWet || d.Type == DiaperTypeMixed
}

// IsDirty checks if this was a dirty diaper.
func (d *DiaperEntry) IsDirty() bool {
	return d.Type == DiaperTypeDirty || d.Type == DiaperTypeMixed
}
