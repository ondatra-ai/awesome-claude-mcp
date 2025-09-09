# Technical Assumptions

## Repository Structure: Monorepo

All services, including the Frontend Service (Next.js), Backend Service (Go Fiber), MCP Service (Mark3Labs MCP-Go), and any future components will be maintained in a single monorepo structure to simplify development and deployment for the single developer team.

## Service Architecture

The system will be implemented as a 3-service containerized architecture deployed on AWS ECS Fargate with Application Load Balancer, utilizing AWS services for infrastructure. The architecture includes Frontend Service (Next.js), Backend Service (Go Fiber), and MCP Service (Go with Mark3Labs MCP-Go library) following the design patterns established in the architecture document.

## Testing Requirements

Comprehensive testing pyramid including:
- Unit tests for all business logic (minimum 80% coverage)
- Integration tests for external API interactions
- End-to-end tests for complete workflows
- All tests must pass before any deployment
- Test implementation required for every user story

## Additional Technical Assumptions and Requests

- **Backend Services:** Go 1.21.5 with Fiber framework for Backend and MCP services
- **Frontend Service:** TypeScript with Next.js 14 (App Router) and modern React patterns
- **MCP Protocol:** Mark3Labs MCP-Go library for MCP protocol implementation
- AWS as cloud infrastructure provider (ECS Fargate, Application Load Balancer, CloudWatch)
- Redis for token caching (AWS ElastiCache or similar)
- Docker containerization with multi-stage builds for all services
- New Relic + CloudWatch for monitoring and observability
- Standard MCP protocol implementation without extensions
- OAuth tokens cached until expiry by default
- No access controls or rate limiting in MVP (add in v2)
- Fail-fast error handling - no automatic retries
- All configuration via environment variables
- Infrastructure as Code using Terraform
- GitHub Actions for CI/CD pipeline with ECR integration
- Markdown parsing using goldmark library
- Structured JSON logging for all operations
