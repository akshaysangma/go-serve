package services

import (
	"context"
	"fmt"

	"github.com/akshaysangma/go-serve/internal/api-gateway/repositories"
	db "github.com/akshaysangma/go-serve/internal/database/postgres/sqlc"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ArticleService handles business logic for articles.
type ArticleService struct {
	articleRepo repositories.ArticleRepository
	logger      *zap.Logger
}

// NewArticleService creates a new ArticleService.
func NewArticleService(articleRepo repositories.ArticleRepository, logger *zap.Logger) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		logger:      logger,
	}
}

// CreateArticle creates a new article.
func (s *ArticleService) CreateArticle(ctx context.Context, title, content string, authorID uuid.UUID) (db.Article, error) {
	article, err := s.articleRepo.CreateArticle(ctx, repositories.CreateArticleParams{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
	})
	if err != nil {
		s.logger.Error("Service: Failed to create article via repository", zap.Error(err), zap.String("title", title))
		return db.Article{}, fmt.Errorf("could not create article: %w", err)
	}
	return article, nil
}

// GetArticleByID retrieves an article by ID.
func (s *ArticleService) GetArticleByID(ctx context.Context, id uuid.UUID) (db.Article, error) {
	article, err := s.articleRepo.GetArticleByID(ctx, id)
	if err != nil {
		s.logger.Error("Service: Failed to get article by ID via repository", zap.Error(err), zap.String("article_id", id.String()))
		return db.Article{}, fmt.Errorf("could not get article: %w", err)
	}
	return article, nil
}

// ListArticles lists all articles.
func (s *ArticleService) ListArticles(ctx context.Context) ([]db.Article, error) {
	articles, err := s.articleRepo.ListArticles(ctx)
	if err != nil {
		s.logger.Error("Service: Failed to list articles via repository", zap.Error(err))
		return nil, fmt.Errorf("could not list articles: %w", err)
	}
	return articles, nil
}

// UpdateArticle updates an existing article.
func (s *ArticleService) UpdateArticle(ctx context.Context, id uuid.UUID, title, content string) (db.Article, error) {
	article, err := s.articleRepo.UpdateArticle(ctx, repositories.UpdateArticleParams{
		ID:      id,
		Title:   title,
		Content: content,
	})
	if err != nil {
		s.logger.Error("Service: Failed to update article via repository", zap.Error(err), zap.String("article_id", id.String()))
		return db.Article{}, fmt.Errorf("could not update article: %w", err)
	}
	return article, nil
}

// DeleteArticle deletes an article by ID.
func (s *ArticleService) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	err := s.articleRepo.DeleteArticle(ctx, id)
	if err != nil {
		s.logger.Error("Service: Failed to delete article via repository", zap.Error(err), zap.String("article_id", id.String()))
		return fmt.Errorf("could not delete article: %w", err)
	}
	return nil
}
