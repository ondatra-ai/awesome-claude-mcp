# MCP Google Docs Editor - Development Makefile
.PHONY: help init dev test-unit test-e2e lint-backend lint-frontend lint-scripts lint-terraform lint-terraform-modules tf-bootstrap tf-init tf-validate tf-plan tf-apply tf-plan-destroy tf-destroy

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

## Terraform helpers (use TF_ENV=dev|staging|prod, or ENV=... as an alias)
# Prefer TF_ENV; fall back to ENV if provided
TF_ENV ?= $(or $(ENV),dev)
TF_DIR=infrastructure/terraform
TF_BOOTSTRAP_DIR=infrastructure/terraform/bootstrap

tf-bootstrap: ## Bootstrap S3 backend (run once to create bucket and DynamoDB table)
	@echo "🔧 Bootstrapping S3 backend infrastructure..."
	terraform -chdir=$(TF_BOOTSTRAP_DIR) init
	terraform -chdir=$(TF_BOOTSTRAP_DIR) plan -out=bootstrap.tfplan
	terraform -chdir=$(TF_BOOTSTRAP_DIR) apply -auto-approve bootstrap.tfplan
	@echo "✅ S3 backend bootstrap completed!"

tf-init: ## Terraform init for ENV (ENV=dev|staging|prod)
	@echo "🔧 Terraform init for $(TF_ENV)..."
	terraform -chdir=$(TF_DIR) init -backend-config=backend-$(TF_ENV).hcl

tf-validate: ## Terraform validate for ENV (ENV=dev|staging|prod)
	@echo "🧪 Terraform validate for $(TF_ENV)..."
	terraform -chdir=$(TF_DIR) validate

TF_PLAN ?= plan.out
tf-plan: ## Terraform plan for ENV (ENV=dev|staging|prod)
	@echo "🗺️  Terraform plan for $(TF_ENV)..."
	terraform -chdir=$(TF_DIR) plan -var-file="environments/$(TF_ENV).tfvars" -out $(TF_PLAN)

tf-apply: ## Terraform apply for ENV (ENV=dev|staging|prod)
	@echo "🚀 Terraform apply for $(TF_ENV)..."
	@if [ ! -f "$(TF_DIR)/$(TF_PLAN)" ]; then \
	  echo "❌ Plan file '$(TF_PLAN)' not found in $(TF_DIR). Run 'make tf-plan TF_ENV=$(TF_ENV)' first."; \
	  exit 1; \
	fi
	terraform -chdir=$(TF_DIR) apply -auto-approve -input=false $(TF_PLAN)

TF_DESTROY_PLAN ?= destroy.out
tf-plan-destroy: ## Terraform plan destroy for ENV (ENV=dev|staging|prod)
	@echo "🗑️  Terraform plan destroy for $(TF_ENV)..."
	terraform -chdir=$(TF_DIR) plan -destroy -var-file="environments/$(TF_ENV).tfvars" -out $(TF_DESTROY_PLAN)

tf-destroy: ## Terraform destroy for ENV (ENV=dev|staging|prod)
	@echo "🗑️  Terraform destroy for $(TF_ENV)..."
	terraform -chdir=$(TF_DIR) destroy -var-file="environments/$(TF_ENV).tfvars" -auto-approve

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

lint-terraform: ## Run tflint on Terraform code (auto-fix when possible)
	@echo "🔍 Running tflint on Terraform infrastructure..."
	@echo "📦 Installing tflint plugins..."
	tflint --init
	@echo "🔧 Running tflint with auto-fix on infrastructure/terraform..."
	tflint --fix --chdir=infrastructure/terraform
	@echo "🔧 Running terraform fmt on infrastructure/terraform..."
	terraform fmt -recursive infrastructure/terraform/
	@echo "✅ Terraform linting completed!"

lint-terraform-modules: ## Run tflint on all Terraform modules
	@echo "🔍 Running tflint on Terraform modules..."
	@echo "📦 Installing tflint plugins..."
	tflint --init
	@for module in infrastructure/terraform/modules/*; do \
		if [ -d "$$module" ]; then \
			echo "🔧 Linting module: $$module"; \
			tflint --fix --chdir="$$module"; \
		fi; \
	done
	@echo "🔧 Running terraform fmt on all modules..."
	terraform fmt -recursive infrastructure/terraform/
	@echo "✅ All Terraform modules linting completed!"
