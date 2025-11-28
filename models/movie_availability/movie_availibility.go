package movieavailability

import "github.com/google/uuid"

type MovieAvailability struct {
	MovieID     uuid.UUID `json:"movie_id"`
	CountryCode string    `json:"country_code"`
}
