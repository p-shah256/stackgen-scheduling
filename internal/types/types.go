package types

import "time"

type DailySlots struct {
	Date  time.Time  `json:"date" bson:"date"` // YYYY-MM-DD format
	Slots []TimeSlot `json:"slots" bson:"slots"`
}

type TimeSlot struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

type User struct {
	ID       int          `json:"id" bson:"_id"`
	Name     string       `json:"name" bson:"name"`
	Timezone string       `json:"timezone" bson:"timezone"`
	Avail    []DailySlots `json:"availability" bson:"availability"`
}

// NOTE: we can use this if the user gets too big in memory, like if they share the whole year worth of avail
// this bad boy could just be a particular days' avail

// type UserAvailability struct {
//     UserID      string     `json:"userId"`
//     Name        string     `json:"name"`
//     Timezone    string     `json:"timezone"`
//     Slots       []TimeSlot `json:"slots"`
// }

type Event struct {
	ID            int       `json:"id" bson:"_id"`
	Title         string       `json:"title" bson:"title"`
	OrganizerID   string       `json:"organizerId" bson:"organizerId"`
	PossibleSlots []DailySlots `json:"possibleSlots" bson:"possibleSlots"`
	DurationMins  int          `json:"durationMins" bson:"durationMins"`
	Participants  []string     `json:"participantIds" bson:"participantIds"`
}

type Rec struct {
	Slot               TimeSlot `json:"slot" bson:"slot"`
	Date               time.Time   `json:"date" bson:"date"`
	AvailableUserIDs   []int `json:"availableUserIds" bson:"availableUserIds"`
	UnavailableUserIDs []int `json:"unavailableUserIds" bson:"unavailableUserIds"`
	ParticipationRate  float64  `json:"participationRate" bson:"participationRate"`
}

type RecResults struct {
	EventID         string `json:"eventId" bson:"eventId"`
	Recommendations []Rec  `json:"recommendations" bson:"recommendations"`
}
