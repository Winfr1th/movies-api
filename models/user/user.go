package user

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
}