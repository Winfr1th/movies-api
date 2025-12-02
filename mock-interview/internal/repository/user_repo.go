package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/winfr1th/mock-interview/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	FindUserByID(ctx context.Context, id string) (model.User, error)
	FindUserByAPIKey(ctx context.Context, apiKey string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	UpdateAPIKey(ctx context.Context, userID uuid.UUID, apiKey string) error
	DeleteUser(ctx context.Context, id string) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (id, name, date_of_birth, api_key_hash) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.DateOfBirth, user.APIKeyHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) FindUserByID(ctx context.Context, id string) (model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return model.User{}, errors.New("invalid user ID format")
	}

	query := `SELECT id, name, date_of_birth, api_key_hash FROM users WHERE id = $1`
	var user model.User
	err = r.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Name, &user.DateOfBirth, &user.APIKeyHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepo) FindUserByAPIKey(ctx context.Context, apiKey string) (model.User, error) {
	query := `SELECT id, name, date_of_birth, api_key_hash FROM users WHERE api_key_hash = $1`
	var user model.User
	err := r.db.QueryRow(ctx, query, apiKey).Scan(&user.ID, &user.Name, &user.DateOfBirth, &user.APIKeyHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("invalid API key")
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user model.User) error {
	query := `UPDATE users SET name = $1, date_of_birth = $2 WHERE id = $3`
	result, err := r.db.Exec(ctx, query, user.Name, user.DateOfBirth, user.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepo) UpdateAPIKey(ctx context.Context, userID uuid.UUID, apiKey string) error {
	query := `UPDATE users SET api_key_hash = $1 WHERE id = $2`
	result, err := r.db.Exec(ctx, query, apiKey, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
