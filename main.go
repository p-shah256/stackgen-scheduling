package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type TimeSlot struct {
	Start_t    time.Time `json:"start_t"`
	End_t      time.Time `json:"end_t"`
	StartStr string    `json:"start,omitempty"`
	EndStr   string    `json:"end,omitempty"`
}

func (ts *TimeSlot) UnmarshalJSON(data []byte) error {
	type TimeSlotAlias TimeSlot
	aux := struct {
		StartStr string `json:"start_str"`
		EndStr   string `json:"end_str"`
		*TimeSlotAlias
	}{
		TimeSlotAlias: (*TimeSlotAlias)(ts),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.StartStr != "" {
		startTime, err := parseTimeString(aux.StartStr)
		if err != nil {
			return err
		}
		ts.Start_t = startTime
		ts.StartStr = aux.StartStr
	}

	if aux.EndStr != "" {
		endTime, err := parseTimeString(aux.EndStr)
		if err != nil {
			return err
		}
		ts.End_t = endTime
		ts.EndStr = aux.EndStr
	}

	return nil
}

func parseTimeString(timeStr string) (time.Time, error) {
	formats := []string{
		"2 Jan 2006, 3:04PM MST", // 12 Jan 2025, 2:00PM EST
		"2 Jan 2006, 3PM MST",    // 12 Jan 2025, 2PM EST
		"2 Jan 2006, 15:04 MST",  // 12 Jan 2025, 14:00 EST
		"2 Jan 2006, 15 MST",     // 12 Jan 2025, 14 EST
	}

	timeStr = strings.TrimSpace(timeStr)

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeStr)
}

type UserAvailability struct {
	UserID string     `json:"user_id"`
	Slots  []TimeSlot `json:"slots"`
}

type Event struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	DurationMins int                `json:"duration_mins"`
	Slots        []TimeSlot         `json:"slots"`
	UserSlots    []UserAvailability `json:"user_slots"`
}

type SlotRecommendation struct {
	Slot             TimeSlot `json:"slot"`
	AvailableUsers   []string `json:"available_users"`
	UnavailableUsers []string `json:"unavailable_users"`
}

var events = make(map[string]*Event)

func formatTimeForDisplay(t time.Time) string {
	return t.Format("2 Jan 2006, 3:04PM MST")
}

func findOptimalSlots(event Event) []SlotRecommendation {
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

func main() {
	router := mux.NewRouter()

	// Event endpoints - PUT or POST
	router.HandleFunc("/events/{id}", handleEvent).Methods("POST", "PUT")
	router.HandleFunc("/events/{id}", getEvent).Methods("GET")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	// User availability endpoints
	router.HandleFunc("/events/{id}/availability/{user_id}", handleUserAvailability).Methods("POST", "PUT")
	router.HandleFunc("/events/{id}/availability/{user_id}", deleteUserAvailability).Methods("DELETE")

	// Recommendation endpoint
	router.HandleFunc("/events/{id}/recommendations", getRecommendations).Methods("GET")

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
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

	_, exists := events[id]

	// POST = create (fail if exists), PUT = update (fail if not exists)
	if r.Method == "POST" && exists {
		http.Error(w, "Event already exists", http.StatusConflict)
		return
	} else if r.Method == "PUT" && !exists {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	events[id] = &event

	message := "Event created successfully"
	if r.Method == "PUT" {
		message = "Event updated successfully"
	}

	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
		"id":      id,
	})
}

func getEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	event, found := events[id]
	if !found {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if _, found := events[id]; !found {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	delete(events, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}

func handleUserAvailability(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	userID := params["user_id"]

	event, found := events[id]
	if !found {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	var userAvail UserAvailability
	if err := json.NewDecoder(r.Body).Decode(&userAvail); err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return
	}

	userAvail.UserID = userID

	userExists := false
	for i, ua := range event.UserSlots {
		if ua.UserID == userID {
			// POST = create (fail if exists), PUT = update (fail if not exists)
			if r.Method == "POST" {
				http.Error(w, "User availability already exists", http.StatusConflict)
				return
			}

			event.UserSlots[i] = userAvail
			userExists = true
			break
		}
	}

	if !userExists && r.Method == "PUT" {
		http.Error(w, "User availability not found", http.StatusNotFound)
		return
	}

	if !userExists {
		event.UserSlots = append(event.UserSlots, userAvail)
	}

	message := "User availability added"
	if userExists {
		message = "User availability updated"
	}

	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" && !userExists {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
	})
}

func deleteUserAvailability(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	userID := params["user_id"]

	event, found := events[id]
	if !found {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	for i, ua := range event.UserSlots {
		if ua.UserID == userID {
			lastIndex := len(event.UserSlots) - 1
			event.UserSlots[i] = event.UserSlots[lastIndex]
			event.UserSlots = event.UserSlots[:lastIndex]

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "User availability deleted",
			})
			return
		}
	}

	http.Error(w, "User availability not found", http.StatusNotFound)
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	event, found := events[id]
	if !found {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	recommendations := findOptimalSlots(*event)

	for i := range recommendations {
		slot := &recommendations[i].Slot
		slot.StartStr = formatTimeForDisplay(slot.Start_t)
		slot.EndStr = formatTimeForDisplay(slot.End_t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
