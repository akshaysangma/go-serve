// Package handlers consist of REST API Handlers
package handlers

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Healthcheck(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
		logger.Info("Health check requested",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
	}
}
