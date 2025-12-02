package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/winfr1th/mock-interview/internal/repository"
	"github.com/winfr1th/mock-interview/internal/utils"
)

// ListMovies handles GET /movies - List movies with filtering, sorting, and pagination
func ListMovies(repo repository.MovieRepository) http.HandlerFunc {
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

		// Parse and validate country filter (optional)
		var countryCode *string
		countryParam := strings.TrimSpace(r.URL.Query().Get("country"))
		if countryParam != "" {
			// Validate country code format (ISO-3166-1 alpha-2: 2 uppercase letters)
			countryParam = strings.ToUpper(countryParam)
			if len(countryParam) != 2 {
				utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_COUNTRY_CODE",
					"Invalid country code: must be ISO-3166-1 alpha-2 format (2 characters)", nil)
				return
			}
			countryCode = &countryParam
		}

		// Parse and validate genre filter (optional)
		var genreID *uuid.UUID
		genreParam := strings.TrimSpace(r.URL.Query().Get("genre"))
		if genreParam != "" {
			parsedGenreID, err := uuid.Parse(genreParam)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_GENRE_ID",
					"Invalid genre ID: must be a valid UUID", nil)
				return
			}
			genreID = &parsedGenreID
		}

		// Parse sort parameter (optional, default: -year)
		sortBy := strings.TrimSpace(r.URL.Query().Get("sort"))
		if sortBy == "" {
			sortBy = "-year" // Default: newest first
		} else {
			// Validate sort parameter
			if sortBy != "year" && sortBy != "-year" {
				utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SORT_PARAMETER",
					"Invalid sort parameter: must be 'year' or '-year'", nil)
				return
			}
		}

		// Get movies from repository
		movies, total, err := repo.ListMovies(r.Context(), countryCode, genreID, page, pageSize, sortBy)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to fetch movies: "+err.Error(), nil)
			return
		}

		// Create simplified movie response (only id, title, year)
		type MovieResponse struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Year  int    `json:"year"`
		}

		movieResponses := make([]MovieResponse, len(movies))
		for i, movie := range movies {
			movieResponses[i] = MovieResponse{
				ID:    movie.ID.String(),
				Title: movie.Title,
				Year:  movie.Year,
			}
		}

		// Create paginated response
		response := utils.CreatePagedResponse(movieResponses, total, page, pageSize)

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
