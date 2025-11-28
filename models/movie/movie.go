package movie

import "github.com/google/uuid"

type Movie struct {
	ID uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Year  int       `json:"year"`
	GenreID uuid.UUID `json:"genre_id"`
}