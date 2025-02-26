package api

import (
	"context"
	"errors"
	"exercise/internal/db"
	"exercise/internal/scheduler"
	"exercise/internal/types"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrDatabase = errors.New("database error")
	ErrInvalid  = errors.New("invalid input")
)

func CreateEvent(ctx context.Context, event types.Event) (types.Event, error) {
	if event.EventId == "" {
		return types.Event{}, ErrInvalid
	}

	slog.Info("Attempting to create event",
		"id", event.EventId,
		"title", event.Title,
		"participants", event.Participants,
		"date", event.Date)

	collection := db.Database.Collection("events")
	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		slog.Error("Failed to insert event", "error", err)
		return types.Event{}, ErrDatabase
	}
	slog.Info("Created event", "id", event.EventId, "title", event.Title)
	return event, nil
}

func UpdateEvent(ctx context.Context, eventID string, event types.Event) (types.Event, error) {
	event.EventId = eventID
	collection := db.Database.Collection("events")
	result, err := collection.ReplaceOne(ctx, bson.M{"_id": eventID}, event)
	if err != nil {
		return types.Event{}, ErrDatabase
	}
	if result.MatchedCount == 0 {
		return types.Event{}, ErrNotFound
	}
	return event, nil
}

func DeleteEvent(ctx context.Context, eventID string) error {
	collection := db.Database.Collection("events")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": eventID})
	if err != nil {
		return ErrDatabase
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func GetEventRecommendations(ctx context.Context, eventID string) (types.RecResults, error) {
	slog.Info("getting recommendation for event")
	collection := db.Database.Collection("events")
	var event types.Event
	err := collection.FindOne(ctx, bson.M{"eventId": eventID}).Decode(&event)
	if err != nil {
		slog.Error("error while finding event", "err", err)
		return types.RecResults{}, err
	}
	date := event.Date

	participants := event.Participants
	userAvails := []types.UsrAvailDate{}

	for _, p := range participants {
		avail, err := GetUserAvail(ctx, p, date)
		if err != nil {
			slog.Error("error while finding participants avail", "err", err, "participant", p)
			return types.RecResults{}, err
		}
		userAvails = append(userAvails, avail)
	}

	// find optimal meeting slot
	recs := scheduler.FindOptimalMeetingSlot(event, userAvails)
	return types.RecResults{EventID: event.EventId, Recommendations: recs}, nil
}

func CreateUser(ctx context.Context, user types.User) (types.User, error) {
	collection := db.Database.Collection("users")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return types.User{}, ErrDatabase
	}
	return user, nil
}

func UpdateUserTimezone(ctx context.Context, userID string, timezone string) error {
	collection := db.Database.Collection("users")
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{"$set": bson.M{"timezone": timezone}},
	)
	if err != nil {
		return ErrDatabase
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func UpdateUserAvail(ctx context.Context, userID string, avail types.UsrAvailDate) error {
	// convert avail to bson
	collection := db.Database.Collection("availabilities")
	_, err := collection.ReplaceOne(ctx, bson.M{"userId": userID}, avail)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func CreateUserAvail(ctx context.Context, userID string, avail types.UsrAvailDate) error {
	// convert avail to bson
	collection := db.Database.Collection("availabilities")
	_, err := collection.InsertOne(ctx, avail)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func DeleteUserAvail(ctx context.Context, userID string, date time.Time) error {
	collection := db.Database.Collection("availabilities")
	result, err := collection.DeleteOne(ctx, bson.M{"userId": userID, "date": date})
	if err != nil {
		return ErrDatabase
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func GetUserAvail(ctx context.Context, userID string, date time.Time) (types.UsrAvailDate, error) {
	collection := db.Database.Collection("availabilities")
	var avail types.UsrAvailDate
	slog.Info("Getting users avail", "userid", userID, "date", date)
	err := collection.FindOne(ctx, bson.M{"userId": userID, "date": date}).Decode(&avail)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			slog.Error("Cant find users avail", "userid", userID, "err", err)
			return types.UsrAvailDate{}, ErrNotFound
		}
		slog.Error("Cant find users avail", "userid", userID, "err", err)
		return types.UsrAvailDate{}, ErrDatabase
	}
	return avail, nil
}
