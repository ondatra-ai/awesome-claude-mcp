# MCP Google Docs Editor - Source Tree

This document summarizes the current monorepo layout after the migration to Railway. It replaces the prior cloud Infrastructure-as-Code structure.

## Root Directory Overview

```text
.
├── .github/
│   └── workflows/
│       └── deploy_to_railway.yml       # GitHub Actions deployment workflow
├── docs/                               # Product, architecture, QA documentation
├── services/
│   ├── frontend/                       # Next.js App Router frontend (React UI)
│   ├── backend/                        # Go REST API service (user management, OAuth)
│   └── mcp-service/                    # Go MCP Protocol Handler (HTTP+SSE, tool execution)
├── tests/                              # Playwright and other cross-service tests
├── scripts/                            # Utility scripts (linting, PR triage, etc.)
├── Makefile                            # Local tasks (lint, test, railway deploy)
├── docker-compose.yml                  # Local dev stack
├── railway.toml                        # Railway service definitions
├── service.toml                        # Railway CLI defaults
├── package.json                        # Root package for tooling
├── turbo.json                          # Turborepo configuration
└── CLAUDE.md                           # Claude guidance
```

## Services Directory

### `services/frontend`
- Next.js 14 App Router structure (`app/`, `components/`, `lib/`, etc.)
- Health check route exposed at `/`
- Environment variables consumed from Railway (e.g., `NEXT_PUBLIC_API_URL`)

### `services/backend`
- Go module with Fiber HTTP server (`cmd/main.go` entry point)
- REST endpoints (`/health`, `/version`) for user management API
- Dockerfile multi-stage build for Railway deployments

### `services/mcp-service`
- Go module with Mark3Labs MCP-Go server (`cmd/main.go` entry point)
- MCP protocol HTTP+SSE endpoints for Claude AI communication (Streamable HTTP)
- Tool registration and Google Docs operations
- Dockerfile multi-stage build for Railway deployments

## Deployment Configuration

- `railway.toml`: Defines four Railway services (production/staging/dev variants) and build contexts
- `service.toml`: CLI defaults for `railway up`
- `Makefile`: Targets `deploy-dev`, `deploy-staging`, `deploy-prod`, `deploy-service`
- `.github/workflows/deploy_to_railway.yml`: Maps branches/environments to Railway deployments

## Environment Strategy

| Environment | Railway Environment | Services | Domains |
|-------------|--------------------|----------|---------|
| Development | `development`      | `frontend-dev`, `backend-dev`, `mcp-service-dev` | `dev.ondatra-ai.xyz`, `api.dev.ondatra-ai.xyz`, `mcp.dev.ondatra-ai.xyz` |
| Staging     | `staging`          | `frontend-staging`, `backend-staging`, `mcp-service-staging` | `staging.ondatra-ai.xyz` (planned), `api.staging.ondatra-ai.xyz` (planned), `mcp.staging.ondatra-ai.xyz` (planned) |
| Production  | `production`       | `frontend`, `backend`, `mcp-service` | `app.ondatra-ai.xyz` (planned), `api.ondatra-ai.xyz` (planned), `mcp.ondatra-ai.xyz` (planned) |

## Detailed Service Structures

### `services/frontend/` - Next.js Frontend
```text
services/frontend/
├── app/                    # Next.js 14 App Router
│   ├── layout.tsx         # Root layout
│   ├── page.tsx           # Home page
│   └── api/               # API routes
├── components/            # React components
├── lib/                   # Utility libraries
├── public/                # Static assets
├── package.json
├── tsconfig.json
├── Dockerfile
└── .env
```

### `services/backend/` - Go REST API Service
```text
services/backend/
├── cmd/
│   └── main.go            # Entry point for REST API server
├── internal/
│   ├── api/               # HTTP endpoints and middleware
│   ├── auth/              # OAuth and JWT handling
│   ├── users/             # User management
│   ├── cache/             # Redis client
│   └── config/            # Configuration
├── pkg/
│   ├── errors/            # Custom errors
│   └── utils/             # Utility functions
├── go.mod
├── go.sum
├── Dockerfile
└── .env
```

### `services/mcp-service/` - Go MCP Protocol Handler
```text
services/mcp-service/
├── cmd/
│   └── main.go            # Entry point for MCP HTTP+SSE server
├── internal/
│   ├── server/            # MCP protocol server and HTTP+SSE handling
│   │   ├── mcp.go        # MCP server setup
│   │   ├── tools.go      # Tool registration
│   │   ├── handlers.go   # Tool handlers
│   │   └── middleware.go # Middleware
│   ├── operations/        # Document operations (replace, append, insert)
│   ├── docs/              # Google Docs API integration
│   ├── auth/              # OAuth for service accounts
│   ├── cache/             # Redis client
│   └── config/            # Configuration
├── pkg/
│   ├── types/             # MCP request/response types
│   ├── errors/            # Custom errors
│   └── utils/             # Utility functions
├── go.mod                 # Includes Mark3Labs MCP-Go dependency
├── go.sum
├── Dockerfile
└── .env
```

**CRITICAL: Service Separation**
- **Frontend** (`services/frontend/`): Next.js UI only
- **Backend** (`services/backend/`): REST API for user management and OAuth
- **MCP Service** (`services/mcp-service/`): MCP Protocol Handler for Claude/LLM communication
- MCP code is **NOT** in `services/backend/` - it has its own separate service

## Legacy Cloud Structure

The previous cloud-based infrastructure topology has been archived. If a future initiative revives it, consult historical commits prior to `2025-09-20` or the `docs/architecture.md` legacy appendix.
