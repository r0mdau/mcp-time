package types

// TimeResult represents the current time information for a specific timezone.
// It matches the Python MCP example output format.
type TimeResult struct {
	Timezone  string `json:"timezone"`
	Datetime  string `json:"datetime"`
	DayOfWeek string `json:"day_of_week"`
	IsDst     bool   `json:"is_dst"`
}

// TimeConversionResult represents a time conversion between two timezones.
type TimeConversionResult struct {
	Source         TimeResult `json:"source"`
	Target         TimeResult `json:"target"`
	TimeDifference string     `json:"time_difference"`
}

// GetCurrentTimeInput represents the input parameters for the get_current_time tool.
type GetCurrentTimeInput struct {
	Timezone string `json:"timezone"`
}

// ConvertTimeInput represents the input parameters for the convert_time tool.
// Time is expected in HH:MM (24-hour) format.
type ConvertTimeInput struct {
	SourceTimezone string `json:"source_timezone"`
	Time           string `json:"time"` // expected HH:MM
	TargetTimezone string `json:"target_timezone"`
}
