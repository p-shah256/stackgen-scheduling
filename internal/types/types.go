package types

import "time"

type TimeSlot struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

type User struct {
	UserID   string `json:"userId" bson:"userId"`
	Timezone string `json:"timezone" bson:"timezone"`
}

// separated availability from user for better querying
type UsrAvailDate struct {
	UserID string     `json:"userId" bson:"userId"`
	Date   time.Time  `json:"date" bson:"date"` // YYYY-MM-DD format
	Slots  []TimeSlot `json:"slots" bson:"slots"`
}

type Event struct {
	EventId       string     `json:"eventId" bson:"eventId"`
	Title         string     `json:"title" bson:"title"`
	PossibleSlots []TimeSlot `json:"possibleSlots" bson:"possibleSlots"`
	DurationMins  int        `json:"durationMins" bson:"durationMins"`
	Participants  []string   `json:"participantIds" bson:"participantIds"`
	Date          time.Time  `json:"date" bson:"date"` // YYYY-MM-DD format
}

type Rec struct {
	Slot               TimeSlot  `json:"slot" bson:"slot"`
	Date               time.Time `json:"date" bson:"date"`
	AvailableUserIDs   []string  `json:"availableUserIds" bson:"availableUserIds"`
	UnavailableUserIDs []string  `json:"unavailableUserIds" bson:"unavailableUserIds"`
	ParticipationRate  float64   `json:"participationRate" bson:"participationRate"`
}

type RecResults struct {
	EventID         string `json:"eventId" bson:"eventId"`
	Recommendations []Rec  `json:"recommendations" bson:"recommendations"`
}
