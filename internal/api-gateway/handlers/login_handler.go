package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/akshaysangma/go-serve/internal/api-gateway/middleware"
	"github.com/akshaysangma/go-serve/internal/api-gateway/services"
	"github.com/akshaysangma/go-serve/internal/common/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(config config.JWTConfig, u *services.UserService, defaultLogger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger := middleware.LoggerFromContext(r.Context(), defaultLogger)
		userIDStr := r.URL.Query().Get("user")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.Error("bad user ID provided", zap.Error(err), zap.String("user_id", userIDStr))
			http.Error(w, "Bad user ID provided", http.StatusBadRequest)
			return
		}

		// mimicing credentials validation and get user obj from third party or database
		user, err := u.GetUserByID(r.Context(), userID)
		if err != nil {
			logger.Error("user not exist", zap.Error(err), zap.String("user_id", userIDStr))
			http.Error(w, "user not exist", http.StatusNotFound)
			return
		}

		expirationTime := time.Now().Add(config.ExpirationDuration)
		claims := &middleware.AuthClaims{
			UserID:   user.ID.String(),
			Email:    user.Email,
			Username: user.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-serve",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Subject:   user.ID.String(),
				Audience:  []string{"web"},
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(config.Secret))
		if err != nil {
			logger.Error("unable to sign token", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		res := LoginResponse{
			Token: tokenStr,
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
