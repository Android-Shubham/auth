package main

import "net/http"

func (apiConfig *ApiConfig) errorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Errors working fine.")
}