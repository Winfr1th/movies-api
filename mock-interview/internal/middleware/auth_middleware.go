package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/winfr1th/mock-interview/internal/repository"
	"github.com/winfr1th/mock-interview/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

// APIKeyAuth middleware validates API key and adds user to context
func APIKeyAuth(repo repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get API key from header
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				// Try Authorization header with Bearer token
				authHeader := r.Header.Get("Authorization")
				if authHeader != "" {
					parts := strings.Split(authHeader, " ")
					if len(parts) == 2 && parts[0] == "Bearer" {
						apiKey = parts[1]
					}
				}
			}

			if apiKey == "" {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED",
					"API key required", nil)
				return
			}

			// Find user by API key (plain comparison)
			user, err := repo.FindUserByAPIKey(r.Context(), apiKey)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED",
					"Invalid API key", nil)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from request context
func GetUserID(r *http.Request) interface{} {
	return r.Context().Value(UserIDKey)
}
