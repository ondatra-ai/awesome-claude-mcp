# Test Naming Conventions

## Overview

This document establishes naming conventions for all automated tests (Unit, Integration, and End-to-End) to ensure consistency, traceability, and maintainability. This document covers only automatable tests that can run without manual intervention.

## Automatable Requirements Criteria

**Requirements are considered automatable if they meet ALL of these criteria:**
- No manual setup or authentication steps required
- Deterministic, measurable outcomes
- No external dependency configuration needed
- No human interaction during execution
- Services must be pre-configured and running
- Executable in CI environment without manual intervention

**Excluded from automation:**
- Requirements with manual configuration steps
- File/directory structure validation requiring filesystem checks
- Infrastructure setup requirements
- Authentication configuration requirements
- External service configuration requirements

## Scenario ID System

### ID Format

All tests must follow the scenario ID system for complete traceability:

- **UT-XXXXX-YY**: Unit Test scenarios
- **IT-XXXXX-YY**: Integration Test scenarios
- **EE-XXXXX-YY**: End-to-End Test scenarios

Where:
- **XXXXX**: FR number (00001, 00002, etc.) linking to functional requirements
- **YY**: Globally sequential number within each FR (01, 02, 03, 04...)

### Test Name Format

**With mapped scenario:**
```typescript
test('EE-00001-04: should access version endpoint directly', async ({ request }) => {
  // Test implementation
});
```

**Without mapped scenario (orphan):**
```typescript
test('ORPHAN: should validate request headers', async ({ request }) => {
  // Test implementation
});
```

## Core Naming Principles

### 1. Scenario ID Integration

All test names **MUST** include their Scenario ID as defined in `docs/requirements.yml`:

```typescript
// ✅ Correct - Mapped scenario
test('EE-00001-04: should access version endpoint directly', async ({ page }) => {
  // Test implementation
});

// ✅ Correct - Orphaned test
test('ORPHAN: should validate error response format', async ({ page }) => {
  // Test implementation
});

// ❌ Incorrect - No scenario identification
test('should access version endpoint directly', async ({ page }) => {
  // Test implementation
});
```

### 2. Descriptive Action-Based Names

Test names should clearly describe the action being tested:

```typescript
// Pattern: [SCENARIO-ID]: should [action] [expected result]
test('EE-00007-06: should fetch and display backend version', async ({ page }) => {
  // Implementation
});

test('EE-00010-01: should load homepage within 2 seconds', async ({ page }) => {
  // Implementation
});
```

## Test Type Naming Patterns

### Unit Tests (UT-XXXXX-YY)

**Backend Unit Tests** - `services/backend/**/*_test.go`

```go
// Version handler tests
func TestVersionHandler_ValidRequest_ReturnsCorrectVersion(t *testing.T) {
    // UT-00001-01: handler should return correct version
}

func TestVersionHandler_InvalidMethod_ReturnsError(t *testing.T) {
    // UT-00001-02: handler should reject invalid methods
}

// Health check tests
func TestHealthHandler_ValidRequest_ReturnsHealthyStatus(t *testing.T) {
    // UT-00002-01: health check handler returns correct status
}
```

**Frontend Unit Tests** - `services/frontend/__tests__/**/*.test.ts`

```typescript
// Component tests
test('UT-00006-01: homepage component should render without errors', () => {
  // Implementation
});

test('UT-00006-02: component should handle props correctly', () => {
  // Implementation
});

// API client tests
test('UT-00007-01: API client should construct correct request', () => {
  // Implementation
});
```

### Integration Tests (IT-XXXXX-YY)

**Backend Integration** - `tests/integration/backend.test.ts`

```typescript
test('IT-00001-03: server should respond with correct status', async () => {
  // Implementation
});
```

**Frontend Integration** - `tests/integration/frontend.test.ts`

```typescript
test('IT-00006-03: Next.js routes should serve pages correctly', async () => {
  // Implementation
});
```

**Full-Stack Integration** - `tests/integration/fullstack.test.ts`

```typescript
test('IT-00007-05: frontend should fetch from backend successfully', async () => {
  // Implementation
});
```

### End-to-End Tests (EE-XXXXX-YY)

**Backend API Tests** - `tests/e2e/backend-api.spec.ts`

```typescript
test.describe('Backend API Endpoints', () => {
  test('EE-00001-04: should access version endpoint directly', async ({ request }) => {
    // Implementation
  });

  test('EE-00002-02: should access health endpoint directly', async ({ request }) => {
    // Implementation
  });

  test('EE-00003-01: should handle 404 for non-existent endpoints', async ({ request }) => {
    // Implementation
  });

  test('EE-00004-01: should handle method not allowed for POST on version endpoint', async ({ request }) => {
    // Implementation
  });

  test('EE-00005-01: should verify CORS headers for frontend requests', async ({ request }) => {
    // Implementation
  });
});
```

**Frontend UI Tests** - `tests/e2e/homepage.spec.ts`

```typescript
test.describe('Frontend Homepage', () => {
  test('EE-00006-04: should load homepage and display title', async ({ page }) => {
    // Implementation
  });

  test('EE-00007-06: should fetch and display backend version', async ({ page }) => {
    // Implementation
  });

  test('EE-00008-01: should display welcome card with features', async ({ page }) => {
    // Implementation
  });

  test('EE-00009-01: should have responsive design', async ({ page }) => {
    // Implementation
  });
});
```

**Performance Tests** - `tests/e2e/performance.spec.ts`

```typescript
test.describe('Performance Requirements', () => {
  test('EE-00010-01: should load homepage within 2 seconds', async ({ page }) => {
    const startTime = Date.now();

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    const loadTime = Date.now() - startTime;
    expect(loadTime).toBeLessThan(2000);
  });
});
```

## File Organization

### Directory Structure

```
tests/
├── e2e/                          # End-to-End tests
│   ├── backend-api.spec.ts       # Backend API tests (EE-00001-04 to EE-00005-01)
│   ├── homepage.spec.ts          # Frontend UI tests (EE-00006-04 to EE-00009-01)
│   ├── performance.spec.ts       # Performance tests (EE-00010-01)
│   └── helpers/                  # Shared utilities and fixtures
│       ├── api-helpers.ts
│       └── page-helpers.ts
├── integration/                  # Integration tests
│   ├── backend.test.ts           # Backend integration (IT-00001-03)
│   ├── frontend.test.ts          # Frontend integration (IT-00006-03)
│   └── fullstack.test.ts         # Full-stack integration (IT-00007-05)
└── unit/                         # Unit test helpers (tests are co-located with source)
    └── helpers/
        ├── mock-helpers.ts
        └── test-utils.ts

services/
├── backend/
│   └── cmd/
│       └── main_test.go          # Backend unit tests (UT-00001-01, UT-00001-02, UT-00002-01)
└── frontend/
    └── __tests__/
        ├── lib/
        │   └── api.test.ts       # API client tests (UT-00007-01)
        └── components/
            ├── HomePage.test.tsx         # Component tests (UT-00006-01, UT-00006-02)
            └── VersionDisplay.test.tsx   # Version component tests (UT-00007-02, UT-00007-03, UT-00007-04)
```

### File Naming Rules

1. **Descriptive Domains**: File names should reflect the domain being tested
2. **Kebab Case**: Use kebab-case for test file names (e.g., `backend-api.spec.ts`)
3. **Test Extension**:
   - E2E tests: `.spec.ts`
   - Unit/Integration tests: `.test.ts` or `_test.go`
4. **Logical Grouping**: Group related Scenario IDs in the same file when they test the same system component

## Test Structure Standards

### Test Block Organization

```typescript
import { test, expect } from '@playwright/test';

// Group related tests with descriptive names
test.describe('Backend API Endpoints', () => {
  test('EE-00001-04: should access version endpoint directly', async ({ request }) => {
    // EE-00001-04: service returns version with headers
    // Source: FR-00001 - Backend /version endpoint returns 1.0.0
    // Requirements: Version must be returned via HTTP GET with proper headers

    // Arrange
    const endpoint = '/version';

    // Act
    const response = await request.get(endpoint);

    // Assert
    expect(response.status()).toBe(200);
    const data = await response.json();
    expect(data.version).toBe('1.0.0');
  });
});
```

### Comment Standards

```typescript
test('EE-00007-06: should fetch and display backend version', async ({ page }) => {
  // EE-00007-06: version appears at bottom of page
  // Source: FR-00007 - Homepage displays backend version at bottom
  // Requirements: Version must be fetched from /version endpoint and displayed

  await page.goto('/');

  // Wait for version to load and be displayed
  const versionElement = await page.waitForSelector('[data-testid="backend-version"]');
  const versionText = await versionElement.textContent();

  expect(versionText).toContain('1.0.0');
});
```

## Orphan Test Handling

### Identifying Orphan Tests

Tests that don't map to documented functional requirements must be marked as ORPHAN:

```typescript
// ✅ Correct orphan marking
test('ORPHAN: should validate request headers', async ({ request }) => {
  // Test implementation for functionality not covered by FR requirements
});

test('ORPHAN: should handle malformed JSON gracefully', async ({ request }) => {
  // Test implementation for edge case not in requirements
});
```

### Orphan Test Guidelines

1. **Mark Clearly**: Always prefix with "ORPHAN:"
2. **Document Reason**: Add comment explaining why it's not mapped
3. **Regular Review**: Periodically review orphan tests for potential requirement mapping
4. **Consider Removal**: Remove orphan tests that don't provide value

```typescript
test('ORPHAN: should validate Content-Type header', async ({ request }) => {
  // Consider: Should this be added as a new functional requirement?

  // Test implementation
});
```

## Validation Rules

### Required Elements

Every test **MUST** include:

1. **Scenario ID or ORPHAN prefix**: Clear identification
2. **Descriptive action**: Clear description of what is being tested
3. **Source comment**: Reference to requirements document
4. **Data test IDs**: Use `data-testid` for reliable element selection (E2E tests)
5. **Proper assertions**: Clear, specific expectations

### Forbidden Patterns

```typescript
// ❌ Missing scenario identification
test('should load homepage', async ({ page }) => {});

// ❌ Vague description
test('EE-00001-04: should work', async ({ page }) => {});

// ❌ Implementation details in name
test('EE-00001-04: should call axios.get with /version', async ({ page }) => {});

// ❌ Multiple responsibilities
test('EE-00001-04: should access version and health endpoints', async ({ page }) => {});
```

## Integration with Requirements

### Traceability Matrix

Each test must maintain bidirectional traceability:

1. **Forward Traceability**: Scenario ID → Test Implementation
2. **Backward Traceability**: Test → Requirements Document

### Requirements Comments

```typescript
test('EE-00010-01: should load homepage within 2 seconds', async ({ page }) => {
  // SCENARIO: EE-00010-01 - page loads under 2000ms
  // REQUIREMENT: FR-00010 - Homepage loads within 2 seconds
  // SOURCE: Story 1.1 (1.1-E2E-009)
  // STATUS: Implemented
  // PRIORITY: High

  const startTime = Date.now();
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  const loadTime = Date.now() - startTime;

  expect(loadTime).toBeLessThan(2000);
});
```

## Maintenance Guidelines

### Updating Test Names

When requirements change:

1. Update `docs/requirements.yml` first
2. Update scenario IDs in test names to match new requirements
3. Update source comments
4. Verify traceability matrix

### Adding New Tests

1. Verify the requirement is automatable (no manual setup required)
2. Assign next available Scenario ID in `docs/requirements.yml`
3. Follow naming conventions exactly
4. Add to appropriate test file
5. Update requirements mapping

### Removing Tests

1. Mark Scenario ID as deprecated in requirements
2. Add deprecation comment to test
3. Schedule removal after validation

## Quality Gates

### Test Naming Validation

All tests must pass these checks:

- Contains valid Scenario ID (UT/IT/EE-XXXXX-YY) or ORPHAN prefix
- Scenario ID exists in `docs/requirements.yml`
- Description is clear and action-based
- File is in correct directory for test type
- No forbidden patterns used

### Coverage Requirements

- All functional requirements must have mapped test scenarios
- All test scenarios must trace back to functional requirements
- Orphan tests should be minimized and regularly reviewed

---

**Last Updated**: Updated to support all test types with scenario ID system
**Next Review**: When new test scenarios are identified
**Maintainer**: Test Architecture Team
