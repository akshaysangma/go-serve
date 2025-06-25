package repositories

import (
	"context"
	"fmt"

	db "github.com/akshaysangma/go-serve/internal/database/postgres/sqlc"
	"github.com/google/uuid"
)

type User = db.User

type CreateUserParams struct {
	Username string
	Email    string
}

type UpdateUserParams struct {
	ID       uuid.UUID
	Username string
	Email    string
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error) // Add this query if not already
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type postgresUserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &postgresUserRepository{
		queries: queries,
	}
}

// Implement methods from the UserRepository interface
func (r *postgresUserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username: arg.Username,
		Email:    arg.Email,
	})
	if err != nil {
		return User{}, fmt.Errorf("repo: failed to create user: %w", err)
	}
	return user, nil
}

func (r *postgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("repo: failed to get user by ID: %w", err)
	}
	return user, nil
}

func (r *postgresUserRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email) // Ensure you have this query in users.sql
	if err != nil {
		return User{}, fmt.Errorf("repo: failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *postgresUserRepository) ListUsers(ctx context.Context) ([]User, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: failed to list users: %w", err)
	}
	return users, nil
}

func (r *postgresUserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	user, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       arg.ID,
		Username: arg.Username,
		Email:    arg.Email,
	})
	if err != nil {
		return User{}, fmt.Errorf("repo: failed to update user: %w", err)
	}
	return user, nil
}

func (r *postgresUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: failed to delete user: %w", err)
	}
	return nil
}
