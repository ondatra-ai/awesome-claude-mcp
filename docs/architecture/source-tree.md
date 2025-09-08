# MCP Google Docs Editor - Source Tree Organization

## Introduction

This document defines the complete source tree structure for the MCP Google Docs Editor project. It establishes a monorepo organization with clear separation between frontend (Next.js) and backend (Go) services, supporting the development requirements outlined in the PRD.

**Structure Type:** Monorepo  
**Primary Languages:** Go (backend), TypeScript (frontend)  
**Target Audience:** AI development agents and human developers  
**Last Updated:** 2025-09-07

## Root Directory Structure

```text
mcp-google-docs-editor/
├── .github/                    # GitHub-specific configuration
│   ├── workflows/              # CI/CD pipeline definitions
│   │   ├── ci.yml              # Main CI/CD workflow
│   │   ├── security-scan.yml   # Security scanning workflow
│   │   └── deploy.yml          # Production deployment workflow
│   ├── ISSUE_TEMPLATE/         # GitHub issue templates
│   ├── PULL_REQUEST_TEMPLATE.md # PR template
│   └── dependabot.yml         # Dependency update configuration
├── frontend/                   # Next.js Frontend Application
├── backend/                    # Go Backend Services
├── infrastructure/             # Infrastructure as Code
├── scripts/                    # Build and deployment scripts
├── docs/                       # Project documentation
├── tests/                      # Cross-service integration tests
├── .gitignore                  # Git ignore patterns
├── .env.example               # Environment variable template
├── docker-compose.yml         # Local development stack
├── Makefile                   # Build automation
├── README.md                  # Project overview and setup
├── LICENSE                    # MIT license file
└── CLAUDE.md                  # Claude Code configuration
```

## Frontend Directory Structure (Next.js 14)

```text
frontend/
├── app/                        # Next.js App Router (primary routing)
│   ├── (auth)/                 # Route group for authentication
│   │   ├── login/              # Login page route
│   │   │   └── page.tsx        # Login page component
│   │   ├── callback/           # OAuth callback route
│   │   │   └── page.tsx        # OAuth callback handler
│   │   └── layout.tsx          # Auth layout wrapper
│   ├── dashboard/              # Dashboard route
│   │   ├── page.tsx            # Dashboard page component
│   │   └── loading.tsx         # Loading UI for dashboard
│   ├── documents/              # Document management routes
│   │   ├── page.tsx            # Document list page
│   │   ├── [id]/               # Dynamic document routes
│   │   │   ├── page.tsx        # Document detail page
│   │   │   ├── edit/           # Document editing route
│   │   │   │   └── page.tsx    # Document editor page
│   │   │   └── loading.tsx     # Document loading UI
│   │   └── new/                # New document route
│   │       └── page.tsx        # New document page
│   ├── settings/               # User settings routes
│   │   ├── page.tsx            # Settings overview
│   │   ├── accounts/           # Account management
│   │   │   └── page.tsx        # Connected accounts page
│   │   └── api-keys/           # API key management
│   │       └── page.tsx        # API keys page
│   ├── api/                    # API Routes (Server-side)
│   │   ├── auth/               # Authentication API routes
│   │   │   ├── signin/         # Sign-in endpoint
│   │   │   │   └── route.ts    # Sign-in API handler
│   │   │   ├── callback/       # OAuth callback endpoint
│   │   │   │   └── route.ts    # OAuth callback handler
│   │   │   └── signout/        # Sign-out endpoint
│   │   │       └── route.ts    # Sign-out API handler
│   │   ├── documents/          # Document proxy API routes
│   │   │   ├── route.ts        # Document list API
│   │   │   └── [id]/           # Document-specific APIs
│   │   │       ├── route.ts    # Single document API
│   │   │       └── operations/ # Document operations
│   │   │           └── route.ts # Operations proxy
│   │   ├── health/             # Health check endpoint
│   │   │   └── route.ts        # Health check handler
│   │   └── webhook/            # Webhook endpoints
│   │       └── route.ts        # Webhook handler
│   ├── globals.css             # Global CSS styles
│   ├── layout.tsx              # Root layout component
│   ├── page.tsx                # Homepage component
│   ├── loading.tsx             # Global loading component
│   ├── error.tsx               # Global error component
│   ├── not-found.tsx           # 404 page component
│   └── favicon.ico             # Application favicon
├── components/                 # React Components
│   ├── ui/                     # Basic UI components (shadcn/ui)
│   │   ├── button.tsx          # Button component
│   │   ├── input.tsx           # Input component
│   │   ├── card.tsx            # Card component
│   │   ├── dialog.tsx          # Dialog/Modal component
│   │   ├── form.tsx            # Form components
│   │   ├── table.tsx           # Table component
│   │   ├── badge.tsx           # Badge component
│   │   ├── alert.tsx           # Alert component
│   │   └── skeleton.tsx        # Loading skeleton component
│   ├── forms/                  # Form-specific components
│   │   ├── LoginForm.tsx       # Login form component
│   │   ├── DocumentForm.tsx    # Document creation/editing form
│   │   ├── OperationForm.tsx   # Document operation form
│   │   └── SettingsForm.tsx    # User settings form
│   ├── layout/                 # Layout components
│   │   ├── Header.tsx          # Page header component
│   │   ├── Footer.tsx          # Page footer component
│   │   ├── Sidebar.tsx         # Navigation sidebar
│   │   ├── Navigation.tsx      # Main navigation component
│   │   └── Breadcrumbs.tsx     # Breadcrumb navigation
│   ├── documents/              # Document-related components
│   │   ├── DocumentList.tsx    # Document list display
│   │   ├── DocumentCard.tsx    # Individual document card
│   │   ├── DocumentEditor.tsx  # Document editing interface
│   │   ├── OperationPanel.tsx  # Document operations panel
│   │   └── StatusIndicator.tsx # Operation status display
│   ├── auth/                   # Authentication components
│   │   ├── AuthProvider.tsx    # Authentication context provider
│   │   ├── LoginButton.tsx     # Google login button
│   │   ├── LogoutButton.tsx    # Logout button
│   │   ├── AuthGuard.tsx       # Route protection component
│   │   └── AccountSelector.tsx # Multi-account selector
│   └── monitoring/             # Monitoring and debugging components
│       ├── StatusPanel.tsx     # System status display
│       ├── ErrorBoundary.tsx   # Error boundary component
│       └── DebugPanel.tsx      # Development debug panel
├── hooks/                      # Custom React Hooks
│   ├── useAuth.ts              # Authentication hook
│   ├── useDocuments.ts         # Document management hook
│   ├── useOperations.ts        # Document operations hook
│   ├── useWebSocket.ts         # WebSocket connection hook
│   ├── useLocalStorage.ts      # Local storage hook
│   └── useApi.ts               # API client hook
├── lib/                        # Utility Libraries and Configurations
│   ├── api.ts                  # API client configuration
│   ├── auth.ts                 # Authentication utilities
│   ├── websocket.ts            # WebSocket client
│   ├── storage.ts              # Local storage utilities
│   ├── utils.ts                # General utility functions
│   ├── constants.ts            # Application constants
│   ├── validators.ts           # Input validation functions
│   └── formatters.ts           # Data formatting utilities
├── types/                      # TypeScript Type Definitions
│   ├── api.ts                  # API request/response types
│   ├── auth.ts                 # Authentication types
│   ├── document.ts             # Document-related types
│   ├── operation.ts            # Document operation types
│   ├── user.ts                 # User profile types
│   └── global.ts               # Global type definitions
├── styles/                     # Additional Style Files
│   ├── globals.css             # Additional global styles
│   └── components.css          # Component-specific styles
├── public/                     # Static Assets
│   ├── images/                 # Image assets
│   │   ├── logo.svg            # Application logo
│   │   ├── hero-bg.jpg         # Homepage hero background
│   │   └── icons/              # Icon assets
│   ├── fonts/                  # Custom font files
│   └── manifest.json           # PWA manifest (future)
├── tests/                      # Frontend Tests
│   ├── __mocks__/              # Jest mocks
│   ├── components/             # Component tests
│   ├── hooks/                  # Hook tests
│   ├── lib/                    # Utility tests
│   └── setup.ts                # Test setup configuration
├── .env.local                  # Local environment variables
├── .env.example               # Environment variables template
├── .eslintrc.json             # ESLint configuration
├── .prettierrc                # Prettier configuration
├── .gitignore                 # Frontend-specific git ignores
├── jest.config.js             # Jest testing configuration
├── next.config.js             # Next.js configuration
├── package.json               # Node.js dependencies and scripts
├── package-lock.json          # Dependency lock file
├── postcss.config.js          # PostCSS configuration
├── tailwind.config.ts         # Tailwind CSS configuration
├── tsconfig.json              # TypeScript configuration
└── README.md                  # Frontend-specific documentation
```

## Backend Directory Structure (Go Services)

```text
backend/
├── cmd/                        # Application Entry Points
│   ├── api/                    # REST API Server
│   │   ├── main.go             # API server entry point
│   │   └── config.go           # API server configuration
│   └── mcp/                    # MCP WebSocket Server
│       ├── main.go             # MCP server entry point
│       └── config.go           # MCP server configuration
├── internal/                   # Internal Packages (Private)
│   ├── api/                    # HTTP API Implementation
│   │   ├── handlers/           # HTTP request handlers
│   │   │   ├── auth.go         # Authentication handlers
│   │   │   ├── documents.go    # Document operation handlers
│   │   │   ├── health.go       # Health check handlers
│   │   │   └── middleware.go   # HTTP middleware
│   │   ├── routes/             # Route definitions
│   │   │   ├── auth.go         # Authentication routes
│   │   │   ├── documents.go    # Document operation routes
│   │   │   └── routes.go       # Main route setup
│   │   └── server.go           # HTTP server setup
│   ├── mcp/                    # MCP Protocol Implementation
│   │   ├── server.go           # MCP WebSocket server
│   │   ├── handlers/           # MCP message handlers
│   │   │   ├── tools.go        # Tool registration handlers
│   │   │   ├── operations.go   # Document operation handlers
│   │   │   └── discovery.go    # Service discovery handlers
│   │   ├── protocol/           # MCP protocol implementation
│   │   │   ├── messages.go     # Message types and validation
│   │   │   ├── transport.go    # Transport layer handling
│   │   │   └── client.go       # Client connection management
│   │   └── tools/              # MCP tool definitions
│   │       ├── replace_all.go  # Replace all operation tool
│   │       ├── append.go       # Append operation tool
│   │       ├── prepend.go      # Prepend operation tool
│   │       ├── replace_match.go # Replace match operation tool
│   │       ├── insert_before.go # Insert before operation tool
│   │       └── insert_after.go # Insert after operation tool
│   ├── auth/                   # Authentication and Authorization
│   │   ├── oauth.go            # OAuth 2.0 implementation
│   │   ├── tokens.go           # Token management
│   │   ├── middleware.go       # Authentication middleware
│   │   ├── google.go           # Google-specific auth logic
│   │   └── cache.go            # Token caching
│   ├── operations/             # Document Operations Business Logic
│   │   ├── processor.go        # Main document processor
│   │   ├── markdown.go         # Markdown parsing and conversion
│   │   ├── replace_all.go      # Replace all operation
│   │   ├── append.go           # Append operation
│   │   ├── prepend.go          # Prepend operation
│   │   ├── replace_match.go    # Replace match operation
│   │   ├── insert_before.go    # Insert before operation
│   │   ├── insert_after.go     # Insert after operation
│   │   └── validator.go        # Operation input validation
│   ├── docs/                   # Google Docs Integration
│   │   ├── client.go           # Google Docs API client
│   │   ├── service.go          # Document service wrapper
│   │   ├── formatter.go        # Document formatting utilities
│   │   ├── batch.go            # Batch operation handler
│   │   └── errors.go           # Google Docs specific errors
│   ├── cache/                  # Caching Layer
│   │   ├── redis.go            # Redis client implementation
│   │   ├── tokens.go           # Token caching logic
│   │   ├── documents.go        # Document metadata caching
│   │   └── interface.go        # Caching interface definition
│   └── config/                 # Configuration Management
│       ├── config.go           # Configuration structure and loading
│       ├── env.go              # Environment variable handling
│       └── validation.go       # Configuration validation
├── pkg/                        # Public Packages (Importable)
│   ├── errors/                 # Custom Error Types
│   │   ├── api.go              # API error types
│   │   ├── auth.go             # Authentication error types
│   │   ├── docs.go             # Document operation error types
│   │   └── common.go           # Common error utilities
│   ├── utils/                  # Utility Functions
│   │   ├── strings.go          # String manipulation utilities
│   │   ├── time.go             # Time handling utilities
│   │   ├── crypto.go           # Cryptographic utilities
│   │   └── validation.go       # Input validation utilities
│   └── logger/                 # Logging Package
│       ├── logger.go           # Main logger implementation
│       ├── context.go          # Context-aware logging
│       └── formatters.go       # Log formatting utilities
├── deployments/                # Deployment Configurations
│   ├── docker/                 # Docker configurations
│   │   ├── Dockerfile.api      # API service Dockerfile
│   │   ├── Dockerfile.mcp      # MCP service Dockerfile
│   │   └── docker-compose.yml  # Local development compose
│   ├── aws/                    # AWS deployment configurations
│   │   ├── template.yaml       # SAM template
│   │   └── buildspec.yml       # AWS CodeBuild specification
│   └── k8s/                    # Kubernetes manifests (future)
│       └── deployment.yaml     # Kubernetes deployment (future)
├── tests/                      # Test Files and Utilities
│   ├── integration/            # Integration Tests
│   │   ├── api_test.go         # API integration tests
│   │   ├── mcp_test.go         # MCP protocol integration tests
│   │   ├── auth_test.go        # Authentication flow tests
│   │   └── documents_test.go   # Document operation tests
│   ├── fixtures/               # Test Data and Fixtures
│   │   ├── documents/          # Sample document data
│   │   │   ├── simple.json     # Simple document fixture
│   │   │   ├── complex.json    # Complex document fixture
│   │   │   └── markdown.md     # Sample Markdown content
│   │   ├── auth/               # Authentication test data
│   │   │   └── tokens.json     # Sample OAuth tokens
│   │   └── responses/          # API response fixtures
│   │       ├── success.json    # Successful operation responses
│   │       └── errors.json     # Error response examples
│   ├── helpers/                # Test Helper Functions
│   │   ├── auth.go             # Authentication test helpers
│   │   ├── documents.go        # Document test helpers
│   │   ├── server.go           # Test server utilities
│   │   └── mocks.go            # Mock generation utilities
│   └── mocks/                  # Generated Mocks
│       ├── auth_mock.go        # Authentication service mocks
│       ├── docs_mock.go        # Google Docs client mocks
│       └── cache_mock.go       # Cache interface mocks
├── scripts/                    # Build and Utility Scripts
│   ├── build.sh                # Build script
│   ├── test.sh                 # Test execution script
│   ├── lint.sh                 # Linting script
│   ├── generate-mocks.sh       # Mock generation script
│   └── migrate.sh              # Database migration script (future)
├── docs/                       # Backend-Specific Documentation
│   ├── api.md                  # API documentation
│   ├── mcp.md                  # MCP protocol implementation
│   └── deployment.md           # Deployment instructions
├── .env.example               # Environment variables template
├── .gitignore                 # Go-specific git ignores
├── .golangci.yml              # golangci-lint configuration
├── .dockerignore              # Docker ignore patterns
├── go.mod                     # Go module definition
├── go.sum                     # Go dependency checksums
├── Makefile                   # Build automation
└── README.md                  # Backend-specific documentation
```

## Infrastructure Directory Structure

```text
infrastructure/
├── terraform/                  # Infrastructure as Code
│   ├── environments/           # Environment-specific configurations
│   │   ├── dev/                # Development environment
│   │   │   ├── main.tf         # Development infrastructure
│   │   │   ├── variables.tf    # Development variables
│   │   │   └── terraform.tfvars # Development values
│   │   ├── staging/            # Staging environment
│   │   │   ├── main.tf         # Staging infrastructure
│   │   │   ├── variables.tf    # Staging variables
│   │   │   └── terraform.tfvars # Staging values
│   │   └── prod/               # Production environment
│   │       ├── main.tf         # Production infrastructure
│   │       ├── variables.tf    # Production variables
│   │       └── terraform.tfvars # Production values
│   ├── modules/                # Reusable Terraform modules
│   │   ├── lambda/             # Lambda function module
│   │   │   ├── main.tf         # Lambda resources
│   │   │   ├── variables.tf    # Lambda variables
│   │   │   └── outputs.tf      # Lambda outputs
│   │   ├── api-gateway/        # API Gateway module
│   │   │   ├── main.tf         # API Gateway resources
│   │   │   ├── variables.tf    # API Gateway variables
│   │   │   └── outputs.tf      # API Gateway outputs
│   │   ├── redis/              # ElastiCache Redis module
│   │   │   ├── main.tf         # Redis resources
│   │   │   ├── variables.tf    # Redis variables
│   │   │   └── outputs.tf      # Redis outputs
│   │   └── monitoring/         # CloudWatch monitoring module
│   │       ├── main.tf         # Monitoring resources
│   │       ├── variables.tf    # Monitoring variables
│   │       └── outputs.tf      # Monitoring outputs
│   ├── shared/                 # Shared Terraform configurations
│   │   ├── providers.tf        # Provider configurations
│   │   ├── backend.tf          # Remote state configuration
│   │   └── versions.tf         # Terraform version constraints
│   ├── main.tf                 # Root Terraform configuration
│   ├── variables.tf            # Global variables
│   ├── outputs.tf              # Global outputs
│   └── README.md              # Infrastructure documentation
├── aws/                       # AWS-specific configurations
│   ├── sam/                   # SAM templates
│   │   ├── template.yaml      # Main SAM template
│   │   └── samconfig.toml     # SAM configuration
│   └── cloudformation/        # CloudFormation templates
│       └── stack.yaml         # CloudFormation stack
├── docker/                    # Docker configurations
│   ├── Dockerfile.dev         # Development Dockerfile
│   ├── Dockerfile.prod        # Production Dockerfile
│   └── docker-compose.yml     # Multi-service composition
├── k8s/                       # Kubernetes manifests (future)
│   ├── namespace.yaml         # Namespace definition
│   ├── deployment.yaml        # Application deployment
│   ├── service.yaml           # Service definitions
│   └── ingress.yaml           # Ingress configuration
├── monitoring/                # Monitoring configurations
│   ├── cloudwatch/            # CloudWatch configurations
│   │   ├── dashboards/        # Custom dashboards
│   │   └── alarms.yaml        # Alarm definitions
│   ├── newrelic/              # New Relic configurations
│   │   └── config.yaml        # New Relic configuration
│   └── grafana/               # Grafana configurations (future)
│       └── dashboard.json     # Custom dashboard
├── scripts/                   # Infrastructure scripts
│   ├── deploy.sh              # Deployment script
│   ├── setup.sh               # Initial setup script
│   ├── destroy.sh             # Infrastructure teardown
│   └── validate.sh            # Configuration validation
└── README.md                  # Infrastructure documentation
```

## Scripts Directory Structure

```text
scripts/
├── build/                     # Build Scripts
│   ├── build-backend.sh       # Go backend build script
│   ├── build-frontend.sh      # Next.js frontend build script
│   ├── build-all.sh           # Full project build script
│   └── optimize.sh            # Build optimization script
├── deploy/                    # Deployment Scripts
│   ├── deploy-dev.sh          # Development deployment
│   ├── deploy-staging.sh      # Staging deployment
│   ├── deploy-prod.sh         # Production deployment
│   └── rollback.sh            # Deployment rollback script
├── test/                      # Testing Scripts
│   ├── test-backend.sh        # Backend test execution
│   ├── test-frontend.sh       # Frontend test execution
│   ├── test-integration.sh    # Integration test execution
│   ├── test-e2e.sh            # End-to-end test execution
│   └── test-all.sh            # Full test suite execution
├── dev/                       # Development Scripts
│   ├── dev-setup.sh           # Development environment setup
│   ├── dev-start.sh           # Start development servers
│   ├── dev-stop.sh            # Stop development servers
│   └── dev-reset.sh           # Reset development environment
├── ci/                        # CI/CD Scripts
│   ├── ci-setup.sh            # CI environment setup
│   ├── ci-test.sh             # CI test execution
│   ├── ci-build.sh            # CI build process
│   └── ci-deploy.sh           # CI deployment process
├── db/                        # Database Scripts (future)
│   ├── migrate.sh             # Database migration
│   ├── seed.sh                # Database seeding
│   └── backup.sh              # Database backup
├── maintenance/               # Maintenance Scripts
│   ├── cleanup.sh             # Cleanup old resources
│   ├── backup.sh              # System backup
│   ├── health-check.sh        # System health verification
│   └── log-rotate.sh          # Log rotation
└── README.md                  # Scripts documentation
```

## Documentation Directory Structure

```text
docs/
├── architecture/              # Architecture Documentation
│   ├── architecture.md        # Main architecture document
│   ├── coding-standards.md    # Development standards
│   ├── tech-stack.md          # Technology stack details
│   ├── source-tree.md         # This document
│   ├── api-design.md          # API design patterns
│   ├── security.md            # Security architecture
│   └── performance.md         # Performance considerations
├── api/                       # API Documentation
│   ├── openapi.yaml           # OpenAPI/Swagger specification
│   ├── rest-api.md            # REST API documentation
│   ├── mcp-protocol.md        # MCP protocol documentation
│   └── webhooks.md            # Webhook documentation
├── deployment/                # Deployment Documentation
│   ├── aws-setup.md           # AWS infrastructure setup
│   ├── local-development.md   # Local development guide
│   ├── ci-cd.md               # CI/CD pipeline documentation
│   └── monitoring.md          # Monitoring and alerting setup
├── user/                      # User Documentation
│   ├── getting-started.md     # User onboarding guide
│   ├── operations-guide.md    # Document operations guide
│   ├── troubleshooting.md     # Common issues and solutions
│   └── faq.md                 # Frequently asked questions
├── developer/                 # Developer Documentation
│   ├── contributing.md        # Contribution guidelines
│   ├── testing.md             # Testing strategy and guidelines
│   ├── debugging.md           # Debugging guide
│   └── code-review.md         # Code review guidelines
├── stories/                   # User Stories and Requirements
│   ├── 1.2.project-repository-setup.md # Current story
│   └── [future-stories].md   # Additional user stories
├── qa/                        # Quality Assurance Documentation
│   ├── assessments/           # QA assessment reports
│   └── gates/                 # Quality gate decisions
├── compliance/                # Compliance Documentation
│   ├── security-audit.md      # Security audit reports
│   ├── privacy-policy.md      # Privacy policy
│   └── terms-of-service.md    # Terms of service
└── README.md                  # Documentation overview
```

## Testing Directory Structure

```text
tests/
├── e2e/                       # End-to-End Tests
│   ├── auth/                  # Authentication flow tests
│   │   ├── login.spec.ts      # Login functionality tests
│   │   └── multi-account.spec.ts # Multi-account tests
│   ├── documents/             # Document operation tests
│   │   ├── replace-all.spec.ts # Replace all operation tests
│   │   ├── append.spec.ts     # Append operation tests
│   │   ├── prepend.spec.ts    # Prepend operation tests
│   │   └── insert-ops.spec.ts # Insert operations tests
│   ├── integration/           # Service integration tests
│   │   ├── api-frontend.spec.ts # API-Frontend integration
│   │   └── mcp-client.spec.ts # MCP-Client integration
│   ├── fixtures/              # E2E test fixtures
│   │   ├── test-documents.json # Test document data
│   │   └── test-users.json    # Test user data
│   ├── helpers/               # E2E test helpers
│   │   ├── auth-helper.ts     # Authentication helpers
│   │   └── page-helper.ts     # Page interaction helpers
│   └── playwright.config.ts   # Playwright configuration
├── load/                      # Load and Performance Tests
│   ├── api-load.js            # API load tests (Artillery)
│   ├── websocket-load.js      # WebSocket load tests
│   └── stress-test.js         # Stress testing scenarios
├── security/                  # Security Tests
│   ├── auth-security.test.js  # Authentication security tests
│   ├── api-security.test.js   # API security tests
│   └── input-validation.test.js # Input validation tests
├── compatibility/             # Cross-platform Tests
│   ├── browser-compat.js      # Browser compatibility tests
│   ├── mobile-responsive.js   # Mobile responsiveness tests
│   └── accessibility.js       # Accessibility compliance tests
├── data/                      # Test Data
│   ├── documents/             # Document test data
│   │   ├── small-doc.json     # Small document samples
│   │   ├── large-doc.json     # Large document samples
│   │   └── markdown-samples/  # Markdown test samples
│   ├── auth/                  # Authentication test data
│   │   └── oauth-tokens.json  # OAuth token samples
│   └── responses/             # API response samples
│       ├── success-responses.json # Successful responses
│       └── error-responses.json # Error responses
├── utils/                     # Test Utilities
│   ├── test-server.js         # Test server setup
│   ├── mock-google-api.js     # Google API mocking
│   ├── test-database.js       # Test database utilities
│   └── assertion-helpers.js   # Custom assertion helpers
├── config/                    # Test Configuration
│   ├── jest.config.js         # Jest configuration
│   ├── playwright.config.js   # Playwright configuration
│   └── test-env.js            # Test environment setup
└── README.md                  # Testing documentation
```

## Configuration Files Organization

### Root Level Configuration Files

```text
# Git Configuration
.gitignore                     # Git ignore patterns
.gitattributes                # Git attributes configuration

# Environment Configuration
.env.example                  # Environment variables template
.env.local                    # Local environment variables (git-ignored)

# Development Tools
.editorconfig                 # Editor configuration
.prettierrc                   # Prettier configuration
.prettierignore              # Prettier ignore patterns

# Docker Configuration
docker-compose.yml           # Local development stack
.dockerignore               # Docker ignore patterns

# Build Configuration
Makefile                    # Build automation
package.json                # Node.js workspace configuration (if using workspaces)

# Documentation
README.md                   # Project overview
LICENSE                     # MIT license
CONTRIBUTING.md             # Contribution guidelines
CHANGELOG.md               # Version change log
CLAUDE.md                  # Claude Code configuration
```

### Language-Specific Configuration

**Go Configuration (backend/):**
```text
go.mod                      # Go module definition
go.sum                      # Dependency checksums
.golangci.yml              # golangci-lint configuration
```

**Node.js Configuration (frontend/):**
```text
package.json               # Dependencies and scripts
package-lock.json          # Dependency lock file
.eslintrc.json            # ESLint configuration
.eslintignore             # ESLint ignore patterns
tsconfig.json             # TypeScript configuration
next.config.js            # Next.js configuration
tailwind.config.ts        # Tailwind CSS configuration
postcss.config.js         # PostCSS configuration
jest.config.js            # Jest testing configuration
```

## File Naming Conventions

### General Conventions

**Directories:** 
- Use lowercase with hyphens for multi-word names: `api-gateway/`, `user-management/`
- Use descriptive, clear names: `authentication/` not `auth/` (except for very common abbreviations)

**Files:**
- Go files: snake_case with `.go` extension: `document_processor.go`, `auth_middleware.go`
- TypeScript/JavaScript: PascalCase for components, camelCase for utilities: `DocumentEditor.tsx`, `apiClient.ts`
- Configuration files: lowercase with dots/hyphens: `.env.example`, `docker-compose.yml`
- Documentation: lowercase with hyphens: `getting-started.md`, `api-documentation.md`

### Specific File Type Conventions

**Go Files:**
```text
# Main entry points
main.go                    # Application entry point
server.go                  # Server setup and configuration
config.go                  # Configuration handling

# Business logic
document_processor.go      # Document processing logic
auth_service.go           # Authentication service
operation_handler.go      # Operation handlers

# Tests
document_processor_test.go # Unit tests
integration_test.go       # Integration tests
```

**TypeScript/React Files:**
```text
# Components (PascalCase)
DocumentEditor.tsx         # React component
AuthProvider.tsx          # Context provider
ApiClient.ts              # Class-based utilities

# Hooks (camelCase with 'use' prefix)
useAuth.ts                # Authentication hook
useDocuments.ts           # Document management hook

# Utilities (camelCase)
apiClient.ts              # API client utilities
formatters.ts             # Data formatting utilities
validators.ts             # Input validation

# Types (camelCase)
api.ts                    # API type definitions
document.ts               # Document types
```

**Configuration Files:**
```text
# Environment
.env.local                # Local environment
.env.example             # Template for environment variables

# Build and deployment
Dockerfile               # Docker container definition
docker-compose.yml       # Multi-service composition
template.yaml           # SAM/CloudFormation template

# Code quality
.eslintrc.json          # ESLint configuration
.prettierrc             # Prettier configuration
.golangci.yml           # Go linter configuration
```

## Development Workflow Integration

### Local Development Structure

When running the application locally, the following structure supports the development workflow:

```text
# Terminal 1: Backend API Server
cd backend && make dev-api

# Terminal 2: Backend MCP Server  
cd backend && make dev-mcp

# Terminal 3: Frontend Development Server
cd frontend && npm run dev

# Terminal 4: Infrastructure Services
docker-compose up redis
```

### Build Artifact Organization

```text
# Build outputs (git-ignored)
frontend/.next/            # Next.js build output
backend/bin/              # Go compiled binaries
dist/                     # Distribution packages
coverage/                 # Test coverage reports
logs/                     # Development logs
```

### IDE Integration

**VS Code Configuration (.vscode/):**
```text
.vscode/
├── settings.json         # Workspace settings
├── tasks.json           # Build tasks
├── launch.json          # Debug configurations
└── extensions.json      # Recommended extensions
```

**Recommended Extensions:**
- Go extension for Go development
- TypeScript and React extensions
- Tailwind CSS IntelliSense
- ESLint and Prettier extensions
- Docker extension
- AWS Toolkit

This source tree organization provides a clear, scalable structure that supports the development requirements outlined in the PRD while maintaining separation of concerns and enabling efficient AI-assisted development.