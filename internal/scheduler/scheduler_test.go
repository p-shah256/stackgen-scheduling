package scheduler

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestHappyPath_AllUsersAvailable(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)


	event := Event{
		ID:           "event1",
		Title:        "Team Meeting",
		DurationMins: 60,
		PossibleSlots: []TimeSlot{
			{
				Start: time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 3, 1, 16, 0, 0, 0, time.UTC),
			},
			{
				Start: time.Date(2025, 3, 2, 10, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 3, 2, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	users := []User{
		{
			ID: "user1",
			Avail: []TimeSlot{
				{
					Start: time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 1, 15, 30, 0, 0, time.UTC),
				},
			},
		},
		{
			ID: "user2",
			Avail: []TimeSlot{
				{
					Start: time.Date(2025, 3, 1, 14, 30, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 1, 16, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			ID: "user3",
			Avail: []TimeSlot{
				{
					Start: time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 1, 16, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	results := FindOptimalMeetingSlot(event, users)

	if results[0].ParticipationRate != 1.0 {
		t.Errorf("Expected 100%% participation rate, got %.2f%%", results[0].ParticipationRate*100)
	}
}

func TestEdgeCase_NoCommonSlot(t *testing.T) {
	event := Event{
		ID:             "event2",
		Title:          "Budget Review",
		DurationMins: 60,
		PossibleSlots: []TimeSlot{
			{
				Start: time.Date(2025, 3, 5, 13, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 3, 5, 17, 0, 0, 0, time.UTC),
			},
		},
	}

	// Define users with conflicting availabilities
	users := []User{
		{
			ID: "user1",
			Avail: []TimeSlot{
				{
					Start: time.Date(2025, 3, 5, 13, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 5, 15, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			ID: "user2",
			Avail: []TimeSlot{
				{
					Start: time.Date(2025, 3, 5, 14, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 5, 16, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			ID: "user3",
			Avail: []TimeSlot{
				{
					Start: time.Date(2025, 3, 5, 15, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 5, 17, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	results := FindOptimalMeetingSlot(event, users)

	// The best we can do is 2/3 participation (no slot works for all 3 users)
	// for a 90-minute meeting
	if results[0].ParticipationRate != 2.0/3.0 {
		t.Errorf("Expected 66.67%% participation rate, got %.2f%%", results[0].ParticipationRate*100)
	}
}
