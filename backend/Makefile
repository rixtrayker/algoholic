.PHONY: help install build run dev test clean lint format db-setup db-migrate db-seed db-reset

# Colors
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

# Variables
BINARY_NAME=algoholic
MAIN_PATH=./main.go
BUILD_DIR=./bin
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Database
DB_USER=leetcode
DB_NAME=leetcode_training
DB_HOST=localhost
DB_PORT=5432

help: ## Show this help message
	@echo "$(BLUE)Backend Makefile - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

install: ## Install Go dependencies
	@echo "$(BLUE)Installing Go dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

build: ## Build the backend binary
	@echo "$(BLUE)Building backend...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

run: build ## Build and run the backend
	@echo "$(BLUE)Starting backend server...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Run backend in development mode (with auto-reload via air if available)
	@echo "$(BLUE)Starting backend in development mode...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)Air not installed. Running with 'go run'...$(NC)"; \
		echo "$(YELLOW)Install air for hot reload: go install github.com/cosmtrek/air@latest$(NC)"; \
		go run $(MAIN_PATH); \
	fi

test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✓ Tests complete$(NC)"

test-coverage: test ## Run tests with coverage report
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test -v -tags=integration ./tests/...

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f *.log
	@echo "$(GREEN)✓ Clean complete$(NC)"

lint: ## Run golangci-lint
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting complete$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not installed$(NC)"; \
		echo "Install: brew install golangci-lint"; \
		echo "Running go vet instead..."; \
		go vet ./...; \
	fi

format: ## Format Go code with gofmt and goimports
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w $(GO_FILES); \
	else \
		echo "$(YELLOW)goimports not installed (optional)$(NC)"; \
		echo "Install: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi
	@echo "$(GREEN)✓ Formatting complete$(NC)"

db-setup: ## Create database and user
	@echo "$(BLUE)Setting up database...$(NC)"
	@psql -U postgres -c "CREATE USER $(DB_USER) WITH PASSWORD 'leetcode_password';" 2>/dev/null || echo "$(YELLOW)User already exists$(NC)"
	@psql -U postgres -c "CREATE DATABASE $(DB_NAME) OWNER $(DB_USER);" 2>/dev/null || echo "$(YELLOW)Database already exists$(NC)"
	@psql -U postgres -d $(DB_NAME) -c "GRANT ALL PRIVILEGES ON DATABASE $(DB_NAME) TO $(DB_USER);" 2>/dev/null || true
	@echo "$(GREEN)✓ Database setup complete$(NC)"

db-migrate: ## Run GORM auto-migrations
	@echo "$(BLUE)Running database migrations...$(NC)"
	@echo "$(YELLOW)Migrations run automatically on backend startup via GORM$(NC)"
	@go run $(MAIN_PATH) migrate

db-seed: ## Seed database with sample data
	@echo "$(BLUE)Seeding database...$(NC)"
	@go run $(MAIN_PATH) seed

db-reset: ## Drop and recreate database
	@echo "$(RED)⚠️  This will delete all data!$(NC)"
	@echo "$(YELLOW)Dropping database...$(NC)"
	@psql -U postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "$(YELLOW)Recreating database...$(NC)"
	@psql -U postgres -c "CREATE DATABASE $(DB_NAME) OWNER $(DB_USER);"
	@echo "$(GREEN)✓ Database reset complete$(NC)"
	@echo "$(YELLOW)Run 'make run' to apply migrations$(NC)"

db-console: ## Open PostgreSQL console
	@psql -U $(DB_USER) -d $(DB_NAME)

check: lint test ## Run linter and tests

watch: ## Watch for changes and run tests (requires entr)
	@echo "$(BLUE)Watching for changes...$(NC)"
	@if command -v entr > /dev/null; then \
		find . -name '*.go' | entr -c go test ./...; \
	else \
		echo "$(RED)entr not installed$(NC)"; \
		echo "Install: brew install entr"; \
	fi

deps-upgrade: ## Upgrade all Go dependencies
	@echo "$(BLUE)Upgrading dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies upgraded$(NC)"

install-tools: ## Install development tools (air, golangci-lint, etc.)
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(YELLOW)Install golangci-lint manually: brew install golangci-lint$(NC)"
	@echo "$(GREEN)✓ Tools installed$(NC)"

.DEFAULT_GOAL := help
