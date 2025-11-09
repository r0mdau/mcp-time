package timezone

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

// FormatISOSeconds formats time to ISO 8601 with seconds precision
// matching Python's .isoformat(timespec="seconds")
func FormatISOSeconds(t time.Time) string {
	// Format: 2006-01-02T15:04:05-07:00 (RFC3339 with seconds only, no fractional seconds)
	return t.Format("2006-01-02T15:04:05-07:00")
}

// IsDST determines if time t is in DST for its location
func IsDST(t time.Time) bool {
	// Compare offsets in Jan and Jul to determine the standard offset
	year := t.Year()
	loc := t.Location()
	jan := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	jul := time.Date(year, time.July, 1, 0, 0, 0, 0, loc)
	_, offJan := jan.Zone()
	_, offJul := jul.Zone()
	standard := offJan
	if offJul < standard {
		standard = offJul
	}
	_, offNow := t.Zone()
	return offNow != standard
}

// GetNowInLocation returns the current time in the requested IANA location.
func GetNowInLocation(tz string) (time.Time, error) {
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown timezone %q: %w", tz, err)
	}
	return time.Now().In(loc), nil
}

// ConvertTimeString converts a time given as a string from a source timezone to a destination timezone.
// The function prefers RFC3339 input. If that fails, it will try the layout "2006-01-02 15:04:05" and
// interpret that naive timestamp in the provided fromTZ (or UTC if empty).
func ConvertTimeString(tstr, fromTZ, toTZ string) (time.Time, error) {
	if tstr == "" {
		return time.Time{}, fmt.Errorf("time string is empty")
	}

	// First try RFC3339 which includes an offset or Z
	parsed, err := time.Parse(time.RFC3339, tstr)
	if err != nil {
		// Try a common naive layout
		parsed, err = time.Parse("2006-01-02 15:04:05", tstr)
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to parse time %q: %w", tstr, err)
		}
		// Assign location based on fromTZ or default to UTC
		if fromTZ == "" {
			parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), parsed.Nanosecond(), time.UTC)
		} else {
			locFrom, err := time.LoadLocation(fromTZ)
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid from timezone %q: %w", fromTZ, err)
			}
			parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), parsed.Nanosecond(), locFrom)
		}
	}

	// Load destination timezone
	locTo, err := time.LoadLocation(toTZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid to timezone %q: %w", toTZ, err)
	}

	return parsed.In(locTo), nil
}

// GetLocalTimezone determines the local timezone to use.
// If override is provided, it returns that. Otherwise, it attempts to detect
// the system's IANA timezone name, falling back to UTC if detection fails.
func GetLocalTimezone(override string) string {
	// If override is provided, use it
	if override != "" {
		return override
	}

	// Try to get the local timezone name
	// On Unix systems, this will return the IANA timezone
	// On Windows or if it can't be determined, it may return "Local"
	localZone := time.Local.String()
	if localZone == "Local" {
		// Fallback to UTC if we can't determine the local timezone
		return "UTC"
	}
	return localZone
}
