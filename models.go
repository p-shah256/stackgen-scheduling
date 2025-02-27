package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TimeSlot represents a time period with start and end times
type TimeSlot struct {
	Start_t  time.Time `json:"start_t" bson:"start_t"`
	End_t    time.Time `json:"end_t" bson:"end_t"`
	StartStr string    `json:"start,omitempty" bson:"start_str,omitempty"`
	EndStr   string    `json:"end,omitempty" bson:"end_str,omitempty"`
	TimeZone string    `json:"timezone,omitempty" bson:"timezone,omitempty"`
}

// UnmarshalJSON custom unmarshaler to handle string time formats
func (ts *TimeSlot) UnmarshalJSON(data []byte) error {
	type TimeSlotAlias TimeSlot
	aux := struct {
		StartStr string `json:"start_str"`
		EndStr   string `json:"end_str"`
		TimeZone string `json:"timezone"`
		*TimeSlotAlias
	}{
		TimeSlotAlias: (*TimeSlotAlias)(ts),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Default timezone to UTC if not provided
	timezone := aux.TimeZone
	if timezone == "" {
		timezone = "UTC"
	}
	ts.TimeZone = timezone

	if aux.StartStr != "" {
		startTime, err := parseTimeString(aux.StartStr, timezone)
		if err != nil {
			return err
		}
		ts.Start_t = startTime
		ts.StartStr = aux.StartStr
	}

	if aux.EndStr != "" {
		endTime, err := parseTimeString(aux.EndStr, timezone)
		if err != nil {
			return err
		}
		ts.End_t = endTime
		ts.EndStr = aux.EndStr
	}

	return nil
}

// parseTimeString parses various time formats into time.Time with timezone
func parseTimeString(timeStr string, timezone string) (time.Time, error) {
	formats := []string{
		"2 Jan 2006, 3:04PM MST", // 12 Jan 2025, 2:00PM EST
		"2 Jan 2006, 3PM MST",    // 12 Jan 2025, 2PM EST
		"2 Jan 2006, 15:04 MST",  // 12 Jan 2025, 14:00 EST
		"2 Jan 2006, 15 MST",     // 12 Jan 2025, 14 EST
	}

	timeStr = strings.TrimSpace(timeStr)
	
	// Try to parse with the timezone already in the string
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}
	
	// If no timezone in string, try to add the specified timezone
	// Extract timezone from string if present, otherwise use provided timezone
	hasTZ := strings.Contains(timeStr, "EST") || strings.Contains(timeStr, "PST") || 
	         strings.Contains(timeStr, "UTC") || strings.Contains(timeStr, "GMT") || 
	         strings.Contains(timeStr, "MST")
			 
	if !hasTZ {
		// Try adding the timezone to the string
		for _, format := range formats {
			tzFormat := strings.Replace(format, "MST", "", -1) // Remove MST from format
			if t, err := time.Parse(tzFormat, timeStr); err == nil {
				// Get location from timezone string
				loc, err := time.LoadLocation(timezone)
				if err != nil {
					return time.Time{}, fmt.Errorf("invalid timezone: %s", timezone)
				}
				return t.In(loc), nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s with timezone: %s", timeStr, timezone)
}

// UserAvailability represents a user's available time slots
type UserAvailability struct {
	UserID string     `json:"user_id" bson:"user_id"`
	Slots  []TimeSlot `json:"slots" bson:"slots"`
}

// Event represents a meeting event with slots and user availabilities
type Event struct {
	ID           string             `json:"id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	DurationMins int                `json:"duration_mins" bson:"duration_mins"`
	Slots        []TimeSlot         `json:"slots" bson:"slots"`
	UserSlots    []UserAvailability `json:"user_slots" bson:"user_slots"`
}

// SlotRecommendation represents a recommended time slot with user availability info
type SlotRecommendation struct {
	Slot             TimeSlot `json:"slot" bson:"slot"`
	AvailableUsers   []string `json:"available_users" bson:"available_users"`
	UnavailableUsers []string `json:"unavailable_users" bson:"unavailable_users"`
}

// Response provides a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// formatTimeForDisplay formats a time.Time for display with timezone
func formatTimeForDisplay(t time.Time, timezone string) string {
	if timezone == "" {
		return t.Format("2 Jan 2006, 3:04PM MST")
	}
	
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return t.Format("2 Jan 2006, 3:04PM MST")
	}
	
	return t.In(loc).Format("2 Jan 2006, 3:04PM MST")
}
