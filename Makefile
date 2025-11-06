.PHONY: help dev build run migrate seed test test-unit test-int docker-build docker-up docker-down clean swagger

# Variables
APP_NAME=greenlight-api
MAIN_PATH=cmd/api/main.go
SEED_PATH=scripts/seed.go
BIN_DIR=bin
DOCKER_COMPOSE=docker-compose

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: ## Run application in development mode
	@echo "Starting application in development mode..."
	@go run $(MAIN_PATH)

build: ## Build the application binary
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Binary created at $(BIN_DIR)/$(APP_NAME)"

run: build ## Build and run the application
	@echo "Running application..."
	@./$(BIN_DIR)/$(APP_NAME)

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run migrations/migrate.go

seed: ## Seed the database with sample data
	@echo "Seeding database..."
	@go run $(SEED_PATH)

test: ## Run all tests
	@echo "Running all tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -race -short ./...

test-int: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v -race -run Integration ./...

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@swag init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated in docs/"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@$(DOCKER_COMPOSE) build

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services are up! API available at http://localhost:8080"
	@echo "Run 'make docker-seed' to seed the database"

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down

docker-logs: ## View Docker container logs
	@$(DOCKER_COMPOSE) logs -f

docker-seed: ## Seed database in Docker
	@echo "Seeding database in Docker..."
	@$(DOCKER_COMPOSE) exec api go run /root/scripts/seed.go || echo "Run seed manually after containers are fully up"

docker-clean: ## Stop containers and remove volumes
	@echo "Cleaning up Docker containers and volumes..."
	@$(DOCKER_COMPOSE) down -v

docker-reset: ## Complete reset - stop, clean, rebuild, and start fresh
	@echo "Performing complete Docker reset..."
	@$(DOCKER_COMPOSE) down -v
	@docker volume prune -f
	@echo "Building fresh containers..."
	@$(DOCKER_COMPOSE) up -d --build
	@echo "Waiting for services to be ready..."
	@sleep 30
	@echo "Services are ready!"
	@echo "Run 'make docker-seed' to add sample data"

clean: ## Clean build artifacts and test files
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@rm -rf uploads/*
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

tidy: ## Tidy go.mod and go.sum
	@echo "Tidying dependencies..."
	@go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

.DEFAULT_GOAL := help

