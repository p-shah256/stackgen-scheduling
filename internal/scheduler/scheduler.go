package scheduler

import (
	"time")

type TimeSlot struct {
	Start time.Time
	End   time.Time
}

type User struct {
	ID           string
	Availability []TimeSlot
}

type Event struct {
	ID             string
	Title          string
	PossibleSlots  []TimeSlot
	DurationMinutes int
}

type RecommendationResult struct {
	Slot              TimeSlot
	AvailableUsers    []User
	UnavailableUsers  []User
	ParticipationRate float64
}


func FindOptimalTimeSlots(event Event, users []User) []RecommendationResult {
	// create a recommendationresult struct
	// 1. for each evnetDuration in Possible slots 
		// 2. for each user in users, update result struct
	// 3. return result struct with max users
	return nil
}
