package scheduler

import (
	"exercise/internal/types"
	"log/slog"
	"sort"
	"time"
)

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

func FindOptimalMeetingSlot(event types.Event, userAvails []types.UsrAvailDate) []types.Rec {
	meetingDuration := time.Duration(event.DurationMins) * time.Minute

	timePoints := []TimePoint{}

	for _, userAvail := range userAvails {
		for _, timeSlot := range userAvail.Slots {

			for _, eventSlot := range event.PossibleSlots {
				start := max(timeSlot.Start, eventSlot.Start)
				end := min(timeSlot.End, eventSlot.End)

				if start.After(end) || start.Equal(end) {
					continue
				}

				timePoints = append(timePoints,
					TimePoint{Time: start, UserID: userAvail.UserID, IsStart: true},
					TimePoint{Time: end, UserID: userAvail.UserID, IsStart: false},
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
		if !tp.Time.Equal(lastTime) && len(currentUsers) > 0 {
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
			currentUsers[tp.UserID] = true
		} else {
			delete(currentUsers, tp.UserID)
		}

		if i == len(timePoints)-1 || !tp.Time.Equal(timePoints[i+1].Time) {
			lastTime = tp.Time
		}
	}
	slog.Debug("INTERVALS::", "intervals", intervals)

	recommendationResults := []types.Rec{}
	userMap := make(map[string]types.UsrAvailDate)
	for _, avail := range userAvails {
		userMap[avail.UserID] = avail
	}
	allUserIDs := make([]string, 0, len(userMap))
	for userID := range userMap {
		allUserIDs = append(allUserIDs, userID)
	}

	for _, interval := range intervals {
		duration := interval.End.Sub(interval.Start)
		if duration >= meetingDuration {
			endTime := interval.Start.Add(meetingDuration)
			candidateInterval := Interval{
				Start:   interval.Start,
				End:     endTime,
				UserIDs: interval.UserIDs,
			}

			availableUserIDs := []string{}
			unavailableUserIDs := []string{}

			for _, userID := range allUserIDs {
				if candidateInterval.UserIDs[userID] {
					availableUserIDs = append(availableUserIDs, userID)
				} else {
					unavailableUserIDs = append(unavailableUserIDs, userID)
				}
			}

			recommendationResults = append(recommendationResults, types.Rec{
				Slot: types.TimeSlot{
					Start: candidateInterval.Start,
					End:   candidateInterval.End,
				},
				Date:               interval.Start.Truncate(24 * time.Hour),
				AvailableUserIDs:   availableUserIDs,
				UnavailableUserIDs: unavailableUserIDs,
				ParticipationRate:  float64(len(availableUserIDs)) / float64(len(allUserIDs)),
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
