# MCP Google Docs Editor

A minimal full-stack application for editing Google Docs via the Model Context Protocol (MCP). This project provides a foundation with Next.js frontend and Go backend services.

## Architecture

- **Frontend**: Next.js 14 with TypeScript and Tailwind CSS
- **Backend**: Go with Fiber framework
- **Deployment**: Docker containerization with docker-compose
- **Testing**: Playwright E2E, Jest for frontend, testify for Go backend

## Quick Start

### Prerequisites

- Go 1.21.5+
- Node.js 18+
- Docker and Docker Compose

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd awesome-claude-mcp
   ```

2. **Start development environment**
   ```bash
   # Option 1: Using Make (recommended)
   make dev

   # Option 2: Using Docker Compose directly
   docker-compose up --build
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Manual Development

If you prefer to run services individually:

1. **Start Backend**
   ```bash
   cd services/backend
   go run ./cmd/main.go
   ```

2. **Start Frontend**
   ```bash
   cd services/frontend
   npm install
   npm run dev
   ```

## API Endpoints

### Backend (Port 8080)

- `GET /version` - Returns API version (1.0.0)
- `GET /health` - Health check endpoint

### Frontend (Port 3000)

- `GET /` - Homepage displaying backend version
- `GET /api/health` - Frontend health check

## Testing

### Backend Tests (Go)
```bash
cd services/backend
go test ./...
```

### Frontend Tests (Jest)
```bash
cd services/frontend
npm run test
```

### End-to-End Tests (Playwright)
Run unified E2E tests for both frontend UI and backend API integration:

```bash
cd tests/e2e
npx playwright test
```

This runs:
- **Frontend UI Tests** (`tests/e2e/frontend/`): Homepage functionality, authentication flows, document management UI
- **Backend API Tests** (`tests/e2e/backend/`): Core API functionality, authentication endpoints, CORS validation

### Run All Tests
```bash
make test
```

## Project Structure

```
├── services/
│   ├── backend/          # Go backend service
│   │   ├── cmd/          # Application entry point
│   │   ├── internal/     # Internal packages
│   │   ├── pkg/          # Public packages
│   │   └── Dockerfile    # Backend container definition
│   └── frontend/         # Next.js frontend service
│       ├── app/          # Next.js app directory
│       ├── components/   # React components (future)
│       ├── lib/          # Utilities (future)
│       └── Dockerfile    # Frontend container definition
├── tests/
│   └── e2e/              # Playwright E2E tests
├── docker-compose.yml    # Development stack
├── Makefile             # Development commands
└── README.md            # This file
```

## Development Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start development environment |
| `make dev-up` | Start services in background |
| `make dev-down` | Stop development services |
| `make dev-logs` | View service logs |
| `make test` | Run all tests |
| `make clean` | Clean development environment |
| `make help` | Show available commands |

## Environment Variables

### Backend
- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment mode (development/production)

### Frontend
- `NEXT_PUBLIC_API_URL` - Backend API URL (default: http://localhost:8080)

## Docker

### Building Images

```bash
# Backend
docker build -t mcp-backend ./services/backend

# Frontend
docker build -t mcp-frontend ./services/frontend
```

### Running with Docker Compose

```bash
# Development
docker-compose up --build

# Background
docker-compose up -d --build

# Stop
docker-compose down
```

## Health Checks

Both services include health check endpoints:

- Backend: `curl http://localhost:8080/health`
- Frontend: `curl http://localhost:3000/api/health`

## Next Steps

This is a foundational setup. Future development will include:

1. MCP protocol implementation
2. Google Docs API integration
3. OAuth authentication
4. Document operation tools
5. Advanced frontend features

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

This project is licensed under the MIT License.