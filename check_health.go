package main

import "net/http"

func (apiConfig *ApiConfig)  checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	responseWithJson(w, http.StatusOK, map[string]string{"status": "healthy"})
}