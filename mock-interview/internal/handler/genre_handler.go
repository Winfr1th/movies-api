package handler

import (
	"encoding/json"
	"net/http"

	"github.com/winfr1th/mock-interview/internal/repository"
	"github.com/winfr1th/mock-interview/internal/utils"
)

// ListGenres handles GET /genres - List all genres with pagination
func ListGenres(repo repository.GenreRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
				"Method not allowed", nil)
			return
		}

		// Parse pagination parameters
		page, pageSize, err := utils.ParsePaginationParams(r)
		if err != nil {
			if pagErr, ok := err.(*utils.PaginationError); ok {
				utils.WriteErrorResponse(w, http.StatusBadRequest, pagErr.Code, pagErr.Message, nil)
			} else {
				utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_PARAMETER", err.Error(), nil)
			}
			return
		}

		// Get genres from repository
		genres, total, err := repo.ListGenres(r.Context(), page, pageSize)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to fetch genres: "+err.Error(), nil)
			return
		}

		// Create paginated response
		response := utils.CreatePagedResponse(genres, total, page, pageSize)

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
