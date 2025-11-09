package handlers

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/r0mdau/mcp-time/internal/types"
)

func TestGetCurrentTimeValid(t *testing.T) {
	input := types.GetCurrentTimeInput{Timezone: "UTC"}
	_, out, err := GetCurrentTime(context.Background(), nil, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Timezone != "UTC" {
		t.Fatalf("expected timezone UTC, got %s", out.Timezone)
	}
	// Datetime should parse as RFC3339
	parsed, err := time.Parse(time.RFC3339, out.Datetime)
	if err != nil {
		t.Fatalf("datetime not RFC3339: %v", err)
	}
	// Check DayOfWeek matches parsed time
	if parsed.Weekday().String() != out.DayOfWeek {
		t.Fatalf("day_of_week mismatch: expected %s got %s", parsed.Weekday().String(), out.DayOfWeek)
	}
}

func TestGetCurrentTimeInvalidTZ(t *testing.T) {
	input := types.GetCurrentTimeInput{Timezone: "Invalid/Zone"}
	_, _, err := GetCurrentTime(context.Background(), nil, input)
	if err == nil {
		t.Fatalf("expected error for invalid timezone")
	}
}

func TestConvertTimeValidFractionalOffset(t *testing.T) {
	// Use UTC -> Asia/Kathmandu (UTC+5:45) which has a fractional offset
	input := types.ConvertTimeInput{SourceTimezone: "UTC", Time: "00:00", TargetTimezone: "Asia/Kathmandu"}
	_, out, err := ConvertTime(context.Background(), nil, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Source.Timezone != "UTC" {
		t.Fatalf("unexpected source timezone: %s", out.Source.Timezone)
	}
	if out.Target.Timezone != "Asia/Kathmandu" {
		t.Fatalf("unexpected target timezone: %s", out.Target.Timezone)
	}
	// TimeDifference should end with 'h' and be parseable as float
	if !strings.HasSuffix(out.TimeDifference, "h") {
		t.Fatalf("time_difference missing 'h' suffix: %s", out.TimeDifference)
	}
	num := strings.TrimSuffix(out.TimeDifference, "h")
	// allow leading + sign
	num = strings.TrimPrefix(num, "+")
	f, ferr := strconv.ParseFloat(num, 64)
	if ferr != nil {
		t.Fatalf("unable to parse time_difference numeric part: %v", ferr)
	}
	// Kathmandu offset is +5.75 hours relative to UTC
	if f < 5.7 || f > 5.8 {
		t.Fatalf("unexpected fractional hours for Kathmandu: %v", f)
	}
}

func TestConvertTimeInvalidFormat(t *testing.T) {
	input := types.ConvertTimeInput{SourceTimezone: "UTC", Time: "badtime", TargetTimezone: "Europe/Paris"}
	_, _, err := ConvertTime(context.Background(), nil, input)
	if err == nil {
		t.Fatalf("expected error for invalid time format")
	}
}

func TestConvertTimeMissingFields(t *testing.T) {
	cases := []types.ConvertTimeInput{
		{SourceTimezone: "", Time: "12:00", TargetTimezone: "Europe/Paris"},
		{SourceTimezone: "UTC", Time: "", TargetTimezone: "Europe/Paris"},
		{SourceTimezone: "UTC", Time: "12:00", TargetTimezone: ""},
	}
	for i, tc := range cases {
		_, _, err := ConvertTime(context.Background(), nil, tc)
		if err == nil {
			t.Fatalf("case %d: expected error for missing field, got nil", i)
		}
	}
}

func TestConvertTimeIntegerOffset(t *testing.T) {
	// UTC -> Europe/Paris should be an integer hour offset (e.g., +1 or +2 depending on date)
	input := types.ConvertTimeInput{SourceTimezone: "UTC", Time: "00:00", TargetTimezone: "Europe/Paris"}
	_, out, err := ConvertTime(context.Background(), nil, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(out.TimeDifference, "h") {
		t.Fatalf("time_difference missing 'h' suffix: %s", out.TimeDifference)
	}
	// integer offsets are formatted with one decimal place (e.g., +1.0h)
	if !strings.Contains(out.TimeDifference, ".0h") {
		t.Fatalf("expected integer-hour formatting (.0h), got %s", out.TimeDifference)
	}
}

func TestConvertTimeInvalidTargetTimezone(t *testing.T) {
	input := types.ConvertTimeInput{SourceTimezone: "UTC", Time: "00:00", TargetTimezone: "Invalid/Zone"}
	_, _, err := ConvertTime(context.Background(), nil, input)
	if err == nil {
		t.Fatalf("expected error for invalid target timezone")
	}
}

func TestConvertTimeInvalidSourceTimezone(t *testing.T) {
	input := types.ConvertTimeInput{SourceTimezone: "Invalid/Zone", Time: "00:00", TargetTimezone: "UTC"}
	_, _, err := ConvertTime(context.Background(), nil, input)
	if err == nil {
		t.Fatalf("expected error for invalid source timezone")
	}
}

func TestGetCurrentTimeEmptyTimezone(t *testing.T) {
	input := types.GetCurrentTimeInput{Timezone: ""}
	_, out, err := GetCurrentTime(context.Background(), nil, input)
	if err != nil {
		t.Fatalf("unexpected error for empty timezone: %v", err)
	}
	if out.Timezone != "UTC" {
		t.Fatalf("expected default timezone UTC, got %s", out.Timezone)
	}
}

func TestRegisterTools(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-time-test", Version: "vtest"}, nil)
	// Should not panic
	RegisterTools(server, "UTC")
}

func TestConvertTimeDSTTransition(t *testing.T) {
	// Test conversion during DST transition period
	input := types.ConvertTimeInput{
		SourceTimezone: "America/New_York",
		Time:           "12:00",
		TargetTimezone: "Europe/London",
	}
	_, out, err := ConvertTime(context.Background(), nil, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify structure is complete
	if out.Source.Timezone == "" || out.Target.Timezone == "" {
		t.Error("timezone fields should not be empty")
	}
	if out.Source.Datetime == "" || out.Target.Datetime == "" {
		t.Error("datetime fields should not be empty")
	}
	if out.Source.DayOfWeek == "" || out.Target.DayOfWeek == "" {
		t.Error("day_of_week fields should not be empty")
	}
	if out.TimeDifference == "" {
		t.Error("time_difference should not be empty")
	}
}

func TestGetCurrentTimeMultipleTimezones(t *testing.T) {
	timezones := []string{
		"UTC",
		"America/New_York",
		"Europe/Paris",
		"Asia/Tokyo",
		"Australia/Sydney",
	}

	for _, tz := range timezones {
		t.Run(tz, func(t *testing.T) {
			input := types.GetCurrentTimeInput{Timezone: tz}
			_, out, err := GetCurrentTime(context.Background(), nil, input)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", tz, err)
			}
			if out.Timezone != tz {
				t.Errorf("expected timezone %s, got %s", tz, out.Timezone)
			}
			// Verify datetime can be parsed
			if _, err := time.Parse(time.RFC3339, out.Datetime); err != nil {
				t.Errorf("datetime not parseable for %s: %v", tz, err)
			}
		})
	}
}
