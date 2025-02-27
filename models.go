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
}

// UnmarshalJSON custom unmarshaler to handle string time formats
func (ts *TimeSlot) UnmarshalJSON(data []byte) error {
	type TimeSlotAlias TimeSlot
	aux := struct {
		StartStr string `json:"start_str"`
		EndStr   string `json:"end_str"`
		*TimeSlotAlias
	}{
		TimeSlotAlias: (*TimeSlotAlias)(ts),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.StartStr != "" {
		startTime, err := parseTimeString(aux.StartStr)
		if err != nil {
			return err
		}
		ts.Start_t = startTime
		ts.StartStr = aux.StartStr
	}

	if aux.EndStr != "" {
		endTime, err := parseTimeString(aux.EndStr)
		if err != nil {
			return err
		}
		ts.End_t = endTime
		ts.EndStr = aux.EndStr
	}

	return nil
}

// parseTimeString parses various time formats into time.Time
func parseTimeString(timeStr string) (time.Time, error) {
	formats := []string{
		"2 Jan 2006, 3:04PM MST", // 12 Jan 2025, 2:00PM EST
		"2 Jan 2006, 3PM MST",    // 12 Jan 2025, 2PM EST
		"2 Jan 2006, 15:04 MST",  // 12 Jan 2025, 14:00 EST
		"2 Jan 2006, 15 MST",     // 12 Jan 2025, 14 EST
	}

	timeStr = strings.TrimSpace(timeStr)

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeStr)
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

// formatTimeForDisplay formats a time.Time for display
func formatTimeForDisplay(t time.Time) string {
	return t.Format("2 Jan 2006, 3:04PM MST")
}
