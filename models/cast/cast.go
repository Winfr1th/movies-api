package cast

import "github.com/google/uuid"

type Cast struct {
	ID uuid.UUID `json:"id"`
	MovieID uuid.UUID `json:"movie_id"`
	ActorID uuid.UUID `json:"actor_id"`
}