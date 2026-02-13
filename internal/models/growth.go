package models

// GrowthEntry represents a single growth measurement record.
type GrowthEntry struct {
	ID                int     `json:"id"`
	Date              string  `json:"date"`      // YYYY-MM-DD
	Weight            float64 `json:"weight"`    // kg
	Height            float64 `json:"height"`    // cm
	HeadCircumference float64 `json:"head_circ"` // cm
	Notes             string  `json:"notes"`
}

// HasWeight checks if weight was recorded.
func (g *GrowthEntry) HasWeight() bool {
	return g.Weight > 0
}

// HasHeight checks if height was recorded.
func (g *GrowthEntry) HasHeight() bool {
	return g.Height > 0
}

// HasHeadCircumference checks if head circumference was recorded.
func (g *GrowthEntry) HasHeadCircumference() bool {
	return g.HeadCircumference > 0
}
