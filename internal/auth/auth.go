package auth

import (
	"fmt"
	"net/http"
)

func GetTokenFromHeader(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Assuming the format is "Bearer <token>"
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return authHeader[7:], nil
}