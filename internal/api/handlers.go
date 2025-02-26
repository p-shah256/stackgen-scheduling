package api

import (
	"encoding/json"
	"exercise/internal/types"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Helper functions to reduce repetition
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func decodeJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func createEvnt(w http.ResponseWriter, r *http.Request) {
	var newEvent types.Event
	if err := decodeJSONBody(r, &newEvent); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	slog.Info("Creating event", "id", newEvent.EventId, "title", newEvent.Title)
	createdEvent, err := CreateEvent(r.Context(), newEvent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, createdEvent)
}

func updateEvntId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	var updatedEvent types.Event
	if err := decodeJSONBody(r, &updatedEvent); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedEvent.EventId = eventID
	updatedEvent, err := UpdateEvent(r.Context(), eventID, updatedEvent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, updatedEvent)
}

func deleteEvntId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	err := DeleteEvent(r.Context(), eventID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Event deleted successfully"})
}

func getEvntRecId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	event, err := GetEventRecommendations(r.Context(), eventID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, event)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser types.User
	if err := decodeJSONBody(r, &newUser); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	createdUser, err := CreateUser(r.Context(), newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, createdUser)
}

func updateUserAvail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var newAvail types.UsrAvailDate
	if err := decodeJSONBody(r, &newAvail); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newAvail.UserID = userID
	err := UpdateUserAvail(r.Context(), userID, newAvail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Availability updated successfully"})
}

func createUserAvail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var newAvail types.UsrAvailDate
	if err := decodeJSONBody(r, &newAvail); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newAvail.UserID = userID
	err := CreateUserAvail(r.Context(), userID, newAvail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Availability created successfully"})
}

func updateUserTimezone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var newTimezone string
	if err := decodeJSONBody(r, &newTimezone); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err := UpdateUserTimezone(r.Context(), userID, newTimezone)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Timezone updated successfully"})
}

func rmUserAvail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	err := DeleteUserAvail(r.Context(), userID, time.Now())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Availability deleted successfully"})
}

// Helper function to handle service errors
func handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		respondWithError(w, http.StatusNotFound, "Resource not found")
	case ErrInvalid:
		respondWithError(w, http.StatusBadRequest, "Invalid input")
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
