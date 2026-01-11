package auth

import (
	"errors"
	"net/http"
)

func GetRefreshToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	if len(authHeader) < 8 || authHeader[:8] != "Refresh " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[8:], nil
}
