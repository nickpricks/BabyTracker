// Package models defines the data structures used throughout the Baby Tracker application.
// These models represent the core entities and provide a consistent data layer
// for storage, retrieval, and manipulation of baby tracking information.
package models

import (
	"time"
)

// FeedEntry represents a single feeding session record
// This structure captures all relevant information about a baby's feeding,
// including timing, quantity, type, and contextual notes for comprehensive tracking.
type FeedEntry struct {
	ID       int       `json:"id"`       // Unique identifier for database storage
	Date     string    `json:"date"`     // Date of the feeding (YYYY-MM-DD)
	Time     time.Time `json:"time"`     // When the feeding occurred
	Type     string    `json:"type"`     // Type of feed (bottle, breast, solid)
	Quantity float64   `json:"quantity"` // Amount consumed (ml/oz), 0 if not applicable
	Notes    string    `json:"notes"`    // Additional observations or comments
	Duration int       `json:"duration"` // Feeding duration in minutes (for breastfeeding)
}

// FeedType constants for consistent feed categorization
// These constants ensure data consistency and enable proper analytics
const (
	FeedTypeBottle      = "Bottle"
	FeedTypeBreastLeft  = "Breast (Left)"
	FeedTypeBreastRight = "Breast (Right)"
	FeedTypeBreastBoth  = "Breast (Both)"
	FeedTypeSolid       = "Solid Food"
)

// IsBottleFeed checks if the feed type involves a bottle
// Helper method for analytics and quantity validation
func (f *FeedEntry) IsBottleFeed() bool {
	return f.Type == FeedTypeBottle
}

// IsBreastFeed checks if the feed type involves breastfeeding
// Helper method for duration tracking and analytics
func (f *FeedEntry) IsBreastFeed() bool {
	return f.Type == FeedTypeBreastLeft || f.Type == FeedTypeBreastRight || f.Type == FeedTypeBreastBoth
}

// HasQuantity checks if this feed type typically has a measurable quantity
// Useful for form validation and data display logic
func (f *FeedEntry) HasQuantity() bool {
	return f.IsBottleFeed() || f.Type == FeedTypeSolid
}
