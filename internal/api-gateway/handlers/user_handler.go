package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/akshaysangma/go-serve/internal/api-gateway/middleware"
	"github.com/akshaysangma/go-serve/internal/api-gateway/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func CreateUserHandler(u *services.UserService, defaultLogger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.LoggerFromContext(r.Context(), defaultLogger)

		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode create user request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, articles, err := u.CreateUserTX(r.Context(), req.Username, req.Email)
		if err != nil {
			logger.Error("Failed to create user and article transactionally", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user":    user,
			"article": articles,
		})

		logger.Info("User and default article created successfully", zap.String("user_id", user.ID.String()))
	}
}

func GetUserByIDHandler(u *services.UserService, defaultLogger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.LoggerFromContext(r.Context(), defaultLogger)
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			logger.Error("Failed to retrive valid User ID from request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := u.GetUserByID(r.Context(), id)
		if err != nil {
			logger.Error("Failed to get user bu ID", zap.Error(err), zap.String("user_id", id.String()))
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "User Not Found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)

		logger.Info("User retrieved successfully", zap.String("user_id", user.ID.String()))
	}
}

func ListUsersHandler(u *services.UserService, defaultLogger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.LoggerFromContext(r.Context(), defaultLogger)

		users, err := u.ListUsers(r.Context())
		if err != nil {
			logger.Error("Fail to fetch all users", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)

		logger.Info("All users retrieved successfully", zap.Int("count", len(users)))
	}
}
