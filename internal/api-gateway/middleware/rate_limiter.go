package middleware

import (
	"net/http"

	"github.com/akshaysangma/go-serve/internal/common/config"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware is a leaking bucket type Ratelimit
func RateLimitMiddleware(config config.RateLimitConfig, defaultLogger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		rl := rate.NewLimiter(rate.Every(config.LimitInterval), config.Burst)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := LoggerFromContext(r.Context(), defaultLogger)
			if rl.Allow() {
				logger.Info("Rate limited API called", zap.Float64("tokens_available", rl.Tokens()))
				next.ServeHTTP(w, r)
				return
			}
			logger.Error("Rate limit reached")
			w.Header().Add("Retry-After", config.LimitInterval.String())
			http.Error(w, "Rate limit reached. Please try after sometime.", http.StatusTooManyRequests)
		})
	}
}
