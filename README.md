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

### Prerequisites Setup
1. **Clone the repository:**
```bash
git clone git@github.com:ondatra-ai/awesome-claude-mcp.git
cd awesome-claude-mcp
```

2. **Initialize project (install dependencies and build Docker images):**
```bash
make init
```

### Development Workflow

1. **Start all services:**
```bash
make dev
```

2. **Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

3. **Stop all services:**
```bash
# Press Ctrl+C to stop, then clean up with:
docker compose down
```

## 🧪 Testing

### E2E Testing

**Prerequisites:**
Create a `.env.test` file in the project root with required environment variables:

```bash
# Copy the example file and add your API key
cp .env.test.example .env.test

# Then edit .env.test and add your Anthropic API key
# ANTHROPIC_API_KEY=your_actual_api_key_here
```

**Running Tests:**
```bash
# Run complete E2E test suite with Docker
make test-e2e
```

### Unit Testing
```bash
# Run all unit tests (backend + frontend)
make test-unit
```

## 🏭 Production Deployment

### Docker Production Build

1. **Initialize and build optimized images:**
```bash
make init
```

2. **Deploy with production configuration:**
```bash
docker compose -f docker-compose.yml up -d
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
make lint-backend

# Frontend linting and formatting
make lint-frontend
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
# Start services in development mode (includes debug logging)
make dev
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes following the coding standards
4. Add tests for new functionality
5. Run the full test suite: `make test-e2e`
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
