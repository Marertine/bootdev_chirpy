package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// MakeRefreshToken generates a secure random refresh token.
func MakeRefreshToken() (string, error) {
	const tokenSize = 32 // 32 bytes = 256 bits

	// Create a byte slice to hold the random bytes
	tokenBytes := make([]byte, tokenSize)

	// Read random bytes into the slice
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to a hex string
	refreshToken := hex.EncodeToString(tokenBytes)

	return refreshToken, nil
}
