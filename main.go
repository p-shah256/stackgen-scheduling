package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var eventsCollection *mongo.Collection

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	err := client.Ping(ctx, nil)
	
	response := Response{
		Success: err == nil,
		Message: "Service health status",
		Data: map[string]string{
			"status": "healthy",
			"db_connection": "connected",
		},
	}
	
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusServiceUnavailable
		response.Success = false
		response.Data = map[string]string{
			"status": "unhealthy",
			"db_connection": "disconnected",
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "meetingScheduler")
	port := getEnv("PORT", "8082")
	
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB at", mongoURI)
	eventsCollection = client.Database(dbName).Collection("events")
	
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	router := mux.NewRouter()
	
	// Health check endpoint for Kubernetes
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Event endpoints
	router.HandleFunc("/events/{id}", handleEvent).Methods("POST", "PUT")
	router.HandleFunc("/events/{id}", getEvent).Methods("GET")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	// User availability endpoints
	router.HandleFunc("/events/{id}/availability/{user_id}", handleUserAvailability).Methods("POST", "PUT")
	router.HandleFunc("/events/{id}/availability/{user_id}", deleteUserAvailability).Methods("DELETE")

	// Recommendation endpoint
	router.HandleFunc("/events/{id}/recommendations", getRecommendations).Methods("GET")

	fmt.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
