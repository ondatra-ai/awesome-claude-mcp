# MCP Google Docs Editor - Development Makefile
.PHONY: help init dev test-unit test-e2e lint-backend lint-frontend lint-scripts lint-docs lint

SUPPORTED_E2E_ENVS := local dev
E2E_ENV ?= local
SKIP_DEV_TARGET ?= 0

E2E_EXTRA_GOALS := $(filter-out test-e2e,$(MAKECMDGOALS))
E2E_CMD_ENV := $(firstword $(filter $(SUPPORTED_E2E_ENVS),$(E2E_EXTRA_GOALS)))

ifneq ($(E2E_CMD_ENV),)
  ifeq ($(firstword $(MAKECMDGOALS)),test-e2e)
    SKIP_DEV_TARGET := 1
  endif
  ifneq ($(origin E2E_ENV),command line)
    E2E_ENV := $(E2E_CMD_ENV)
  endif
endif

# Default target
help: ## Show available commands
	@echo "MCP Google Docs Editor - Development Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

init: ## Install dependencies and build Docker images with caching
	@echo "🚀 Initializing project dependencies and Docker builds..."
	@echo "📦 Installing backend Go dependencies..."
	go mod download -C services/backend
	go mod tidy -C services/backend
	@echo "📦 Installing frontend Node.js dependencies..."
	npm install --prefix services/frontend
	@echo "📦 Installing test dependencies..."
	npm install --prefix tests
	@echo "🐳 Building all Docker images with cached dependencies..."
	docker compose -f docker-compose.test.yaml build --parallel
	@echo "✅ All dependencies and Docker images ready with caching optimized!"

dev: ## Start all services with Docker Compose
	@if [ "$(SKIP_DEV_TARGET)" = "1" ]; then \
		echo "[warn] Skipping dev target (interpreted as E2E environment flag)."; \
	else \
		docker compose up --build; \
	fi

test-unit: ## Run unit tests for both services (via Docker)
	@echo "🧪 Running unit tests..."
	@echo "🔧 Running Go backend tests (via Docker)..."
	docker build -t mcp-backend-test --target test -f services/backend/Dockerfile services/backend
	docker run --rm mcp-backend-test go test ./...
	@echo "🔧 Running Node.js frontend tests (via Docker)..."
	docker build -t mcp-frontend-test --target test -f services/frontend/Dockerfile services/frontend
	docker run --rm mcp-frontend-test npm test
	@echo "✅ Unit tests completed!"

test-e2e: ## Run E2E tests (default local; append environment name e.g. `make test-e2e dev`)
	@E2E_ENV=$(E2E_ENV); \
	printf "🚀 Starting E2E Test Pipeline for '%s'...\n" "$$E2E_ENV"; \
	docker compose -f docker-compose.test.yaml down --remove-orphans >/dev/null 2>&1 || true; \
	if [ "$$E2E_ENV" = "local" ]; then \
		echo "🔧 Starting backend and frontend services..."; \
		echo "⏳ Waiting for services to be healthy (docker compose --wait)..."; \
		docker compose -f docker-compose.test.yaml up -d --wait backend frontend; \
	else \
		echo "🌐 Using remote endpoints; skipping local service startup."; \
	fi; \
	echo "🧪 Running E2E tests..."; \
	docker compose -f docker-compose.test.yaml run --build --no-deps --rm \
		-e E2E_ENV=$$E2E_ENV \
		playwright-test; \
	TEST_EXIT_CODE=$$?; \
	echo "🧹 Cleaning up containers..."; \
	docker compose -f docker-compose.test.yaml down --remove-orphans >/dev/null 2>&1 || true; \
	if [ $$TEST_EXIT_CODE -eq 0 ]; then \
		echo "✅ All tests passed!"; \
	else \
		echo "❌ Tests exited with $$TEST_EXIT_CODE (frontend failures may be expected for remote envs)."; \
	fi; \
	exit $$TEST_EXIT_CODE

lint-backend: ## Run Go linter on backend code (auto-fix when possible)
	@echo "🔧 Running go fmt to fix formatting on backend..."
	find services/backend/ -name "*.go" -exec gofmt -l -w {} \;
	@echo "🔍 Running golangci-lint on backend..."
	cd services/backend && golangci-lint run --fix ./...
	@echo "✅ Backend linting completed!"

lint-frontend: ## Run ESLint and Prettier on frontend code (via Docker)
	@echo "🔍 Running Next.js ESLint on frontend (via Docker)..."
	docker build -t mcp-frontend-test --target test -f services/frontend/Dockerfile services/frontend
	docker run --rm mcp-frontend-test npm run lint
	@echo "✅ Frontend linting completed!"

lint-scripts: ## Run Go linter on scripts with Go code (auto-fix when possible)
	@echo "🔧 Running go fmt to fix formatting on all Go scripts..."
	find scripts/ -name "*.go" -exec gofmt -l -w {} \;
	@echo "🔍 Running golangci-lint on bmad-cli..."
	cd scripts/bmad-cli && golangci-lint run --fix ./...
	@echo "✅ Scripts linting completed!"

lint-docs: ## Validate requirements.yaml and epic YAML files against Yamale schemas
	@echo "🔍 Validating requirements.yaml against schema (strict mode)..."
	yamale -s docs/requirements-schema.yaml docs/requirements.yaml
	@echo "🔍 Validating epic YAML files against schema (strict mode)..."
	yamale -s docs/epics/jsons/epics-schema.yaml docs/epics/jsons/epic-*.yaml
	@echo "✅ Documentation validation completed!"

lint: lint-backend lint-frontend lint-scripts lint-docs ## Run all linting checks
