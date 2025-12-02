package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/winfr1th/mock-interview/internal/models"
)

type SaveMoviesRepository interface {
	ListSavedMovies(ctx context.Context, userID uuid.UUID, countryCode string, page, pageSize int, sortBy string) ([]model.Movie, int, error)
	SaveMovie(ctx context.Context, userID, movieID uuid.UUID) error
	RemoveSavedMovie(ctx context.Context, userID, movieID uuid.UUID) error
	IsMovieSaved(ctx context.Context, userID, movieID uuid.UUID) (bool, error)
}

type saveMoviesRepo struct {
	db *pgxpool.Pool
}

func NewSaveMoviesRepository(db *pgxpool.Pool) SaveMoviesRepository {
	return &saveMoviesRepo{
		db: db,
	}
}

func (r *saveMoviesRepo) ListSavedMovies(ctx context.Context, userID uuid.UUID, countryCode string, page, pageSize int, sortBy string) ([]model.Movie, int, error) {
	// Validate and set sort order
	var sortClause string
	switch sortBy {
	case "date_added":
		sortClause = "ORDER BY sm.date_added ASC"
	case "-date_added":
		sortClause = "ORDER BY sm.date_added DESC"
	default:
		sortClause = "ORDER BY sm.date_added DESC" // Default: -date_added (newest first)
	}

	// Build WHERE clause - filter by user_id and country
	whereClause := "WHERE sm.user_id = $1 AND ma.country_code = $2"
	args := []interface{}{userID, strings.ToUpper(countryCode)}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT sm.movie_id)
		FROM save_movies sm
		INNER JOIN movie_availability ma ON sm.movie_id = ma.movie_id
		%s
	`, whereClause)

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Build main query - only return movies available in the specified country
	query := fmt.Sprintf(`
		SELECT DISTINCT m.id, m.title, m.year, m.genre_id
		FROM save_movies sm
		INNER JOIN movies m ON sm.movie_id = m.id
		INNER JOIN movie_availability ma ON m.id = ma.movie_id
		%s
		%s
		LIMIT $3 OFFSET $4
	`, whereClause, sortClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.GenreID); err != nil {
			return nil, 0, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (r *saveMoviesRepo) SaveMovie(ctx context.Context, userID, movieID uuid.UUID) error {
	// Check if movie is already saved (duplicate prevention)
	isSaved, err := r.IsMovieSaved(ctx, userID, movieID)
	if err != nil {
		return err
	}
	if isSaved {
		return errors.New("movie already saved")
	}

	// Insert the saved movie
	query := `INSERT INTO save_movies (user_id, movie_id, date_added) VALUES ($1, $2, CURRENT_TIMESTAMP)`
	_, err = r.db.Exec(ctx, query, userID, movieID)
	if err != nil {
		// Check for unique constraint violation (additional safety)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return errors.New("movie already saved")
		}
		return err
	}

	return nil
}

func (r *saveMoviesRepo) RemoveSavedMovie(ctx context.Context, userID, movieID uuid.UUID) error {
	query := `DELETE FROM save_movies WHERE user_id = $1 AND movie_id = $2`
	result, err := r.db.Exec(ctx, query, userID, movieID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("movie not saved")
	}

	return nil
}

func (r *saveMoviesRepo) IsMovieSaved(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM save_movies WHERE user_id = $1 AND movie_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, movieID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
