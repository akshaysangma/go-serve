package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	RequestLoggerKey contextKey = "requestLogger"
)

type wrapperResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (ww *wrapperResponseWriter) WriteHeader(statusCode int) {
	ww.ResponseWriter.WriteHeader(statusCode)
	ww.statusCode = statusCode
}

func RequestLoggerMiddleware(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &wrapperResponseWriter{ResponseWriter: w}

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			requestLogger := logger.With(zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_ip", r.RemoteAddr),
			)

			ctx := context.WithValue(r.Context(), RequestLoggerKey, requestLogger)
			next.ServeHTTP(ww, r.WithContext(ctx))

			requestLogger.Info("Request Completed",
				zap.Int("status_code", ww.statusCode),
				zap.Duration("duration", time.Since(start)))

		})
	}
}

func LoggerFromContext(ctx context.Context, defaultLogger *zap.Logger) *zap.Logger {
	if logger, ok := ctx.Value(RequestLoggerKey).(*zap.Logger); ok {
		return logger
	}

	defaultLogger.Warn("Contextual logger not found, using default logger.")
	return defaultLogger
}
