# E2E Test Requirements Mapping (Automatable Only)

## Overview

This document provides a mapping of End-to-End (E2E) test requirements that can be fully automated without manual intervention. Only requirements that can run programmatically with pre-configured services are included.

## Automation Criteria

**Requirements included in this document must:**
- Be executable without manual setup or authentication
- Have deterministic, measurable outcomes
- Not require external dependency configuration
- Not require human interaction during execution

## Automatable Requirements (10 total)

### Backend API Requirements (5/5 - 100% implemented)

**FR-00001**: Backend /version endpoint returns 1.0.0
- **Source**: Story 1.1 (1.1-E2E-001)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00001 should access version endpoint directly"
- **Automation**: HTTP GET request validation
- **Scenarios**:
  ```gherkin
  # Unit Test: 1.1-UNIT-002 (P0)
  Given version handler function exists
  When handler is called with valid HTTP request
  Then JSON response {"version": "1.0.0"} is returned

  # Unit Test: 1.1-UNIT-003 (P0)
  Given version endpoint receives invalid HTTP method
  When non-GET request is made to /version
  Then appropriate HTTP error status is returned

  # Integration Test: 1.1-INT-002 (P0)
  Given Fiber server is running with version endpoint
  When HTTP GET request is made to /version
  Then server responds with correct JSON and 200 status

  # E2E Test: 1.1-E2E-001 (P0)
  Given full backend service is deployed and running
  When external HTTP client calls /version endpoint
  Then service responds with "1.0.0" and proper HTTP headers
  ```

**FR-00002**: Backend /health endpoint returns healthy status
- **Source**: Story 1.1 (1.1-E2E-006)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00002 should access health endpoint directly"
- **Automation**: HTTP GET request validation
- **Scenarios**:
  ```gherkin
  # Unit Test: 1.1-UNIT-015 (P0)
  Given health check handler functions exist
  When handlers are called
  Then correct health status responses are returned

  # E2E Test: 1.1-E2E-006 (P0)
  Given services are running with health endpoints
  When health check URLs are accessed via HTTP
  Then both services respond with healthy status
  ```

**FR-00003**: Backend handles 404 for non-existent endpoints
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00003 should handle 404 for non-existent endpoints"
- **Automation**: HTTP GET request with error validation
- **Scenario**:
  ```gherkin
  Given the backend service is running
  When I send a GET request to a non-existent endpoint
  Then I should receive a 404 status code
  And the response should indicate the endpoint was not found
  ```

**FR-00004**: Backend rejects invalid HTTP methods
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00004 should handle method not allowed for POST on version endpoint"
- **Automation**: HTTP POST request with error validation
- **Scenario**:
  ```gherkin
  Given the backend service is running
  When I send a POST request to the /version endpoint
  Then I should receive a 405 Method Not Allowed status code
  And the response should indicate the method is not allowed
  ```

**FR-00005**: Backend provides CORS headers for frontend
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00005 should verify CORS headers for frontend requests"
- **Automation**: HTTP request with Origin header validation
- **Scenario**:
  ```gherkin
  Given the backend service is running
  When I send a request with an Origin header from the frontend domain
  Then I should receive appropriate CORS headers in the response
  And the Access-Control-Allow-Origin header should allow the frontend domain
  ```

### Infrastructure Requirements (referenced from QA assessments)

**Monorepo Structure**: Services directory with cross-communication
- **Scenarios**:
  ```gherkin
  # Unit Test: 1.1-UNIT-001 (P1)
  Given project root directory exists
  When monorepo structure creation is validated
  Then services directory and subdirectories exist with correct permissions

  # Integration Test: 1.1-INT-001 (P1)
  Given services directory structure is created
  When services attempt cross-communication
  Then services can discover and communicate with each other
  ```

**Docker Configuration**: Local development setup
- **Scenarios**:
  ```gherkin
  # Integration Test: 1.1-INT-005 (P0)
  Given backend Dockerfile and configuration exist
  When backend container is built and started
  Then container starts successfully and exposes correct ports

  # Integration Test: 1.1-INT-006 (P0)
  Given frontend Dockerfile and configuration exist
  When frontend container is built and started
  Then container starts successfully and serves application

  # Integration Test: 1.1-INT-007 (P1)
  Given both services are running in Docker network
  When frontend attempts to communicate with backend
  Then cross-service communication works via Docker networking

  # E2E Test: 1.1-E2E-004 (P1)
  Given docker-compose.yml configuration exists
  When full stack is started with docker-compose up
  Then all services start and function together correctly
  ```

**Playwright Testing Framework**: E2E testing setup
- **Scenarios**:
  ```gherkin
  # Unit Test: 1.1-UNIT-010 (P2)
  Given Playwright configuration file exists
  When configuration is validated
  Then all required settings and browser configurations are correct

  # E2E Test: 1.1-E2E-005 (P2)
  Given Playwright is installed and configured
  When basic test scenario is executed
  Then test runs successfully and reports results correctly
  ```

### Frontend UI Requirements (4/4 - 100% implemented)

**FR-00006**: Frontend single-page application loads successfully
- **Source**: Story 1.1 (1.1-E2E-002)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00006 should load homepage and display title"
- **Automation**: Playwright page load and element validation
- **Scenarios**:
  ```gherkin
  # Unit Test: 1.1-UNIT-004 (P1)
  Given homepage React component is defined
  When component is rendered in test environment
  Then component renders without errors and displays expected content

  # Unit Test: 1.1-UNIT-005 (P1)
  Given component receives props for configuration
  When props are passed to component
  Then component handles props correctly and renders accordingly

  # Integration Test: 1.1-INT-003 (P1)
  Given Next.js App Router is configured
  When application routes are accessed
  Then pages are served correctly by Next.js framework

  # E2E Test: 1.1-E2E-002 (P1)
  Given frontend application is built and served
  When browser navigates to application URL
  Then single-page application loads and renders in browser
  ```

**FR-00007**: Homepage displays backend version at bottom
- **Source**: Story 1.1 (1.1-E2E-003)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00007 should fetch and display backend version"
- **Automation**: Playwright element content validation
- **Scenarios**:
  ```gherkin
  # Unit Test: 1.1-UNIT-006 (P0)
  Given API client utility exists
  When API client constructs version request
  Then correct HTTP request structure is created

  # Unit Test: 1.1-UNIT-007 (P0)
  Given version display component receives version data
  When component is rendered with version "1.0.0"
  Then version text is displayed at bottom of page

  # Unit Test: 1.1-UNIT-008 (P1)
  Given API call fails with network error
  When version fetch encounters failure
  Then error state is handled gracefully with user-friendly message

  # Unit Test: 1.1-UNIT-009 (P1)
  Given API call is in progress
  When version is being fetched from backend
  Then loading state is displayed to user

  # Integration Test: 1.1-INT-004 (P0)
  Given frontend and backend services are running
  When frontend makes API call to backend version endpoint
  Then version "1.0.0" is successfully retrieved and displayed

  # E2E Test: 1.1-E2E-003 (P0)
  Given full stack application is running
  When user loads homepage in browser
  Then backend version "1.0.0" appears at bottom of page
  ```

**FR-00008**: Homepage displays welcome card with features
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00008 should display welcome card with features"
- **Automation**: Playwright element visibility validation
- **Scenario**:
  ```gherkin
  Given the frontend application is loaded
  When I view the homepage
  Then a welcome card should be visible
  And the welcome card should display key features
  ```

**FR-00009**: Homepage has responsive design
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00009 should have responsive design"
- **Automation**: Playwright viewport testing
- **Scenario**:
  ```gherkin
  Given the frontend application is loaded
  When I resize the browser to different viewport sizes
  Then the homepage layout should adapt responsively
  And all elements should remain accessible and properly positioned
  ```

### Performance Requirements (0/1 - 0% implemented)

**FR-00014**: Homepage loads within 2 seconds
- **Source**: Story 1.1 (1.1-E2E-009)
- **Status**: ❌ Not Implemented
- **Implementation**: None
- **Automation**: Playwright performance measurement with timing assertions
- **Scenarios**:
  ```gherkin
  # Suggested E2E Test: 1.1-E2E-009 (P1)
  Given frontend application is deployed
  When user navigates to homepage
  Then page fully loads in under 2000ms
  And all critical resources are loaded within time limit

  # Additional Performance Validation
  Given full stack application is running
  When user accesses homepage with backend version fetch
  Then complete page interaction is responsive within 2 seconds
  ```

## Coverage Analysis

### ✅ Automated Tests (9/10 - 90% implemented)

| FR ID | Test Name | File | Automation Type |
|-------|-----------|------|-----------------|
| FR-00001 | should access version endpoint directly | backend-api.spec.ts | API Testing |
| FR-00002 | should access health endpoint directly | backend-api.spec.ts | API Testing |
| FR-00003 | should handle 404 for non-existent endpoints | backend-api.spec.ts | Error Testing |
| FR-00004 | should handle method not allowed for POST on version endpoint | backend-api.spec.ts | Error Testing |
| FR-00005 | should verify CORS headers for frontend requests | backend-api.spec.ts | Header Testing |
| FR-00006 | should load homepage and display title | homepage.spec.ts | UI Testing |
| FR-00007 | should fetch and display backend version | homepage.spec.ts | Integration Testing |
| FR-00008 | should display welcome card with features | homepage.spec.ts | UI Testing |
| FR-00009 | should have responsive design | homepage.spec.ts | Responsive Testing |

### ❌ Missing Automated Tests (1/10 - 10% remaining)

| FR ID | Description | Priority | Implementation Plan |
|-------|-------------|----------|-------------------|
| FR-00014 | Performance (2s load time) | High | Add to homepage.spec.ts with timing assertions |

## Test File Organization

```
tests/e2e/
├── backend-api.spec.ts       # Backend API tests (FR-00001 to FR-00005)
├── homepage.spec.ts          # Frontend UI tests (FR-00006 to FR-00009)
└── performance.spec.ts       # Performance tests (FR-00014) - planned
```

## Implementation Priority

### 🔴 High Priority (Immediate)
1. **FR-00014**: Performance testing (2-second load requirement)
   - Add timing assertions to homepage.spec.ts or create performance.spec.ts
   - Validate load time under 2000ms

### ✅ Complete Coverage
- All API endpoints (5/5)
- All UI components (4/4)
- Error handling scenarios
- Cross-browser responsiveness

## Automation Quality Gates

- **Coverage**: 90% of automatable requirements implemented
- **Test Execution**: All tests run without manual intervention
- **Environment**: Tests assume services are pre-configured and running
- **Reliability**: Tests are idempotent and deterministic
- **Maintenance**: Clear FR-ID traceability for all tests

## Success Metrics

- **Total Automatable Requirements**: 10
- **Currently Automated**: 9 (90%)
- **Remaining**: 1 (10%)
- **Test Execution Time**: < 5 minutes for full suite
- **Test Reliability**: > 95% pass rate in CI environment

---

**Last Updated**: Refactored to include only automatable requirements
**Next Review**: When new automatable test scenarios are identified
**Maintainer**: Test Architecture Team
