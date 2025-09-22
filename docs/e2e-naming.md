# E2E Test Naming Conventions

## Overview

This document establishes comprehensive naming conventions for End-to-End (E2E) tests to ensure consistency, traceability, and maintainability across the test suite.

## Core Naming Principles

### 1. FR-ID Integration

All test names **MUST** include their Functional Requirement ID as defined in `docs/requirements.md`:

```typescript
// ✅ Correct
test('FR-00001 should access version endpoint directly', async ({ page }) => {
  // Test implementation
});

// ❌ Incorrect
test('should access version endpoint directly', async ({ page }) => {
  // Test implementation
});
```

### 2. Descriptive Action-Based Names

Test names should clearly describe the action being tested using consistent patterns:

```typescript
// Pattern: FR-XXXXX should [action] [expected result]
test('FR-00007 should fetch and display backend version', async ({ page }) => {
  // Implementation
});

test('FR-00014 should load homepage within 2 seconds', async ({ page }) => {
  // Implementation
});
```

## Naming Patterns

### Backend API Tests

**File**: `tests/e2e/backend-api.spec.ts`

```typescript
// Endpoint access tests
test('FR-00001 should access version endpoint directly', async ({ request }) => {});
test('FR-00002 should access health endpoint directly', async ({ request }) => {});

// Error handling tests
test('FR-00003 should handle 404 for non-existent endpoints', async ({ request }) => {});
test('FR-00004 should handle method not allowed for POST on version endpoint', async ({ request }) => {});

// Infrastructure tests
test('FR-00005 should verify CORS headers for frontend requests', async ({ request }) => {});
```

### Frontend UI Tests

**File**: `tests/e2e/homepage.spec.ts`

```typescript
// Page loading tests
test('FR-00006 should load homepage and display title', async ({ page }) => {});
test('FR-00007 should fetch and display backend version', async ({ page }) => {});

// UI component tests
test('FR-00008 should display welcome card with features', async ({ page }) => {});
test('FR-00009 should have responsive design', async ({ page }) => {});

// Performance tests
test('FR-00014 should load homepage within 2 seconds', async ({ page }) => {});
```

### Infrastructure Tests

**File**: `tests/e2e/docker-compose.spec.ts` (planned)

```typescript
// Service orchestration tests
test('FR-00010 should start all services correctly with docker-compose', async () => {});
test('FR-00011 should execute Playwright test framework successfully', async () => {});
```

### Railway Infrastructure Tests

**File**: `tests/e2e/railway.spec.ts` (planned)

```typescript
// Authentication tests
test('FR-00015 should link Railway project and authenticate CLI', async () => {});

// Deployment tests
test('FR-00016 should deploy successfully via GitHub Actions workflow', async () => {});

// Environment tests
test('FR-00017 should create services for each Railway environment', async () => {});
test('FR-00019 should configure environment variables per service', async () => {});

// Domain tests
test('FR-00018 should map and verify custom domains', async () => {});
```

### Documentation Tests

**File**: `tests/e2e/documentation.spec.ts` (planned)

```typescript
// Setup validation tests
test('FR-00012 should validate README setup instructions work step-by-step', async () => {});
test('FR-00013 should validate full development environment operational from fresh setup', async () => {});
```

## File Organization

### Directory Structure

```
tests/e2e/
├── backend-api.spec.ts       # Backend API tests (FR-00001 to FR-00005)
├── homepage.spec.ts          # Frontend UI tests (FR-00006 to FR-00009, FR-00014)
├── docker-compose.spec.ts    # Infrastructure tests (FR-00010, FR-00011)
├── railway.spec.ts           # Railway infrastructure tests (FR-00015 to FR-00019)
├── documentation.spec.ts     # Documentation tests (FR-00012, FR-00013)
└── helpers/                  # Shared utilities and fixtures
    ├── api-helpers.ts
    ├── page-helpers.ts
    └── railway-helpers.ts
```

### File Naming Rules

1. **Descriptive Domains**: File names should reflect the domain being tested
2. **Kebab Case**: Use kebab-case for file names (e.g., `backend-api.spec.ts`)
3. **Spec Suffix**: All test files must end with `.spec.ts`
4. **Logical Grouping**: Group related FR-IDs in the same file when they test the same system component

## Test Structure Standards

### Test Block Organization

```typescript
import { test, expect } from '@playwright/test';

// Group related tests with descriptive names
test.describe('Backend API Endpoints', () => {
  test('FR-00001 should access version endpoint directly', async ({ request }) => {
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

test.describe('Frontend Homepage', () => {
  test('FR-00006 should load homepage and display title', async ({ page }) => {
    // Implementation
  });
});
```

### Comment Standards

```typescript
test('FR-00007 should fetch and display backend version', async ({ page }) => {
  // FR-00007: Homepage displays backend version at bottom
  // Source: Story 1.1 (1.1-E2E-003)
  // Requirements: Version must be fetched from /version endpoint and displayed

  await page.goto('/');

  // Wait for version to load and be displayed
  const versionElement = await page.waitForSelector('[data-testid="backend-version"]');
  const versionText = await versionElement.textContent();

  expect(versionText).toContain('1.0.0');
});
```

## Validation Rules

### Required Elements

Every E2E test **MUST** include:

1. **FR-ID in test name**: `FR-XXXXX` prefix
2. **Descriptive action**: Clear description of what is being tested
3. **Source comment**: Reference to requirements document
4. **Data test IDs**: Use `data-testid` for reliable element selection
5. **Proper assertions**: Clear, specific expectations

### Forbidden Patterns

```typescript
// ❌ Missing FR-ID
test('should load homepage', async ({ page }) => {});

// ❌ Vague description
test('FR-00001 should work', async ({ page }) => {});

// ❌ Implementation details in name
test('FR-00001 should call axios.get with /version', async ({ page }) => {});

// ❌ Multiple responsibilities
test('FR-00001 should access version and health endpoints', async ({ page }) => {});
```

## Performance Test Naming

```typescript
// Performance tests should specify timing requirements
test('FR-00014 should load homepage within 2 seconds', async ({ page }) => {
  const startTime = Date.now();

  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const loadTime = Date.now() - startTime;
  expect(loadTime).toBeLessThan(2000);
});
```

## Security Test Naming

```typescript
// Security tests should specify the security concern
test('FR-00005 should verify CORS headers for frontend requests', async ({ request }) => {
  const response = await request.get('/version', {
    headers: { 'Origin': 'https://example.com' }
  });

  expect(response.headers()['access-control-allow-origin']).toBeDefined();
});
```

## Integration with Requirements

### Traceability Matrix

Each test must maintain bidirectional traceability:

1. **Forward Traceability**: FR-ID → Test Implementation
2. **Backward Traceability**: Test → Requirements Document

### Requirements Comments

```typescript
test('FR-00010 should start all services correctly with docker-compose', async () => {
  // REQUIREMENT: FR-00010 - Docker-compose starts all services correctly
  // SOURCE: Story 1.1 (1.1-E2E-004)
  // STATUS: Not Implemented
  // PRIORITY: High

  // Test implementation when created
});
```

## Maintenance Guidelines

### Updating Test Names

When requirements change:

1. Update `docs/requirements.md` first
2. Update test names to match new FR-IDs
3. Update source comments
4. Verify traceability matrix

### Adding New Tests

1. Assign next available FR-ID in `docs/requirements.md`
2. Follow naming conventions exactly
3. Add to appropriate test file
4. Update requirements mapping

### Removing Tests

1. Mark FR-ID as deprecated in requirements
2. Add deprecation comment to test
3. Schedule removal after validation

---

**Last Updated**: Created with comprehensive E2E naming conventions
**Next Review**: When new test files are added or naming patterns change
**Maintainer**: Test Architecture Team
