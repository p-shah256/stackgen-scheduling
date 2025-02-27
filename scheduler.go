package main

import (
	"sort"
	"time"
)

// TimePoint represents a single point in time where availability changes
type TimePoint struct {
	Time    time.Time
	UserID  string
	IsStart bool
}

// findOptimalSlots finds optimal meeting slots using a line sweep algorithm
func findOptimalSlots(event Event) []SlotRecommendation {
	meetingDuration := time.Duration(event.DurationMins) * time.Minute

	// If no users or slots, return empty recommendations
	if len(event.UserSlots) == 0 || len(event.Slots) == 0 {
		return []SlotRecommendation{}
	}

	allUsers := make(map[string]bool)
	for _, user := range event.UserSlots {
		allUsers[user.UserID] = true
	}

	// Generate time points
	points := []TimePoint{}

	// For each user's availability
	for _, user := range event.UserSlots {
		for _, userSlot := range user.Slots {
			// For each event slot
			for _, eventSlot := range event.Slots {
				// Find overlap
				start := later(userSlot.Start_UTC, eventSlot.Start_UTC)
				end := earlier(userSlot.End_UTC, eventSlot.End_UTC)

				// Skip if no overlap
				if start.After(end) || start.Equal(end) {
					continue
				}

				// Add start and end points
				points = append(points,
					TimePoint{Time: start, UserID: user.UserID, IsStart: true},
					TimePoint{Time: end, UserID: user.UserID, IsStart: false},
				)
			}
		}
	}

	// Sort time points
	sort.Slice(points, func(i, j int) bool {
		if points[i].Time.Equal(points[j].Time) {
			return !points[i].IsStart && points[j].IsStart
		}
		return points[i].Time.Before(points[j].Time)
	})

	// Process time points to find intervals
	var recommendations []SlotRecommendation
	activeUsers := make(map[string]bool)
	var startTime time.Time

	for i, point := range points {
		// Update active users
		if point.IsStart {
			// If this is our first user becoming available, mark the start time
			if len(activeUsers) == 0 {
				startTime = point.Time
			}
			activeUsers[point.UserID] = true
		} else {
			delete(activeUsers, point.UserID)
		}

		// Check if we can create a valid interval
		intervalEnd := point.Time
		intervalDuration := intervalEnd.Sub(startTime)

		// If we have a valid interval and it's long enough for the meeting
		if len(activeUsers) > 0 && intervalDuration >= meetingDuration && i > 0 {
			// Create list of available and unavailable users
			available := []string{}
			unavailable := []string{}

			for user := range allUsers {
				if activeUsers[user] {
					available = append(available, user)
				} else {
					unavailable = append(unavailable, user)
				}
			}

			// Add recommendation
			recommendations = append(recommendations, SlotRecommendation{
				Slot: TimeSlot{
					Start_UTC: startTime,
					End_UTC:   startTime.Add(meetingDuration),
				},
				AvailableUsers:   available,
				UnavailableUsers: unavailable,
			})
		}

		// If no users are active anymore, reset the start time
		if len(activeUsers) == 0 {
			startTime = time.Time{}
		} else if i > 0 && !points[i-1].Time.Equal(point.Time) {
			// Only update start time if we've moved to a new time point
			startTime = point.Time
		}
	}

	// Sort recommendations by number of available users (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return len(recommendations[i].AvailableUsers) > len(recommendations[j].AvailableUsers)
	})

	return recommendations
}

// earlier returns the earlier of two time points
func earlier(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// later returns the later of two time points
func later(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
