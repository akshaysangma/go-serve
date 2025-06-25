package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

const (
	AuthContextKey contextKey = "authtoken"
	authHeader     string     = "Authorization"
)

// Claims struct that extends jwt.RegisteredClaims
type AuthClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret []byte, defaultLogger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := LoggerFromContext(r.Context(), defaultLogger)
			token := GetTokenFromHeader(r)
			if token == "" {
				logger.Error("Invalid auth token received")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			jwtToken, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return jwtSecret, nil
			})
			if err != nil {
				logger.Error("Token validation error", zap.Error(err))
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			if claims, ok := jwtToken.Claims.(*AuthClaims); ok && jwtToken.Valid {
				reqWithAuth := r.WithContext(context.WithValue(r.Context(), AuthContextKey, claims.ID))
				next.ServeHTTP(w, reqWithAuth)
			} else {
				logger.Error("Invalid or expired auth token received")
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			}

		})
	}
}

func GetTokenFromHeader(r *http.Request) string {
	token := r.Header.Get(authHeader)
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:]
	}

	return ""
}
