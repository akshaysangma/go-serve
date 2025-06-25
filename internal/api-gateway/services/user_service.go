package services

import (
	"context"
	"fmt"

	"github.com/akshaysangma/go-serve/internal/api-gateway/repositories"
	db "github.com/akshaysangma/go-serve/internal/database/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo    repositories.UserRepository
	articleRepo repositories.ArticleRepository
	dBConn      *pgxpool.Pool
	logger      *zap.Logger
}

func NewUserService(
	userRepo repositories.UserRepository,
	articleRepo repositories.ArticleRepository,
	dBConn *pgxpool.Pool,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		articleRepo: articleRepo,
		dBConn:      dBConn,
		logger:      logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username, email string) (db.User, error) {
	user, err := s.userRepo.CreateUser(ctx, repositories.CreateUserParams{
		Username: username,
		Email:    email,
	})
	if err != nil {
		s.logger.Error("Service: Failed to create user via repository", zap.Error(err), zap.String("username", username))
	}
	return user, nil
}

func (s *UserService) CreateUserTX(ctx context.Context, username, email string) (db.User, db.Article, error) {
	tx, err := s.dBConn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		s.logger.Error("Service: Failed to create transactions", zap.Error(err))
		return db.User{}, db.Article{}, fmt.Errorf("could not begin transaction %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			s.logger.Error("Service: Failed to rollback transactions", zap.Error(err))
		}
	}()

	txQueries := db.New(tx)
	txUserRepo := repositories.NewUserRepository(txQueries)
	txArcticleRepo := repositories.NewArticleRepository(txQueries)

	// Create user
	user, err := txUserRepo.CreateUser(ctx, repositories.CreateUserParams{
		Username: username,
		Email:    email,
	})
	if err != nil {
		s.logger.Error("Service: Failed to create user within transaction", zap.Error(err), zap.String("username", username))
		return db.User{}, db.Article{}, fmt.Errorf("fail to create user in transaction: %w", err)
	}

	article, err := txArcticleRepo.CreateArticle(ctx, repositories.CreateArticleParams{
		Title:    fmt.Sprintf("Welcome %s", username),
		Content:  fmt.Sprintf("Thank you for joining our platform, %s! This is your first article.", user.Username),
		AuthorID: user.ID,
	})
	if err != nil {
		s.logger.Error("Service: Failed to create default article within transaction", zap.Error(err), zap.String("user_id", user.ID.String()))
		return db.User{}, db.Article{}, fmt.Errorf("fail to create article: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("Service: Failed to commit transaction for user and article creation", zap.Error(err))
		return db.User{}, db.Article{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	s.logger.Info("User and default article created successfully in transaction",
		zap.String("user_id", user.ID.String()),
		zap.String("article_id", article.ID.String()),
	)

	return user, article, nil
}

// GetUserByID retrieves a user by ID.
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("Service: Failed to get user by ID via repository", zap.Error(err), zap.String("user_id", id.String()))
		return db.User{}, fmt.Errorf("could not get user: %w", err)
	}
	return user, nil
}

// ListUsers lists all users.
func (s *UserService) ListUsers(ctx context.Context) ([]db.User, error) {
	users, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		s.logger.Error("Service: Failed to list users via repository", zap.Error(err))
		return nil, fmt.Errorf("could not list users: %w", err)
	}
	return users, nil
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, username, email string) (db.User, error) {
	user, err := s.userRepo.UpdateUser(ctx, repositories.UpdateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
	})
	if err != nil {
		s.logger.Error("Service: Failed to update user via repository", zap.Error(err), zap.String("user_id", id.String()))
		return db.User{}, fmt.Errorf("could not update user: %w", err)
	}
	return user, nil
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		s.logger.Error("Service: Failed to delete user via repository", zap.Error(err), zap.String("user_id", id.String()))
		return fmt.Errorf("could not delete user: %w", err)
	}
	return nil
}
