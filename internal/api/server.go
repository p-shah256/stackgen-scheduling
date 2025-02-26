package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"exercise/internal/db"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

func RunServer() {
    ctx := context.Background()
    if err := db.Connect(ctx); err != nil {
        slog.Error("Failed to connect to database", "error", err)
        os.Exit(1)
    }

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

    slog.Info("Starting server on 0.0.0.0:8081")
    err := http.ListenAndServe("0.0.0.0:8081", router)
    if err != nil {
        slog.Error("Server error", "error", err)
        os.Exit(1)
    }
}
