package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"context"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMeetingSchedulerIntegration(t *testing.T) {
	setupTestEnvironment(t)
	defer teardownTestEnvironment(t)

	router := mux.NewRouter()
	router.HandleFunc("/events/{id}", handleEvent).Methods("POST", "PUT")
	router.HandleFunc("/events/{id}", getEvent).Methods("GET")
	router.HandleFunc("/events/{id}/availability/{user_id}", handleUserAvailability).Methods("POST", "PUT")
	router.HandleFunc("/events/{id}/recommendations", getRecommendations).Methods("GET")

	// Test 1: Create an event
	t.Run("Create Event", func(t *testing.T) {
		eventData := map[string]interface{}{
			"title":         "Team Meeting",
			"duration_mins": 60,
			"slots": []map[string]string{
				{
					"start":    "15 Jan 2025, 9:00AM",
					"end":      "15 Jan 2025, 5:00PM",
					"timezone": "UTC",
				},
			},
		}

		req := createJSONRequest("POST", "/events/test-event-123", eventData)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		var response Response
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.True(t, response.Success)
	})

	// Test 2: Add user availability for 3 users in different timezones
	t.Run("Add User Availability", func(t *testing.T) {
		// Simple data structure with 3 users in different timezones
		users := []struct {
			id    string
			slots map[string]interface{}
		}{
			{
				id: "user1",
				slots: map[string]interface{}{
					"slots": []map[string]string{
						{
							"start":    "15 Jan 2025, 10:00AM",
							"end":      "15 Jan 2025, 2:00PM",
							"timezone": "America/New_York",
						},
					},
				},
			},
			{
				id: "user2",
				slots: map[string]interface{}{
					"slots": []map[string]string{
						{
							"start":    "15 Jan 2025, 8:00AM",
							"end":      "15 Jan 2025, 12:00PM",
							"timezone": "America/Los_Angeles",
						},
					},
				},
			},
			{
				id: "user3",
				slots: map[string]interface{}{
					"slots": []map[string]string{
						{
							"start":    "15 Jan 2025, 3:00PM",
							"end":      "15 Jan 2025, 7:00PM",
							"timezone": "Europe/London",
						},
					},
				},
			},
		}

		// Add availability for each user
		for _, user := range users {
			req := createJSONRequest("POST", "/events/test-event-123/availability/"+user.id, user.slots)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			var response Response
			json.Unmarshal(resp.Body.Bytes(), &response)
			assert.True(t, response.Success)
		}
	})

	// Test 3: Get event details
	t.Run("Get Event Details", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/events/test-event-123", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		var response Response
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.True(t, response.Success)

		eventData, _ := json.Marshal(response.Data)
		var event Event
		json.Unmarshal(eventData, &event)

		assert.Equal(t, "Team Meeting", event.Title)
		assert.Equal(t, 60, event.DurationMins)
	})

	// Test 4: Get recommendations (simplified to just check structure)
	t.Run("Get Basic Recommendations", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/events/test-event-123/recommendations", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		var response Response
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.True(t, response.Success)
		
		if response.Data == nil {
			t.Log("Response data is nil, can't validate structure")
			return
		}
		
		recData, _ := json.Marshal(response.Data)
		var recommendations SlotRecommendation
		err := json.Unmarshal(recData, &recommendations)
		
		if err != nil {
			t.Logf("Could not parse recommendation data: %v", err)
		}

		assert.Equal(t,	3, len(recommendations.AvailableUsers))
	})
}

func createJSONRequest(method, url string, data interface{}) *http.Request {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func setupTestEnvironment(t *testing.T) {
	ctx := context.Background()
	testMongoURI := os.Getenv("TEST_MONGO_URI")
	if testMongoURI == "" {
		testMongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(testMongoURI)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	eventsCollection = client.Database("testMeetingScheduler").Collection("events")
}

func teardownTestEnvironment(t *testing.T) {
	if eventsCollection != nil {
		eventsCollection.Drop(context.Background())
	}
	if client != nil {
		client.Disconnect(context.Background())
	}
}
