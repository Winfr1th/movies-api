package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/winfr1th/mock-interview/internal/repository"
	"github.com/winfr1th/mock-interview/internal/utils"
)

const (
	ErrorCodeDuplicateSave        = "DUPLICATE_SAVE"
	ErrorCodeUnavailableInCountry = "UNAVAILABLE_IN_COUNTRY"
	ErrorCodeNotSaved             = "NOT_SAVED"
)

// ListSavedMovies handles GET /users/{user_id}/movies - List saved movies by user
func ListSavedMovies(saveRepo repository.SaveMoviesRepository, movieRepo repository.MovieRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
				"Method not allowed", nil)
			return
		}

		// Parse user_id from URL path
		vars := mux.Vars(r)
		userIDStr := vars["user_id"]
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID",
				"Invalid user ID format", nil)
			return
		}

		// Parse and validate country parameter (required)
		countryCode := strings.TrimSpace(r.URL.Query().Get("country"))
		if countryCode == "" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_COUNTRY",
				"Country parameter is required", nil)
			return
		}
		countryCode = strings.ToUpper(countryCode)
		if len(countryCode) != 2 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_COUNTRY_CODE",
				"Invalid country code: must be ISO-3166-1 alpha-2 format (2 characters)", nil)
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

		// Parse sort parameter (optional, default: -date_added)
		sortBy := strings.TrimSpace(r.URL.Query().Get("sort"))
		if sortBy == "" {
			sortBy = "-date_added" // Default: newest first
		} else {
			// Validate sort parameter
			if sortBy != "date_added" && sortBy != "-date_added" {
				utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SORT_PARAMETER",
					"Invalid sort parameter: must be 'date_added' or '-date_added'", nil)
				return
			}
		}

		// Get saved movies from repository
		movies, total, err := saveRepo.ListSavedMovies(r.Context(), userID, countryCode, page, pageSize, sortBy)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to fetch saved movies: "+err.Error(), nil)
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

// SaveMovie handles POST /users/{user_id}/movies - Save a movie for a user
func SaveMovie(saveRepo repository.SaveMoviesRepository, movieRepo repository.MovieRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
				"Method not allowed", nil)
			return
		}

		// Parse user_id from URL path
		vars := mux.Vars(r)
		userIDStr := vars["user_id"]
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID",
				"Invalid user ID format", nil)
			return
		}

		// Parse and validate country parameter (required)
		countryCode := strings.TrimSpace(r.URL.Query().Get("country"))
		if countryCode == "" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_COUNTRY",
				"Country parameter is required", nil)
			return
		}
		countryCode = strings.ToUpper(countryCode)
		if len(countryCode) != 2 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_COUNTRY_CODE",
				"Invalid country code: must be ISO-3166-1 alpha-2 format (2 characters)", nil)
			return
		}

		// Parse request body
		var req struct {
			MovieID string `json:"movie_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST",
				"Invalid request body", nil)
			return
		}

		// Parse and validate movie_id
		movieID, err := uuid.Parse(req.MovieID)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_MOVIE_ID",
				"Invalid movie_id: must be a valid UUID", nil)
			return
		}

		// Validate movie exists
		_, err = movieRepo.GetMovieByID(r.Context(), movieID)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "MOVIE_NOT_FOUND",
				"Movie not found", nil)
			return
		}

		// Validate movie is available in country
		available, err := movieRepo.IsMovieAvailableInCountry(r.Context(), movieID, countryCode)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to check movie availability: "+err.Error(), nil)
			return
		}
		if !available {
			// Return 422 Unprocessable Entity with error code
			utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, ErrorCodeUnavailableInCountry,
				"Movie is not available in the specified country", nil)
			return
		}

		// Save the movie
		err = saveRepo.SaveMovie(r.Context(), userID, movieID)
		if err != nil {
			// Check for duplicate save error
			if strings.Contains(err.Error(), "already saved") {
				utils.WriteErrorResponse(w, http.StatusConflict, ErrorCodeDuplicateSave,
					"Movie is already saved", nil)
				return
			}
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to save movie: "+err.Error(), nil)
			return
		}

		// Get movie details to return
		movie, err := movieRepo.GetMovieByID(r.Context(), movieID)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to fetch movie details: "+err.Error(), nil)
			return
		}

		// Return movie detail
		response := map[string]interface{}{
			"id":    movie.ID.String(),
			"title": movie.Title,
			"year":  movie.Year,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// RemoveSavedMovie handles DELETE /users/{user_id}/movies/{movie_id} - Remove a saved movie
func RemoveSavedMovie(saveRepo repository.SaveMoviesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
				"Method not allowed", nil)
			return
		}

		// Parse user_id and movie_id from URL path
		vars := mux.Vars(r)
		userIDStr := vars["user_id"]
		movieIDStr := vars["movie_id"]

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID",
				"Invalid user ID format", nil)
			return
		}

		movieID, err := uuid.Parse(movieIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_MOVIE_ID",
				"Invalid movie ID format", nil)
			return
		}

		// Remove the saved movie
		err = saveRepo.RemoveSavedMovie(r.Context(), userID, movieID)
		if err != nil {
			// Check for "not saved" error
			if strings.Contains(err.Error(), "not saved") {
				utils.WriteErrorResponse(w, http.StatusNotFound, ErrorCodeNotSaved,
					"Movie is not saved", nil)
				return
			}
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
				"Failed to remove saved movie: "+err.Error(), nil)
			return
		}

		// Return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}
