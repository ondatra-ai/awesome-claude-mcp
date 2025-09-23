# E2E Test Requirements (Automatable Only)

This document maps End-to-End test requirements using a Requirement→Scenario structure with flat numbering for complete traceability.

**FR-00001**: Story 1.1 (1.1-E2E-001) Backend /version endpoint returns 1.0.0

**Scenarios**:
```gherkin
# UT-00001-01: handler should return correct version (services/backend/cmd/main_test.go)
Given version handler function exists
When handler is called with valid HTTP request
Then JSON response {"version": "1.0.0"} is returned

# UT-00001-02: handler should reject invalid methods (services/backend/cmd/main_test.go)
Given version endpoint receives invalid HTTP method
When non-GET request is made to /version
Then appropriate HTTP error status is returned

# IT-00001-03: server responds with correct status (tests/integration/backend.test.ts)
Given Fiber server is running with version endpoint
When HTTP GET request is made to /version
Then server responds with correct JSON and 200 status

# EE-00001-04: service returns version with headers (tests/e2e/backend-api.spec.ts)
Given full backend service is deployed and running
When external HTTP client calls /version endpoint
Then service responds with "1.0.0" and proper HTTP headers
```

**FR-00002**: Story 1.1 (1.1-E2E-006) Backend /health endpoint returns healthy status

**Scenarios**:
```gherkin
# UT-00002-01: health check handler returns correct status (services/backend/cmd/main_test.go)
Given health check handler functions exist
When handlers are called
Then correct health status responses are returned

# EE-00002-02: health endpoint accessible via HTTP (tests/e2e/backend-api.spec.ts)
Given services are running with health endpoints
When health check URLs are accessed via HTTP
Then both services respond with healthy status
```

**FR-00003**: Backend handles 404 for non-existent endpoints

**Scenarios**:
```gherkin
# EE-00003-01: returns 404 for non-existent endpoints (tests/e2e/backend-api.spec.ts)
Given backend service is running
When GET request is sent to non-existent endpoint
Then service returns 404 status code
And response indicates endpoint was not found
```

**FR-00004**: Backend rejects invalid HTTP methods

**Scenarios**:
```gherkin
# EE-00004-01: returns 405 for POST on version endpoint (tests/e2e/backend-api.spec.ts)
Given backend service is running
When POST request is sent to /version endpoint
Then service returns 405 Method Not Allowed status code
And response indicates method is not allowed
```

**FR-00005**: Backend provides CORS headers for frontend

**Scenarios**:
```gherkin
# EE-00005-01: includes CORS headers for frontend requests (tests/e2e/backend-api.spec.ts)
Given backend service is running
When request with Origin header from frontend domain is sent
Then response includes appropriate CORS headers
And Access-Control-Allow-Origin header allows frontend domain
```

**FR-00006**: Story 1.1 (1.1-E2E-002) Frontend single-page application loads successfully

**Scenarios**:
```gherkin
# UT-00006-01: homepage component renders without errors (services/frontend/__tests__/components/HomePage.test.tsx)
Given homepage React component is defined
When component is rendered in test environment
Then component renders without errors and displays expected content

# UT-00006-02: component handles props correctly (services/frontend/__tests__/components/HomePage.test.tsx)
Given component receives props for configuration
When props are passed to component
Then component handles props correctly and renders accordingly

# IT-00006-03: Next.js routes serve pages correctly (tests/integration/frontend.test.ts)
Given Next.js App Router is configured
When application routes are accessed
Then pages are served correctly by Next.js framework

# EE-00006-04: SPA loads in browser (tests/e2e/homepage.spec.ts)
Given frontend application is built and served
When browser navigates to application URL
Then single-page application loads and renders in browser
```

**FR-00007**: Story 1.1 (1.1-E2E-003) Homepage displays backend version at bottom

**Scenarios**:
```gherkin
# UT-00007-01: API client constructs correct request (services/frontend/__tests__/lib/api.test.ts)
Given API client utility exists
When API client constructs version request
Then correct HTTP request structure is created

# UT-00007-02: version component displays "1.0.0" (services/frontend/__tests__/components/VersionDisplay.test.tsx)
Given version display component receives version data
When component is rendered with version "1.0.0"
Then version text is displayed at bottom of page

# UT-00007-03: handles API errors gracefully (services/frontend/__tests__/components/VersionDisplay.test.tsx)
Given API call fails with network error
When version fetch encounters failure
Then error state is handled gracefully with user-friendly message

# UT-00007-04: shows loading state during fetch (services/frontend/__tests__/components/VersionDisplay.test.tsx)
Given API call is in progress
When version is being fetched from backend
Then loading state is displayed to user

# IT-00007-05: frontend fetches from backend successfully (tests/integration/fullstack.test.ts)
Given frontend and backend services are running
When frontend makes API call to backend version endpoint
Then version "1.0.0" is successfully retrieved and displayed

# EE-00007-06: version appears at bottom of page (tests/e2e/homepage.spec.ts)
Given full stack application is running
When user loads homepage in browser
Then backend version "1.0.0" appears at bottom of page
```

**FR-00008**: Homepage displays welcome card with features

**Scenarios**:
```gherkin
# EE-00008-01: welcome card is visible (tests/e2e/homepage.spec.ts)
Given frontend application is loaded
When user views homepage
Then welcome card is visible
And welcome card displays key features
```

**FR-00009**: Homepage has responsive design

**Scenarios**:
```gherkin
# EE-00009-01: adapts to different viewport sizes (tests/e2e/homepage.spec.ts)
Given frontend application is loaded
When browser is resized to different viewport sizes
Then homepage layout adapts responsively
And all elements remain accessible and properly positioned
```

**FR-00010**: Story 1.1 (1.1-E2E-009) Homepage loads within 2 seconds

**Scenarios**:
```gherkin
# EE-00010-01: page loads under 2000ms (tests/e2e/performance.spec.ts)
Given frontend application is deployed
When user navigates to homepage
Then page fully loads in under 2000ms
And all critical resources are loaded within time limit
```