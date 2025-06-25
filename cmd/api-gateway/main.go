package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/akshaysangma/go-serve/internal/api-gateway/handlers"
	"github.com/akshaysangma/go-serve/internal/api-gateway/middleware"
	"github.com/akshaysangma/go-serve/internal/api-gateway/repositories"
	"github.com/akshaysangma/go-serve/internal/api-gateway/services"
	"github.com/akshaysangma/go-serve/internal/common/config"
	"github.com/akshaysangma/go-serve/internal/common/logging"
	database "github.com/akshaysangma/go-serve/internal/database/postgres"
	db "github.com/akshaysangma/go-serve/internal/database/postgres/sqlc"
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
	logger.Info("Successfully connected to Database")

	dBQueries := db.New(dB)

	router := http.NewServeMux()
	router.Handle("GET /health", handlers.Healthcheck(logger))

	// V1 API Group
	v1 := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", v1))

	// User V1
	userRepo := repositories.NewUserRepository(dBQueries)
	articleRepo := repositories.NewArticleRepository(dBQueries)
	userService := services.NewUserService(userRepo, articleRepo, dB, logger)
	userMiddlewareChain := middleware.ChainMiddleware(middleware.RequestLoggerMiddleware(logger))
	v1.Handle("GET /users/{id}", userMiddlewareChain(handlers.GetUserByIDHandler(userService, logger)))
	v1.Handle("POST /users", userMiddlewareChain(handlers.CreateUserHandler(userService, logger)))
	v1.Handle("GET /users", userMiddlewareChain(handlers.ListUsersHandler(userService, logger)))

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
