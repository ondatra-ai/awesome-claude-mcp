# MCP Google Docs Editor

A Model Context Protocol integration for seamless Google Docs editing with Claude Code and ChatGPT.

## 🚀 Overview

This project provides a full-stack foundation for MCP (Model Context Protocol) Google Docs integration featuring:

- **Go Backend Service**: High-performance REST API with structured logging and health monitoring
- **Next.js Frontend**: Modern React application with TypeScript and Tailwind CSS
- **Docker Support**: Complete containerization with multi-stage builds
- **E2E Testing**: Comprehensive Playwright test suite
- **Production Ready**: Structured logging, health checks, and graceful shutdown

## 📋 Prerequisites

- **Go**: 1.21.5 or higher
- **Node.js**: 18.x or higher
- **Docker**: Latest version (for containerized deployment)
- **Git**: For version control

## 🏗️ Architecture

### Monorepo Structure
```
awesome-claude-mcp/
├── services/
│   ├── backend/           # Go backend service
│   │   ├── cmd/main.go    # Application entry point
│   │   ├── go.mod         # Go dependencies
│   │   └── Dockerfile     # Multi-stage Docker build
│   └── frontend/          # Next.js frontend application
│       ├── app/           # App Router pages
│       ├── components/    # Reusable UI components
│       ├── lib/           # Utilities and API client
│       └── Dockerfile     # Production Docker build
├── tests/
│   └── e2e/              # Playwright E2E tests
├── docker-compose.yml    # Local development stack
├── playwright.config.ts  # E2E testing configuration
└── README.md
```

### Services Overview

#### Backend Service (Go + Fiber)
- **Framework**: Fiber 2.52.0 for high-performance HTTP server
- **Logging**: Structured JSON logging with zerolog
- **Health Monitoring**: Built-in health check endpoints
- **Production Features**: Graceful shutdown, CORS support, security headers

**API Endpoints:**
- `GET /version` - Returns application version
- `GET /health` - Service health status

#### Frontend Service (Next.js + TypeScript)
- **Framework**: Next.js 14.1.0 with App Router
- **Styling**: Tailwind CSS 3.4.0 with shadcn/ui components
- **Type Safety**: TypeScript 5.3.3 with strict configuration
- **API Integration**: Type-safe client for backend communication

## 🚀 Quick Start

### Docker Development (Recommended)

1. **Clone the repository:**
```bash
git clone git@github.com:ondatra-ai/awesome-claude-mcp.git
cd awesome-claude-mcp
```

2. **Start all services with Docker Compose:**
```bash
# Build and start in foreground (recommended for development)
docker-compose up --build

# Or start in background
docker-compose up -d --build
```

3. **Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

4. **Stop all services:**
```bash
docker-compose down
```

5. **View logs:**
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Local Development (Alternative)

For local development without Docker:

1. **Start the backend service:**
```bash
cd services/backend
go mod download
go run ./cmd/main.go
```

2. **Start the frontend service:** *(in a new terminal)*
```bash
cd services/frontend
npm install
npm run dev
```

## 🧪 Testing

### Docker E2E Testing (Recommended)

**Full automated testing pipeline:**
```bash
# Run E2E test suite with Docker (recommended)
make test-e2e
# or
./scripts/test.sh

# Or run tests directly with Docker Compose
docker-compose -f docker-compose.test.yml run --rm playwright-test
```

### Local E2E Testing (Playwright)
```bash
# First, ensure services are running (docker-compose up)

# Install Playwright browsers (first time only)
npm run install:browsers --prefix tests

# Run E2E tests
npm test --prefix tests

# Run tests with UI
npm run test:ui --prefix tests

# Run tests in headed mode (visible browser)
npm run test:headed --prefix tests

# Show test report
npm run test:report --prefix tests
```

### Unit Testing

**All Tests:**
```bash
make test           # Run all unit tests (backend + frontend)
```

**Individual Service Tests:**
```bash
make test-be        # Run backend unit tests only
make test-fe        # Run frontend unit tests only
```

**Manual Testing:**
```bash
# Backend Testing (Go)
go test ./services/backend/...

# Frontend Testing (Node.js)
npm test --prefix services/frontend
```

## 🏭 Production Deployment

### Docker Production Build

1. **Build optimized images:**
```bash
docker-compose build
```

2. **Deploy with production configuration:**
```bash
docker-compose -f docker-compose.yml up -d
```

### Manual Production Build

#### Backend
```bash
cd services/backend
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go
```

#### Frontend
```bash
cd services/frontend
npm run build
npm start
```

## ⚙️ Configuration

### Environment Variables

#### Backend Service
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

#### Frontend Service
- `NEXT_PUBLIC_API_URL`: Backend API URL (default: http://localhost:8080)
- `NODE_ENV`: Environment (development, production)

### Docker Environment
Environment variables are configured in `docker-compose.yml`:
- Backend connects to frontend via Docker network
- CORS configured for cross-service communication
- Health checks ensure service reliability

## 📊 Monitoring & Health Checks

### Health Endpoints
- **Backend**: `GET /health` - Returns service status and timestamp
- **Frontend**: Docker health check on port 3000

### Logging
- **Structured Logging**: JSON format with contextual information
- **Log Levels**: Configurable via LOG_LEVEL environment variable
- **Request Tracking**: Automatic HTTP request/response logging

## 🛠️ Development Workflow

### Code Quality
```bash
# Backend linting and formatting
cd services/backend
go fmt ./...
go vet ./...

# Frontend linting and type checking
cd services/frontend
npm run lint
npm run type-check
```

### Git Workflow
1. Create feature branches from `main`
2. Run tests locally before committing
3. Use conventional commit messages
4. Submit pull requests for code review

## 🔧 Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Kill processes using ports 3000 or 8080
lsof -ti:3000 | xargs kill -9
lsof -ti:8080 | xargs kill -9
```

#### Docker Issues
```bash
# Clean up Docker resources
docker-compose down
docker system prune -f

# Rebuild from scratch
docker-compose build --no-cache
```

#### CORS Issues
- Ensure backend CORS is configured for frontend origin
- Check NEXT_PUBLIC_API_URL environment variable
- Verify Docker network connectivity

### Debug Mode
```bash
# Backend with debug logging
cd services/backend
LOG_LEVEL=debug go run ./cmd/main.go

# Frontend with detailed errors
cd services/frontend
NODE_ENV=development npm run dev
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes following the coding standards
4. Add tests for new functionality
5. Run the full test suite: `npm run test:e2e`
6. Commit your changes: `git commit -m "feat: add your feature"`
7. Push to your branch: `git push origin feature/your-feature-name`
8. Submit a pull request

## 📝 License

This project is licensed under the ISC License.

## 📞 Support

For questions and support:
- Create an issue in the GitHub repository
- Check the troubleshooting section above
- Review the E2E tests for usage examples

---

**Built with ❤️ for Model Context Protocol integration**
