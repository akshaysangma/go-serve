package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/akshaysangma/go-serve/internals/api-gateway/handlers"
	"github.com/akshaysangma/go-serve/internals/common/config"
	"github.com/akshaysangma/go-serve/internals/common/logging"
	"go.uber.org/zap"
)

func main() {
	config := config.LoadConfig()
	logger, err := logging.InitLogger(config.Log.Level, config.Log.Encoding)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	// flush all buffer before exiting
	defer logging.Sync(logger)

	logger.Info("Configuration loaded successfully", zap.Int("port", config.App.Port), zap.String("log_Level", config.Log.Level))

	router := http.NewServeMux()
	router.Handle("GET /healthcheck", handlers.Healthcheck(logger))

	v1 := http.NewServeMux()

	router.Handle("/v1/", http.StripPrefix("/v1", v1))

	apiServer := &http.Server{
		Addr:    ":" + strconv.Itoa(config.App.Port),
		Handler: router,
	}

	go func() {
		logger.Info("Starting API Server", zap.Int("port", config.App.Port))
		if err := apiServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("Failed to start API Server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("Shutting down Server...",
		zap.String("signal", sig.String()),
		zap.Duration("graceful_shutdown_period", config.App.GracefulShutdownPeriod))

	ctx, cancel := context.WithTimeout(context.Background(), config.App.GracefulShutdownPeriod)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully.")
}
