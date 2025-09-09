# MCP Google Docs Editor - Development Commands

.PHONY: dev dev-build dev-up dev-down dev-logs test clean help

# Start development environment
dev:
	@echo "Starting development environment..."
	docker-compose up --build

# Build development images
dev-build:
	@echo "Building development images..."
	docker-compose build

# Start development services (detached)
dev-up:
	@echo "Starting services in background..."
	docker-compose up -d --build

# Stop development services
dev-down:
	@echo "Stopping development services..."
	docker-compose down

# View development logs
dev-logs:
	@echo "Showing logs..."
	docker-compose logs -f

# Run tests
test:
	@echo "Running backend tests..."
	cd services/backend && go test ./...
	@echo "Running frontend tests..."
	cd services/frontend && npm run test

# Clean development environment
clean:
	@echo "Cleaning development environment..."
	docker-compose down -v --rmi local
	docker system prune -f

# Show help
help:
	@echo "Available commands:"
	@echo "  dev       - Start development environment (interactive)"
	@echo "  dev-up    - Start development services (background)"
	@echo "  dev-down  - Stop development services"
	@echo "  dev-build - Build development images"
	@echo "  dev-logs  - View service logs"
	@echo "  test      - Run all tests"
	@echo "  clean     - Clean development environment"
	@echo "  help      - Show this help message"