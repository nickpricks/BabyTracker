package models

import "testing"

func TestSleepEntry_IsNap(t *testing.T) {
	s := SleepEntry{Type: SleepTypeNap}
	if !s.IsNap() {
		t.Error("expected IsNap to be true")
	}
	s.Type = SleepTypeNight
	if s.IsNap() {
		t.Error("expected IsNap to be false for Night")
	}
}

func TestSleepEntry_IsNightSleep(t *testing.T) {
	s := SleepEntry{Type: SleepTypeNight}
	if !s.IsNightSleep() {
		t.Error("expected IsNightSleep to be true")
	}
	s.Type = SleepTypeNap
	if s.IsNightSleep() {
		t.Error("expected IsNightSleep to be false for Nap")
	}
}
