package timeutil

import (
	"strings"
	"testing"
	"time"

	"github.com/r0mdau/mcp-time/internal/types"
)

func TestFormatTimeDifference(t *testing.T) {
	tests := []struct {
		name            string
		offsetSource    int
		offsetTarget    int
		wantContains    string
		checkInteger    bool
		checkFractional bool
	}{
		{
			name:         "integer offset positive",
			offsetSource: 0,
			offsetTarget: 3600, // +1 hour
			wantContains: ".0h",
			checkInteger: true,
		},
		{
			name:         "integer offset negative",
			offsetSource: 3600,
			offsetTarget: 0,
			wantContains: "-1",
			checkInteger: true,
		},
		{
			name:            "fractional offset",
			offsetSource:    0,
			offsetTarget:    20700, // +5:45 (Nepal)
			wantContains:    "5.75h",
			checkFractional: true,
		},
		{
			name:         "zero offset",
			offsetSource: 0,
			offsetTarget: 0,
			wantContains: "0.0h",
			checkInteger: true,
		},
		{
			name:            "fractional negative",
			offsetSource:    20700,
			offsetTarget:    0,
			wantContains:    "-5.75h",
			checkFractional: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTimeDifference(tt.offsetSource, tt.offsetTarget)
			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("FormatTimeDifference() = %v, want to contain %v", got, tt.wantContains)
			}
			if tt.checkInteger && !strings.Contains(got, ".0h") {
				t.Errorf("expected integer format (.0h), got %v", got)
			}
			if tt.checkFractional && strings.Contains(got, ".0h") {
				t.Errorf("expected fractional format, got integer %v", got)
			}
		})
	}
}

func TestBuildTimeResult(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	testTime := time.Date(2025, 7, 15, 12, 30, 45, 0, loc)
	result := BuildTimeResult(testTime, "America/New_York")

	if result.Timezone != "America/New_York" {
		t.Errorf("expected timezone America/New_York, got %s", result.Timezone)
	}
	if result.DayOfWeek != "Tuesday" {
		t.Errorf("expected Tuesday, got %s", result.DayOfWeek)
	}
	if !strings.Contains(result.Datetime, "2025-07-15") {
		t.Errorf("expected date 2025-07-15 in datetime, got %s", result.Datetime)
	}
	// July in New York is DST
	if !result.IsDst {
		t.Error("expected DST to be true in July for New York")
	}
}

func TestParseTimeInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantHour  int
		wantMin   int
		wantError bool
	}{
		{"valid morning", "09:30", 9, 30, false},
		{"valid afternoon", "15:45", 15, 45, false},
		{"valid midnight", "00:00", 0, 0, false},
		{"valid end of day", "23:59", 23, 59, false},
		{"valid single digit hour", "9:30", 9, 30, false}, // Go's time.Parse accepts this
		{"invalid format", "25:00", 0, 0, true},
		{"invalid format alpha", "ab:cd", 0, 0, true},
		{"empty string", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, min, err := ParseTimeInput(tt.input)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if hour != tt.wantHour {
				t.Errorf("hour = %v, want %v", hour, tt.wantHour)
			}
			if min != tt.wantMin {
				t.Errorf("minute = %v, want %v", min, tt.wantMin)
			}
		})
	}
}

func TestValidateConvertTimeInput(t *testing.T) {
	tests := []struct {
		name      string
		input     types.ConvertTimeInput
		wantError bool
		errMsg    string
	}{
		{
			name: "valid input",
			input: types.ConvertTimeInput{
				SourceTimezone: "UTC",
				Time:           "12:00",
				TargetTimezone: "America/New_York",
			},
			wantError: false,
		},
		{
			name: "missing source timezone",
			input: types.ConvertTimeInput{
				SourceTimezone: "",
				Time:           "12:00",
				TargetTimezone: "America/New_York",
			},
			wantError: true,
			errMsg:    "source_timezone is required",
		},
		{
			name: "missing target timezone",
			input: types.ConvertTimeInput{
				SourceTimezone: "UTC",
				Time:           "12:00",
				TargetTimezone: "",
			},
			wantError: true,
			errMsg:    "target_timezone is required",
		},
		{
			name: "missing time",
			input: types.ConvertTimeInput{
				SourceTimezone: "UTC",
				Time:           "",
				TargetTimezone: "America/New_York",
			},
			wantError: true,
			errMsg:    "time is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConvertTimeInput(tt.input)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error = %v, want to contain %v", err, tt.errMsg)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func BenchmarkFormatTimeDifference(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatTimeDifference(0, 3600)
	}
}

func BenchmarkBuildTimeResult(b *testing.B) {
	loc, _ := time.LoadLocation("UTC")
	t := time.Now().In(loc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildTimeResult(t, "UTC")
	}
}
