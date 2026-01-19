# Epic 4: OAuth Authentication

**Goal:** Implement complete Google OAuth 2.0 authentication flow with multi-account support

**Context:** Story 4.1 provides MVP authentication via Service Account for shared documents. Stories 4.2-4.6 extend to full user OAuth for personal document access.

## User Stories

### Story 4.1: Shared Document Editing
**As a** Claude User
**I want** Claude to edit Google Docs that I've shared with provided service account
**So that** I can edit my documents conversationally by just sharing them, without needing OAuth logins or copy-pasting content back and forth

**Acceptance Criteria:**
- When I share a Google Doc with the service account, Claude can modify it
- When I ask Claude to update shared document content, the changes appear in Google Docs
- When I try to edit an unshared document, I receive a clear message explaining how to share it
- The service account email is visible so I know what to share with

### Story 4.2: OAuth Configuration
**As a** Developer/Maintainer
**I want** to configure Google OAuth application
**So that** users can authenticate with Google

**Acceptance Criteria:**
- Google Cloud project created
- OAuth 2.0 credentials generated
- Redirect URIs configured
- Scopes defined for Google Docs access
- Credentials stored as Railway environment variables (or external secret manager if required)
- Environment variables configured

### Story 4.3: OAuth Flow Implementation
**As a** Claude User
**I want** to authenticate with my Google account
**So that** the service can access my documents

**Acceptance Criteria:**
- OAuth initiation endpoint created
- Google consent screen displayed with permissions
- Authorization code exchange implemented
- Tokens received and validated
- Error handling for auth failures
- Success redirect to settings page

### Story 4.4: Token Management
**As a** Developer/Maintainer
**I want** to securely manage OAuth tokens
**So that** user sessions remain valid

**Acceptance Criteria:**
- Tokens encrypted before storage
- Redis cache integration for tokens
- Token refresh logic implemented
- Expiry handling automated
- Token retrieval by user ID
- Secure token deletion capability

### Story 4.5: Multi-Account Support
**As a** Claude User
**I want** to connect multiple Google accounts
**So that** I can edit documents from different accounts

**Acceptance Criteria:**
- Support multiple tokens per user
- Account selection mechanism
- Account switching capability
- Display connected accounts
- Remove account functionality
- Default account designation

### Story 4.6: Authentication Error Handling
**As a** Claude User
**I want** clear authentication error messages
**So that** I can resolve authentication issues

**Acceptance Criteria:**
- Structured error responses for auth failures
- Token refresh failure handling
- Clear error codes and messages
- Guidance for resolution steps
- Logging of auth errors
- Immediate error return without retry

## Technical Notes

- Story 4.1 uses Service Account for MVP - no user OAuth flow required
- Service Account requires documents to be explicitly shared by user
- Migration to OAuth (4.2-4.6) uses same Google Docs API calls - only auth layer changes
