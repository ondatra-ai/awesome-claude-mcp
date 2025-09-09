# MCP Google Docs Editor - Development Makefile
.PHONY: help dev test-unit test-e2e lint-backend lint-frontend

# Default target
help: ## Show available commands
	@echo "MCP Google Docs Editor - Development Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

dev: ## Start all services with Docker Compose
	docker-compose up --build

test-unit: ## Run unit tests for both services
	@echo "🧪 Running unit tests..."
	go test ./services/backend/...
	npm test --prefix services/frontend
	@echo "✅ Unit tests completed!"

test-e2e: ## Run E2E tests with Docker
	@echo "🚀 Starting E2E Test Pipeline..."
	@echo "🧹 Cleaning up existing containers..."
	@docker-compose -f docker-compose.test.yml down --remove-orphans || true
	@echo "🔧 Starting backend and frontend services..."
	@docker-compose -f docker-compose.test.yml up -d backend frontend
	@echo "⏳ Waiting for services to be healthy..."
	@docker-compose -f docker-compose.test.yml exec backend wget --no-verbose --tries=1 --spider http://localhost:8080/health
	@docker-compose -f docker-compose.test.yml exec frontend wget --no-verbose --tries=1 --spider http://localhost:3000
	@echo "🧪 Running E2E tests..."
	@docker-compose -f docker-compose.test.yml run --rm playwright-test; \
	TEST_EXIT_CODE=$$?; \
	echo "🧹 Cleaning up containers..."; \
	docker-compose -f docker-compose.test.yml down --remove-orphans; \
	if [ $$TEST_EXIT_CODE -eq 0 ]; then \
		echo "✅ All tests passed!"; \
	else \
		echo "❌ Tests failed with exit code: $$TEST_EXIT_CODE"; \
	fi; \
	exit $$TEST_EXIT_CODE

lint-backend: ## Run Go linter on backend code (auto-fix when possible)
	@echo "🔍 Running Go lint on backend..."
	GOWORK=off go run -C services/backend golang.org/x/lint/golint@latest ./cmd
	@echo "🔧 Running go fmt to fix formatting..."
	gofmt -l -w services/backend/cmd
	@echo "✅ Backend linting completed!"

lint-frontend: ## Run ESLint and Prettier on frontend code (auto-fix when possible)
	@echo "🔍 Running Next.js ESLint with --fix on frontend..."
	npm run lint --prefix services/frontend -- --fix
	@echo "🎨 Running Prettier with --write on frontend..."
	npx prettier --write services/frontend/ --ignore-path services/frontend/.prettierignore --config services/frontend/.prettierrc.json
	@echo "✅ Frontend linting completed!"
