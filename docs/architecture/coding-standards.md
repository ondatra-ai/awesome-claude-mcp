# MCP Google Docs Editor - Coding Standards

## Introduction

This document establishes comprehensive coding standards for the MCP Google Docs Editor project. These standards ensure code quality, maintainability, and consistency across the full-stack application comprising Next.js frontend and Go backend services.

**Scope:** All source code including Go backend services, TypeScript/React frontend, configuration files, tests, and documentation.

## General Principles

### Code Quality Fundamentals

1. **Readability First**: Code should be self-documenting through clear naming and structure
2. **Consistency**: Follow established patterns within each language ecosystem
3. **Simplicity**: Prefer simple, explicit solutions over clever abstractions
4. **Testability**: Write code that's easy to test with clear dependencies
5. **Error Handling**: Explicit error handling with informative messages
6. **Performance**: Balance readability with performance, optimize based on actual metrics

### Project-Specific Guidelines

- **Fail-Fast Philosophy**: Return errors immediately without retry attempts (per PRD)
- **Structured Logging**: All operations must be logged with structured JSON format
- **MCP Compliance**: Strictly follow MCP protocol standards without custom extensions
- **OAuth Security**: Secure token handling with encryption at rest
- **Configuration**: All configuration via environment variables, no hardcoded values

## Go Backend Standards

### File Organization

```text
backend/
├── cmd/                    # Application entry points
│   ├── api/main.go        # REST API server
│   └── mcp/main.go        # MCP WebSocket server
├── internal/              # Internal packages (not importable externally)
│   ├── api/               # HTTP handlers and middleware
│   ├── mcp/               # MCP protocol implementation
│   ├── auth/              # OAuth and authentication
│   ├── operations/        # Document operations
│   ├── docs/              # Google Docs integration
│   ├── cache/             # Redis caching
│   └── config/            # Configuration management
├── pkg/                   # Public packages (importable)
│   ├── errors/            # Custom error types
│   └── utils/             # Utility functions
└── deployments/           # Docker and deployment configs
```

### Naming Conventions

**Packages:**
- Use lowercase, single words when possible: `auth`, `cache`, `docs`
- Avoid generic names: `util`, `common`, `base`
- Package names should match directory names

**Functions and Methods:**
```go
// Public functions: PascalCase
func ProcessDocument(docID string) error

// Private functions: camelCase
func validateToken(token string) bool

// Interface names: append 'er' or descriptive suffix
type DocumentProcessor interface
type TokenValidator interface
```

**Variables and Constants:**
```go
// Variables: camelCase
var documentCache map[string]*Document
var oauthConfig *oauth2.Config

// Constants: SCREAMING_SNAKE_CASE for exported, camelCase for private
const MAX_DOCUMENT_SIZE = 10 * 1024 * 1024
const defaultTimeout = 30 * time.Second

// Error variables: Err prefix
var ErrTokenExpired = errors.New("oauth token expired")
var ErrDocumentNotFound = errors.New("document not found")
```

### Code Structure

**Function Design:**
```go
// Good: Clear function signature with context
func ReplaceAllContent(ctx context.Context, docID, content string) (*OperationResult, error) {
    if docID == "" {
        return nil, ErrInvalidDocumentID
    }

    // Implementation
    return result, nil
}

// Avoid: Generic parameters, unclear return types
func DoOperation(args ...interface{}) (interface{}, error)
```

**Error Handling:**
```go
// Use custom error types with context
type DocumentError struct {
    DocumentID string
    Operation  string
    Cause      error
}

func (e *DocumentError) Error() string {
    return fmt.Sprintf("document %s: operation %s failed: %v",
        e.DocumentID, e.Operation, e.Cause)
}

// Wrap errors with context
if err := validateDocument(doc); err != nil {
    return nil, &DocumentError{
        DocumentID: doc.ID,
        Operation:  "validate",
        Cause:      err,
    }
}
```

**Structured Logging:**
```go
// Use structured logging with consistent fields
logger.Info("document operation started",
    "document_id", docID,
    "operation", "replace_all",
    "user_id", userID,
    "content_length", len(content),
)

// Log errors with full context
logger.Error("document operation failed",
    "document_id", docID,
    "operation", "replace_all",
    "error", err.Error(),
    "duration_ms", time.Since(start).Milliseconds(),
)
```

### Testing Standards

**Test File Organization:**
- Test files alongside implementation: `document_test.go`
- Integration tests: `tests/integration/`
- E2E tests: `tests/e2e/` (organized by service: frontend/, backend/, mcp-service/)

**Test Naming:**
```go
func TestReplaceAllContent_ValidDocument_Success(t *testing.T) {
    // Arrange, Act, Assert pattern
}

func TestReplaceAllContent_EmptyDocID_ReturnsError(t *testing.T) {
    // Test error conditions
}
```

**Test Structure (AAA Pattern):**
```go
func TestDocumentProcessor_ReplaceAll(t *testing.T) {
    // Arrange
    processor := NewDocumentProcessor(mockClient)
    docID := "test-doc-123"
    content := "# New Content\n\nTest document"

    // Act
    result, err := processor.ReplaceAll(context.Background(), docID, content)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, docID, result.DocumentID)
}
```

**Mock Usage:**
```go
// Generate mocks using gomock
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_interfaces.go

// Use mocks in tests
func TestAuthHandler_WithMockedTokenValidator(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockValidator := mocks.NewMockTokenValidator(ctrl)
    mockValidator.EXPECT().Validate("valid-token").Return(nil)

    // Test implementation
}
```

## TypeScript Frontend Standards

### File Organization

```text
frontend/
├── app/                    # Next.js App Router
│   ├── (auth)/            # Route groups
│   ├── api/               # API routes
│   ├── globals.css        # Global styles
│   ├── layout.tsx         # Root layout
│   └── page.tsx           # Home page
├── components/            # React components
│   ├── ui/                # Basic UI components (shadcn/ui)
│   ├── forms/             # Form components
│   └── layout/            # Layout components
├── lib/                   # Utilities and configurations
│   ├── api.ts             # API client
│   ├── auth.ts            # Authentication utilities
│   └── utils.ts           # General utilities
├── hooks/                 # Custom React hooks
├── types/                 # TypeScript type definitions
└── public/                # Static assets
```

### Naming Conventions

**Components:**
```tsx
// PascalCase for components and interfaces
interface DocumentEditorProps {
    documentId: string;
    onSave: (content: string) => void;
}

export function DocumentEditor({ documentId, onSave }: DocumentEditorProps) {
    // Component implementation
}
```

**Files and Directories:**
- Components: `PascalCase.tsx` (e.g., `DocumentEditor.tsx`)
- Utilities: `camelCase.ts` (e.g., `apiClient.ts`)
- Pages: Next.js conventions (`page.tsx`, `layout.tsx`)
- Hooks: `use` prefix (e.g., `useDocuments.ts`)

**Variables and Functions:**
```tsx
// camelCase for variables, functions
const documentCache = new Map<string, Document>();
const isAuthenticated = useAuth();

// Async functions: clear naming
const fetchDocument = async (id: string): Promise<Document> => {
    // Implementation
};

// Event handlers: handle prefix
const handleDocumentSave = (content: string) => {
    // Implementation
};
```

### Component Standards

**Component Structure:**
```tsx
'use client'; // Mark client components explicitly

import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import type { Document } from '@/types/document';

interface DocumentEditorProps {
    document: Document;
    onSave: (content: string) => Promise<void>;
    readOnly?: boolean;
}

export function DocumentEditor({
    document,
    onSave,
    readOnly = false
}: DocumentEditorProps) {
    const [content, setContent] = useState(document.content);
    const [isSaving, setIsSaving] = useState(false);

    const handleSave = async () => {
        setIsSaving(true);
        try {
            await onSave(content);
        } catch (error) {
            console.error('Failed to save document:', error);
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <div className="document-editor">
            {/* Component JSX */}
        </div>
    );
}
```

**Custom Hooks:**
```tsx
// Custom hook naming and structure
export function useDocumentOperations(documentId: string) {
    const [document, setDocument] = useState<Document | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const replaceAll = async (content: string) => {
        try {
            setIsLoading(true);
            const result = await apiClient.replaceAll(documentId, content);
            setDocument(result.document);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Unknown error');
        } finally {
            setIsLoading(false);
        }
    };

    return { document, isLoading, error, replaceAll };
}
```

### TypeScript Standards

**Type Definitions:**
```tsx
// Explicit interfaces over implicit types
interface ApiResponse<T> {
    success: boolean;
    data: T;
    error?: string;
}

interface DocumentOperation {
    type: 'replace_all' | 'append' | 'prepend' | 'replace_match' | 'insert_before' | 'insert_after';
    documentId: string;
    content: string;
    anchor?: string; // For insertion operations
}

// Use utility types when appropriate
type DocumentUpdate = Partial<Pick<Document, 'title' | 'content'>>;
```

**API Client Standards:**
```tsx
// Typed API client with proper error handling
class ApiClient {
    private baseURL: string;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    async replaceAll(documentId: string, content: string): Promise<ApiResponse<OperationResult>> {
        const response = await fetch(`${this.baseURL}/api/documents/${documentId}/replace-all`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ content }),
        });

        if (!response.ok) {
            throw new Error(`API request failed: ${response.statusText}`);
        }

        return response.json();
    }
}
```

## Configuration and Environment

### Environment Variables

**Go Backend:**
```env
# Server Configuration
PORT=8080
MCP_PORT=8081
ENVIRONMENT=development

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000/auth/callback

# Railway Deployment (managed via CLI/UI)
# RAILWAY_PROJECT_ID=801ad5e0-95bf-4ce6-977e-6f2fa37529fd
# Railway injects environment variables per service (configured in dashboard)

# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_TTL_HOURS=24

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

**Next.js Frontend:**
```env
# Next.js Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_MCP_URL=ws://localhost:8081

# Authentication
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key
```

### Configuration Management

**Go Configuration:**
```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Google   GoogleConfig   `mapstructure:"google"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
    Port     int    `mapstructure:"port" default:"8080"`
    MCPPort  int    `mapstructure:"mcp_port" default:"8081"`
    Environment string `mapstructure:"environment" default:"development"`
}

// Load configuration using Viper
func LoadConfig() (*Config, error) {
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

## Testing Standards

### Test Coverage Requirements

- **Minimum Coverage**: 80% overall, 85% for business logic
- **Critical Paths**: 100% coverage for authentication, document operations
- **Error Scenarios**: All error paths must be tested
- **Integration Points**: External API interactions must be tested with mocks

### Test Categories

**Unit Tests (60% of test pyramid):**
- Individual functions and methods
- Business logic validation
- Error condition handling
- Mock external dependencies

**Integration Tests (30% of test pyramid):**
- Service-to-service communication
- Database interactions (Redis)
- Authentication flows
- MCP protocol compliance

**End-to-End Tests (10% of test pyramid):**
- Complete user workflows across all three services
- Service-to-service integration testing
- Document operation workflows (frontend → backend → MCP service)
- Authentication workflows (OAuth flows through all services)
- **Framework**: Playwright for cross-browser testing
- **Organization**: Service-specific test directories (frontend/, backend/, mcp-service/)

### Test Data Management

```go
// Test fixtures in organized structure
tests/
├── fixtures/
│   ├── documents/
│   │   ├── simple.json
│   │   └── complex.json
│   └── auth/
│       └── tokens.json
└── helpers/
    ├── auth_helper.go
    └── document_helper.go

// Helper functions for test data
func LoadTestDocument(t *testing.T, filename string) *Document {
    data, err := os.ReadFile(filepath.Join("fixtures", "documents", filename))
    require.NoError(t, err)

    var doc Document
    err = json.Unmarshal(data, &doc)
    require.NoError(t, err)

    return &doc
}
```

## Documentation Standards

### Code Documentation

**Go Documentation:**
```go
// Package documentation
// Package auth provides OAuth 2.0 authentication for Google services.
//
// This package implements the complete OAuth flow including token acquisition,
// refresh, and secure storage. It supports multiple Google accounts per user
// and handles token expiration gracefully.
package auth

// Function documentation
// ValidateToken checks if the provided OAuth token is valid and not expired.
//
// It performs the following validations:
//   - Token format and structure
//   - Expiration time
//   - Token signature (if applicable)
//
// Returns nil if the token is valid, otherwise returns a descriptive error.
func ValidateToken(token *oauth2.Token) error {
    // Implementation
}
```

**TypeScript Documentation:**
```tsx
/**
 * Custom hook for managing document operations.
 *
 * Provides methods for all six document operations (replace_all, append, prepend,
 * replace_match, insert_before, insert_after) with proper error handling and
 * loading states.
 *
 * @param documentId - The Google Docs document ID
 * @returns Object containing document state and operation methods
 */
export function useDocumentOperations(documentId: string) {
    // Implementation
}
```

### API Documentation

**OpenAPI/Swagger for REST APIs:**
```yaml
paths:
  /api/documents/{documentId}/replace-all:
    post:
      summary: Replace entire document content
      parameters:
        - name: documentId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ReplaceAllRequest'
      responses:
        '200':
          description: Operation successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/OperationResult'
        '400':
          description: Invalid request
        '404':
          description: Document not found
```

## Security Standards

### Authentication and Authorization

**OAuth Token Handling:**
```go
// Encrypt tokens before storage
func StoreToken(userID string, token *oauth2.Token) error {
    encrypted, err := encryptToken(token)
    if err != nil {
        return err
    }

    return redisClient.Set(
        fmt.Sprintf("token:%s", userID),
        encrypted,
        tokenTTL,
    ).Err()
}

// Always validate tokens before use
func GetValidToken(userID string) (*oauth2.Token, error) {
    token, err := retrieveToken(userID)
    if err != nil {
        return nil, err
    }

    if token.Expiry.Before(time.Now()) {
        return nil, ErrTokenExpired
    }

    return token, nil
}
```

### Input Validation

**Go Input Validation:**
```go
// Validate all inputs explicitly
func ValidateReplaceAllRequest(req *ReplaceAllRequest) error {
    if req.DocumentID == "" {
        return ErrInvalidDocumentID
    }

    if len(req.Content) > MAX_DOCUMENT_SIZE {
        return ErrDocumentTooLarge
    }

    if !isValidMarkdown(req.Content) {
        return ErrInvalidMarkdown
    }

    return nil
}
```

**TypeScript Input Validation:**
```tsx
// Use Zod for runtime validation
import { z } from 'zod';

const ReplaceAllSchema = z.object({
    documentId: z.string().min(1, 'Document ID required'),
    content: z.string().max(10485760, 'Content too large'), // 10MB limit
});

export function validateReplaceAllRequest(data: unknown) {
    return ReplaceAllSchema.parse(data);
}
```

### Error Handling Security

```go
// Don't expose internal errors to clients
func HandleError(err error) *APIError {
    if errors.Is(err, ErrTokenExpired) {
        return &APIError{
            Code:    "TOKEN_EXPIRED",
            Message: "Please re-authenticate",
            Status:  401,
        }
    }

    // Log internal error but return generic message
    logger.Error("internal error", "error", err)
    return &APIError{
        Code:    "INTERNAL_ERROR",
        Message: "An unexpected error occurred",
        Status:  500,
    }
}
```

## Performance Standards

### Response Time Requirements

- **Document Operations**: < 2 seconds (95th percentile)
- **Authentication**: < 1 second
- **Homepage Load**: < 2 seconds
- **API Health Check**: < 100ms

### Resource Usage Guidelines

**Go Services:**
- Memory usage: < 128MB per Lambda function
- Cold start: < 1 second
- Concurrent requests: Support 10+ per function instance

**Frontend Performance:**
- First Contentful Paint: < 1.5 seconds
- Largest Contentful Paint: < 2.5 seconds
- Cumulative Layout Shift: < 0.1
- First Input Delay: < 100ms

### Caching Strategy

```go
// Cache frequently accessed data
type DocumentCache struct {
    cache  *sync.Map
    ttl    time.Duration
}

func (c *DocumentCache) Get(docID string) (*Document, bool) {
    if val, ok := c.cache.Load(docID); ok {
        if entry := val.(*CacheEntry); time.Now().Before(entry.Expiry) {
            return entry.Document, true
        }
        c.cache.Delete(docID) // Clean expired entry
    }
    return nil, false
}
```

## Deployment and Operations

### Build Standards

**Go Build Pipeline:**
```makefile
# Consistent build commands
.PHONY: build test lint deploy

build:
    CGO_ENABLED=0 GOOS=linux go build -o bin/api ./cmd/api
    CGO_ENABLED=0 GOOS=linux go build -o bin/mcp ./cmd/mcp

test:
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

lint:
    golangci-lint run --config .golangci.yml  # Runs 70+ linters including cyclop, funlen, gosec

deploy: build
    sam deploy --guided
```

**Next.js Build Pipeline:**
```json
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "type-check": "tsc --noEmit",
    "test": "jest",
    "test:e2e": "playwright test"
  }
}
```

### Monitoring and Observability

**Structured Logging Format:**
```json
{
  "timestamp": "2025-01-07T10:30:00Z",
  "level": "info",
  "service": "mcp-server",
  "operation": "replace_all",
  "document_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
  "user_id": "user-123",
  "duration_ms": 1250,
  "success": true,
  "content_length": 1024
}
```

**Health Check Standards:**
```go
// Comprehensive health checks
type HealthCheck struct {
    Service   string    `json:"service"`
    Version   string    `json:"version"`
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Checks    []Check   `json:"checks"`
}

type Check struct {
    Name   string `json:"name"`
    Status string `json:"status"`
    Error  string `json:"error,omitempty"`
}

func PerformHealthCheck() *HealthCheck {
    return &HealthCheck{
        Service:   "mcp-google-docs-editor",
        Version:   version.Get(),
        Status:    "healthy",
        Timestamp: time.Now(),
        Checks: []Check{
            checkRedisConnection(),
            checkGoogleAPIAccess(),
            checkTokenValidation(),
        },
    }
}
```

## Compliance and Quality Gates

### Pre-commit Checks

**Automated via lint-staged (.lintstagedrc.json):**
- **TypeScript files (`*.{ts,tsx}`)**:
  - `eslint --fix` - Auto-fix linting issues
  - `prettier --write` - Auto-format code
- **Config/Docs (`*.{json,md,yml,yaml}`)**:
  - `prettier --write` - Auto-format configuration files

**Required checks before any commit:**
1. All tests pass
2. Linting passes (golangci-lint with 70+ linters, eslint with strict rules)
3. Type checking passes (Go build, TypeScript compiler)
4. Security scanning passes (gosec, eslint security rules)
5. Code coverage meets minimums (80%+ overall)
6. Pre-commit hooks pass (lint-staged)
7. Documentation updated if needed

### Code Review Checklist

**Mandatory Review Items:**
- [ ] Follows naming conventions
- [ ] Has appropriate test coverage
- [ ] Handles errors properly
- [ ] Logs operations with structured format
- [ ] Validates all inputs
- [ ] Documentation updated
- [ ] No hardcoded configuration
- [ ] Security implications considered
- [ ] Performance implications considered

### Quality Metrics

**Automated Quality Gates (Go):**
- **Cyclomatic complexity**: ≤15 per function (configured in .golangci.yml cyclop linter)
- **Function length**: ≤80 lines / ≤40 statements (configured in .golangci.yml funlen linter)
- **No critical security vulnerabilities**: gosec linter enabled
- **All linting rules pass**: 70+ linters in golangci-lint configuration

**Automated Quality Gates (TypeScript):**
- **Complexity limit**: ≤10 per function
- **Function length**: ≤50 lines (≤100 for container/handler files, ≤250 for tests)
- **Line length**: ≤80 characters
- **File length**: ≤300 lines
- **Max parameters**: ≤5 per function
- **Nesting depth**: ≤4 levels
- **Interface naming**: Must use 'I' prefix

**Common Quality Gates:**
- Code coverage: ≥80% overall, ≥85% business logic
- No critical security vulnerabilities
- All linting rules pass
- Build time: <5 minutes for full pipeline
- Test execution time: <2 minutes for full suite

This coding standards document ensures consistent, maintainable, and secure code across the MCP Google Docs Editor project while supporting the development requirements outlined in the PRD.
