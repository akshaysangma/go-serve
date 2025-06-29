// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteArticle(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetArticleByID(ctx context.Context, id uuid.UUID) (Article, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	ListArticles(ctx context.Context) ([]Article, error)
	ListArticlesByAuthorID(ctx context.Context, authorID uuid.UUID) ([]Article, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateArticle(ctx context.Context, arg UpdateArticleParams) (Article, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
