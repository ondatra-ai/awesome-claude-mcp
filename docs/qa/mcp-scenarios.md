# MCP Service QA Test Scenarios

## Overview

This document contains comprehensive test scenarios for the MCP (Model Context Protocol) service, covering Claude Code integration, document operations, and MCP protocol compliance.

**Last Updated:** 2025-09-09  
**Coverage:** Go MCP Service (Mark3Labs MCP-Go Framework)  
**Framework:** Go testify (Unit Tests) + MCP Protocol Tests  
**Status:** Future Implementation (Post-MVP)

## MCP Protocol Compliance

### Scenario: MCP Server Initialization
**Test ID:** MCP-001  
**Priority:** Critical  
**Type:** Protocol Compliance Test

**Description:** Verify that the MCP server initializes correctly with proper protocol compliance.

**Endpoint:** MCP stdio transport

**Preconditions:**
- MCP service compiled and ready
- Mark3Labs MCP-Go library integrated
- Configuration parameters set

**Test Steps:**
1. Initialize MCP server with stdio transport
2. Verify server capabilities announcement
3. Validate protocol version compatibility

**Expected Results:**
- Server starts without errors
- Capabilities include required tool support
- Protocol version matches client expectations
- Server ready to accept tool requests

**Test Implementation:**
- **Go Unit Test:** `services/mcp-service/internal/server/mcp_test.go`
- **Protocol Test:** MCP client integration tests

### Scenario: Tool Registration and Discovery
**Test ID:** MCP-002  
**Priority:** Critical  
**Type:** Tool Discovery Test

**Description:** Verify that document operation tools are properly registered and discoverable.

**Tool Categories:**
- Document Content Operations
- Document Formatting Operations
- Document Metadata Operations

**Test Steps:**
1. Query server for available tools
2. Validate tool schemas and parameters
3. Verify tool descriptions and examples

**Expected Results:**
- All document operation tools listed
- Tool schemas are valid JSON Schema
- Parameter validation rules defined
- Tool descriptions are clear and actionable

**Registered Tools (Future Implementation):**
- `replace_all_content`
- `append_content`
- `prepend_content`
- `insert_content`
- `replace_text_range`
- `format_document`
- `get_document_info`

## Document Operation Tools

### Scenario: Replace All Content Operation
**Test ID:** MCP-003  
**Priority:** Critical  
**Type:** Document Operation Test

**Description:** Verify that the replace_all_content tool correctly replaces entire document content.

**Tool:** `replace_all_content`

**Input Schema:**
```json
{
  "document_id": "string (required)",
  "content": "string (required)",
  "format": "markdown|plain (optional, default: markdown)"
}
```

**Test Cases:**

1. **Valid Operation:**
   - Input: Valid document ID and markdown content
   - Expected: Document content completely replaced
   - Response: Success with operation details

2. **Invalid Document ID:**
   - Input: Non-existent document ID
   - Expected: Error response with clear message
   - Response: MCP error result

3. **Permission Denied:**
   - Input: Document user doesn't have edit access to
   - Expected: Permission error
   - Response: MCP error with permission details

4. **Large Content:**
   - Input: Content exceeding Google Docs limits
   - Expected: Size limit error or chunked processing
   - Response: Appropriate handling strategy

**Test Implementation:**
- **Go Unit Test:** Tool handler function tests
- **Integration Test:** Mock Google Docs API responses

### Scenario: Append Content Operation
**Test ID:** MCP-004  
**Priority:** High  
**Type:** Document Operation Test

**Description:** Verify that the append_content tool correctly adds content to document end.

**Tool:** `append_content`

**Input Schema:**
```json
{
  "document_id": "string (required)",
  "content": "string (required)",
  "separator": "string (optional, default: \\n\\n)"
}
```

**Test Cases:**
1. Simple text append
2. Markdown formatting preservation
3. Custom separator handling
4. Multiple consecutive appends

### Scenario: Insert Content Operation
**Test ID:** MCP-005  
**Priority:** High  
**Type:** Document Operation Test

**Description:** Verify that the insert_content tool correctly inserts content at specified position.

**Tool:** `insert_content`

**Input Schema:**
```json
{
  "document_id": "string (required)",
  "content": "string (required)",
  "position": "number (required)",
  "reference": "start|end|cursor (optional)"
}
```

**Test Cases:**
1. Insert at document start
2. Insert at document end
3. Insert at specific character position
4. Insert with invalid position handling

## Google Docs API Integration

### Scenario: Authentication and Authorization
**Test ID:** MCP-006  
**Priority:** Critical  
**Type:** Authentication Test

**Description:** Verify that the MCP service correctly authenticates with Google Docs API.

**Authentication Flow:**
1. Service account authentication (server-to-server)
2. OAuth token validation from frontend
3. User permission verification

**Test Cases:**
1. **Service Account Auth:**
   - Valid service account credentials
   - Token refresh handling
   - Permission scope validation

2. **User Token Validation:**
   - Valid user OAuth tokens from frontend
   - Expired token handling
   - Invalid token rejection

3. **Permission Checks:**
   - Document read permissions
   - Document edit permissions
   - Ownership verification

### Scenario: Document Access and Modification
**Test ID:** MCP-007  
**Priority:** Critical  
**Type:** Google Docs Integration Test

**Description:** Verify that document operations correctly interact with Google Docs API.

**Test Cases:**
1. **Document Retrieval:**
   - Get document metadata
   - Retrieve document content
   - Handle non-existent documents

2. **Document Updates:**
   - Apply text replacements
   - Insert formatted content
   - Batch operation handling

3. **Error Handling:**
   - API rate limiting
   - Network connectivity issues
   - Google API error responses

**Test Implementation:**
- **Integration Test:** Mock Google Docs API
- **E2E Test:** Real Google Docs interaction (with test documents)

## Content Formatting and Conversion

### Scenario: Markdown to Google Docs Conversion
**Test ID:** MCP-008  
**Priority:** High  
**Type:** Content Conversion Test

**Description:** Verify that markdown content is correctly converted to Google Docs format.

**Conversion Features:**
- Headers (H1-H6)
- Bold and italic text
- Lists (ordered and unordered)
- Links and images
- Code blocks and inline code
- Tables
- Blockquotes

**Test Cases:**
1. **Basic Formatting:**
   - `**bold**` → Bold text style
   - `*italic*` → Italic text style
   - `# Header` → Heading 1 style

2. **Complex Structures:**
   - Nested lists
   - Tables with formatting
   - Mixed formatting within paragraphs

3. **Edge Cases:**
   - Malformed markdown
   - Unsupported markdown features
   - Very large content blocks

**Test Implementation:**
- **Unit Test:** Markdown parser and converter functions
- **Integration Test:** End-to-end conversion validation

### Scenario: Google Docs to Markdown Export
**Test ID:** MCP-009  
**Priority:** Medium  
**Type:** Content Export Test

**Description:** Verify that Google Docs content can be exported to markdown format.

**Export Features:**
- Text formatting preservation
- Structure hierarchy maintenance
- Link and image handling
- Table export
- Special character handling

## Error Handling and Recovery

### Scenario: Network and Connectivity Issues
**Test ID:** MCP-010  
**Priority:** High  
**Type:** Error Handling Test

**Description:** Verify graceful handling of network and connectivity problems.

**Error Scenarios:**
1. **Google API Unavailable:**
   - Timeout handling
   - Retry logic with exponential backoff
   - Fallback error responses

2. **Rate Limiting:**
   - Rate limit detection
   - Queue management
   - Client notification

3. **Authentication Failures:**
   - Token expiry handling
   - Permission denied responses
   - User notification strategies

### Scenario: Data Validation and Sanitization
**Test ID:** MCP-011  
**Priority:** High  
**Type:** Security Test

**Description:** Verify that all input data is properly validated and sanitized.

**Validation Categories:**
1. **Document IDs:** Format validation, existence checks
2. **Content:** Size limits, character encoding, HTML sanitization
3. **Positions:** Range validation, boundary checks
4. **Parameters:** Type validation, required field checks

## Performance and Scalability

### Scenario: Large Document Handling
**Test ID:** MCP-012  
**Priority:** Medium  
**Type:** Performance Test

**Description:** Verify that the service handles large documents efficiently.

**Test Parameters:**
- Document sizes: 1MB, 5MB, 10MB+ content
- Complex formatting with many elements
- Multiple simultaneous operations

**Expected Results:**
- Operations complete within reasonable time
- Memory usage remains stable
- No timeout errors
- Proper progress feedback

### Scenario: Concurrent Operation Handling
**Test ID:** MCP-013  
**Priority:** Medium  
**Type:** Concurrency Test

**Description:** Verify that multiple concurrent operations are handled correctly.

**Test Cases:**
1. Multiple tools operating on same document
2. Different users operating on different documents
3. Rate limiting under high load
4. Resource contention management

## MCP Client Integration

### Scenario: Claude Code Integration
**Test ID:** MCP-014  
**Priority:** Critical  
**Type:** Client Integration Test

**Description:** Verify seamless integration with Claude Code client.

**Integration Points:**
- Tool discovery and schema validation
- Request/response format compliance
- Error message clarity for Claude
- Operation result formatting

**Test Cases:**
1. **Tool Discovery:**
   - Claude can list available tools
   - Tool schemas are properly formatted
   - Examples are helpful and accurate

2. **Operation Execution:**
   - Tool calls execute correctly
   - Results are properly formatted
   - Errors are user-friendly

3. **Session Management:**
   - Multiple operations in sequence
   - State management between calls
   - Session cleanup

### Scenario: Multi-Client Support
**Test ID:** MCP-015  
**Priority:** Low (Future)  
**Type:** Multi-Client Test

**Description:** Verify that the service can handle multiple MCP clients simultaneously.

**Supported Clients:**
- Claude Code
- ChatGPT (future)
- Custom MCP clients

**Test Cases:**
- Concurrent client connections
- Client-specific configurations
- Resource isolation between clients

## Test Execution Commands

### MCP Protocol Tests
```bash
cd services/mcp-service
go test ./...
go test -v ./internal/server/
go test -cover ./internal/operations/
```

### MCP Integration Tests
```bash
cd services/mcp-service
go test -tags=integration ./tests/integration/
```

### Mock Google Docs Tests
```bash
cd services/mcp-service
go test ./internal/docs/ -run TestWithMocks
```

## Test Data Requirements

### MCP Test Data
- **Test Documents:** Google Docs with known content and IDs
- **Mock API Responses:** Realistic Google Docs API responses
- **Tool Schemas:** Valid and invalid MCP tool parameter combinations
- **Authentication Tokens:** Valid and expired OAuth tokens
- **Content Samples:** Various markdown and plain text samples

## Automation Status

| Scenario | Test Type | Status | Framework | Priority |
|----------|-----------|--------|-----------|----------|
| MCP-001 | Protocol | ❌ Not Implemented | Go testify | Critical |
| MCP-002 | Tool Discovery | ❌ Not Implemented | Go testify | Critical |
| MCP-003 | Replace Content | ❌ Not Implemented | Go testify | Critical |
| MCP-004 | Append Content | ❌ Not Implemented | Go testify | High |
| MCP-005 | Insert Content | ❌ Not Implemented | Go testify | High |
| MCP-006 | Authentication | ❌ Not Implemented | Go testify | Critical |
| MCP-007 | Google Docs API | ❌ Not Implemented | Integration | Critical |
| MCP-008 | Markdown Conversion | ❌ Not Implemented | Go testify | High |
| MCP-009 | Export Conversion | ❌ Not Implemented | Go testify | Medium |
| MCP-010 | Error Handling | ❌ Not Implemented | Go testify | High |
| MCP-011 | Data Validation | ❌ Not Implemented | Go testify | High |
| MCP-012 | Large Documents | ❌ Not Implemented | Performance | Medium |
| MCP-013 | Concurrency | ❌ Not Implemented | Load Test | Medium |
| MCP-014 | Claude Integration | ❌ Not Implemented | E2E | Critical |
| MCP-015 | Multi-Client | ❌ Not Implemented | Integration | Low |

## Coverage Metrics

### Current Coverage (Story 1.1)
- **MCP Service:** 0% (not implemented in MVP)
- **Protocol Compliance:** 0% (future implementation)
- **Document Operations:** 0% (future implementation)
- **Google Docs Integration:** 0% (future implementation)

### Target Coverage (Future Implementation)
- **Protocol Compliance:** 100% (all MCP features)
- **Core Operations:** 100% (document manipulation tools)
- **Error Handling:** 95% (comprehensive error scenarios)
- **Performance:** 80% (load and stress testing)

## Implementation Roadmap

### Phase 1: Core MCP Server (Epic 2)
- MCP server initialization with Mark3Labs library
- Basic tool registration and discovery
- Simple document operations (replace_all_content)

### Phase 2: Document Operations (Epic 3)
- Full set of document manipulation tools
- Markdown to Google Docs conversion
- Content validation and sanitization

### Phase 3: Advanced Features (Epic 4)
- Performance optimization
- Multi-client support
- Advanced error recovery

This document serves as the comprehensive MCP service testing specification and will be implemented as the MCP functionality is developed in future epics.