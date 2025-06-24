package main

import (
	"net/http"

	"github.com/akshaysangma/go-serve/internals/api-gateway/handlers"
	"go.uber.org/zap"
)

func initRoutes(logger *zap.Logger) *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("GET /health", handlers.Healthcheck(logger))

	// V1 API Group
	v1 := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", v1))

	return router
}
