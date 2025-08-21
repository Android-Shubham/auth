package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Android-Shubham/auth/internal/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apiConfig *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var p payload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = validate(p.Name, p.Email, p.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = apiConfig.db.GetUserByEmail(r.Context(), p.Email)
	if err == nil {
		respondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user, err := apiConfig.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      p.Name,
		Email:     p.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	secret := apiConfig.secret
	token, err := generateJWT([]byte(secret), user)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	response := map[string]interface{}{
		"token": token,
		"name":  user.Name,
		"email": user.Email,
	}
	responseWithJson(w, http.StatusCreated, response)
}

func generateJWT(secret []byte, user database.User) (string, error) {
	clains := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return clains.SignedString(secret)
}

func (apiConfig *ApiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var p payload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := apiConfig.db.GetUserByEmail(r.Context(), p.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	secret := apiConfig.secret
	token, err := generateJWT([]byte(secret), user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := map[string]interface{}{
		"token": token,
		"name":  user.Name,
		"email": user.Email,
	}
	responseWithJson(w, http.StatusOK, response)
}

func (apiConfig *ApiConfig) getUserDetails(w http.ResponseWriter, r *http.Request, user database.User) {
	response := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"id":    user.ID.String(),
	}
	responseWithJson(w, http.StatusOK, response)
}
