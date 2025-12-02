package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/winfr1th/mock-interview/internal/auth"
	model "github.com/winfr1th/mock-interview/internal/models"
	"github.com/winfr1th/mock-interview/internal/repository"
	"github.com/winfr1th/mock-interview/internal/utils"
)

// Register handles user registration and returns an API key
func Register(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
				"Method not allowed", nil)
			return
		}

		var req model.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST",
				"Invalid request body", nil)
			return
		}

		// Validate required fields
		if req.Name == "" || req.DateOfBirth == "" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_FIELDS",
				"Name and date_of_birth are required", nil)
			return
		}

		// Generate API key
		apiKey, err := auth.GenerateAPIKey()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to generate API key", nil)
			return
		}

		// Create user
		user := model.User{
			ID:          uuid.New(),
			Name:        req.Name,
			DateOfBirth: req.DateOfBirth,
			APIKeyHash:  apiKey, // Store plain API key
		}

		if err := repo.CreateUser(r.Context(), user); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to create user: "+err.Error(), nil)
			return
		}

		// Return response with API key (only time it's returned)
		response := model.RegisterResponse{
			UserID: user.ID,
			APIKey: apiKey,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
