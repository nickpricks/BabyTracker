package models

import "testing"

func TestGrowthEntry_HasWeight(t *testing.T) {
	g := GrowthEntry{Weight: 3.5}
	if !g.HasWeight() {
		t.Error("expected HasWeight to be true")
	}
	g.Weight = 0
	if g.HasWeight() {
		t.Error("expected HasWeight to be false for 0")
	}
}

func TestGrowthEntry_HasHeight(t *testing.T) {
	g := GrowthEntry{Height: 50.0}
	if !g.HasHeight() {
		t.Error("expected HasHeight to be true")
	}
	g.Height = 0
	if g.HasHeight() {
		t.Error("expected HasHeight to be false for 0")
	}
}

func TestGrowthEntry_HasHeadCircumference(t *testing.T) {
	g := GrowthEntry{HeadCircumference: 35.0}
	if !g.HasHeadCircumference() {
		t.Error("expected HasHeadCircumference to be true")
	}
	g.HeadCircumference = 0
	if g.HasHeadCircumference() {
		t.Error("expected HasHeadCircumference to be false for 0")
	}
}
