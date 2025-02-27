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

// findOptimalSlots placeholder that will be filled later with actual logic
func findOptimalSlots(event Event) []SlotRecommendation {
	// Placeholder implementation - to be filled in later as mentioned
	return []SlotRecommendation{
		{
			Slot: TimeSlot{
				Start_t: time.Now().Add(24 * time.Hour),
				End_t:   time.Now().Add(25 * time.Hour),
			},
			AvailableUsers:   []string{"user1", "user2"},
			UnavailableUsers: []string{"user3"},
		},
	}
}

// handleEvent handles both creation (POST) and updates (PUT) of events
func handleEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return
	}

	event.ID = id
	
	// Always initialize UserSlots to prevent nil array issues
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
		http.Error(w, "Event already exists", http.StatusConflict)
		return
	} else if r.Method == "PUT" && !exists {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	// Simple upsert operation
	opts := options.Update().SetUpsert(true)
	result, err := eventsCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": event}, opts)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Event created successfully"
	statusCode := http.StatusCreated
	if r.Method == "PUT" {
		message = "Event updated successfully"
		statusCode = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": message,
		"id":      id,
		"matched": result.MatchedCount,
		"upserted": result.UpsertedCount,
	})
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
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// deleteEvent removes an event by ID
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := eventsCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}

// handleUserAvailability adds or updates user availability for an event
func handleUserAvailability(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	userID := params["user_id"]

	// First, get the entire event
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var event Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse the new user availability
	var userAvail UserAvailability
	if err := json.NewDecoder(r.Body).Decode(&userAvail); err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
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
				http.Error(w, "User availability already exists", http.StatusConflict)
				return
			}
			// Update existing entry
			event.UserSlots[i] = userAvail
			break
		}
	}

	if r.Method == "PUT" && !userExists {
		http.Error(w, "User availability not found", http.StatusNotFound)
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

	// Save the full updated event
	_, err = eventsCollection.ReplaceOne(ctx, bson.M{"_id": id}, event)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
	})
}

// deleteUserAvailability removes a user's availability from an event
func deleteUserAvailability(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	userID := params["user_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First, get the entire event
	var event Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Find and remove the user
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
		http.Error(w, "User availability not found", http.StatusNotFound)
		return
	}

	// Update the event with the user removed
	event.UserSlots = newUserSlots
	_, err = eventsCollection.ReplaceOne(ctx, bson.M{"_id": id}, event)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User availability deleted",
	})
}

// getRecommendations returns recommended time slots for an event
func getRecommendations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var event Event
	err := eventsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	recommendations := findOptimalSlots(event)

	// Format times for display
	for i := range recommendations {
		slot := &recommendations[i].Slot
		slot.StartStr = formatTimeForDisplay(slot.Start_t)
		slot.EndStr = formatTimeForDisplay(slot.End_t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
