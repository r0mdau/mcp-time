package timeutil

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/r0mdau/mcp-time/internal/timezone"
	"github.com/r0mdau/mcp-time/internal/types"
)

// FormatTimeDifference formats the hour difference between two timezones
// Returns format like "+5.0h" for integer hours or "+5.75h" for fractional
func FormatTimeDifference(offsetSource, offsetTarget int) string {
	hoursDifference := float64(offsetTarget-offsetSource) / 3600.0

	if math.Mod(hoursDifference, 1.0) == 0 {
		return fmt.Sprintf("%+.1fh", hoursDifference)
	}

	s := fmt.Sprintf("%+.2f", hoursDifference)
	// Strip trailing zeros and possible trailing dot
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s + "h"
}

// BuildTimeResult creates a TimeResult from a time.Time
func BuildTimeResult(t time.Time, tz string) types.TimeResult {
	return types.TimeResult{
		Timezone:  tz,
		Datetime:  timezone.FormatISOSeconds(t),
		DayOfWeek: t.Weekday().String(),
		IsDst:     timezone.IsDST(t),
	}
}

// ParseTimeInput parses HH:MM format time string
func ParseTimeInput(timeStr string) (hour, minute int, err error) {
	parsed, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid time format. Expected HH:MM [24-hour format]")
	}
	return parsed.Hour(), parsed.Minute(), nil
}

// ValidateConvertTimeInput validates the ConvertTimeInput fields
func ValidateConvertTimeInput(input types.ConvertTimeInput) error {
	if input.SourceTimezone == "" {
		return fmt.Errorf("source_timezone is required")
	}
	if input.TargetTimezone == "" {
		return fmt.Errorf("target_timezone is required")
	}
	if input.Time == "" {
		return fmt.Errorf("time is required")
	}
	return nil
}
