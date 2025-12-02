package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/winfr1th/mock-interview/internal/database"
	"github.com/winfr1th/mock-interview/internal/handler"
	"github.com/winfr1th/mock-interview/internal/middleware"
	"github.com/winfr1th/mock-interview/internal/repository"
)

func main() {
	// Create context
	ctx := context.Background()

	// Establish database connection
	db, err := database.NewConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseConnection(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	saveMoviesRepo := repository.NewSaveMoviesRepository(db)

	// Setup router
	router := mux.NewRouter()

	// Public endpoints (no auth required)
	router.HandleFunc("/register", handler.Register(userRepo)).Methods("POST")
	router.HandleFunc("/genres", handler.ListGenres(genreRepo)).Methods("GET")
	router.HandleFunc("/movies", handler.ListMovies(movieRepo)).Methods("GET")

	// Protected endpoints - require API key authentication
	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.APIKeyAuth(userRepo))

	// User endpoints
	protectedRouter.HandleFunc("/users", handler.CreateUser(userRepo)).Methods("POST")
	protectedRouter.HandleFunc("/users/{id}", handler.GetUserByID(userRepo)).Methods("GET")

	// Saved movies endpoints
	protectedRouter.HandleFunc("/users/{user_id}/movies", handler.ListSavedMovies(saveMoviesRepo, movieRepo)).Methods("GET")
	protectedRouter.HandleFunc("/users/{user_id}/movies", handler.SaveMovie(saveMoviesRepo, movieRepo)).Methods("POST")
	protectedRouter.HandleFunc("/users/{user_id}/movies/{movie_id}", handler.RemoveSavedMovie(saveMoviesRepo)).Methods("DELETE")

	// Start server
	log.Println("Server starting on :8080")

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down server...")
}
