-- Seed data for development/testing
-- This script creates a test user with a known API key

-- Insert test user with known API key
-- API Key: 550e8400-e29b-41d4-a716-446655440000
INSERT INTO users (id, name, date_of_birth, api_key_hash)
VALUES (
    '550e8400-e29b-41d4-a716-446655440001',
    'Test User',
    '1990-01-01',
    '550e8400-e29b-41d4-a716-446655440000'  -- Plain API key (not hashed)
)
ON CONFLICT (id) DO NOTHING;

-- Insert sample genres
INSERT INTO genres (id, name)
VALUES
    ('550e8400-e29b-41d4-a716-446655440010', 'Action'),
    ('550e8400-e29b-41d4-a716-446655440011', 'Comedy'),
    ('550e8400-e29b-41d4-a716-446655440012', 'Drama'),
    ('550e8400-e29b-41d4-a716-446655440013', 'Horror'),
    ('550e8400-e29b-41d4-a716-446655440014', 'Sci-Fi')
ON CONFLICT (id) DO NOTHING;

-- Insert sample countries
INSERT INTO countries (code, name)
VALUES
    ('US', 'United States'),
    ('GB', 'United Kingdom'),
    ('CA', 'Canada'),
    ('BR', 'Brazil'),
    ('FR', 'France')
ON CONFLICT (code) DO NOTHING;

-- Insert sample movies
INSERT INTO movies (id, title, year, genre_id)
VALUES
    ('550e8400-e29b-41d4-a716-446655440020', 'The Matrix', 1999, '550e8400-e29b-41d4-a716-446655440013'),
    ('550e8400-e29b-41d4-a716-446655440021', 'Inception', 2010, '550e8400-e29b-41d4-a716-446655440013'),
    ('550e8400-e29b-41d4-a716-446655440022', 'The Dark Knight', 2008, '550e8400-e29b-41d4-a716-446655440010'),
    ('550e8400-e29b-41d4-a716-446655440023', 'Pulp Fiction', 1994, '550e8400-e29b-41d4-a716-446655440010'),
    ('550e8400-e29b-41d4-a716-446655440024', 'The Shawshank Redemption', 1994, '550e8400-e29b-41d4-a716-446655440012')
ON CONFLICT (id) DO NOTHING;

-- Insert movie availability (all movies available in US and GB)
INSERT INTO movie_availability (movie_id, country_code)
SELECT m.id, c.code
FROM movies m
CROSS JOIN countries c
WHERE c.code IN ('US', 'GB')
ON CONFLICT (movie_id, country_code) DO NOTHING;

