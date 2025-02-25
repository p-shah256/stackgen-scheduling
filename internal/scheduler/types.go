package scheduler

import (
	"time"
)

type TimeSlot struct {
	Start time.Time
	End   time.Time
}

type User struct {
	ID           string
	Avail []TimeSlot
}

type Event struct {
	ID              string
	Title           string
	PossibleSlots   []TimeSlot
	DurationMins int
}

type RecommendationResult struct {
	Slot              TimeSlot
	AvailableUsers    []User
	UnavailableUsers  []User
	ParticipationRate float64
}

type TimePoint struct {
	Time    time.Time
	UserID  string
	IsStart bool
}

type Interval struct {
	Start   time.Time
	End     time.Time
	UserIDs map[string]bool
}


