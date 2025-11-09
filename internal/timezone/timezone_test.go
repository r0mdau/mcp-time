package timezone

import (
	"testing"
	"time"
	_ "time/tzdata"
)

func TestFormatISOSeconds(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "UTC time",
			time:     time.Date(2025, 11, 9, 12, 30, 45, 0, time.UTC),
			expected: "2025-11-09T12:30:45+00:00",
		},
		{
			name:     "New York time with negative offset",
			time:     time.Date(2025, 11, 9, 12, 30, 45, 0, mustLoadLocation(t, "America/New_York")),
			expected: "2025-11-09T12:30:45-05:00",
		},
		{
			name:     "Paris time with positive offset",
			time:     time.Date(2025, 11, 9, 12, 30, 45, 0, mustLoadLocation(t, "Europe/Paris")),
			expected: "2025-11-09T12:30:45+01:00",
		},
		{
			name:     "Tokyo time with positive offset",
			time:     time.Date(2025, 7, 15, 14, 20, 30, 0, mustLoadLocation(t, "Asia/Tokyo")),
			expected: "2025-07-15T14:20:30+09:00",
		},
		{
			name:     "Sydney DST time",
			time:     time.Date(2025, 1, 15, 8, 45, 10, 0, mustLoadLocation(t, "Australia/Sydney")),
			expected: "2025-01-15T08:45:10+11:00",
		},
		{
			name:     "midnight UTC",
			time:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2025-01-01T00:00:00+00:00",
		},
		{
			name:     "end of day",
			time:     time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: "2025-12-31T23:59:59+00:00",
		},
		{
			name:     "time with nanoseconds truncated",
			time:     time.Date(2025, 6, 15, 10, 20, 30, 123456789, time.UTC),
			expected: "2025-06-15T10:20:30+00:00",
		},
		{
			name:     "DST summer time in Paris",
			time:     time.Date(2025, 7, 15, 14, 30, 45, 0, mustLoadLocation(t, "Europe/Paris")),
			expected: "2025-07-15T14:30:45+02:00",
		},
		{
			name:     "DST winter time in Paris",
			time:     time.Date(2025, 1, 15, 14, 30, 45, 0, mustLoadLocation(t, "Europe/Paris")),
			expected: "2025-01-15T14:30:45+01:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatISOSeconds(tt.time)
			if result != tt.expected {
				t.Errorf("FormatISOSeconds() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// mustLoadLocation is a helper that loads a timezone or fails the test
func mustLoadLocation(t *testing.T, name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("failed to load location %s: %v", name, err)
	}
	return loc
}

func TestConvertUTCtoParis(t *testing.T) {
	// 2025-11-09 is in CET (UTC+1)
	input := "2025-11-09T12:00:00Z"
	out, err := ConvertTimeString(input, "", "Europe/Paris")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Format("15:04") != "13:00" {
		t.Fatalf("expected 13:00, got %s", out.Format("15:04"))
	}
}

func TestParseNoOffsetWithFromTZ(t *testing.T) {
	input := "2025-11-09 12:00:00"
	out, err := ConvertTimeString(input, "America/New_York", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// New York is UTC-5 on 2025-11-09 => 17:00 UTC
	want := "2025-11-09T17:00:00Z"
	if out.Format(time.RFC3339) != want {
		t.Fatalf("expected %s, got %s", want, out.Format(time.RFC3339))
	}
}

func TestGetNowInLocationEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		tz        string
		wantError bool
	}{
		{"empty defaults to UTC", "", false},
		{"valid UTC", "UTC", false},
		{"valid timezone", "America/New_York", false},
		{"invalid timezone", "Invalid/Timezone", true},
		{"malformed timezone", "Not-A-Real-TZ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetNowInLocation(tt.tz)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result.IsZero() {
				t.Error("expected non-zero time")
			}
		})
	}
}

func TestConvertTimeStringEmptyInput(t *testing.T) {
	_, err := ConvertTimeString("", "UTC", "America/New_York")
	if err == nil {
		t.Error("expected error for empty time string")
	}
}

func TestConvertTimeStringInvalidFromTZ(t *testing.T) {
	input := "2025-11-09 12:00:00"
	_, err := ConvertTimeString(input, "Invalid/TZ", "UTC")
	if err == nil {
		t.Error("expected error for invalid from timezone")
	}
}

func TestConvertTimeStringInvalidToTZ(t *testing.T) {
	input := "2025-11-09T12:00:00Z"
	_, err := ConvertTimeString(input, "", "Invalid/TZ")
	if err == nil {
		t.Error("expected error for invalid to timezone")
	}
}

func TestConvertTimeStringInvalidFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid format", "not-a-time"},
		{"partial date", "2025-11-09"},
		{"only time", "12:00:00"},
		{"wrong format", "11/09/2025 12:00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConvertTimeString(tt.input, "UTC", "America/New_York")
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestConvertTimeStringWithNaiveTimeAndEmptyFromTZ(t *testing.T) {
	input := "2025-11-09 12:00:00"
	out, err := ConvertTimeString(input, "", "Europe/Paris")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should default to UTC when fromTZ is empty
	// UTC 12:00 -> Paris should be 13:00 or 14:00 depending on DST
	hour := out.Hour()
	if hour != 13 && hour != 14 {
		t.Errorf("expected hour 13 or 14, got %d", hour)
	}
}

func TestIsDST(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}
	summer := time.Date(2025, time.July, 1, 12, 0, 0, 0, loc)
	winter := time.Date(2025, time.January, 1, 12, 0, 0, 0, loc)

	if !IsDST(summer) {
		t.Fatalf("expected summer to be DST in Europe/Paris")
	}
	if IsDST(winter) {
		t.Fatalf("expected winter to NOT be DST in Europe/Paris")
	}
}

func TestIsDSTEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		tz      string
		month   time.Month
		wantDST bool
	}{
		{"NYC summer", "America/New_York", time.July, true},
		{"NYC winter", "America/New_York", time.January, false},
		{"UTC always no DST", "UTC", time.July, false},
		{"Sydney summer", "Australia/Sydney", time.January, true},
		{"Sydney winter", "Australia/Sydney", time.July, false},
		{"Tokyo no DST", "Asia/Tokyo", time.July, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := time.LoadLocation(tt.tz)
			if err != nil {
				t.Fatalf("failed to load location: %v", err)
			}
			testTime := time.Date(2025, tt.month, 15, 12, 0, 0, 0, loc)
			got := IsDST(testTime)
			if got != tt.wantDST {
				t.Errorf("IsDST(%s in %s) = %v, want %v", tt.month, tt.tz, got, tt.wantDST)
			}
		})
	}
}

func TestGetLocalTimezone(t *testing.T) {
	tests := []struct {
		name     string
		override string
		want     string
	}{
		{
			name:     "with override",
			override: "America/New_York",
			want:     "America/New_York",
		},
		{
			name:     "with UTC override",
			override: "UTC",
			want:     "UTC",
		},
		{
			name:     "empty override returns non-empty",
			override: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLocalTimezone(tt.override)
			if tt.want != "" && got != tt.want {
				t.Errorf("GetLocalTimezone(%q) = %v, want %v", tt.override, got, tt.want)
			}
			if tt.want == "" && got == "" {
				t.Error("GetLocalTimezone should never return empty string")
			}
		})
	}
}

func BenchmarkGetNowInLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetNowInLocation("America/New_York")
	}
}

func BenchmarkFormatISOSeconds(b *testing.B) {
	loc, _ := time.LoadLocation("America/New_York")
	testTime := time.Date(2025, 11, 9, 12, 30, 45, 123456789, loc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatISOSeconds(testTime)
	}
}

func BenchmarkConvertTimeString(b *testing.B) {
	input := "2025-11-09T12:00:00Z"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertTimeString(input, "", "Europe/Paris")
	}
}
