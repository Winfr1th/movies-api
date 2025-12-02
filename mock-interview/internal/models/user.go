package model

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth string    `json:"date_of_birth"`
	APIKeyHash  string    `json:"-"` // Don't expose hash in JSON responses
}

type RegisterRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
}

type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id"`
	APIKey string    `json:"api_key"` // Only returned once during registration
}
