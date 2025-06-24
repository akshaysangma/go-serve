package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/akshaysangma/go-serve/internals/common/config"
	"github.com/akshaysangma/go-serve/internals/common/logging"
	database "github.com/akshaysangma/go-serve/internals/database/postgres"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() {
	config := config.LoadConfig()
	logger, err := logging.InitLogger(config.Log.Level, config.Log.Encoding)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	// flush all buffer before exiting
	defer logger.Sync()

	logger.Info("Configuration loaded successfully", zap.Int("port", config.App.Port), zap.String("log_Level", config.Log.Level))

	// creating Db Connection Pool
	dB, err := database.ConnectDB(config.Database.URL, config.Database.MaxConnections)
	if err != nil {
		logger.Fatal("Unable to connect to Database", zap.Error(err))
	}
	defer dB.Close()

	apiServer := &http.Server{
		Addr:    ":" + strconv.Itoa(config.App.Port),
		Handler: initRoutes(logger),
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
