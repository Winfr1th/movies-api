package save

import (
	"time"

	"github.com/google/uuid"
)

type SaveMovies struct {
	UserID    uuid.UUID `json:"user_id"`
	MovieID   uuid.UUID `json:"movie_id"`
	DateAdded time.Time `json:"date_added"`
}
