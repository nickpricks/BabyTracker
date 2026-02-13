package models

import "testing"

func TestDiaperEntry_IsWet(t *testing.T) {
	d := DiaperEntry{Type: DiaperTypeWet}
	if !d.IsWet() {
		t.Error("expected IsWet to be true for Wet")
	}
	d.Type = DiaperTypeMixed
	if !d.IsWet() {
		t.Error("expected IsWet to be true for Mixed")
	}
	d.Type = DiaperTypeDirty
	if d.IsWet() {
		t.Error("expected IsWet to be false for Dirty")
	}
}

func TestDiaperEntry_IsDirty(t *testing.T) {
	d := DiaperEntry{Type: DiaperTypeDirty}
	if !d.IsDirty() {
		t.Error("expected IsDirty to be true for Dirty")
	}
	d.Type = DiaperTypeMixed
	if !d.IsDirty() {
		t.Error("expected IsDirty to be true for Mixed")
	}
	d.Type = DiaperTypeWet
	if d.IsDirty() {
		t.Error("expected IsDirty to be false for Wet")
	}
}
