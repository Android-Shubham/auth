package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, status int, message string) {
	if status > 499 {
		log.Printf("Server Error: %s", message)
	}

	type ErrorResponse struct {
		Message string `json:"message"`
	}

	responseWithJson(w, status, ErrorResponse{Message: message})
}
func responseWithJson(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
