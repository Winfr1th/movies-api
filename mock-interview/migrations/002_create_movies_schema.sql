-- Create genres table based on Genre model
-- Model fields: ID (uuid.UUID), Name (string)
CREATE TABLE IF NOT EXISTS genres (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

-- Create countries table based on Country model
-- Model fields: Code (string), Name (string)
-- Code is ISO-3166-1 alpha-2 format (2 characters)
CREATE TABLE IF NOT EXISTS countries (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- Create movies table based on Movie model
-- Model fields: ID (uuid.UUID), Title (string), Year (int), GenreID (uuid.UUID)
CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    year INTEGER NOT NULL,
    genre_id UUID NOT NULL,
    FOREIGN KEY (genre_id) REFERENCES genres(id)
);

-- Create movie_availability table based on MovieAvailability model
-- Junction table: MovieID (uuid.UUID), CountryCode (string)
CREATE TABLE IF NOT EXISTS movie_availability (
    movie_id UUID NOT NULL,
    country_code TEXT NOT NULL,
    PRIMARY KEY (movie_id, country_code),
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (country_code) REFERENCES countries(code) ON DELETE CASCADE
);

-- Create save_movies table based on SaveMovies model
-- Junction table: UserID (uuid.UUID), MovieID (uuid.UUID), DateAdded (time.Time)
CREATE TABLE IF NOT EXISTS save_movies (
    user_id UUID NOT NULL,
    movie_id UUID NOT NULL,
    date_added TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, movie_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);

-- Create indexes for better query performance

-- Index on movies.genre_id for filtering by genre
CREATE INDEX IF NOT EXISTS idx_movies_genre_id ON movies(genre_id);

-- Index on movies.year for sorting by year
CREATE INDEX IF NOT EXISTS idx_movies_year ON movies(year);

-- Index on movie_availability.country_code for filtering by country
CREATE INDEX IF NOT EXISTS idx_movie_availability_country_code ON movie_availability(country_code);

-- Index on save_movies.date_added for sorting saved movies
CREATE INDEX IF NOT EXISTS idx_save_movies_date_added ON save_movies(date_added);

-- Index on save_movies.user_id for faster lookups of user's saved movies
CREATE INDEX IF NOT EXISTS idx_save_movies_user_id ON save_movies(user_id);

