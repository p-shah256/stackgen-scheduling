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
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Database connection failed"})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "meetingScheduler")
	
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

	fmt.Println("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
