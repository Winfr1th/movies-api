-- Create users table based on User model
-- Model fields: ID (uuid.UUID), Name (string), DateOfBirth (string), APIKeyHash (string)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    api_key_hash TEXT
);

-- Create index on api_key_hash for faster lookups during authentication
CREATE INDEX IF NOT EXISTS idx_users_api_key_hash ON users(api_key_hash);