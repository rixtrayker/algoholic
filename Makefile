.PHONY: help install build test run clean dev stop backend frontend db-setup api-test all

# Colors for terminal output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Algoholic - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Quick Start:$(NC)"
	@echo "  make install    # Install all dependencies"
	@echo "  make dev        # Run backend + frontend in dev mode"
	@echo "  make test       # Run all tests"
	@echo ""

install: ## Install all dependencies (backend + frontend)
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@$(MAKE) -C backend install
	@$(MAKE) -C frontend install
	@echo "$(GREEN)✓ All dependencies installed$(NC)"

build: ## Build backend + frontend for production
	@echo "$(BLUE)Building all services...$(NC)"
	@$(MAKE) -C backend build
	@$(MAKE) -C frontend build
	@echo "$(GREEN)✓ Build complete$(NC)"

test: ## Run all tests (backend + frontend + API)
	@echo "$(BLUE)Running all tests...$(NC)"
	@$(MAKE) -C backend test
	@$(MAKE) -C frontend test
	@$(MAKE) api-test
	@echo "$(GREEN)✓ All tests passed$(NC)"

dev: ## Run backend + frontend in development mode (parallel)
	@echo "$(BLUE)Starting development servers...$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:4000$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:5173$(NC)"
	@trap 'kill 0' EXIT; \
		$(MAKE) -C backend dev & \
		$(MAKE) -C frontend dev & \
		wait

run: dev ## Alias for 'make dev'

backend: ## Run backend only
	@echo "$(BLUE)Starting backend server...$(NC)"
	@$(MAKE) -C backend run

frontend: ## Run frontend only
	@echo "$(BLUE)Starting frontend server...$(NC)"
	@$(MAKE) -C frontend run

backend-dev: ## Run backend in development mode
	@$(MAKE) -C backend dev

frontend-dev: ## Run frontend in development mode
	@$(MAKE) -C frontend dev

backend-build: ## Build backend only
	@$(MAKE) -C backend build

frontend-build: ## Build frontend only
	@$(MAKE) -C frontend build

backend-test: ## Run backend tests only
	@$(MAKE) -C backend test

frontend-test: ## Run frontend tests only
	@$(MAKE) -C frontend test

api-test: ## Run API tests via Postman/Newman
	@echo "$(BLUE)Running API tests...$(NC)"
	@cd postman && ./run-tests.sh

clean: ## Clean build artifacts and dependencies
	@echo "$(BLUE)Cleaning...$(NC)"
	@$(MAKE) -C backend clean
	@$(MAKE) -C frontend clean
	@echo "$(GREEN)✓ Clean complete$(NC)"

db-setup: ## Setup PostgreSQL database
	@echo "$(BLUE)Setting up database...$(NC)"
	@$(MAKE) -C backend db-setup
	@echo "$(GREEN)✓ Database setup complete$(NC)"

db-migrate: ## Run database migrations
	@$(MAKE) -C backend db-migrate

db-seed: ## Seed database with sample data
	@$(MAKE) -C backend db-seed

db-reset: ## Reset database (drop + recreate + migrate)
	@$(MAKE) -C backend db-reset

lint: ## Run linters (backend + frontend)
	@echo "$(BLUE)Running linters...$(NC)"
	@$(MAKE) -C backend lint
	@$(MAKE) -C frontend lint
	@echo "$(GREEN)✓ Linting complete$(NC)"

format: ## Format code (backend + frontend)
	@echo "$(BLUE)Formatting code...$(NC)"
	@$(MAKE) -C backend format
	@$(MAKE) -C frontend format
	@echo "$(GREEN)✓ Formatting complete$(NC)"

check: lint test ## Run linters + tests

logs: ## Show running service logs
	@echo "$(BLUE)Checking running services...$(NC)"
	@ps aux | grep -E "(go run main.go|npm run dev)" | grep -v grep || echo "$(YELLOW)No services running$(NC)"

stop: ## Stop all running services
	@echo "$(BLUE)Stopping services...$(NC)"
	@pkill -f "go run main.go" || true
	@pkill -f "npm run dev" || true
	@pkill -f "vite" || true
	@echo "$(GREEN)✓ Services stopped$(NC)"

status: ## Show service status
	@echo "$(BLUE)Service Status:$(NC)"
	@echo ""
	@echo "$(YELLOW)Backend (port 4000):$(NC)"
	@curl -s http://localhost:4000/health | jq '.' 2>/dev/null || echo "$(RED)✗ Not running$(NC)"
	@echo ""
	@echo "$(YELLOW)Frontend (port 5173):$(NC)"
	@curl -s -o /dev/null -w "%{http_code}" http://localhost:5173 2>/dev/null | grep -q "200" && echo "$(GREEN)✓ Running$(NC)" || echo "$(RED)✗ Not running$(NC)"
	@echo ""
	@echo "$(YELLOW)Database (PostgreSQL):$(NC)"
	@psql -U leetcode -d leetcode_training -c "SELECT version();" 2>/dev/null | head -3 | tail -1 || echo "$(RED)✗ Not connected$(NC)"

health: status ## Alias for 'make status'

all: clean install build test ## Clean, install, build, and test everything

docker-up: ## Start all services with Docker Compose (when available)
	@echo "$(YELLOW)Docker Compose not yet configured$(NC)"
	@echo "Services are running natively. See RUNNING.md"

docker-down: ## Stop all Docker services
	@echo "$(YELLOW)Docker Compose not yet configured$(NC)"

.DEFAULT_GOAL := help
