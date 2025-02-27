package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func sendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// handleEvent handles both creation (POST) and updates (PUT) of events
func handleEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		sendResponse(w, http.StatusBadRequest, false, "Invalid request format: "+err.Error(), nil)
		return
	}
	event.ID = id

	if event.UserSlots == nil {
		event.UserSlots = []UserAvailability{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the event exists
	var existingEvent Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingEvent)
	exists := err == nil

	// POST = create (fail if exists), PUT = update (fail if not exists)
	if r.Method == "POST" && exists {
		sendResponse(w, http.StatusConflict, false, "Event already exists", nil)
		return
	} else if r.Method == "PUT" && !exists {
		sendResponse(w, http.StatusNotFound, false, "Event not found", nil)
		return
	}

	// Simple upsert operation
	opts := options.Update().SetUpsert(true)
	_, err = eventsCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": event}, opts)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		return
	}

	message := "Event created successfully"
	statusCode := http.StatusCreated
	if r.Method == "PUT" {
		message = "Event updated successfully"
		statusCode = http.StatusOK
	}

	data := event
	sendResponse(w, statusCode, true, message, data)
}

// getEvent retrieves an event by ID
func getEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var event Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			sendResponse(w, http.StatusNotFound, false, "Event not found", nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		}
		return
	}
	sendResponse(w, http.StatusOK, true, "Event retrieved successfully", event)
}

// deleteEvent removes an event by ID
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := eventsCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		return
	}
	if result.DeletedCount == 0 {
		sendResponse(w, http.StatusNotFound, false, "Event not found", nil)
		return
	}
	sendResponse(w, http.StatusOK, true, "Event deleted successfully", nil)
}

// handleUserAvailability adds or updates user availability for an event
func handleUserAvailability(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	userID := params["user_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var event Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			sendResponse(w, http.StatusNotFound, false, "Event not found", nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		}
		return
	}

	var userAvail UserAvailability
	if err := json.NewDecoder(r.Body).Decode(&userAvail); err != nil {
		sendResponse(w, http.StatusBadRequest, false, "Invalid request format: "+err.Error(), nil)
		return
	}
	userAvail.UserID = userID

	// Find if user already exists and determine operation type
	userExists := false
	for i, ua := range event.UserSlots {
		if ua.UserID == userID {
			userExists = true
			// POST = create (fail if exists), PUT = update
			if r.Method == "POST" {
				sendResponse(w, http.StatusConflict, false, "User availability already exists", nil)
				return
			}
			// Update existing entry
			event.UserSlots[i] = userAvail
			break
		}
	}

	if r.Method == "PUT" && !userExists {
		sendResponse(w, http.StatusNotFound, false, "User availability not found", nil)
		return
	}

	// Add new entry if it doesn't exist
	if !userExists {
		// Initialize array if nil
		if event.UserSlots == nil {
			event.UserSlots = []UserAvailability{}
		}
		event.UserSlots = append(event.UserSlots, userAvail)
	}

	_, err = eventsCollection.ReplaceOne(ctx, bson.M{"_id": id}, event)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		return
	}

	message := "User availability updated"
	statusCode := http.StatusOK
	if !userExists {
		message = "User availability added"
		if r.Method == "POST" {
			statusCode = http.StatusCreated
		}
	}
	sendResponse(w, statusCode, true, message, userAvail)
}

func deleteUserAvailability(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	userID := params["user_id"]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var event Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			sendResponse(w, http.StatusNotFound, false, "Event not found", nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		}
		return
	}

	found := false
	newUserSlots := []UserAvailability{}
	for _, ua := range event.UserSlots {
		if ua.UserID != userID {
			newUserSlots = append(newUserSlots, ua)
		} else {
			found = true
		}
	}
	if !found {
		sendResponse(w, http.StatusNotFound, false, "User availability not found", nil)
		return
	}

	event.UserSlots = newUserSlots
	_, err = eventsCollection.ReplaceOne(ctx, bson.M{"_id": id}, event)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		return
	}
	sendResponse(w, http.StatusOK, true, "User availability deleted", nil)
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	timezone := r.URL.Query().Get("timezone")

	if timezone == "" {
		timezone = "UTC" // Default timezone
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var event Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			sendResponse(w, http.StatusNotFound, false, "Event not found", nil)
		} else {
			sendResponse(w, http.StatusInternalServerError, false, "Database error: "+err.Error(), nil)
		}
		return
	}

	recommendations := findOptimalSlots(event)

		slot := &recommendations[0].Slot
		slot.StartStr = formatTimeForDisplay(slot.Start_UTC, timezone)
		slot.EndStr = formatTimeForDisplay(slot.End_UTC, timezone)

	sendResponse(w, http.StatusOK, true, "Recommendations retrieved successfully", recommendations[0])
}
