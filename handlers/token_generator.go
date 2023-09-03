package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// GenerateRandomToken generates a random token of the specified length and hashes it with SHA-256.
func GenerateRandomToken(length int) (string, error) {
	// Ensure the token length is at least 16 characters
	if length < 32 {
		return "", errors.New("token length is too short")
	}

	// Generates random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Hashing the random bytes with SHA-256
	hashedBytes := sha256.Sum256(randomBytes)

	// Encoding the hashed bytes as a base64 string
	token := base64.RawURLEncoding.EncodeToString(hashedBytes[:])

	return token, nil
}
