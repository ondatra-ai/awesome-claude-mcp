<!-- Powered by BMAD™ Core -->

# Generate Playwright Test from BDD Scenario

## Your Task

Generate a Playwright test for scenario **{{.ScenarioID}}** and update the requirements registry.

---

## Step 0: Understand Project Architecture & Best Practices

**CRITICAL**: Before generating tests, gather both local project knowledge and official framework best practices.

### Local Project Knowledge

1. **Read Architecture Documents**:
  - Read(`docs/architecture.md`) - Architecture Document
  - Read(`docs/architecture/source-tree.md`) - Source Tree Structure
  - Read(`docs/architecture/coding-standards.md`) - Coding Standards
  - Read(`docs/architecture/tech-stack.md`) - Tech Stack
  - Read(`docs/architecture/mcp-e2e-testing.md`) - MCP E2E Testing with Claude SDK

2. **Determine Correct Service**:
   - Review which service this scenario actually tests (frontend, backend, mcp-service)
   - The scenario's service field indicates the target service
   - Example: If scenario tests MCP WebSocket → service should be `mcp-service`

3. **Verify Test File Location**:
   - Ensure test file path aligns with service structure from source-tree.md
   - Cross-reference with story tasks to see specified file paths

### Official Framework & Library Documentation (Context7)

4. **Fetch Latest Documentation for Relevant Libraries**:

   Based on the test level (`{{.Level}}`) and service (`{{.Service}}`), fetch official documentation from Context7 for the libraries and frameworks you'll use.

   **Step 4.1: Identify Relevant Libraries**

   Review `docs/architecture/tech-stack.md` (already read in step 1) to identify which libraries are relevant for this test:

   **For ALL tests:**
   - **Playwright** (`/microsoft/playwright`) - Testing framework
   - **TypeScript** patterns for type-safe test code

   **For Integration (INT) tests:**
   - **Fetch API** or **Axios** - HTTP client patterns (if testing REST APIs)
   - **@modelcontextprotocol/sdk** - If testing MCP protocol
   - **@anthropic-ai/sdk** - If testing with Claude API simulation

   **For UI E2E tests:**
   - **Next.js** (`/vercel/next.js`) - If testing frontend routing/rendering
   - **React Hook Form** - If testing forms
   - **Zod** - If validating form schemas
   - **Zustand** - If testing state management

   **For MCP E2E tests (Claude SDK):**
   - **@anthropic-ai/sdk** - Claude API client for LLM-driven tool selection
   - **@modelcontextprotocol/sdk** - MCP client for tool execution
   - Claude decides which tools to call based on prompts
   - Tests real LLM → MCP Server → Google Docs flow
   - No browser required (Playwright for test structure and assertions only)
   - See `docs/architecture/mcp-e2e-testing.md` for implementation patterns

   **For Backend/MCP service tests:**
   - **Go** patterns for backend testing (if applicable)
   - **WebSocket** patterns for MCP protocol

   **Step 4.2: Fetch Documentation via Context7**

   Use Context7 to get official, up-to-date patterns. **Fetch docs in order of importance:**

   ```
   # 1. ALWAYS fetch Playwright (primary testing framework)
   context7.get-library-docs(
     libraryID: "/microsoft/playwright",
     topic: "test assertions, error handling, and best practices",
     tokens: 2000
   )
   ```

   ```
   # 2. Fetch additional libraries based on test type
   # Example for MCP E2E tests:
   context7.get-library-docs(
     libraryID: "/anthropic-ai/sdk",  # Use resolver if needed
     topic: "testing with Claude API, message handling, tool calling",
     tokens: 1500
   )
   ```

   ```
   # Example for Frontend E2E tests:
   context7.get-library-docs(
     libraryID: "/vercel/next.js",
     topic: "testing Next.js applications, routing, server components",
     tokens: 1500
   )
   ```

   **Why use Context7:**
   - ✅ Official documentation from library maintainers (not hallucinated patterns)
   - ✅ Latest best practices from framework teams
   - ✅ Real code examples from official docs
   - ✅ Version-specific guidance for current library versions
   - ✅ Ensures tests follow recommended patterns for each library

   **What to extract from Context7 results:**

   *From Playwright docs:*
   - How to use `expect()` vs `expect.soft()` vs `expect.poll()`
   - Custom error message patterns: `expect(locator, 'custom message').toBeVisible()`
   - Web-first assertions: `await expect(page.getByText()).toBeVisible()`
   - Async event handling patterns (for WebSocket/event-driven tests)
   - Assertion best practices and anti-patterns

   *From library-specific docs:*
   - Recommended testing patterns for that library
   - Common test scenarios and how to implement them
   - Integration points with Playwright
   - API patterns and error handling

   **Apply these patterns** in Step 3 when generating test code.

---

## Scenario Details

**ID**: `{{.ScenarioID}}`
**Description**: {{.Description}}
**Level**: {{.Level}}
**Service**: {{.Service}}
**Priority**: {{.Priority}}

**Test Steps (Given-When-Then)**:
```
{{.FormatSteps}}
```

---

## Step 1: Determine Test File Path

Based on the scenario metadata:
- Level: `{{.Level}}`
- Service: `{{.Service}}`

**Determine target file:**
- If level=`e2e` AND service=`mcp-service` → `tests/e2e/mcp-claude-sdk.spec.ts`
- Otherwise → `tests/{{.Level}}/{{.Service}}.spec.ts`

**⚠️ IMPORTANT**: MCP E2E tests go to `mcp-claude-sdk.spec.ts`, NOT `mcp-service.spec.ts`

---

## Step 2: Read Existing Test File (if exists)

**For MCP E2E tests (level=e2e, service=mcp-service):**
```
Read tests/e2e/mcp-claude-sdk.spec.ts
```

**For all other tests:**
```
Read tests/{{.Level}}/{{.Service}}.spec.ts
```

If the file exists, analyze:
- Import statements pattern
- test.describe structure
- Existing test naming convention
- Assertion style (expect patterns)
- Variable naming conventions

If the file doesn't exist, you'll create it with proper structure.

---

## Step 3: Generate Test Code

### CRITICAL: Determine Test Pattern Based on Level and Service

**Before writing any code, determine which pattern to use:**

| Level | Service | Pattern | DO NOT USE |
|-------|---------|---------|------------|
| `integration` | any | Playwright Request API | - |
| `e2e` | `frontend` | Playwright Browser API | - |
| `e2e` | `mcp-service` | **Claude SDK + MCP SDK** | ❌ Playwright Request API |

---

### For Integration Tests (API)
Use Playwright Request API (`{ request }`):

```typescript
test('{{.ScenarioID}}: {{.Description}}', async ({ request }) => {
  // Given: [map Given steps to setup]

  // When: [map When steps to API call]
  const response = await request.get(`${backendUrl}/endpoint`);

  // Then: [map Then steps to assertions]
  expect(response.status()).toBe(200);
  const data = await response.json();
  expect(data).toHaveProperty('key', 'value');
});
```

---

### For MCP E2E Tests (Claude SDK) - MANDATORY FOR e2e + mcp-service

**⚠️ CRITICAL: When level=`e2e` AND service=`mcp-service`, you MUST use Claude SDK.**

**❌ DO NOT:**
- Use Playwright Request API (`{ request }`)
- Use direct HTTP calls to MCP endpoints
- Copy patterns from existing HTTP-based tests

**✅ MUST USE:**
- `@anthropic-ai/sdk` for Claude API client
- `@modelcontextprotocol/sdk` for MCP client
- Claude decides which tools to call (LLM-driven tool selection)

**Target file**: `tests/e2e/mcp-claude-sdk.spec.ts` (NOT `mcp-service.spec.ts`)

**Required imports:**
```typescript
import { test, expect } from '@playwright/test';
import Anthropic from '@anthropic-ai/sdk';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { SSEClientTransport } from '@modelcontextprotocol/sdk/client/sse.js';
import { getEnvironmentConfig } from '../config/environments';
```

**Test pattern:**
```typescript
test('{{.ScenarioID}}: {{.Description}}', async () => {
  const { mcpServiceUrl } = getEnvironmentConfig(process.env.E2E_ENV);

  // Given: MCP server runs with document operation tools
  const mcpClient = new Client({ name: 'e2e-test', version: '1.0.0' });
  const transport = new SSEClientTransport(new URL(`${mcpServiceUrl}/mcp`));
  await mcpClient.connect(transport);

  // Get tools from MCP server and convert to Claude format
  const { tools } = await mcpClient.listTools();
  const claudeTools = tools.map(tool => ({
    name: tool.name,
    description: tool.description || '',
    input_schema: tool.inputSchema as Anthropic.Tool.InputSchema,
  }));

  // When: Claude API client performs workflow
  const anthropic = new Anthropic();
  const response = await anthropic.messages.create({
    model: 'claude-sonnet-4-20250514',
    max_tokens: 1024,
    tools: claudeTools,
    messages: [{
      role: 'user',
      content: 'Your prompt here based on scenario'
    }]
  });

  // Then: Verify Claude selected correct tool and MCP returned valid result
  const toolUse = response.content.find(block => block.type === 'tool_use');
  expect(toolUse, 'Claude should use a tool').toBeDefined();

  // Execute tool via MCP
  const result = await mcpClient.callTool({
    name: toolUse.name,
    arguments: toolUse.input as Record<string, unknown>,
  });

  expect(result.isError, 'Tool call should succeed').toBeFalsy();

  // Cleanup
  await mcpClient.close();
});
```

**See `docs/architecture/mcp-e2e-testing.md` for complete patterns including:**
- Multi-tool workflows
- Error handling
- Helper utilities (ClaudeMcpSession)

---

### For UI E2E Tests (Browser)
Use Playwright Browser API (`{ page }`):

```typescript
test('{{.ScenarioID}}: {{.Description}}', async ({ page }) => {
  // Given: [map Given steps to navigation/setup]
  await page.goto(frontendUrl);

  // When: [map When steps to interactions]
  await page.click('button');

  // Then: [map Then steps to visibility/content checks]
  await expect(page.locator('.result')).toBeVisible();
  await expect(page.locator('.result')).toHaveText('Expected');
});
```

---

## Step 4: Add Test to File

### If file exists:
- Append test inside the existing `test.describe()` block
- Match the existing code style and formatting
- Preserve all existing tests

### If file doesn't exist:
Create complete file structure:
```typescript
import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('{{.Service | title}} {{.Level | title}} Tests', () => {

  test('{{.ScenarioID}}: {{.Description}}', async ({ request }) => {
    // Implementation here
  });

});
```

---

## Step 5: Update Requirements Registry

Update `{{.RequirementsFile}}` for scenario `{{.ScenarioID}}`:

```yaml
{{.ScenarioID}}:
  # ... keep existing fields ...
  implementation_status:
    status: "implemented"  # Change from "pending"
    file_path: "tests/{{.Level}}/{{.Service}}.spec.ts"
  # ... keep other fields unchanged ...
```

Use the Edit tool to update only the `implementation_status` section.

---

## Mapping Guide: Given-When-Then → Playwright Code

### Given (Preconditions)
- "service runs" → Setup base URL/config
- "user authenticated" → Add auth headers/cookies
- "data exists in database" → Pre-populate via API
- "page loaded" → `await page.goto(url)`

### When (Actions)
- "client sends GET request" → `await request.get(url)`
- "user clicks button" → `await page.click('button')`
- "form submitted" → `await page.fill() + page.click('submit')`
- "API called" → `await request.post(url, { data })`

### Then (Assertions)
- "returns 200 status" → `expect(response.status()).toBe(200)`
- "element visible" → `await expect(locator).toBeVisible()`
- "contains text" → `await expect(locator).toHaveText('text')`
- "property exists" → `expect(data).toHaveProperty('key')`

---

## Important Rules

1. **Test ID in name**: Always prefix test name with `{{.ScenarioID}}:`
2. **Match existing style**: If file exists, follow its patterns exactly
3. **Comments**: Add Given/When/Then comments for clarity
4. **No mocks**: Generate real API/browser interactions
5. **Update requirements**: Always update implementation_status after test creation

---

## CRITICAL: Explicit Failure Handling - Avoid Hiding Problems

**Tests MUST fail fast and explicitly. Any issue must immediately fail the test with a clear error message.**

_Note: The patterns below align with official Playwright documentation fetched via Context7 in Step 0. Always prefer the official docs from Context7 over these examples if there are conflicts._

### ❌ ANTI-PATTERN: Conditional Logic That Hides Failures

**DO NOT** write tests with conditional guards that silently skip logic:

```typescript
// ❌ BAD: Hides failures with state machine logic
if (!handshakeComplete && message.id === 1 && message.result) {
  handshakeComplete = true;
  return;  // Silently skips if condition fails
}

// Later:
expect(handshakeComplete).toBeTruthy();  // Fails, but WHY?
```

**Problems**:
- If `message.error` instead of `message.result` → condition fails silently
- If wrong `message.id` → condition fails silently
- Test fails at final assertion with unclear error
- Root cause is hidden in conditional logic

### ✅ CORRECT PATTERN: Explicit Error Checking with expect()

**DO** write tests that fail immediately with descriptive expect() assertions:

```typescript
// ✅ GOOD: Explicit error checking with expect()
if (message.id === 1) {
  expect(message.error, 'Handshake should not return error').toBeUndefined();
  expect(message.result, 'Handshake response must include result').toBeDefined();
  expect(handshakeComplete, 'Should not receive duplicate handshake response').toBe(false);

  handshakeComplete = true;
  return;
}
```

**Benefits**:
- Clear error: `Expected: undefined, Received: { code: -32600, message: "Invalid request" }`
- Immediate failure at root cause with Playwright error formatting
- Custom messages explain WHY the assertion exists
- Easy debugging with stack traces
- No hidden state

### Principles for Test Generation

1. **Use Playwright expect() for All Assertions**: Follow Playwright's web-first assertion pattern
   ```typescript
   // ✅ GOOD: Use expect() with custom messages
   await expect(response, 'API response should succeed').toHaveStatus(200);
   expect(response.status(), 'Expected successful status').toBe(200);

   // ❌ BAD: Don't throw errors for assertions
   if (response.status() !== 200) {
     throw new Error(`Expected 200, got ${response.status()}`);
   }
   ```

2. **Explicit Error Checks with expect()**: Check errors explicitly, don't filter them silently
   ```typescript
   // ❌ BAD: Hides errors with conditional filtering
   if (message.result) {
     // Process result - but what if message.error exists?
   }

   // ✅ GOOD: Check for errors explicitly
   expect(message.error, 'Server should not return error').toBeUndefined();
   expect(message.result, 'Response must have result field').toBeDefined();
   ```

3. **Custom Messages for Context**: Add descriptive messages to expect() calls
   ```typescript
   // ✅ GOOD: Clear context in assertion
   expect(tools.find(t => t.name === 'replace_all'),
     'replace_all tool should be in tools list').toBeDefined();

   // ❌ BAD: Generic assertion without context
   expect(tools.find(t => t.name === 'replace_all')).toBeDefined();
   ```

4. **Validate State Transitions Explicitly**: Don't rely on implicit state from previous steps
   ```typescript
   // ❌ BAD: Step 2 silently skips if step 1 failed
   if (!step1Complete && condition) { step1Complete = true; }
   if (!step2Complete && step1Complete) { step2Complete = true; }
   // Later: expect(step2Complete).toBeTruthy(); // Unclear why it failed

   // ✅ GOOD: Explicit assertions at each transition
   expect(step1Complete, 'Step 1 must complete before Step 2').toBe(true);
   if (step2Condition) {
     step2Complete = true;
   }
   expect(step2Complete, 'Step 2 should complete successfully').toBe(true);
   ```

5. **Direct Assertions Over Complex Result Objects**: Keep tests simple
   ```typescript
   // ❌ BAD: Complex result object with boolean flags
   const result = { success: false, handshakeComplete: false };
   // ... 100 lines of complex state machine logic ...
   expect(result.success).toBeTruthy(); // Unclear what failed

   // ✅ GOOD: Assert immediately after each operation
   const response = await request.get(url);
   expect(response.status(), 'Initialize request should succeed').toBe(200);
   const data = await response.json();
   expect(data.connectionId, 'Response should include connectionId').toBeDefined();
   ```

6. **Use expect.soft() for Multiple Related Checks**: When you need to check multiple things
   ```typescript
   // ✅ GOOD: Check all fields even if some fail
   await expect.soft(healthData, 'Health endpoint returns status').toHaveProperty('status');
   await expect.soft(healthData, 'Health endpoint returns connections').toHaveProperty('connections');
   await expect.soft(healthData, 'Health endpoint returns dependencies').toHaveProperty('dependencies');
   ```

### When Complex Flow Is Necessary (WebSocket, Async Events)

For WebSocket or event-driven tests, **collect state then assert with expect()**:

```typescript
test('WebSocket flow with explicit assertions', async ({ page }) => {
  const result = await page.evaluate(async () => {
    return new Promise((resolve) => {
      const ws = new WebSocket(url);
      const state = {
        connected: false,
        messageReceived: false,
        error: null as string | null,
        data: null as any
      };

      ws.onopen = () => {
        state.connected = true;
        ws.send(JSON.stringify(request));
      };

      ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        state.messageReceived = true;

        if (msg.error) {
          state.error = `Server error: ${msg.error.message}`;
          ws.close();
          resolve(state);
          return;
        }

        if (!msg.result) {
          state.error = 'Response missing result field';
          ws.close();
          resolve(state);
          return;
        }

        state.data = msg.result;
        ws.close();
        resolve(state);
      };

      ws.onerror = () => {
        state.error = 'WebSocket connection error';
        resolve(state);
      };

      setTimeout(() => {
        state.error = 'WebSocket timeout';
        ws.close();
        resolve(state);
      }, 5000);
    });
  });

  // ✅ Assert each aspect with descriptive messages
  expect(result.connected, 'WebSocket should establish connection').toBe(true);
  expect(result.messageReceived, 'Should receive message from server').toBe(true);
  expect(result.error, 'Should not have errors').toBeNull();
  expect(result.data, 'Response should include data').toBeDefined();
  expect(result.data, 'Data should have expected field').toHaveProperty('expectedField');
});
```

**Key difference**: Instead of throwing errors inside callbacks, **collect state and assert outside**:
- ✅ `expect()` provides better error messages
- ✅ Multiple assertions give complete picture
- ✅ Custom messages explain what failed
- ❌ Don't throw errors - let expect() handle failures

### CRITICAL: WebSocket Reconnection Timeout Handling

**When testing WebSocket reconnections, ALWAYS add explicit connection timeouts to prevent hanging tests.**

**Problem**: Browser WebSocket connections in page.evaluate() can hang indefinitely if:
- Server doesn't respond to connection attempt
- Network issue prevents connection completion
- Connection cleanup from previous test hasn't finished

**❌ BAD: Reconnection without timeout**
```typescript
const ws2 = new WebSocket(url);

ws2.onopen = () => {
  state.reconnected = true;  // May never fire!
  ws2.send(JSON.stringify(request));
};

ws2.onerror = () => {
  state.error = 'Connection failed';  // May never fire!
  resolve(state);
};
```

**Why this fails:**
- If connection hangs, neither `onopen` nor `onerror` fires
- Promise never resolves
- Test hangs until Playwright timeout (30+ seconds)
- `state.reconnected` stays false with no explanation

**✅ GOOD: Reconnection with explicit timeout**
```typescript
const ws2 = new WebSocket(url);
let connectionTimeout: NodeJS.Timeout;

// Set timeout BEFORE waiting for connection
connectionTimeout = setTimeout(() => {
  if (!state.reconnected) {
    state.error = 'Reconnection attempt timed out after 3000ms';
    ws2.close();
    resolve(state);
  }
}, 3000); // 3 second timeout for connection attempt

ws2.onopen = () => {
  clearTimeout(connectionTimeout);  // Cancel timeout on success
  state.reconnected = true;
  ws2.send(JSON.stringify(request));
};

ws2.onerror = (error) => {
  clearTimeout(connectionTimeout);  // Cancel timeout on explicit error
  state.error = `Reconnection WebSocket error: ${error}`;
  resolve(state);
};

ws2.onclose = () => {
  clearTimeout(connectionTimeout);  // Cancel timeout on close
  if (!state.reconnected) {
    state.error = 'Connection closed before reconnection completed';
    resolve(state);
  }
};
```

**Why this works:**
- ✅ Test fails explicitly after 3 seconds if connection hangs
- ✅ Clear error message: "Reconnection attempt timed out after 3000ms"
- ✅ Timeout is cleared if connection succeeds or explicitly fails
- ✅ No hanging tests
- ✅ Fast feedback on connection issues

**Reconnection delay guidelines:**
- Wait 500-2000ms between close and reconnect (allows server cleanup)
- Set 3-5 second timeout for connection attempt
- Clear all timeouts in success/error/close handlers
- Always resolve() promise in timeout handler

**Example: Complete reconnection flow**
```typescript
// Close first connection
ws.close(1006, 'Simulated disconnect');

// Wait for server cleanup
setTimeout(() => {
  const ws2 = new WebSocket(url);
  let connectionTimeout: NodeJS.Timeout;

  connectionTimeout = setTimeout(() => {
    if (!state.reconnected) {
      state.error = 'Reconnection timeout - server may not be accepting new connections';
      if (ws2.readyState === WebSocket.CONNECTING || ws2.readyState === WebSocket.OPEN) {
        ws2.close();
      }
      resolve(state);
    }
  }, 3000);

  ws2.onopen = () => {
    clearTimeout(connectionTimeout);
    state.reconnected = true;
    // ... continue test
  };

  ws2.onerror = () => {
    clearTimeout(connectionTimeout);
    state.error = 'Reconnection failed';
    resolve(state);
  };

  ws2.onclose = () => {
    clearTimeout(connectionTimeout);
    if (!state.reconnected && !state.error) {
      state.error = 'Connection closed unexpectedly';
      resolve(state);
    }
  };
}, 1000); // Wait 1 second before reconnect attempt
```

---

## Output Requirements

After completing the generation, provide:

```
Test Generation Summary:
- Test ID: {{.ScenarioID}}
- Description: {{.Description}}
- File: tests/{{.Level}}/{{.Service}}.spec.ts
- Status: [CREATED | APPENDED]
- Requirements Updated: YES
```

**If any errors occur**, provide clear error message with:
- What step failed
- Error details
- Suggested fix
