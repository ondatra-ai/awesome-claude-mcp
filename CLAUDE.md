# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Information

- **Owner**: ondatra-ai
- **Repository**: awesome-claude-mcp
- **Host**: GitHub (github.com)
- **Clone URL**: git@github.com:ondatra-ai/awesome-claude-mcp.git
- **GitHub User**: @killev

## Repository Overview

This repository contains documentation and specifications for an MCP (Model Context Protocol) Google Docs Editor integration.

## Testing Approach

This project uses Playwright for both Integration (INT) and End-to-End (E2E) testing:

- **Integration (INT)**: Direct API/protocol testing using Playwright Request API
  - No UI required
  - Fast execution (seconds)
  - Tests: HTTP endpoints, WebSocket connections, MCP protocol

- **End-to-End (E2E)**: Complete user journeys using Playwright Browser API
  - Requires browser
  - Slower execution (minutes)
  - Tests: UI workflows, Claude.ai interactions, full-stack operations

- **Unit Tests**: Traditional code-level testing without BDD scenarios
  - Tested directly in code files
  - No Given-When-Then scenarios

**BDD Scenarios**: Only generated for Integration and E2E tests (not unit tests)

See `docs/architecture/bdd-guidelines.md` for scenario writing standards.

## Development Setup

### Quick Start
All development commands use the Makefile in the repository root:

```bash
make help  # Show all available commands
```

### Common Commands

#### Testing
```bash
make test-unit    # Run unit tests (Go backend + Jest frontend)
make test-e2e     # Run integration & E2E tests (starts Docker services automatically)
```

#### Development
```bash
make init         # Install dependencies and build Docker images
make dev          # Start all services with Docker Compose
```

#### Linting
```bash
make lint         # Run all linting checks
make lint-backend # Lint Go backend code
make lint-frontend # Lint TypeScript/React frontend code
make lint-scripts # Lint Go scripts (bmad-cli)
make lint-docs    # Validate YAML documentation
```

### BMAD CLI Usage
**CRITICAL**: Always run BMAD CLI from the repository root directory, not from `scripts/bmad-cli/`

#### Command Timeouts
Different commands require different timeout values based on their complexity:

**User Story Commands:**
```bash
# Create user story - 10 minutes (600 seconds)
go build -C scripts/bmad-cli -o ./bmad-cli && timeout 600 scripts/bmad-cli/bmad-cli us create 3.1

# Implement user story - 2 hours (7200 seconds)
go build -C scripts/bmad-cli -o ./bmad-cli && timeout 7200 scripts/bmad-cli/bmad-cli us implement 3.1

# Implement with force flag - 2 hours (7200 seconds)
go build -C scripts/bmad-cli -o ./bmad-cli && timeout 7200 scripts/bmad-cli/bmad-cli us implement 3.1 --force
```

**Pull Request Commands:**
```bash
# PR triage - 5 minutes (300 seconds)
go build -C scripts/bmad-cli -o ./bmad-cli && timeout 300 scripts/bmad-cli/bmad-cli pr triage
```

**Timeout Guidelines:**
- `us create`: 10 minutes (600s) - Story generation with AI processing
- `us implement`: 2 hours (7200s) - Full implementation with code generation
- `pr triage`: 5 minutes (300s) - PR analysis and triage

**Important Notes:**
- Always run from repository root for proper path resolution
- Timeouts ensure commands don't hang indefinitely during AI processing
- Commands respect configured engine type (see config files)

## Project Structure

Currently empty - update this section as the codebase develops.

## Railway Deployment

- **Railway Project ID**: `801ad5e0-95bf-4ce6-977e-6f2fa37529fd`
- **Environments**: `development`, `staging`, `production`
- **Primary Services**:
  - `frontend`, `backend` (production)
  - `frontend-staging`, `backend-staging`
  - `frontend-dev`, `backend-dev`
- **Custom Domains**:
  - Development: `dev.ondatra-ai.xyz`, `api.dev.ondatra-ai.xyz`
  - Staging (planned): `staging.ondatra-ai.xyz`, `api.staging.ondatra-ai.xyz`
  - Production (planned): `app.ondatra-ai.xyz`, `api.ondatra-ai.xyz`
- **CLI Basics**:
  - Authenticate: `railway login`
  - Link project: `railway link --project 801ad5e0-95bf-4ce6-977e-6f2fa37529fd`
  - Switch environment: `railway environment <development|staging|production>`
  - Deploy service: `railway up --service <name> --path-as-root services/<frontend|backend>`

## BMAD CLI Architecture Principles

### Quality Over Cost Principle
**QUALITY IS PARAMOUNT - TIME, PRICE, AND TOKEN USAGE ARE LOWEST PRIORITY** 🎯

When making decisions about BMAD CLI implementation:
- ✅ **Prioritize output quality**: Always choose the approach that produces the best results
- ✅ **Multi-stage generation is acceptable**: If it takes 3x tokens to get perfect output, do it
- ✅ **Take time for quality**: Generation time is not a concern if results are better
- ✅ **Token usage is not a constraint**: Use as many tokens as needed for comprehensive prompts
- ❌ **Never compromise quality for speed**: Fast but mediocre output is unacceptable
- ❌ **Never optimize for token cost**: Cutting corners on prompts to save tokens is wrong

**Examples:**
- Two-stage generation with critique? ✅ Do it
- Embed full articles in prompts? ✅ Do it
- Multiple validation passes? ✅ Do it
- Self-critique and revision loops? ✅ Do it

### Core Data Flow Principle
**NO CACHING, NO LOADERS, NO UNNECESSARY INTERFACES - JUST DIRECT DATA FLOW!** 🎉

#### The Principle
- ✅ **Direct data flow**: `Epic File → StoryFactory → StoryDocument → Generators`
- ✅ **Single source of truth**: StoryDocument contains all needed data
- ✅ **No abstraction layers**: Components work directly with concrete data
- ✅ **No caching complexity**: Load once, use directly
- ✅ **Simple interfaces**: Functions take concrete types, not abstractions

#### Implementation Guidelines
- **Extend domain models** (like StoryDocument) with required data instead of creating loaders
- **Pass complete data structures** to functions instead of IDs that require loading
- **Load data once** at the factory level and populate the complete structure
- **Avoid interfaces** unless there's a genuine need for multiple implementations
- **Question every layer** - if it doesn't add real value, remove it

#### Example: BMAD CLI Story Generation
```go
// ✅ GOOD: Direct data flow
type StoryDocument struct {
    Story            Story
    Tasks            []Task
    DevNotes         DevNotes
    Testing          Testing
    QAResults        *QAResults
    ArchitectureDocs *docs.ArchitectureDocs  // All data included
}

func (g *Generator) Generate(ctx context.Context, storyDoc *StoryDocument) (Result, error) {
    // Direct access to all needed data
    return processData(storyDoc.Story, storyDoc.ArchitectureDocs)
}

// ❌ BAD: Unnecessary abstractions
type StoryLoader interface { LoadStory(id string) (*Story, error) }
type DataCache struct { /* caching complexity */ }
```

#### When This Principle Was Established
- **Date**: 2025-09-28
- **Context**: BMAD CLI refactoring session
- **Result**: Eliminated 200+ lines of unnecessary abstraction code
- **Verification**: Story generation still works perfectly with much simpler code

## Go File Naming Convention

### Single Entity Single File Principle

**Rule**: Each Go file should contain one primary entity (struct/interface/type), and the filename must match the entity name in snake_case.

**Examples:**
- `GitHubService` struct → `github_service.go`
- `ClaudeClient` struct → `claude_client.go`
- `BranchManager` struct → `branch_manager.go`
- `ThreadProcessor` interface → `thread_processor.go`

**Important**: Treat "GitHub" as a single word (not "Git" + "Hub"):
- ✅ `GitHubService` → `github_service.go` (correct)
- ❌ `GitHubService` → `git_hub_service.go` (wrong)

**Exceptions (Acceptable):**
- Files with only functions/constants (e.g., `logging.go`, `errors.go`)
- Re-export files that aggregate types from internal packages
- Multiple closely related types (e.g., `ExecutionMode` + `ModeFactory` in `execution_mode.go`)
- Data structures bundled with their primary entity (e.g., `ArchitectureLoader` + `ArchitectureDoc`)

**Enforcement:**
- Currently enforced through code review
- No automated linter rule in `.golangci.yaml` yet
- See audit report for compliance status

**Current Status (as of 2025-10-10):**
- 113/115 files (98.3%) comply with this convention
- 2 files need renaming:
  - `internal/adapters/github/client.go` → `github_cli_client.go` (GitHubCLIClient)
  - `internal/adapters/github/queries.go` → `graphql_builder.go` (GraphQLBuilder)

## Notes

- The .gitignore is configured for Go projects
- Environment variables should be stored in .env files (excluded from git)
- Never Update @services/frontend/.eslintrc.json and @.golangci.yaml without my permission
- **CRITICAL**: NEVER merge pull requests without explicit user command to merge
- **CRITICAL**: NEVER use `git commit --amend` or `git push --force`/`--force-with-lease`. Always create new commits.
