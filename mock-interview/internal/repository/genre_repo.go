package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/winfr1th/mock-interview/internal/models"
)

type GenreRepository interface {
	ListGenres(ctx context.Context, page, pageSize int) ([]model.Genre, int, error)
}

type genreRepo struct {
	db *pgxpool.Pool
}

func NewGenreRepository(db *pgxpool.Pool) GenreRepository {
	return &genreRepo{
		db: db,
	}
}

func (r *genreRepo) ListGenres(ctx context.Context, page, pageSize int) ([]model.Genre, int, error) {
	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM genres`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated genres
	query := `
		SELECT id, name 
		FROM genres 
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var genres []model.Genre
	for rows.Next() {
		var genre model.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, 0, err
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return genres, total, nil
}
