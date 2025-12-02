package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/winfr1th/mock-interview/internal/models"
)

type MovieRepository interface {
	ListMovies(ctx context.Context, countryCode *string, genreID *uuid.UUID, page, pageSize int, sortBy string) ([]model.Movie, int, error)
	GetMovieByID(ctx context.Context, movieID uuid.UUID) (model.Movie, error)
	IsMovieAvailableInCountry(ctx context.Context, movieID uuid.UUID, countryCode string) (bool, error)
}

type movieRepo struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) MovieRepository {
	return &movieRepo{
		db: db,
	}
}

func (r *movieRepo) ListMovies(ctx context.Context, countryCode *string, genreID *uuid.UUID, page, pageSize int, sortBy string) ([]model.Movie, int, error) {
	// Build WHERE clause dynamically based on filters
	var whereConditions []string
	var args []interface{}
	argIndex := 1

	// Filter by country (requires join with movie_availability)
	needsJoin := countryCode != nil && *countryCode != ""

	// Filter by genre
	if genreID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("m.genre_id = $%d", argIndex))
		args = append(args, *genreID)
		argIndex++
	}

	// Filter by country
	if needsJoin {
		whereConditions = append(whereConditions, fmt.Sprintf("ma.country_code = $%d", argIndex))
		args = append(args, strings.ToUpper(*countryCode))
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Validate and set sort order
	var sortClause string
	switch sortBy {
	case "year":
		sortClause = "ORDER BY m.year ASC"
	case "-year":
		sortClause = "ORDER BY m.year DESC"
	default:
		sortClause = "ORDER BY m.year DESC" // Default: -year (newest first)
	}

	// Build FROM clause with optional join
	fromClause := "FROM movies m"
	if needsJoin {
		fromClause = "FROM movies m INNER JOIN movie_availability ma ON m.id = ma.movie_id"
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT m.id) %s %s", fromClause, whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Build main query
	query := fmt.Sprintf(`
		SELECT DISTINCT m.id, m.title, m.year, m.genre_id 
		%s 
		%s 
		%s 
		LIMIT $%d OFFSET $%d
	`, fromClause, whereClause, sortClause, argIndex, argIndex+1)

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

func (r *movieRepo) GetMovieByID(ctx context.Context, movieID uuid.UUID) (model.Movie, error) {
	query := `SELECT id, title, year, genre_id FROM movies WHERE id = $1`
	var movie model.Movie
	err := r.db.QueryRow(ctx, query, movieID).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.GenreID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Movie{}, errors.New("movie not found")
		}
		return model.Movie{}, err
	}

	return movie, nil
}

func (r *movieRepo) IsMovieAvailableInCountry(ctx context.Context, movieID uuid.UUID, countryCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM movie_availability WHERE movie_id = $1 AND country_code = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, movieID, strings.ToUpper(countryCode)).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
