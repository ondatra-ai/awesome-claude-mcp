# ADR-001: Unified End-to-End Testing with Playwright

**Status:** Accepted  
**Date:** 2025-09-09  
**Authors:** Claude Code (James - Full Stack Developer)  
**Reviewers:** TBD  

## Context

The MCP Google Docs Editor project requires comprehensive testing coverage across multiple layers:
1. **Frontend UI/UX Testing**: User interface interactions, authentication flows, document management
2. **Backend API Integration Testing**: REST API endpoints, authentication, CORS validation, error handling
3. **Cross-Service Integration**: Frontend-backend communication and data flow validation

Initially, the project used separate testing approaches:
- Go unit tests for backend logic (`testify`)
- Jest for frontend unit tests
- Playwright E2E tests only for frontend UI automation

This approach had limitations:
- **Test Fragmentation**: Three different testing frameworks and configurations
- **Limited API Coverage**: No comprehensive backend API integration testing
- **Maintenance Overhead**: Multiple test environments and reporting systems
- **CI/CD Complexity**: Separate test execution pipelines
- **Coverage Gaps**: Missing backend API validation in isolation

## Decision

We will adopt **Playwright as the unified end-to-end testing framework** for both frontend and backend testing, while maintaining unit tests for service-specific logic.

### New Testing Architecture

```text
tests/
├── e2e/                       # Unified Playwright Tests
│   ├── frontend/              # Frontend UI/UX Tests
│   │   ├── homepage.spec.ts   # Homepage functionality
│   │   ├── auth/              # Authentication flow tests
│   │   └── documents/         # Document management UI
│   ├── backend/               # Backend API Integration Tests  
│   │   ├── api.spec.ts        # Core API functionality
│   │   ├── auth/              # Authentication API tests
│   │   └── documents/         # Document operations API tests
│   ├── playwright.config.ts   # Unified configuration
│   └── package.json           # Single dependency set
```

### Key Benefits

1. **Single Framework**: One testing tool for both frontend and backend
2. **API Testing Capabilities**: Playwright's `request` fixture for HTTP API testing
3. **Unified Configuration**: Single `playwright.config.ts` with project separation
4. **Consistent Reporting**: Single test report covering all integration scenarios
5. **Service Orchestration**: Built-in web server management for test dependencies
6. **TypeScript Support**: Strongly-typed test development across all layers

## Implementation Details

### Playwright Configuration

```typescript
export default defineConfig({
  projects: [
    {
      name: 'frontend',
      testDir: './frontend',
      use: { 
        ...devices['Desktop Chrome'],
        baseURL: 'http://localhost:3000',
      },
    },
    {
      name: 'backend-api',
      testDir: './backend',
      use: { 
        ...devices['Desktop Chrome'],
        baseURL: 'http://localhost:8080',
      },
    },
  ],
  webServer: [
    {
      command: 'cd ../../services/backend && go run ./cmd/main.go',
      port: 8080,
    },
    {
      command: 'cd ../../services/frontend && npm run dev',
      port: 3000,
    },
  ],
});
```

### Backend API Testing Pattern

```typescript
test('API Integration Tests', async ({ request }) => {
  // Test API endpoints directly
  const response = await request.get('/version');
  expect(response.status()).toBe(200);
  
  const data = await response.json();
  expect(data.version).toBe('1.0.0');
  
  // Test CORS headers
  expect(response.headers()['access-control-allow-origin']).toBe('*');
  
  // Test error handling
  const notFoundResponse = await request.get('/nonexistent');
  expect(notFoundResponse.status()).toBe(404);
});
```

### Test Execution Command

```bash
cd tests/e2e
npx playwright test  # Runs both frontend UI and backend API tests
```

## Rationale

### Why Playwright for API Testing?

1. **HTTP Request Fixture**: Native support for API testing via `request` fixture
2. **Service Management**: Automatic startup/shutdown of backend services
3. **Unified Reporting**: Single test report for UI and API results
4. **TypeScript Integration**: Strongly-typed request/response handling
5. **Debugging Tools**: Browser devtools for both UI and API debugging
6. **CI/CD Integration**: Single test execution pipeline

### Compared to Alternatives

| Aspect | Current (Playwright Unified) | Alternative (Separate Tools) |
|--------|----------------------------|----------------------------|
| **Test Frameworks** | 1 (Playwright) | 3 (Playwright + Jest + Go tests) |
| **Configuration Files** | 1 (playwright.config.ts) | 3 (separate configs) |
| **CI/CD Steps** | 1 test execution | Multiple test execution steps |
| **Reporting** | Unified HTML report | Separate reports per framework |
| **API Coverage** | Full HTTP testing | Limited or no API integration tests |
| **Maintenance** | Single tool updates | Multiple tool version management |

## Consequences

### Positive

1. **Simplified Testing Stack**: Single framework reduces complexity
2. **Comprehensive Coverage**: Both UI and API integration testing
3. **Better Developer Experience**: Single command runs all E2E tests
4. **Consistent Test Patterns**: Similar syntax for UI and API tests
5. **Unified CI/CD**: Single test execution and reporting pipeline
6. **Service Integration Testing**: Validates actual frontend-backend communication

### Negative

1. **Learning Curve**: Team needs to learn Playwright API testing capabilities
2. **Tool Dependencies**: All E2E testing depends on Playwright (single point of failure)
3. **Unit Test Separation**: Need clear boundaries between unit tests (Go/Jest) and integration tests (Playwright)

### Mitigation Strategies

1. **Clear Test Boundaries**:
   - **Unit Tests** (Go `testify`, Frontend `Jest`): Business logic, component behavior
   - **Integration Tests** (Playwright): API endpoints, UI workflows, cross-service communication

2. **Documentation**: Create comprehensive testing guide explaining when to use each approach

3. **Training**: Provide examples and patterns for both UI and API testing with Playwright

## References

- [Playwright API Testing Documentation](https://playwright.dev/docs/api-testing)
- [Story 1.1 Implementation](../../../docs/stories/1.1-minimal-frontend-backend-integration.md)
- [Source Tree Architecture](../source-tree.md)
- [Tech Stack Documentation](../tech-stack.md)

## Revision History

- **2025-09-09**: Initial decision - Unified E2E testing with Playwright for both frontend and backend integration testing