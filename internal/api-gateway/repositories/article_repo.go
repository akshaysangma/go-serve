package repositories

import (
	"context"
	"fmt"

	db "github.com/akshaysangma/go-serve/internal/database/postgres/sqlc"
	"github.com/google/uuid"
)

type Article = db.Article

type CreateArticleParams struct {
	Title    string
	Content  string
	AuthorID uuid.UUID
}

type UpdateArticleParams struct {
	ID      uuid.UUID
	Title   string
	Content string
}

type ArticleRepository interface {
	CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error)
	GetArticleByID(ctx context.Context, id uuid.UUID) (Article, error)
	ListArticles(ctx context.Context) ([]Article, error)
	ListArticlesByAuthorID(ctx context.Context, authorID uuid.UUID) ([]Article, error) // If you have this query
	UpdateArticle(ctx context.Context, arg UpdateArticleParams) (Article, error)
	DeleteArticle(ctx context.Context, id uuid.UUID) error
}

type postgresArticleRepository struct {
	queries *db.Queries
}

func NewArticleRepository(queries *db.Queries) ArticleRepository {
	return &postgresArticleRepository{
		queries: queries,
	}
}

// Implement methods from the ArticleRepository interface
func (r *postgresArticleRepository) CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error) {
	article, err := r.queries.CreateArticle(ctx, db.CreateArticleParams{
		Title:    arg.Title,
		Content:  arg.Content,
		AuthorID: arg.AuthorID,
	})
	if err != nil {
		return Article{}, fmt.Errorf("repo: failed to create article: %w", err)
	}
	return article, nil
}

func (r *postgresArticleRepository) GetArticleByID(ctx context.Context, id uuid.UUID) (Article, error) {
	article, err := r.queries.GetArticleByID(ctx, id)
	if err != nil {
		return Article{}, fmt.Errorf("repo: failed to get article by ID: %w", err)
	}
	return article, nil
}

func (r *postgresArticleRepository) ListArticles(ctx context.Context) ([]Article, error) {
	articles, err := r.queries.ListArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: failed to list articles: %w", err)
	}
	return articles, nil
}

func (r *postgresArticleRepository) ListArticlesByAuthorID(ctx context.Context, authorID uuid.UUID) ([]Article, error) {
	articles, err := r.queries.ListArticlesByAuthorID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("repo: failed to list articles by author ID: %w", err)
	}
	return articles, nil
}

func (r *postgresArticleRepository) UpdateArticle(ctx context.Context, arg UpdateArticleParams) (Article, error) {
	article, err := r.queries.UpdateArticle(ctx, db.UpdateArticleParams{
		ID:      arg.ID,
		Title:   arg.Title,
		Content: arg.Content,
	})
	if err != nil {
		return Article{}, fmt.Errorf("repo: failed to update article: %w", err)
	}
	return article, nil
}

func (r *postgresArticleRepository) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteArticle(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: failed to delete article: %w", err)
	}
	return nil
}
