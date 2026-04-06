package models

import (
	"fmt"
	"strings"
	"time"
)

// FlexTime wraps time.Time with flexible JSON unmarshaling.
// Accepts both RFC3339 ("2006-01-02T15:04:05Z") and timezone-less ("2006-01-02T15:04:05") formats.
type FlexTime struct {
	time.Time
}

var flexTimeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
}

func (ft *FlexTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		ft.Time = time.Time{}
		return nil
	}
	for _, layout := range flexTimeFormats {
		if t, err := time.Parse(layout, s); err == nil {
			ft.Time = t
			return nil
		}
	}
	return fmt.Errorf("FlexTime: cannot parse %q (expected RFC3339 or 2006-01-02T15:04:05)", s)
}

func (ft FlexTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf("%q", ft.Time.Format(time.RFC3339))), nil
}
