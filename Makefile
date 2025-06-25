# --- Project Configuration ---
# Define variables for commonly used paths/commands
PROJECT_ROOT := $(shell pwd)
MIGRATIONS_DIR := internal/database/postgres/migrations
SQLC_CONFIG := sqlc.yaml
API_GATEWAY_CMD := ./cmd/api-gateway
API_GATEWAY_BINARY := $(API_GATEWAY_CMD)/api-gateway-app # Binary name
DOCKER_COMPOSE_FILE := deployments/docker-compose.yaml

# --- Environment Variables ---
# Default environment variables for local execution if not already set.
# These will be read by viper in your Go applications.
# For Docker Compose services, env vars are set directly in docker-compose.yaml.
export APP_PORT ?= 8080
export LOG_LEVEL ?= debug
export LOG_ENCODING ?= console
# Ensure DATABASE_URL is set correctly based on your Docker Compose setup for local runs.
# When running inside Docker Compose, `postgres` resolves to the container.
# When running Go app directly on host, `localhost` (if DB also on host).
# When running Go app in Docker but DB on host, `host.docker.internal` (Docker Desktop) or host IP (Linux).
# For this Makefile, we'll assume the DB is running via `docker-compose up` on localhost for migration/testing
# and that your Go app will be run directly on the host using this DATABASE_URL.
export DATABASE_URL ?= postgresql://user:password@localhost:5432/goserve_db?sslmode=disable
# Environment variable for viper prefix (must match internal/common/config/config.go)
export GOBSERVE_APP_PORT := $(APP_PORT)
export GOBSERVE_LOG_LEVEL := $(LOG_LEVEL)
export GOBSERVE_LOG_ENCODING := $(LOG_ENCODING)
export GOBSERVE_DATABASE_URL := $(DATABASE_URL)


# --- Phony Targets ---
# .PHONY indicates that these are not actual files.
.PHONY: all up down migrate-up migrate-down generate build run test clean docker-build-api-gateway

# --- Default Target ---
all: build test

# --- Infrastructure Management ---
up: ## Start all infrastructure services (PostgreSQL, Redis, Kafka, etc.)
	@echo "Starting infrastructure services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

down: ## Stop and remove all infrastructure services
	@echo "Stopping and removing infrastructure services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v # -v removes volumes too (clean slate)

# --- Database Migrations ---
migrate-up: ## Apply all pending database migrations
	@echo "Applying database migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down: ## Rollback the last applied database migration
	@echo "Rolling back last database migration..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

# --- Code Generation ---
generate: ## Run sqlc to generate Go code from SQL queries
	@echo "Running sqlc generate..."
	sqlc generate
	@echo "sqlc generate complete."

# --- Go Application Management ---
build: ## Build all Go service binaries (e.g., api-gateway)
	@echo "Building Go binaries..."
	go build -o $(API_GATEWAY_BINARY) $(API_GATEWAY_CMD)/main.go # Builds api-gateway binary

run: build ## Run the api-gateway Go application directly on the host
	@echo "Running API Gateway..."
	$(API_GATEWAY_BINARY)

docker-build-api-gateway: ## Build the Docker image for the api-gateway Go application
	@echo "Building Docker image for API Gateway..."
	docker build -t realtime-platform/api-gateway -f $(API_GATEWAY_CMD)/Dockerfile .

# --- Testing ---
test: ## Run all Go tests with race detection and coverage
	@echo "Running all tests with race detector..."
	go test -race ./...

test-coverage: ## Generate and view HTML test coverage report
	@echo "Generating test coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	@echo "Coverage report generated: coverage.out"

# --- Cleanup ---
clean: ## Clean up build artifacts and caches
	@echo "Cleaning up..."
	go clean ./...
	rm -f $(API_GATEWAY_BINARY)
	rm -rf internal/database/postgres/sqlc # Remove generated sqlc code
	rm -f coverage.out # Remove coverage report
	@echo "Cleanup complete."

# --- Help Target (for documentation) ---
help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
