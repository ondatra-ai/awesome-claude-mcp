# Epic 1: Foundation & Infrastructure

**Goal:** Establish deployable application foundation with proper architecture, testing framework, and monitoring

## User Stories

### Story 1.1: Minimal Frontend-Backend Integration Setup
**As a** Developer/Maintainer
**I want** to create a minimal Next.js frontend connected to a Go backend
**So that** I have a working full-stack foundation to build upon

**Acceptance Criteria:**
- Monorepo structure with services/ directory created
- Go backend service with `/version` endpoint returning "1.0.0"
- Next.js 14 single-page frontend application
- Frontend homepage displays backend version "1.0.0" at the bottom
- Docker and docker-compose configuration for local development
- Playwright E2E testing framework configured
- Go modules initialized for backend service
- Next.js project initialized with TypeScript
- Dockerfiles created for both services
- Basic project structure follows architecture document
- Health check endpoints for both services
- README.md with setup and run instructions

### Story 1.2: Railway Infrastructure Setup
**As a** Developer/Maintainer
**I want** to configure Railway environments and services
**So that** I have a deployable environment for the application

**Acceptance Criteria:**
- Railway project created and linked to repository
- Development, Staging, and Production environments defined
- Railway services created for frontend/backend (including environment-specific variants)
- Custom domains mapped (`dev.ondatra-ai.xyz`, `api.dev.ondatra-ai.xyz`, etc.)
- Environment variables configured per service (API URLs, OAuth secrets)
- Deployment verified via Railway dashboard for each environment

### Story 1.3: Frontend Service Implementation
**As a** Claude User
**I want** to access a web interface for service management
**So that** I can configure authentication and monitor service status

**Acceptance Criteria:**
- Next.js 14 frontend service deployed and accessible
- Homepage displays "MCP Google Docs Editor" title
- Service status dashboard (operational/degraded/down)
- OAuth authentication management interface
- Connected Google accounts display
- Mobile responsive design with modern UI
- Page loads in under 2 seconds
- Health check endpoint returns proper status

### Story 1.4: CI/CD Pipeline
**As a** Developer/Maintainer
**I want** automated build and deployment
**So that** code changes are safely deployed

**Acceptance Criteria:**
- GitHub Actions workflow (`deploy_to_railway.yml`) configured for all environments
- Railway CLI installed and authenticated via GitHub secret (`RAILWAY_GITHUB_ACTIONS`)
- Automated tests run on pull requests for each service
- Successful builds deploy targeted services to Railway environments based on branch naming
- Manual workflow dispatch supports explicit environment selection
- Build status badges in README (optional)
- Deployment notifications to Slack (optional enhancement)

### Story 1.5: Monitoring Setup
**As a** Developer/Maintainer
**I want** comprehensive monitoring
**So that** I can track system health and performance

**Acceptance Criteria:**
- Application logs accessible via Railway dashboard for each environment
- Health checks exposed on backend services (`/health`) and validated by Railway
- Optional third-party monitoring (e.g., Better Stack, Sentry) documented for future adoption
- Alerts established for service downtime using Railway notifications or external tooling
- Basic dashboard/reporting documented (Railway metrics, custom tooling TBD)

### Story 1.6: Testing Framework
**As a** Developer/Maintainer
**I want** comprehensive testing infrastructure
**So that** I can ensure code quality across all services

**Acceptance Criteria:**
- Unit test framework configured (testify for Go services, Jest for Next.js)
- Integration test environment setup for service-to-service communication
- E2E test framework ready (Playwright for frontend workflows)
- Docker Compose for local testing environment
- Test coverage reporting enabled for all services
- Container-based testing for deployment validation
- Pre-commit hooks for testing all services
- Example tests for each service and test type
