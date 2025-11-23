# MCP E2E Testing with Claude SDK

## Purpose

This document defines the architecture and patterns for End-to-End (E2E) testing of MCP services using the Claude SDK. Unlike integration tests that use direct HTTP requests, E2E tests simulate realistic LLM behavior by having Claude decide which tools to use.

**Last Updated**: 2025-11-22

---

## Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Test Suite    │───▶│   Claude SDK     │───▶│   MCP Server    │
│  (Playwright)   │    │ (@anthropic-ai)  │    │  (HTTP+SSE)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │                        │
                              │ tool_use               │ Google Docs API
                              ▼                        ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │   MCP Client     │───▶│   Google Docs   │
                       │ (@mcp/sdk)       │    │   (Real Docs)   │
                       └──────────────────┘    └─────────────────┘
```

### Flow

1. **Test sends prompt** to Claude API with MCP tools registered
2. **Claude decides** which tool(s) to use based on the prompt
3. **MCP Client** executes tool calls against MCP Server
4. **MCP Server** performs operations on Google Docs
5. **Claude returns** final response with tool results
6. **Test validates** the complete flow succeeded

---

## Key Differences: INT vs E2E Tests

| Aspect | Integration (INT) | End-to-End (E2E) |
|--------|-------------------|------------------|
| Client | Direct HTTP/JSON-RPC | Claude SDK |
| Tool Selection | Test specifies tool | Claude decides |
| Validation | Protocol compliance | Business outcome |
| Speed | Fast (milliseconds) | Slower (seconds) |
| Cost | Free | Claude API tokens |
| Realism | Low | High |

---

## Required Dependencies

```json
{
  "devDependencies": {
    "@playwright/test": "^1.56.1",
    "@anthropic-ai/sdk": "^0.39.0",
    "@modelcontextprotocol/sdk": "^1.0.0"
  }
}
```

---

## Environment Configuration

### Required Environment Variables

```bash
# .env.test
ANTHROPIC_API_KEY=sk-ant-...          # Claude API key (required)
MCP_SERVICE_URL=http://localhost:8081  # MCP server endpoint
E2E_ENV=local                          # Environment: local, dev, staging
TEST_DOCUMENT_ID=1abc...               # Google Doc ID for testing
```

### Environment Parity

Tests must work identically in:
- **Local**: Developer machine with `make dev`
- **CI**: GitHub Actions with Docker services

---

## Available MCP Tools

The MCP server exposes 6 document operation tools:

| Tool | Description | Required Arguments |
|------|-------------|-------------------|
| `replace_all` | Replace entire document content | `documentId`, `content` |
| `append` | Append content to end | `documentId`, `content` |
| `prepend` | Prepend content to beginning | `documentId`, `content` |
| `replace_match` | Find and replace text | `documentId`, `searchText`, `replaceText` |
| `insert_before` | Insert before anchor text | `documentId`, `anchorText`, `content` |
| `insert_after` | Insert after anchor text | `documentId`, `anchorText`, `content` |

---

## Test Patterns

### Pattern 1: Basic Tool Call via Claude

```typescript
import { test, expect } from '@playwright/test';
import Anthropic from '@anthropic-ai/sdk';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { SSEClientTransport } from '@modelcontextprotocol/sdk/client/sse.js';

test('E2E-016: Claude replaces document content', async () => {
  // Given: MCP server runs with document tools
  const mcpClient = new Client({ name: 'e2e-test', version: '1.0.0' });
  const transport = new SSEClientTransport(
    new URL(`${process.env.MCP_SERVICE_URL}/mcp`)
  );
  await mcpClient.connect(transport);

  // Get available tools from MCP server
  const { tools } = await mcpClient.listTools();

  // Convert MCP tools to Claude format
  const claudeTools = tools.map(tool => ({
    name: tool.name,
    description: tool.description,
    input_schema: tool.inputSchema,
  }));

  // When: User asks Claude to edit document
  const anthropic = new Anthropic();
  const response = await anthropic.messages.create({
    model: 'claude-sonnet-4-20250514',
    max_tokens: 1024,
    tools: claudeTools,
    messages: [{
      role: 'user',
      content: `Replace all content in document ${TEST_DOC_ID} with "# Hello World"`
    }]
  });

  // Then: Claude calls replace_all tool
  const toolUse = response.content.find(block => block.type === 'tool_use');
  expect(toolUse, 'Claude should use a tool').toBeDefined();
  expect(toolUse.name).toBe('replace_all');

  // Execute tool call via MCP
  const result = await mcpClient.callTool({
    name: toolUse.name,
    arguments: toolUse.input,
  });

  // Verify success
  expect(result.content[0].type).toBe('text');
  expect(result.isError).toBeFalsy();
});
```

### Pattern 2: Multi-Tool Workflow

```typescript
test('E2E-018: Claude executes multi-step document workflow', async () => {
  // Given: Active MCP session with Claude
  const { mcpClient, claudeTools } = await setupMcpSession();

  // When: User requests complex operation
  const messages = [{
    role: 'user',
    content: `First replace all content in doc ${DOC_ID} with a title,
              then append a paragraph about testing.`
  }];

  // Claude may make multiple tool calls
  let continueLoop = true;
  while (continueLoop) {
    const response = await anthropic.messages.create({
      model: 'claude-sonnet-4-20250514',
      max_tokens: 1024,
      tools: claudeTools,
      messages,
    });

    if (response.stop_reason === 'tool_use') {
      // Execute each tool call
      for (const block of response.content) {
        if (block.type === 'tool_use') {
          const result = await mcpClient.callTool({
            name: block.name,
            arguments: block.input,
          });

          // Add result to conversation
          messages.push({ role: 'assistant', content: response.content });
          messages.push({
            role: 'user',
            content: [{
              type: 'tool_result',
              tool_use_id: block.id,
              content: result.content[0].text,
            }]
          });
        }
      }
    } else {
      continueLoop = false;
    }
  }

  // Then: Both operations completed
  expect(messages.length).toBeGreaterThan(2);
});
```

### Pattern 3: Error Handling

```typescript
test('E2E-013: Claude handles tool errors gracefully', async () => {
  // Given: MCP session ready
  const { mcpClient, claudeTools } = await setupMcpSession();

  // When: User requests operation on non-existent document
  const response = await anthropic.messages.create({
    model: 'claude-sonnet-4-20250514',
    max_tokens: 1024,
    tools: claudeTools,
    messages: [{
      role: 'user',
      content: 'Replace content in document INVALID_DOC_ID_12345 with "test"'
    }]
  });

  // Execute tool (will fail)
  const toolUse = response.content.find(b => b.type === 'tool_use');
  const result = await mcpClient.callTool({
    name: toolUse.name,
    arguments: toolUse.input,
  });

  // Then: MCP returns structured error
  expect(result.isError).toBeTruthy();
  // Error should have MCP error format
});
```

---

## Helper Utilities

### Recommended: `tests/e2e/helpers/claude-mcp-client.ts`

```typescript
import Anthropic from '@anthropic-ai/sdk';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { SSEClientTransport } from '@modelcontextprotocol/sdk/client/sse.js';

export interface ClaudeMcpSession {
  anthropic: Anthropic;
  mcpClient: Client;
  claudeTools: Anthropic.Tool[];
  cleanup: () => Promise<void>;
}

export async function createClaudeMcpSession(
  mcpUrl: string
): Promise<ClaudeMcpSession> {
  const anthropic = new Anthropic();
  const mcpClient = new Client({ name: 'e2e-test', version: '1.0.0' });

  const transport = new SSEClientTransport(new URL(`${mcpUrl}/mcp`));
  await mcpClient.connect(transport);

  const { tools } = await mcpClient.listTools();
  const claudeTools = tools.map(tool => ({
    name: tool.name,
    description: tool.description || '',
    input_schema: tool.inputSchema as Anthropic.Tool.InputSchema,
  }));

  return {
    anthropic,
    mcpClient,
    claudeTools,
    cleanup: async () => {
      await mcpClient.close();
    },
  };
}

export async function executeClaudeWithTools(
  session: ClaudeMcpSession,
  prompt: string
): Promise<{ finalResponse: string; toolCalls: string[] }> {
  const toolCalls: string[] = [];
  const messages: Anthropic.MessageParam[] = [
    { role: 'user', content: prompt }
  ];

  let response = await session.anthropic.messages.create({
    model: 'claude-sonnet-4-20250514',
    max_tokens: 1024,
    tools: session.claudeTools,
    messages,
  });

  while (response.stop_reason === 'tool_use') {
    const assistantContent = response.content;
    messages.push({ role: 'assistant', content: assistantContent });

    const toolResults: Anthropic.ToolResultBlockParam[] = [];

    for (const block of assistantContent) {
      if (block.type === 'tool_use') {
        toolCalls.push(block.name);
        const result = await session.mcpClient.callTool({
          name: block.name,
          arguments: block.input as Record<string, unknown>,
        });

        toolResults.push({
          type: 'tool_result',
          tool_use_id: block.id,
          content: result.content[0].type === 'text'
            ? result.content[0].text
            : JSON.stringify(result.content),
        });
      }
    }

    messages.push({ role: 'user', content: toolResults });

    response = await session.anthropic.messages.create({
      model: 'claude-sonnet-4-20250514',
      max_tokens: 1024,
      tools: session.claudeTools,
      messages,
    });
  }

  const textBlock = response.content.find(b => b.type === 'text');
  return {
    finalResponse: textBlock?.text || '',
    toolCalls,
  };
}
```

---

## BDD Scenario Writing for Claude SDK Tests

### Allowed Actors

- `Claude API client` - The Claude SDK making decisions
- `MCP server` - The tool execution endpoint
- `User` - The human providing prompts

### Example Scenarios

```gherkin
# E2E-016: Complete MCP workflow
Given: MCP server runs with document operation tools
When: Claude API client performs complete workflow from initialize to tool call
Then: Server processes entire flow and returns valid tool result

# E2E-017: Latency SLA
Given: MCP server accepts HTTP connections on configured endpoint
When: Claude API client executes document operation via MCP tools
Then: Server returns tool result within 2 seconds round-trip time

# E2E-018: JSON-RPC sequence
Given: MCP server accepts HTTP connections via Streamable HTTP transport
When: Claude API client sends JSON-RPC message sequence
Then: Server processes HTTP+SSE transport and returns valid responses
```

---

## CI Configuration

### GitHub Actions Example

```yaml
e2e-tests:
  runs-on: ubuntu-latest
  services:
    mcp-service:
      image: mcp-service:latest
      ports:
        - 8081:8081
  env:
    ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
    MCP_SERVICE_URL: http://localhost:8081
    TEST_DOCUMENT_ID: ${{ secrets.TEST_DOCUMENT_ID }}
  steps:
    - uses: actions/checkout@v4
    - run: npm ci
      working-directory: tests
    - run: npx playwright test tests/e2e/mcp-*.spec.ts
      working-directory: tests
```

---

## References

- [Claude API Documentation](https://docs.anthropic.com/claude/docs)
- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [MCP TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)
- [BDD Guidelines](./bdd-guidelines.md)
- [Coding Standards](./coding-standards.md)
