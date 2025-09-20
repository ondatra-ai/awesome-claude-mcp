# MCP Google Docs Editor - Source Tree (Railway Deployment)

This document summarizes the current monorepo layout after the migration to Railway. It replaces the prior AWS Infrastructure-as-Code structure.

## Root Directory Overview

```text
.
├── .github/
│   └── workflows/
│       └── deploy_to_railway.yml       # GitHub Actions deployment workflow
├── docs/                               # Product, architecture, QA documentation
├── services/
│   ├── frontend/                       # Next.js App Router frontend
│   └── backend/                        # Go API (MCP tooling co-located)
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
- REST endpoints (`/health`, `/version`) plus MCP tooling integration points
- Dockerfile multi-stage build for Railway deployments

## Deployment Configuration

- `railway.toml`: Defines four Railway services (production/staging/dev variants) and build contexts
- `service.toml`: CLI defaults for `railway up`
- `Makefile`: Targets `deploy-dev`, `deploy-staging`, `deploy-prod`, `deploy-service`
- `.github/workflows/deploy_to_railway.yml`: Maps branches/environments to Railway deployments

## Environment Strategy

| Environment | Railway Environment | Services | Domains |
|-------------|--------------------|----------|---------|
| Development | `development`      | `frontend-dev`, `backend-dev` | `dev.ondatra-ai.xyz`, `api.dev.ondatra-ai.xyz` |
| Staging     | `staging`          | `frontend-staging`, `backend-staging` | `staging.ondatra-ai.xyz` (planned), `api.staging.ondatra-ai.xyz` (planned) |
| Production  | `production`       | `frontend`, `backend` | `app.ondatra-ai.xyz` (planned), `api.ondatra-ai.xyz` (planned) |

## Legacy AWS Structure

The previous AWS infrastructure topology has been archived. If a future initiative revives it, consult historical commits prior to `2025-09-20` or the `docs/architecture.md` legacy appendix.
