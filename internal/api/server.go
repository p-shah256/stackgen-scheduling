package api

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

type Response struct {
    Message string `json:"message"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{Message: "Welcome to the API!"}
    json.NewEncoder(w).Encode(response)
}

func RunServer() {
    router := mux.NewRouter()

    router.HandleFunc("/", homeHandler).Methods("GET")

    // Add more routes here

    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}
