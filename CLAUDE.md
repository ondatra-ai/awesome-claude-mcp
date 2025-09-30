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

## Development Setup

Since this is a new repository without code, common setup tasks will depend on the project type:

### For Go Projects
- Run tests: `go test ./...`
- Build: `go build`
- Format code: `go fmt ./...`
- Lint: `golangci-lint run` (if installed)

### BMAD CLI Usage
**CRITICAL**: Always run BMAD CLI from the repository root directory, not from `scripts/bmad-cli/`
- Build: `go build -C scripts/bmad-cli -o ./bmad-cli`
- Run: `scripts/bmad-cli/bmad-cli sm us-create 3.1` (from root directory)
- This ensures proper path resolution for config files and tmp directories

### For Other Project Types
Update this file once the project structure is established with:
- Build commands
- Test commands
- Linting/formatting commands
- Local development setup

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

## Notes

- The .gitignore is configured for Go projects
- Environment variables should be stored in .env files (excluded from git)
- Never Update @services/frontend/.eslintrc.json and @.golangci.yml without my permission
