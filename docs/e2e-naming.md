# E2E Test Naming Conventions

## Overview

This document establishes naming conventions for automated End-to-End (E2E) tests to ensure consistency, traceability, and maintainability. This document covers only automatable tests that can run without manual intervention.

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

test('FR-00010 should load homepage within 2 seconds', async ({ page }) => {
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
```

### Performance Tests

**File**: `tests/e2e/performance.spec.ts` (planned)

```typescript
// Performance tests
test('FR-00010 should load homepage within 2 seconds', async ({ page }) => {});
```

## Manual Test Documentation (Not Automated)

For requirements that cannot be automated (infrastructure setup, authentication flows, documentation validation), manual test procedures are documented separately with FR-M prefixes.

### Manual Test Categories

1. **Infrastructure Setup (FR-M001, FR-M002)**: Docker orchestration, framework setup
2. **Documentation Validation (FR-M003, FR-M004)**: README verification, environment setup
3. **Railway Infrastructure (FR-M005 through FR-M009)**: Authentication, deployment, domain configuration

### Manual Test Documentation Format

```markdown
# Manual Test: FR-M001 - Docker-compose Service Orchestration

## Prerequisites
- Docker and docker-compose installed
- Project repository cloned

## Test Steps
1. Navigate to project root
2. Execute `docker-compose up`
3. Verify all services start successfully
4. Validate service health endpoints

## Expected Results
- All containers start without errors
- Services respond to health checks
- Frontend can communicate with backend

## Verification
- [ ] Backend container running on port 8080
- [ ] Frontend container running on port 3000
- [ ] Services communicate successfully
```

## File Organization

### Automated Test Directory Structure

```
tests/e2e/
├── backend-api.spec.ts       # Backend API tests (FR-00001 to FR-00005)
├── homepage.spec.ts          # Frontend UI tests (FR-00006 to FR-00009)
├── performance.spec.ts       # Performance tests (FR-00010) - planned
└── helpers/                  # Shared utilities and fixtures
    ├── api-helpers.ts
    └── page-helpers.ts
```

### Manual Test Documentation Structure

```
docs/manual-tests/
├── infrastructure/
│   ├── FR-M001-docker-compose.md
│   └── FR-M002-playwright-setup.md
├── documentation/
│   ├── FR-M003-readme-validation.md
│   └── FR-M004-fresh-setup.md
└── railway/
    ├── FR-M005-auth.md
    ├── FR-M006-github-actions.md
    ├── FR-M007-environments.md
    ├── FR-M008-domains.md
    └── FR-M009-variables.md
```

### Automated Test File Naming Rules

1. **Descriptive Domains**: File names should reflect the domain being tested
2. **Kebab Case**: Use kebab-case for file names (e.g., `backend-api.spec.ts`)
3. **Spec Suffix**: All automated test files must end with `.spec.ts`
4. **Logical Grouping**: Group related FR-IDs in the same file when they test the same system component

### Manual Test File Naming Rules

1. **FR-M Prefix**: All manual test files use FR-M prefix format (e.g., `FR-M001-docker-compose.md`)
2. **Descriptive Names**: Clear indication of what is being tested manually
3. **Markdown Format**: All manual test documentation in `.md` format
4. **Category Grouping**: Organize by category in subdirectories

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
test('FR-00010 should load homepage within 2 seconds', async ({ page }) => {
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
test('FR-00010 should load homepage within 2 seconds', async ({ page }) => {
  // REQUIREMENT: FR-00010 - Homepage loads within 2 seconds
  // SOURCE: Story 1.1 (1.1-E2E-009) - renumbered from FR-00014
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

### Adding New Automated Tests

1. Verify the requirement is automatable (no manual setup required)
2. Assign next available FR-ID in `docs/requirements.md`
3. Follow naming conventions exactly
4. Add to appropriate automated test file
5. Update requirements mapping

### Adding New Manual Tests

1. Assign FR-M prefix ID for manual requirements
2. Create manual test documentation in `docs/manual-tests/`
3. Follow manual test documentation format
4. Organize by category in appropriate subdirectory

### Removing Tests

1. Mark FR-ID as deprecated in requirements
2. Add deprecation comment to test
3. Schedule removal after validation

---

**Last Updated**: Updated to focus on automatable tests only
**Next Review**: When new automatable test scenarios are identified
**Maintainer**: Test Architecture Team
