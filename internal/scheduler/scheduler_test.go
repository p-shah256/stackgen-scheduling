package scheduler

import (
	"exercise/internal/types"
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

	event := types.Event{
		EventId:           "1",
		Title:        "Team Meeting",
		DurationMins: 60,
		PossibleSlots: []types.TimeSlot{
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

	usersAvailability := []types.UsrAvailDate{
		{
			UserID: "1",
			Date:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			Slots: []types.TimeSlot{
				{
					Start: time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 1, 15, 30, 0, 0, time.UTC),
				},
			},
		},
		{
			UserID: "2",
			Date:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			Slots: []types.TimeSlot{
				{
					Start: time.Date(2025, 3, 1, 14, 30, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 1, 16, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			UserID: "3",
			Date:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			Slots: []types.TimeSlot{
				{
					Start: time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 1, 16, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	results := FindOptimalMeetingSlot(event, usersAvailability)

	if results[0].ParticipationRate != 1.0 {
		t.Errorf("Expected 100%% participation rate, got %.2f%%", results[0].ParticipationRate*100)
	}
}

func TestEdgeCase_NoCommonSlot(t *testing.T) {
	event := types.Event{
		EventId:           "2",
		Title:        "Budget Review",
		DurationMins: 60,
		PossibleSlots: []types.TimeSlot{
			{
				Start: time.Date(2025, 3, 5, 13, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 3, 5, 17, 0, 0, 0, time.UTC),
			},
		},
	}

	usersAvailability := []types.UsrAvailDate{
		{
			UserID: "1",
			Date:   time.Date(2025, 3, 5, 0, 0, 0, 0, time.UTC),
			Slots: []types.TimeSlot{
				{
					Start: time.Date(2025, 3, 5, 13, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 5, 15, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			UserID: "2",
			Date:   time.Date(2025, 3, 5, 0, 0, 0, 0, time.UTC),
			Slots: []types.TimeSlot{
				{
					Start: time.Date(2025, 3, 5, 14, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 5, 16, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			UserID: "3",
			Date:   time.Date(2025, 3, 5, 0, 0, 0, 0, time.UTC),
			Slots: []types.TimeSlot{
				{
					Start: time.Date(2025, 3, 5, 15, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 3, 5, 17, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	results := FindOptimalMeetingSlot(event, usersAvailability)

	if len(results) == 0 || results[0].ParticipationRate != 2.0/3.0 {
		t.Errorf("Expected 66.67%% participation rate, got ")
		if len(results) > 0 {
			t.Errorf("%.2f%%", results[0].ParticipationRate*100)
		} else {
			t.Errorf("No results")
		}

	}
}
