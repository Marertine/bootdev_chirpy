package auth

import (
	"errors"
	"net/http"
)

func getAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	if len(authHeader) < 7 || authHeader[:7] != "ApiKey " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[7:], nil
}
