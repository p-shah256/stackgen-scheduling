package scheduler

import (
	"time"
)

type TimeSlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type User struct {
	ID           string     `json:"id"`
	Timezone     string     `json:"timezone"`     // IANA timezone format
	Avail []TimeSlot `json:"availability"` // Stored in UTC
}

type Event struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	OrganizerID   string     `json:"organizerId"`
	PossibleSlots []TimeSlot `json:"possibleSlots"` // In UTC
	DurationMins  int        `json:"durationMins"`
	Participants  []string   `json:"participantIds"`
}

type Rec struct {
	Slot               TimeSlot `json:"slot"`
	AvailableUserIDs   []string `json:"availableUserIds"`
	UnavailableUserIDs []string `json:"unavailableUserIds"`
	ParticipationRate  float64  `json:"participationRate"`
}

type TimePoint struct {
	Time    time.Time
	Id      string
	IsStart bool
}

type Interval struct {
	Start   time.Time
	End     time.Time
	UserIDs map[string]bool
}
