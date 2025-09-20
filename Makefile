# MCP Google Docs Editor - Development Makefile
.PHONY: help init dev test-unit test-e2e lint-backend lint-frontend lint-scripts railway-login railway-link deploy deploy-dev deploy-staging deploy-prod deploy-service

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
	docker compose -f docker-compose.test.yml build --parallel
	@echo "✅ All dependencies and Docker images ready with caching optimized!"

dev: ## Start all services with Docker Compose
	docker compose up --build

RAILWAY_PROJECT_ID=801ad5e0-95bf-4ce6-977e-6f2fa37529fd
ENV ?= development

railway-login: ## Authenticate Railway CLI
	railway login

railway-link: ## Link repository to Railway project
	railway link --project $(RAILWAY_PROJECT_ID)

deploy: ## Deploy services to Railway environment (ENV=development|staging|production)
	@if [ -z "$(ENV)" ]; then \
		echo "❌ ENV must be set to development, staging, or production"; \
		exit 1; \
	fi
	@if [ "$(ENV)" = "development" ]; then \
		services="frontend-dev backend-dev"; \
	elif [ "$(ENV)" = "staging" ]; then \
		services="frontend-staging backend-staging"; \
	elif [ "$(ENV)" = "production" ]; then \
		services="frontend backend"; \
	else \
		echo "❌ Unknown ENV: $(ENV)"; exit 1; \
	fi; \
	railway environment $(ENV); \
	for svc in $$services; do \
		if echo $$svc | grep -q "frontend"; then \
			path="services/frontend"; \
		else \
			path="services/backend"; \
		fi; \
		echo "🚀 Deploying $$svc from $$path"; \
		railway up --service $$svc --path-as-root $$path; \
	done

deploy-dev: ## Deploy development environment to Railway
	$(MAKE) deploy ENV=development

deploy-staging: ## Deploy staging environment to Railway
	$(MAKE) deploy ENV=staging

deploy-prod: ## Deploy production environment to Railway
	$(MAKE) deploy ENV=production

deploy-service: ## Deploy a single Railway service (SERVICE=frontend|backend|...)
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ SERVICE must be set"; \
		exit 1; \
	fi
	@if [ -z "$(ENV)" ]; then \
		echo "❌ ENV must be set"; \
		exit 1; \
	fi
	railway environment $(ENV)
	@if echo $(SERVICE) | grep -q "frontend"; then \
		path="services/frontend"; \
	else \
		path="services/backend"; \
	fi; \
	echo "🚀 Deploying $(SERVICE) to $(ENV)"; \
	railway up --service $(SERVICE) --path-as-root $$path

test-unit: ## Run unit tests for both services
	@echo "🧪 Running unit tests..."
	@echo "🔧 Running Go backend tests..."
	go test -C services/backend ./...
	@echo "🔧 Running Node.js frontend tests..."
	npm test --prefix services/frontend
	@echo "✅ Unit tests completed!"

test-e2e: ## Run E2E tests with Docker
	@echo "🚀 Starting E2E Test Pipeline..."
	@echo "🧹 Cleaning up existing containers..."
	@docker compose -f docker-compose.test.yml down --remove-orphans || true
	@echo "🔧 Starting backend and frontend services..."
	@docker compose -f docker-compose.test.yml up -d backend frontend
	@echo "⏳ Waiting for services to be healthy..."
	@for i in $$(seq 1 30); do \
		if docker compose -f docker-compose.test.yml exec -T backend wget --no-verbose --tries=1 --spider http://localhost:8080/health > /dev/null 2>&1; then \
			echo "✅ Backend is healthy"; \
			break; \
		fi; \
		echo "Waiting for backend... ($$i/30)"; \
		sleep 2; \
	done
	@for i in $$(seq 1 30); do \
		if docker compose -f docker-compose.test.yml exec -T frontend wget --no-verbose --tries=1 --spider http://0.0.0.0:3000 > /dev/null 2>&1; then \
			echo "✅ Frontend is healthy"; \
			break; \
		fi; \
		echo "Waiting for frontend... ($$i/30)"; \
		sleep 2; \
	done
	@echo "🧪 Running E2E tests..."
	@docker compose -f docker-compose.test.yml run --rm playwright-test; \
	TEST_EXIT_CODE=$$?; \
	echo "🧹 Cleaning up containers..."; \
	docker compose -f docker-compose.test.yml down --remove-orphans; \
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

lint-scripts: ## Run Go linter on scripts/pr-triage code (auto-fix when possible)
	@echo "🔧 Running go fmt to fix formatting..."
	gofmt -l -w scripts/pr-triage
	@echo "🔍 Running golangci-lint on scripts/pr-triage..."
	cd scripts/pr-triage && golangci-lint run --fix .
	@echo "✅ Scripts linting completed!"
