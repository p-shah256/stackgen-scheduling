package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

func RunServer() {
	router := mux.NewRouter()

	router.HandleFunc("/events", getAllEvnts).Methods("GET")
	router.HandleFunc("/events", createEvnt).Methods("POST")
	router.HandleFunc("/events/{eventId}", getEvntId).Methods("GET")
	router.HandleFunc("/events/{eventId}", updateEvntId).Methods("PUT")
	router.HandleFunc("/events/{eventId}", deleteEvntId).Methods("DELETE")
	router.HandleFunc("/events/{eventId}/recommendation", getEvntRecId).Methods("GET")

	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{userId}", getUserId).Methods("GET")
	router.HandleFunc("/users/{userId}", updateUserId).Methods("PUT")
	router.HandleFunc("/users/{userId}", rmUserId).Methods("DELETE")

	slog.Info("Server starting", "port", 8080)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
