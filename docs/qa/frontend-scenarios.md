# Frontend QA Test Scenarios

## Overview

This document contains comprehensive test scenarios for the frontend Next.js application, covering UI functionality, user interactions, and frontend-backend integration from the user's perspective.

**Last Updated:** 2025-09-09  
**Coverage:** Next.js 14 Frontend Service  
**Framework:** Playwright E2E Tests (UI Automation)

## Homepage Functionality

### Scenario: Homepage Loads Successfully
**Test ID:** FE-001  
**Priority:** Critical  
**Type:** UI Functional Test

**Description:** Verify that the homepage loads with all required elements and content.

**Preconditions:**
- Backend service running on http://localhost:8080
- Frontend service running on http://localhost:3000

**Test Steps:**
1. Navigate to http://localhost:3000
2. Wait for page to fully load

**Expected Results:**
- Page loads without errors
- Page title contains "MCP Google Docs Editor"
- Main heading displays "MCP Google Docs Editor"
- Page layout renders correctly with proper styling
- No JavaScript errors in console

**Playwright Test Location:** `tests/e2e/frontend/homepage.spec.ts`

### Scenario: Backend Version Display
**Test ID:** FE-002  
**Priority:** Critical  
**Type:** Integration Test (Frontend-Backend)

**Description:** Verify that the frontend successfully fetches and displays the backend version.

**Preconditions:**
- Backend service running and returning version "1.0.0"
- Frontend service running
- Network connectivity between services

**Test Steps:**
1. Navigate to homepage
2. Wait for version information to load
3. Locate version display element at bottom of page

**Expected Results:**
- Version display shows "Backend version: 1.0.0"
- Version is visible and properly styled
- No loading errors or timeout issues
- Version updates if backend version changes

**Playwright Test Location:** `tests/e2e/frontend/homepage.spec.ts`

## Component Functionality

### Scenario: Version Display Component States
**Test ID:** FE-003  
**Priority:** High  
**Type:** Component Unit Test

**Description:** Verify that the VersionDisplay component handles all states correctly.

**Test Framework:** Jest with React Testing Library

**Test Cases:**
1. **Loading State:**
   - Shows "Loading version..." text initially
   - Loading indicator displays while fetching

2. **Success State:**
   - Displays "Backend version: 1.0.0" after successful fetch
   - Proper styling applied to version text

3. **Error State:**
   - Shows "Error loading version" on fetch failure
   - Error styling applied (red text)
   - Component doesn't crash on network errors

4. **API URL Configuration:**
   - Uses NEXT_PUBLIC_API_URL environment variable
   - Falls back to http://localhost:8080 if not set
   - Makes request to correct endpoint (/version)

**Jest Test Location:** `services/frontend/app/components/VersionDisplay.test.tsx`

## Responsive Design

### Scenario: Mobile Responsiveness
**Test ID:** FE-004  
**Priority:** Medium  
**Type:** Responsive Design Test

**Description:** Verify that the homepage displays correctly on mobile devices.

**Test Steps:**
1. Set viewport to mobile dimensions (375x667)
2. Navigate to homepage
3. Verify layout adaptation

**Expected Results:**
- Content adapts to mobile viewport
- Text remains readable
- Version display remains visible
- No horizontal scrolling required

## Error Handling

### Scenario: Backend Service Unavailable
**Test ID:** FE-005  
**Priority:** High  
**Type:** Error Handling Test

**Description:** Verify graceful handling when backend service is unavailable.

**Preconditions:**
- Frontend service running
- Backend service stopped or unreachable

**Test Steps:**
1. Navigate to homepage
2. Wait for version fetch attempt

**Expected Results:**
- Page loads successfully despite backend unavailability
- Version display shows error message
- No application crashes or white screens
- User can still interact with other page elements

## Performance

### Scenario: Page Load Performance
**Test ID:** FE-006  
**Priority:** Medium  
**Type:** Performance Test

**Description:** Verify that the homepage loads within acceptable time limits.

**Metrics:**
- First Contentful Paint (FCP) < 2s
- Largest Contentful Paint (LCP) < 2.5s
- Time to Interactive (TTI) < 3s
- Backend version fetch < 1s

**Tools:** Lighthouse CI, Playwright performance metrics

## Future Frontend Scenarios

### Authentication Flow (Post-MVP)
**Test ID:** FE-100  
**Priority:** Critical (Future)  
**Type:** Authentication Test

**Description:** User login/logout flow with Google OAuth.

**Test Cases:**
- Login button redirect to Google OAuth
- Successful authentication callback handling
- User profile display
- Logout functionality
- Session management

### Document Management UI (Post-MVP)
**Test ID:** FE-101  
**Priority:** Critical (Future)  
**Type:** Document Management Test

**Description:** Document list, creation, and editing interfaces.

**Test Cases:**
- Document list displays user's Google Docs
- Create new document flow
- Edit existing document interface
- Document operation status indicators
- Real-time operation feedback

## Test Execution Commands

### Run All Frontend Tests
```bash
# E2E Frontend Tests (Playwright)
cd tests/e2e
npx playwright test --project=frontend

# Unit Tests (Jest)
cd services/frontend
npm test

# Component Tests with Coverage
cd services/frontend  
npm test -- --coverage
```

### Debug Frontend Tests
```bash
# Playwright Debug Mode
cd tests/e2e
npx playwright test --project=frontend --debug

# Jest Watch Mode
cd services/frontend
npm test -- --watch
```

## Test Data Requirements

### Frontend Test Data
- **Valid API Responses:** Mock backend responses for version endpoint
- **Error Scenarios:** Network timeouts, 500 errors, invalid JSON
- **Environment Configs:** Different API URL configurations
- **Browser States:** Clean browser state, cached data scenarios

## Automation Status

| Scenario | Test Type | Status | Framework | Location |
|----------|-----------|--------|-----------|-----------|
| FE-001 | UI Functional | ✅ Automated | Playwright | `tests/e2e/frontend/homepage.spec.ts` |
| FE-002 | Integration | ✅ Automated | Playwright | `tests/e2e/frontend/homepage.spec.ts` |
| FE-003 | Unit Test | ✅ Automated | Jest | `services/frontend/app/components/VersionDisplay.test.tsx` |
| FE-004 | Responsive | 📋 Manual | - | - |
| FE-005 | Error Handling | ✅ Automated | Jest | `services/frontend/app/components/VersionDisplay.test.tsx` |
| FE-006 | Performance | 📋 Manual | Lighthouse | - |

## Coverage Metrics

### Current Coverage (Story 1.1)
- **UI Components:** 100% (1/1 components tested)
- **API Integration:** 100% (version fetch tested)
- **Error Handling:** 100% (network errors covered)
- **Browser Compatibility:** Chrome only (Playwright default)

### Target Coverage (Future)
- **Authentication:** 0% (not implemented)
- **Document Operations:** 0% (not implemented)
- **Multi-browser:** Planned (Chrome, Firefox, Safari)
- **Mobile Testing:** Manual only

This document serves as the comprehensive frontend testing specification and should be updated as new features and test scenarios are added to the application.