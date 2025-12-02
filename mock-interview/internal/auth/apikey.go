package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

// GenerateAPIKey generates a UUID-based API key
func GenerateAPIKey() (string, error) {
	return uuid.New().String(), nil
}

// HashAPIKey creates a SHA-256 hash of the API key for storage
func HashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

// VerifyAPIKey compares a plain API key with a stored hash
func VerifyAPIKey(apiKey, storedHash string) bool {
	hash := HashAPIKey(apiKey)
	return hash == storedHash
}
