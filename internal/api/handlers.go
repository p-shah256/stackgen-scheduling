package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func echoParams(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()

	response := map[string]any{
		"method": r.Method,
		"path":   r.URL.Path,
		"params": params,
		"query":  query,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAllEvnts(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func createEvnt(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func getEvntId(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func updateEvntId(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func deleteEvntId(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func getEvntRecId(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func getUserId(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func updateUserId(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}

func rmUserId(w http.ResponseWriter, r *http.Request) {
	echoParams(w, r)
}
