.PHONY: help install build run dev test clean lint format preview

# Colors
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

# Variables
BUILD_DIR=./dist
NODE_MODULES=./node_modules

help: ## Show this help message
	@echo "$(BLUE)Frontend Makefile - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

install: ## Install npm dependencies
	@echo "$(BLUE)Installing npm dependencies...$(NC)"
	@npm install
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

build: ## Build for production
	@echo "$(BLUE)Building frontend for production...$(NC)"
	@npm run build
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)$(NC)"

run: dev ## Alias for 'make dev'

dev: ## Run development server with hot reload
	@echo "$(BLUE)Starting frontend development server...$(NC)"
	@echo "$(YELLOW)URL: http://localhost:5173$(NC)"
	@npm run dev

preview: build ## Preview production build locally
	@echo "$(BLUE)Starting preview server...$(NC)"
	@npm run preview

test: ## Run tests with Vitest
	@echo "$(BLUE)Running tests...$(NC)"
	@npm test
	@echo "$(GREEN)✓ Tests complete$(NC)"

test-watch: ## Run tests in watch mode
	@echo "$(BLUE)Running tests in watch mode...$(NC)"
	@npm run test:watch

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@npm run test:coverage
	@echo "$(GREEN)✓ Coverage report generated$(NC)"

test-ui: ## Run tests with Vitest UI
	@echo "$(BLUE)Opening Vitest UI...$(NC)"
	@npm run test:ui

clean: ## Clean build artifacts and dependencies
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(NODE_MODULES)
	@rm -rf coverage
	@rm -f *.log
	@echo "$(GREEN)✓ Clean complete$(NC)"

clean-build: ## Clean build artifacts only (keep node_modules)
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf coverage
	@echo "$(GREEN)✓ Build artifacts cleaned$(NC)"

lint: ## Run ESLint
	@echo "$(BLUE)Running ESLint...$(NC)"
	@npm run lint
	@echo "$(GREEN)✓ Linting complete$(NC)"

lint-fix: ## Run ESLint with auto-fix
	@echo "$(BLUE)Running ESLint with auto-fix...$(NC)"
	@npm run lint -- --fix
	@echo "$(GREEN)✓ Linting and fixes complete$(NC)"

format: ## Format code with Prettier
	@echo "$(BLUE)Formatting code...$(NC)"
	@if command -v prettier > /dev/null; then \
		npx prettier --write "src/**/*.{js,jsx,ts,tsx,json,css,md}"; \
		echo "$(GREEN)✓ Formatting complete$(NC)"; \
	else \
		echo "$(YELLOW)Prettier not installed$(NC)"; \
		echo "Install: npm install -D prettier"; \
	fi

type-check: ## Run TypeScript type checking
	@echo "$(BLUE)Running TypeScript type check...$(NC)"
	@npx tsc --noEmit
	@echo "$(GREEN)✓ Type check complete$(NC)"

check: lint type-check test ## Run linter, type check, and tests

deps-upgrade: ## Upgrade npm dependencies
	@echo "$(BLUE)Upgrading dependencies...$(NC)"
	@npm update
	@echo "$(GREEN)✓ Dependencies upgraded$(NC)"
	@echo "$(YELLOW)Run 'npm outdated' to see remaining outdated packages$(NC)"

deps-outdated: ## Check for outdated dependencies
	@echo "$(BLUE)Checking for outdated dependencies...$(NC)"
	@npm outdated

deps-audit: ## Run security audit
	@echo "$(BLUE)Running security audit...$(NC)"
	@npm audit
	@echo "$(YELLOW)Run 'npm audit fix' to auto-fix vulnerabilities$(NC)"

deps-audit-fix: ## Fix security vulnerabilities
	@echo "$(BLUE)Fixing security vulnerabilities...$(NC)"
	@npm audit fix

analyze: build ## Analyze bundle size
	@echo "$(BLUE)Analyzing bundle size...$(NC)"
	@echo "$(YELLOW)Bundle analysis requires 'rollup-plugin-visualizer'$(NC)"
	@du -sh $(BUILD_DIR)/* | sort -h

serve: preview ## Alias for 'make preview'

watch: dev ## Alias for 'make dev'

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@npm install -D prettier @typescript-eslint/parser @typescript-eslint/eslint-plugin
	@echo "$(GREEN)✓ Tools installed$(NC)"

.DEFAULT_GOAL := help
