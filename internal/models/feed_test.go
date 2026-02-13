package models

import "testing"

func TestFeedEntry_IsBottleFeed(t *testing.T) {
	f := FeedEntry{Type: FeedTypeBottle}
	if !f.IsBottleFeed() {
		t.Error("expected IsBottleFeed to be true for Bottle type")
	}
	f.Type = FeedTypeBreastLeft
	if f.IsBottleFeed() {
		t.Error("expected IsBottleFeed to be false for Breast type")
	}
}

func TestFeedEntry_IsBreastFeed(t *testing.T) {
	for _, ft := range []string{FeedTypeBreastLeft, FeedTypeBreastRight, FeedTypeBreastBoth} {
		f := FeedEntry{Type: ft}
		if !f.IsBreastFeed() {
			t.Errorf("expected IsBreastFeed to be true for %s", ft)
		}
	}
	f := FeedEntry{Type: FeedTypeBottle}
	if f.IsBreastFeed() {
		t.Error("expected IsBreastFeed to be false for Bottle type")
	}
}

func TestFeedEntry_HasQuantity(t *testing.T) {
	f := FeedEntry{Type: FeedTypeBottle}
	if !f.HasQuantity() {
		t.Error("expected HasQuantity to be true for Bottle")
	}
	f.Type = FeedTypeSolid
	if !f.HasQuantity() {
		t.Error("expected HasQuantity to be true for Solid")
	}
	f.Type = FeedTypeBreastLeft
	if f.HasQuantity() {
		t.Error("expected HasQuantity to be false for Breast")
	}
}
