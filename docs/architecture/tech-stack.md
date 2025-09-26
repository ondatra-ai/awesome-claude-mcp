# MCP Google Docs Editor – Technology Stack

This document captures the current technology stack after migrating from the legacy cloud-based deployment to Railway. It focuses on the components that are actively maintained today.

- **Last Updated:** 2025-09-20
- **Repository Layout:** Monorepo (`services/frontend`, `services/backend`)
- **Runtime Environments:** Railway Development, Staging, Production

## Frontend Stack (Next.js 14)

| Category | Technology / Tool | Version | Purpose | Notes |
|----------|-------------------|---------|---------|-------|
| Framework | Next.js | 14.x | App Router, SSR/ISR | `services/frontend` |
| Language | TypeScript | 5.x | Type safety | Enforced via `tsconfig.json` |
| Styling | Tailwind CSS | 3.x | Utility-first CSS | See `tailwind.config.ts` |
| UI Components | shadcn/ui | latest | Accessible component primitives | |
| State | Zustand | 4.x | Client state management | Lightweight, TypeScript-friendly |
| Forms | React Hook Form | 7.x | Form handling + validation | Integrated with Zod |
| Validation | Zod | 3.x | Schema validation | Shared between client/server |
| HTTP Client | Fetch API / Axios | 1.x | REST calls to backend | `lib/api.ts` |
| Testing | Jest + Playwright | latest | Unit + E2E | Playwright config under `tests/` |
| Linting | ESLint + Prettier | latest | Code quality | `npm run lint`, `npm run format` |

## Backend Stack (Go)

| Category | Technology / Tool | Version | Purpose | Notes |
|----------|-------------------|---------|---------|-------|
| Language | Go | 1.21 | Backend services | `services/backend` |
| Framework | Fiber | 2.x | HTTP routing | `cmd/main.go` |
| MCP Integration | Mark3Labs MCP-Go | latest | Model Context Protocol tools | Embedded in backend |
| Logging | zerolog | 1.x | Structured JSON logs | Output to stdout for Railway |
| Config | viper | 1.x | Environment/config mgmt | Reads Railway env vars |
| OAuth | golang.org/x/oauth2 | latest | Google OAuth flows | |
| Google APIs | google.golang.org/api/docs/v1 | latest | Doc operations | |
| Markdown | goldmark | 1.6 | Markdown parsing | |
| Testing | Go test + testify | latest | Unit/integration tests | `go test ./...` |

## Platform & Deployment

| Category | Technology / Tool | Purpose | Notes |
|----------|-------------------|---------|-------|
| Hosting | Railway | Managed container runtime, TLS, custom domains | Project ID `801ad5e0-95bf-4ce6-977e-6f2fa37529fd` |
| Environments | Railway environments | `development`, `staging`, `production` | Services: `frontend[-dev|-staging]`, `backend[-dev|-staging]` |
| CI/CD | GitHub Actions | Workflow `.github/workflows/deploy_to_railway.yml` | Branch-driven environment selection |
| CLI | Railway CLI (`@railway/cli`) | Local deployments, CI commands | `railway login`, `railway up --service ... --path-as-root ...` |
| Config | `railway.toml`, `service.toml` | Service definitions & defaults | Maintained at repo root |
| Secrets | Railway env vars | API URLs, OAuth secrets | Managed per service/environment |
| Domains | CNAMEs (`dev.ondatra-ai.xyz`, `api.dev.ondatra-ai.xyz`, …) | Map to Railway `*.up.railway.app` endpoints | Staging/prod domains in progress |

## Supporting Tooling

- **Makefile targets:** `make deploy-dev`, `make deploy-staging`, `make deploy-prod`, plus lint/test helpers
- **Local development:** `docker compose up` (optional) or run services individually
- **Testing pipelines:**
  - Unit: `make test-unit`
  - E2E: `make test-e2e` (uses Playwright + Docker)
- **Linting:** `make lint-frontend`, `make lint-backend`, `make lint-scripts`

## Legacy Note

The original cloud infrastructure stack has been archived. Refer to historical commits prior to `2025-09-20` if that infrastructure needs to be revisited.
