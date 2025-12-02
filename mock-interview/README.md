# Mock Interview API

A RESTful API built with Go for managing users with API key-based authentication.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Database Setup](#database-setup)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [Project Structure](#project-structure)

## Features

- ✅ User registration with API key generation
- ✅ API key-based authentication
- ✅ PostgreSQL database integration
- ✅ Connection pooling for optimal performance
- ✅ Graceful shutdown handling
- ✅ Simple UUID-based API keys

## Prerequisites

- Go 1.25.4 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd mock-interview
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o api main.go
```

## Database Setup

1. Create a PostgreSQL database:
```bash
createdb mock_interview
```

Or using psql:
```sql
CREATE DATABASE mock_interview;
```

2. Run the migrations:
```bash
psql -U postgres -d mock_interview -f migrations/001_create_users_table.sql
psql -U postgres -d mock_interview -f migrations/002_create_movies_schema.sql
```

3. (Optional) Load seed data for testing:
```bash
psql -U postgres -d mock_interview -f scripts/seed_data.sql
```

## Configuration

The application uses environment variables for configuration:

### Database Connection

Set the `DATABASE_URL` environment variable:

```bash
export DATABASE_URL="postgres://username:password@localhost:5432/mock_interview?sslmode=disable"
```

If `DATABASE_URL` is not set, the application defaults to:
```
postgres://postgres:postgres@localhost:5432/mock_interview?sslmode=disable
```

### Server Port

The server runs on port `8080` by default. To change it, modify `main.go`.

## Running the Application

1. Set the database URL (if different from default):
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/mock_interview?sslmode=disable"
```

2. Run the application:
```bash
go run main.go
```

Or if you built it:
```bash
./api
```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Endpoints

#### Register User
Register a new user and receive an API key.

**Endpoint:** `POST /register`

**Authentication:** Not required

#### List Genres
Get a paginated list of all genres.

**Endpoint:** `GET /genres`

**Authentication:** Not required

**Query Parameters:**
- `page` (optional, default: 1) - Page number (1-based)
- `page_size` (optional, default: 20, max: 100) - Number of items per page

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "name": "Action"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 5
}
```

#### List Movies
Get a paginated list of movies with optional filtering and sorting.

**Endpoint:** `GET /movies`

**Authentication:** Not required

**Query Parameters:**
- `country` (optional) - ISO-3166-1 alpha-2 country code (e.g., US, BR)
- `genre` (optional) - Genre UUID
- `page` (optional, default: 1) - Page number (1-based)
- `page_size` (optional, default: 20, max: 100) - Number of items per page
- `sort` (optional, default: "-year") - Sort order: `"year"` (ascending) or `"-year"` (descending)

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "title": "The Matrix",
      "year": 1999
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 100
}
```

**Error Responses:**
- `400 Bad Request` - Invalid query parameters

### Protected Endpoints

All protected endpoints require API key authentication. See [Authentication](#authentication) section.

#### Register User (Alternative)
Register a new user and receive an API key.

**Endpoint:** `POST /register`

**Request Body:**
```json
{
  "name": "John Doe",
  "date_of_birth": "1990-01-01"
}
```

**Response:** `201 Created`
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "api_key": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body or missing required fields
- `500 Internal Server Error` - Server error during user creation

### Protected Endpoints

All protected endpoints require API key authentication. See [Authentication](#authentication) section.

#### Get User by ID
Retrieve a user by their ID.

**Endpoint:** `GET /users/{id}`

**Headers:**
```
X-API-Key: <your-api-key-uuid>
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "date_of_birth": "1990-01-01"
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid API key
- `404 Not Found` - User not found

#### Get All Users
Retrieve all users (currently returns not implemented).

**Endpoint:** `GET /users`

**Headers:**
```
X-API-Key: <your-api-key-uuid>
```

#### Create User
Create a new user (currently returns not implemented).

**Endpoint:** `POST /users`

**Headers:**
```
X-API-Key: <your-api-key-uuid>
```

## Authentication

The API uses API key-based authentication for protected endpoints.

### Seed Credentials (Development/Testing)

For development and testing, you can use the seed data which includes a test user:

**API Key:** `550e8400-e29b-41d4-a716-446655440000`  
**User ID:** `550e8400-e29b-41d4-a716-446655440001`

To load seed data:
```bash
psql -U postgres -d mock_interview -f scripts/seed_data.sql
```

**Note:** This seed credential is for development only. In production, always register new users and keep API keys secure.

### Getting an API Key

1. Register a new user via `POST /register`
2. Save the `api_key` from the response (it's only shown once!)

### Using the API Key

You can provide the API key in two ways:

#### Method 1: X-API-Key Header (Recommended)
```bash
curl -H "X-API-Key: <your-api-key-uuid>" http://localhost:8080/users/{id}
```

#### Method 2: Authorization Bearer Token
```bash
curl -H "Authorization: Bearer <your-api-key-uuid>" http://localhost:8080/users/{id}
```

### Security Notes

- The API key is only returned once during registration
- Store your API key securely - if lost, you'll need to register a new user
- API keys are validated on every protected request

## Examples

### Register a New User

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "date_of_birth": "1985-05-15"
  }'
```

**Response:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "api_key": "550e8400-e29b-41d4-a716-446655440001"
}
```

### Get User by ID

```bash
curl -X GET http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "X-API-Key: 550e8400-e29b-41d4-a716-446655440001"
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Jane Smith",
  "date_of_birth": "1985-05-15"
}
```

### Using Bearer Token

```bash
curl -X GET http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer 550e8400-e29b-41d4-a716-446655440001"
```

#### List Saved Movies
Get a paginated list of movies saved by a user.

**Endpoint:** `GET /users/{user_id}/movies`

**Authentication:** Required

**Query Parameters:**
- `country` (required) - ISO-3166-1 alpha-2 country code (e.g., US, BR)
- `page` (optional, default: 1) - Page number (1-based)
- `page_size` (optional, default: 20, max: 100) - Number of items per page
- `sort` (optional, default: "-date_added") - Sort order: `"date_added"` (ascending) or `"-date_added"` (descending)

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "title": "The Matrix",
      "year": 1999
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 10
}
```

#### Save Movie
Save a movie for a user.

**Endpoint:** `POST /users/{user_id}/movies?country={country_code}`

**Authentication:** Required

**Query Parameters:**
- `country` (required) - ISO-3166-1 alpha-2 country code

**Request Body:**
```json
{
  "movie_id": "550e8400-e29b-41d4-a716-446655440020"
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440020",
  "title": "The Matrix",
  "year": 1999
}
```

**Error Responses:**
- `409 Conflict` - Movie already saved (error code: `DUPLICATE_SAVE`)
- `422 Unprocessable Entity` - Movie not available in country (error code: `UNAVAILABLE_IN_COUNTRY`)
- `404 Not Found` - Movie not found

#### Remove Saved Movie
Remove a saved movie for a user.

**Endpoint:** `DELETE /users/{user_id}/movies/{movie_id}`

**Authentication:** Required

**Response:** `204 No Content`

**Error Responses:**
- `404 Not Found` - Movie not saved (error code: `NOT_SAVED`)

## Error Handling

The API returns standard HTTP status codes:

| Status Code | Description |
|------------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid API key |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Duplicate resource |
| 422 | Unprocessable Entity - Business logic error |
| 500 | Internal Server Error |

### Error Response Format

All errors follow a standardized format:

```json
{
  "error": {
    "code": "STRING_CODE",
    "message": "Human readable error message",
    "details": {}
  }
}
```

**Example Error Responses:**

```json
{
  "error": {
    "code": "DUPLICATE_SAVE",
    "message": "Movie is already saved",
    "details": {}
  }
}
```

```json
{
  "error": {
    "code": "UNAVAILABLE_IN_COUNTRY",
    "message": "Movie is not available in the specified country",
    "details": {}
  }
}
```

```json
{
  "error": {
    "code": "NOT_SAVED",
    "message": "Movie is not saved",
    "details": {}
  }
}
```

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid API key",
    "details": {}
  }
}
```

## Project Structure

```
mock-interview/
├── main.go                          # Application entry point
├── go.mod                           # Go module file
├── go.sum                           # Go dependencies checksum
├── README.md                        # This file
├── migrations/                      # Database migrations
│   ├── 001_create_users_table.sql
└── internal/
    ├── auth/                        # Authentication utilities
    │   └── apikey.go                # API key generation
    ├── database/                    # Database connection
    │   └── database.go              # Connection pool management
    ├── handler/                     # HTTP handlers
    │   ├── auth_handler.go          # Registration handler
    │   └── user_handler.go          # User CRUD handlers
    ├── middleware/                  # HTTP middleware
    │   └── auth_middleware.go       # API key authentication middleware
    ├── models/                      # Data models
    │   ├── user.go                  # User model
    │   └── ...                      # Other models
    └── repository/                  # Data access layer
        └── user_repo.go            # User repository implementation
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o api main.go
```

### Code Structure

- **Models**: Define data structures
- **Repository**: Database operations (PostgreSQL)
- **Handlers**: HTTP request/response handling
- **Middleware**: Request processing (authentication, etc.)
- **Auth**: Security utilities (API key generation)

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]

