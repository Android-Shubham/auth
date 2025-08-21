package main

import (
	"net/http"

	"github.com/Android-Shubham/auth/internal/auth"
	"github.com/Android-Shubham/auth/internal/database"
	"github.com/golang-jwt/jwt/v4"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *ApiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetTokenFromHeader(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		secret := apiConfig.secret
		claims := &jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !tkn.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		email, ok := (*claims)["email"].(string)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "invalid user ID")
			return
		}

		user, err := apiConfig.db.GetUserByEmail(r.Context(), email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}
		handler(w, r, user)
	}
}
