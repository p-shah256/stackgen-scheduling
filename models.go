package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TimeSlot struct {
	Start_UTC time.Time `json:"start_utc" bson:"start_utc"`
	End_UTC   time.Time `json:"end_utc" bson:"end_utc"`
	StartStr  string    `json:"start,omitempty" bson:"start,omitempty"`
	EndStr    string    `json:"end,omitempty" bson:"end,omitempty"`
	TimeZone  string    `json:"timezone,omitempty" bson:"timezone,omitempty"` // Add timezone field
}

// UnmarshalJSON handles JSON parsing for TimeSlot
func (ts *TimeSlot) UnmarshalJSON(data []byte) error {
	// Temporary struct to avoid recursion
	userInput := struct {
		StartStr string `json:"start"`
		EndStr   string `json:"end"`
		TimeZone string `json:"timezone"`
	}{}

	if err := json.Unmarshal(data, &userInput); err != nil {
		return err
	}

	if userInput.StartStr == "" || userInput.EndStr == "" {
		return fmt.Errorf("start or end cannot be empty")
	}
	tzName := "UTC"
	if userInput.TimeZone != "" {
		tzName = userInput.TimeZone
	}

	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return fmt.Errorf("invalid timezone: %s", tzName)
	}

	start, err := parseTimeInLocation(userInput.StartStr, loc)
	if err != nil {
		return err
	}
	end, err := parseTimeInLocation(userInput.EndStr, loc)
	if err != nil {
		return err
	}

	// Set fields
	ts.Start_UTC = start.UTC()
	ts.End_UTC = end.UTC()
	ts.StartStr = userInput.StartStr
	ts.EndStr = userInput.EndStr
	ts.TimeZone = tzName

	return nil
}

// parseTimeInLocation parses time string in the given location
func parseTimeInLocation(timeStr string, loc *time.Location) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)

	formats := []string{
		"2 Jan 2006, 3:04PM",
		"2 Jan 2006, 3PM",
		"2 Jan 2006, 15:04",
		"2 Jan 2006, 15",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, loc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse time: %s", timeStr)
}

type UserAvailability struct {
	UserID string     `json:"user_id" bson:"user_id"`
	Slots  []TimeSlot `json:"slots" bson:"slots"`
}

type Event struct {
	ID           string             `json:"id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	DurationMins int                `json:"duration_mins" bson:"duration_mins"`
	Slots        []TimeSlot         `json:"slots" bson:"slots"`
	UserSlots    []UserAvailability `json:"user_slots" bson:"user_slots"`
}

type SlotRecommendation struct {
	Slot             TimeSlot `json:"slot" bson:"slot"`
	AvailableUsers   []string `json:"available_users" bson:"available_users"`
	UnavailableUsers []string `json:"unavailable_users" bson:"unavailable_users"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

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
