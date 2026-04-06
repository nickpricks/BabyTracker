package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlexTimeUnmarshalRFC3339(t *testing.T) {
	var ft FlexTime
	err := json.Unmarshal([]byte(`"2026-04-06T10:30:00Z"`), &ft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ft.Hour() != 10 || ft.Minute() != 30 {
		t.Errorf("expected 10:30, got %d:%d", ft.Hour(), ft.Minute())
	}
}

func TestFlexTimeUnmarshalNoTimezone(t *testing.T) {
	var ft FlexTime
	err := json.Unmarshal([]byte(`"2026-04-06T10:30:00"`), &ft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ft.Hour() != 10 || ft.Minute() != 30 {
		t.Errorf("expected 10:30, got %d:%d", ft.Hour(), ft.Minute())
	}
}

func TestFlexTimeUnmarshalShortFormat(t *testing.T) {
	var ft FlexTime
	err := json.Unmarshal([]byte(`"2026-04-06T10:30"`), &ft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ft.Hour() != 10 || ft.Minute() != 30 {
		t.Errorf("expected 10:30, got %d:%d", ft.Hour(), ft.Minute())
	}
}

func TestFlexTimeUnmarshalEmpty(t *testing.T) {
	var ft FlexTime
	err := json.Unmarshal([]byte(`""`), &ft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ft.IsZero() {
		t.Error("expected zero time for empty string")
	}
}

func TestFlexTimeUnmarshalInvalid(t *testing.T) {
	var ft FlexTime
	err := json.Unmarshal([]byte(`"not-a-time"`), &ft)
	if err == nil {
		t.Error("expected error for invalid time string")
	}
}

func TestFlexTimeMarshal(t *testing.T) {
	ft := FlexTime{Time: time.Date(2026, 4, 6, 10, 30, 0, 0, time.UTC)}
	data, err := json.Marshal(ft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `"2026-04-06T10:30:00Z"`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestFlexTimeMarshalZero(t *testing.T) {
	ft := FlexTime{}
	data, err := json.Marshal(ft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `""` {
		t.Errorf("expected empty string for zero time, got %s", string(data))
	}
}

func TestFlexTimeRoundTrip(t *testing.T) {
	original := FlexTime{Time: time.Date(2026, 4, 6, 14, 45, 0, 0, time.UTC)}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var decoded FlexTime
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !original.Time.Equal(decoded.Time) {
		t.Errorf("round-trip failed: %v != %v", original.Time, decoded.Time)
	}
}
