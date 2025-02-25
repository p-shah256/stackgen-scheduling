package scheduler

import (
	"log/slog"
	"sort"
	"time"
)

func FindOptimalMeetingSlot(event Event, allUsers []User) []RecommendationResult {
	meetingDuration := time.Duration(event.DurationMins) * time.Minute

	timePoints := []TimePoint{}

	for _, user := range allUsers {
		// get overlap
		for _, timeSlot := range user.Avail {
			for _, eventSlot := range event.PossibleSlots {
				start := max(timeSlot.Start, eventSlot.Start)
				end := min(timeSlot.End, eventSlot.End)

				if start.After(end) || start.Equal(end) { // no overlap
					continue
				}

				timePoints = append(timePoints,
					TimePoint{Time: start, Id: user.ID, IsStart: true},
					TimePoint{Time: end, Id: user.ID, IsStart: false},
				)
			}
		}
	}

	sort.Slice(timePoints, func(i, j int) bool {
		if timePoints[i].Time.Equal(timePoints[j].Time) {
			return !timePoints[i].IsStart && timePoints[j].IsStart
		}
		return timePoints[i].Time.Before(timePoints[j].Time)
	})
	slog.Debug("sorted timepoints", "timepoints", timePoints)

	intervals := []Interval{}
	currentUsers := make(map[string]bool)
	var lastTime time.Time
	if len(timePoints) > 0 {
		lastTime = timePoints[0].Time
	}

	for i, tp := range timePoints {
		// interval if time has advanced
		if !tp.Time.Equal(lastTime) && len(currentUsers) > 0 {
			// Deep copy the current users
			users := make(map[string]bool)
			for id := range currentUsers {
				users[id] = true
			}

			intervals = append(intervals, Interval{
				Start:   lastTime,
				End:     tp.Time,
				UserIDs: users,
			})
		}

		if tp.IsStart {
			currentUsers[tp.Id] = true
		} else {
			delete(currentUsers, tp.Id)
		}

		if i == len(timePoints)-1 || !tp.Time.Equal(timePoints[i+1].Time) {
			lastTime = tp.Time
		}
	}
	slog.Debug("INTERVALS::", "intervals", intervals)


	recommendationResults := []RecommendationResult{}
	for _, interval := range intervals {
		duration := interval.End.Sub(interval.Start)
		if duration >= meetingDuration {
			endTime := interval.Start.Add(meetingDuration)
			candidateInterval := Interval{
				Start:   interval.Start,
				End:     endTime,
				UserIDs: interval.UserIDs,
			}

			availableUsers := []User{}
			unavailableUsers := []User{}

			for _, user := range allUsers {
				if candidateInterval.UserIDs[user.ID] {
					availableUsers = append(availableUsers, user)
				} else {
					unavailableUsers = append(unavailableUsers, user)
				}
			}

			participationRate := float64(len(availableUsers)) / float64(len(allUsers))

			recommendationResults = append(recommendationResults, RecommendationResult{
				Slot: TimeSlot{
					Start: candidateInterval.Start,
					End:   candidateInterval.End,
				},
				AvailableUsers:    availableUsers,
				UnavailableUsers:  unavailableUsers,
				ParticipationRate: participationRate,
			})
		}
	}

	sort.Slice(recommendationResults, func(i, j int) bool {
		return recommendationResults[i].ParticipationRate > recommendationResults[j].ParticipationRate
	})

	return recommendationResults
}

func min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
