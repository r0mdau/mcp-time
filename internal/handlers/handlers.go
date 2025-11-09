package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/r0mdau/mcp-time/internal/timeutil"
	"github.com/r0mdau/mcp-time/internal/timezone"
	"github.com/r0mdau/mcp-time/internal/types"
)

// GetCurrentTime implements the get_current_time MCP tool handler.
// It returns the current time in the specified timezone.
func GetCurrentTime(ctx context.Context, req *mcp.CallToolRequest, input types.GetCurrentTimeInput) (
	*mcp.CallToolResult,
	types.TimeResult,
	error,
) {
	tz := input.Timezone
	if tz == "" {
		tz = "UTC"
	}
	now, err := timezone.GetNowInLocation(tz)
	if err != nil {
		// Return error for invalid timezone - SDK will handle it properly
		return nil, types.TimeResult{}, fmt.Errorf("invalid timezone: %w", err)
	}
	return nil, timeutil.BuildTimeResult(now, tz), nil
}

// ConvertTime implements the convert_time MCP tool handler.
// It converts a time specified in HH:MM format from one timezone to another.
func ConvertTime(ctx context.Context, req *mcp.CallToolRequest, input types.ConvertTimeInput) (
	*mcp.CallToolResult,
	types.TimeConversionResult,
	error,
) {
	// Validate input
	if err := timeutil.ValidateConvertTimeInput(input); err != nil {
		return nil, types.TimeConversionResult{}, err
	}

	// Parse time input
	hour, minute, err := timeutil.ParseTimeInput(input.Time)
	if err != nil {
		return nil, types.TimeConversionResult{}, err
	}

	// Get source location and build source time
	sourceNow, err := timezone.GetNowInLocation(input.SourceTimezone)
	if err != nil {
		return nil, types.TimeConversionResult{}, fmt.Errorf("invalid source timezone %q: %w", input.SourceTimezone, err)
	}
	// GetNowInLocation already validated the timezone, so this LoadLocation will succeed
	locFrom, _ := time.LoadLocation(input.SourceTimezone)
	sourceTime := time.Date(sourceNow.Year(), sourceNow.Month(), sourceNow.Day(), hour, minute, 0, 0, locFrom)

	// Convert to target timezone
	locTo, err := time.LoadLocation(input.TargetTimezone)
	if err != nil {
		return nil, types.TimeConversionResult{}, fmt.Errorf("invalid target timezone %q: %w", input.TargetTimezone, err)
	}
	targetTime := sourceTime.In(locTo)

	// Calculate time difference
	_, offSource := sourceTime.Zone()
	_, offTarget := targetTime.Zone()
	timeDiffStr := timeutil.FormatTimeDifference(offSource, offTarget)

	return nil, types.TimeConversionResult{
		Source:         timeutil.BuildTimeResult(sourceTime, input.SourceTimezone),
		Target:         timeutil.BuildTimeResult(targetTime, input.TargetTimezone),
		TimeDifference: timeDiffStr,
	}, nil
}

// RegisterTools attaches the tool handlers to the given server. Extracted for testability.
func RegisterTools(server *mcp.Server, localTZ string) {

	// Tool 1: get_current_time with complete JSON schema
	getCurrentTimeSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"timezone": map[string]any{
				"type":        "string",
				"description": fmt.Sprintf("IANA timezone name (e.g., 'America/New_York', 'Europe/London'). Use '%s' as local timezone if no timezone provided by the user.", localTZ),
			},
		},
		"required": []string{"timezone"},
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_time",
		Description: "Get current time in a specific timezone",
		InputSchema: getCurrentTimeSchema,
	}, GetCurrentTime)

	// Tool 2: convert_time with complete JSON schema
	convertTimeSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"source_timezone": map[string]any{
				"type":        "string",
				"description": fmt.Sprintf("Source IANA timezone name (e.g., 'America/New_York', 'Europe/London'). Use '%s' as local timezone if no source timezone provided by the user.", localTZ),
			},
			"time": map[string]any{
				"type":        "string",
				"description": "Time to convert in 24-hour format (HH:MM)",
			},
			"target_timezone": map[string]any{
				"type":        "string",
				"description": fmt.Sprintf("Target IANA timezone name (e.g., 'Asia/Tokyo', 'America/San_Francisco'). Use '%s' as local timezone if no target timezone provided by the user.", localTZ),
			},
		},
		"required": []string{"source_timezone", "time", "target_timezone"},
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "convert_time",
		Description: "Convert time between timezones",
		InputSchema: convertTimeSchema,
	}, ConvertTime)
}
