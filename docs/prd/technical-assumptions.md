# Technical Assumptions

## Repository Structure: Monorepo

All services, including the Frontend Service (Next.js), Backend Service (Go Fiber), MCP Service (Mark3Labs MCP-Go), and any future components will be maintained in a single monorepo structure to simplify development and deployment for the single developer team.

## Service Architecture

The system will be implemented as a Dockerized architecture deployed to Railway. Railway provides managed environments (Development, Staging, Production), automatic TLS, and container orchestration. The monorepo ships two primary services—Frontend (Next.js) and Backend (Go Fiber)—with environment-specific variants (`-dev`, `-staging`) hosted as separate Railway services. MCP tooling is delivered via the backend service and extended as needed.

## Testing Requirements

Comprehensive testing pyramid including:
- Unit tests for all business logic (minimum 80% coverage)
- Integration tests for external API interactions
- End-to-end tests for complete workflows
- All tests must pass before any deployment
- Test implementation required for every user story

## Additional Technical Assumptions and Requests

- **Backend Services:** Go 1.21.5 with Fiber framework (MCP tooling embedded in backend service)
- **Frontend Service:** TypeScript with Next.js 14 (App Router) and modern React patterns
- **MCP Protocol:** Mark3Labs MCP-Go library for MCP protocol implementation
- Railway as the hosting provider (Docker deploys, managed TLS, custom domains)
- Redis cache (Railway add-on or external provider) for token/session storage
- Docker containerization with multi-stage builds for all services
- Application logging handled within services; Railway logs used for deployment visibility
- Standard MCP protocol implementation without extensions
- OAuth tokens cached until expiry by default
- No access controls or rate limiting in MVP (add in v2)
- Fail-fast error handling - no automatic retries
- All configuration via environment variables managed per Railway environment
- GitHub Actions + Railway CLI for CI/CD (`deploy_to_railway.yml`)
- Markdown parsing using goldmark library
- Structured JSON logging for all operations
