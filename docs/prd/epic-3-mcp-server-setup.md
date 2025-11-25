# Epic 3: MCP Server Setup

**Goal:** Create functional MCP protocol server with tool registration and bidirectional communication

## User Stories

### Story 3.1: MCP Server Implementation
**As a** Developer/Maintainer
**I want** to implement MCP protocol server
**So that** Claude can communicate with the service

**Acceptance Criteria:**
- WebSocket server implemented
- HTTP endpoint for MCP available
- Message parsing and validation
- Response formatting to MCP standard
- Connection management handled
- Concurrent connection support

### Story 3.2: MCP Integration Testing with LLM Simulation
**As a** Developer/Maintainer
**I want** to test MCP server with realistic Claude API client simulation
**So that** MCP tool calling works correctly with real LLM behavior

**Acceptance Criteria:**
- MCP server responds correctly to tool discovery requests from Claude API client
- MCP server processes tool invocations from Claude API client and returns valid responses
- E2E test verifies complete Claude → MCP Server → Tool Execution → Response flow
- MCP server tool invocation protocol tested with realistic LLM request/response patterns
- MCP server error handling verified with invalid tool names, parameters, and message formats

**Technical Approach:**
- Install and configure `@modelcontextprotocol/sdk` for MCP E2E testing
- Use `@playwright/test` as test runner framework (no browser automation)
- Use `McpClient` helper at `tests/e2e/helpers/mcp-client.ts`
- Create test client that connects to MCP server via HTTP+SSE transport
- Test MCP protocol patterns (list tools → call tool → process response)
- Build test fixtures for different tool schemas and expected responses
- Create example tests showing complete Claude ↔ MCP Server ↔ Tool flow
- Document LLM simulation patterns for future test development

### Story 3.3: Tool Registration
**As a** Claude User
**I want** to discover available tools
**So that** I know what operations are available

**Acceptance Criteria:**
- Tool definition schema created
- Tool registration endpoint working
- Tool capabilities described clearly
- Parameter schemas defined
- Version information included
- Dynamic or static registration (TBD after testing)

### Story 3.4: Message Protocol Handler
**As a** Developer/Maintainer
**I want** to process MCP messages correctly
**So that** operations are executed properly

**Acceptance Criteria:**
- Request message parsing implemented
- Command routing to handlers
- Response message formatting
- Error message standards followed
- Message validation complete
- Correlation ID tracking

### Story 3.5: MCP Error Handling
**As a** Claude User
**I want** standard MCP error responses
**So that** I can handle failures appropriately

**Acceptance Criteria:**
- Standard MCP error format used
- Error codes properly categorized
- Descriptive error messages provided
- Stack traces excluded from production
- Error logging implemented
- Rate limit errors handled

### Story 3.6: Connection Management
**As a** Developer/Maintainer
**I want** robust connection handling
**So that** communication remains stable

**Acceptance Criteria:**
- WebSocket heartbeat implemented
- Connection timeout handling
- Reconnection logic for clients
- Graceful shutdown handling
- Connection pooling if needed
- Connection state tracking
