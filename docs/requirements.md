# E2E Test Requirements Mapping (Automatable Only)

## Overview

This document provides a mapping of End-to-End (E2E) test requirements that can be fully automated without manual intervention. Only requirements that can run programmatically with pre-configured services are included.

## Automation Criteria

**Requirements included in this document must:**
- Be executable without manual setup or authentication
- Have deterministic, measurable outcomes
- Not require external dependency configuration
- Not require human interaction during execution

**Requirements excluded from automation:**
- Infrastructure setup (docker-compose up, service deployment)
- Authentication flows (OAuth, CLI login, API keys)
- Environment configuration (Railway setup, DNS, custom domains)
- Documentation validation (README accuracy, setup instructions)

## Automatable Requirements (10 total)

### Backend API Requirements (5/5 - 100% implemented)

**FR-00001**: Backend /version endpoint returns 1.0.0
- **Source**: Story 1.1 (1.1-E2E-001)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00001 should access version endpoint directly"
- **Automation**: HTTP GET request validation

**FR-00002**: Backend /health endpoint returns healthy status
- **Source**: Story 1.1 (1.1-E2E-006)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00002 should access health endpoint directly"
- **Automation**: HTTP GET request validation

**FR-00003**: Backend handles 404 for non-existent endpoints
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00003 should handle 404 for non-existent endpoints"
- **Automation**: HTTP GET request with error validation

**FR-00004**: Backend rejects invalid HTTP methods
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00004 should handle method not allowed for POST on version endpoint"
- **Automation**: HTTP POST request with error validation

**FR-00005**: Backend provides CORS headers for frontend
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/backend-api.spec.ts` - "FR-00005 should verify CORS headers for frontend requests"
- **Automation**: HTTP request with Origin header validation

### Frontend UI Requirements (4/4 - 100% implemented)

**FR-00006**: Frontend single-page application loads successfully
- **Source**: Story 1.1 (1.1-E2E-002)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00006 should load homepage and display title"
- **Automation**: Playwright page load and element validation

**FR-00007**: Homepage displays backend version at bottom
- **Source**: Story 1.1 (1.1-E2E-003)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00007 should fetch and display backend version"
- **Automation**: Playwright element content validation

**FR-00008**: Homepage displays welcome card with features
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00008 should display welcome card with features"
- **Automation**: Playwright element visibility validation

**FR-00009**: Homepage has responsive design
- **Source**: Not in original requirements (orphaned test)
- **Status**: ✅ Implemented
- **Implementation**: `tests/e2e/homepage.spec.ts` - "FR-00009 should have responsive design"
- **Automation**: Playwright viewport testing

### Performance Requirements (0/1 - 0% implemented)

**FR-00010**: Homepage loads within 2 seconds
- **Source**: Story 1.1 (1.1-E2E-009) - renumbered from FR-00014
- **Status**: ❌ Not Implemented
- **Implementation**: None
- **Automation**: Playwright performance measurement with timing assertions

## Manual Verification Requirements (Not Automated)

The following requirements require manual setup or human intervention and are tracked separately with FR-M prefixes:

### Infrastructure Setup (Manual)

**FR-M001**: Docker-compose starts all services correctly
- **Source**: Story 1.1 (1.1-E2E-004)
- **Manual Steps**: `docker-compose up`, service health verification
- **Why Manual**: Requires infrastructure orchestration

**FR-M002**: Playwright test framework executes successfully
- **Source**: Story 1.1 (1.1-E2E-005)
- **Manual Steps**: Framework installation, configuration validation
- **Why Manual**: Environment setup dependency

### Documentation Validation (Manual)

**FR-M003**: README setup instructions work step-by-step
- **Source**: Story 1.1 (1.1-E2E-007)
- **Manual Steps**: Fresh environment setup following instructions
- **Why Manual**: Human interpretation of documentation required

**FR-M004**: Full development environment operational from fresh setup
- **Source**: Story 1.1 (1.1-E2E-008)
- **Manual Steps**: Complete environment setup from scratch
- **Why Manual**: Complex multi-tool installation and configuration

### Railway Infrastructure (Manual)

**FR-M005**: Railway project linked and CLI authenticated
- **Source**: Story 1.2 (RLY-001, RLY-002)
- **Manual Steps**: CLI authentication, project linking
- **Why Manual**: Authentication flow requires user interaction

**FR-M006**: GitHub Actions workflow deploys successfully
- **Source**: Story 1.2 (RLY-003)
- **Manual Steps**: Workflow trigger, deployment monitoring
- **Why Manual**: CI/CD pipeline requires push or manual dispatch

**FR-M007**: Services created for each Railway environment
- **Source**: Story 1.2 (RLY-004)
- **Manual Steps**: Railway service configuration
- **Why Manual**: Platform-specific setup required

**FR-M008**: Custom domains mapped and verified
- **Source**: Story 1.2 (RLY-005, RLY-006)
- **Manual Steps**: DNS configuration, certificate validation
- **Why Manual**: External DNS dependency

**FR-M009**: Environment variables configured per service
- **Source**: Story 1.2 (RLY-007)
- **Manual Steps**: Railway dashboard configuration
- **Why Manual**: Platform-specific configuration

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
| FR-00010 | Performance (2s load time) | High | Add to homepage.spec.ts with timing assertions |

### 📋 Manual Verification Items (9 requirements)

Manual verification requirements (FR-M001 through FR-M009) are documented for completeness but not included in automation coverage metrics.

## Test File Organization

```
tests/e2e/
├── backend-api.spec.ts       # Backend API tests (FR-00001 to FR-00005)
├── homepage.spec.ts          # Frontend UI tests (FR-00006 to FR-00009)
└── performance.spec.ts       # Performance tests (FR-00010) - planned
```

## Implementation Priority

### 🔴 High Priority (Immediate)
1. **FR-00010**: Performance testing (2-second load requirement)
   - Add timing assertions to homepage.spec.ts
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
- **Manual Verification Items**: 9 (tracked separately)
- **Test Execution Time**: < 5 minutes for full suite
- **Test Reliability**: > 95% pass rate in CI environment

---

**Last Updated**: Refactored to include only automatable requirements
**Next Review**: When new automatable test scenarios are identified
**Maintainer**: Test Architecture Team
