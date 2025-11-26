# MCP E2E Testing Guide

This guide covers end-to-end testing for MCP (Model Context Protocol) servers using `@modelcontextprotocol/sdk`.

## Overview

MCP E2E tests verify MCP protocol compliance by:
- Connecting to the MCP server using standard transports
- Listing available tools
- Executing tool calls
- Validating responses follow the MCP specification

### Architecture

```
┌─────────────────┐     ┌─────────────────┐
│  Playwright     │────▶│   MCP Server    │
│  Test Runner    │     │   (Go/HTTP)     │
└─────────────────┘     └─────────────────┘
        │
        ▼
┌─────────────────┐
│   McpClient     │
│  (SDK wrapper)  │
└─────────────────┘
```

Tests use Playwright as the test runner with `@modelcontextprotocol/sdk` for MCP protocol communication.

## Prerequisites

### Dependencies

```json
{
  "devDependencies": {
    "@modelcontextprotocol/sdk": "^1.12.0",
    "@playwright/test": "^1.56.1"
  }
}
```

### Docker Services

The test pipeline requires these services running:
- `mcp-service-test` - MCP server on port 8081
- `mcp-backend-test` - Backend API on port 8080
- `mcp-frontend-test` - Frontend on port 3000

## MCP Client Helper

### Location

```
tests/e2e/helpers/mcp-client.ts
```

### Interfaces

```typescript
interface IMcpClientOptions {
  name?: string;    // Client identifier (default: 'e2e-test-client')
  version?: string; // Client version (default: '1.0.0')
}

interface IMcpTool {
  name: string;
  description?: string;
  inputSchema?: Record<string, unknown>;
}

interface IToolResult {
  content: Array<{ type: string; text?: string }>;
  isError?: boolean;
}
```

### McpClient Class

```typescript
import { McpClient, getResultText } from './helpers/mcp-client';

const client = new McpClient('http://localhost:8081');
await client.connect();

// List tools
const tools = await client.listTools();

// Call a tool
const result = await client.callTool('tool_name', { arg1: 'value' });

// Extract text from result
const text = getResultText(result);

// Cleanup
await client.close();
```

### Connection Strategy

The client attempts connection using:
1. **StreamableHTTPClientTransport** (preferred)
2. **SSEClientTransport** (fallback)

This ensures compatibility with different MCP server implementations.

### Helper Functions

```typescript
// Create multiple clients for concurrent testing
const clients = await createMultipleClients(baseUrl, 3);

// Close all clients
await closeAllClients(clients);

// Extract text from IToolResult
const text = getResultText(result);
```

## Writing Tests

### Test File Location

```
tests/e2e/mcp-service.spec.ts
```

### Test Structure

```typescript
import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';
import {
  McpClient,
  createMultipleClients,
  closeAllClients,
  getResultText,
} from './helpers/mcp-client';

const { mcpServiceUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('MCP Service E2E Tests', () => {
  test('E2E-XXX: Test description', async () => {
    const client = new McpClient(mcpServiceUrl);
    await client.connect();

    // Test logic here

    await client.close();
  });
});
```

### Naming Convention

Test IDs follow the pattern: `E2E-XXX` where XXX is a three-digit sequence number.

### Example: Connect and List Tools

```typescript
test('E2E-011: Connect and list tools', async () => {
  const client = new McpClient(mcpServiceUrl);
  await client.connect();

  expect(client.isConnected()).toBe(true);

  const tools = await client.listTools();

  expect(tools.length).toBeGreaterThan(0);
  expect(tools[0]).toHaveProperty('name');

  await client.close();
});
```

### Example: Call a Tool

```typescript
test('E2E-012: Call replace_all tool', async () => {
  const client = new McpClient(mcpServiceUrl);
  await client.connect();

  const result = await client.callTool('replace_all', {
    documentId: 'test-doc-123',
    content: '# Hello World',
  });

  expect(result.isError).toBeFalsy();
  expect(getResultText(result)).toContain('success');

  await client.close();
});
```

### Example: Concurrent Connections

```typescript
test('E2E-014: Concurrent client connections', async () => {
  const NUM_CLIENTS = 3;
  const clients = await createMultipleClients(mcpServiceUrl, NUM_CLIENTS);

  expect(clients.length).toBe(NUM_CLIENTS);

  for (const client of clients) {
    expect(client.isConnected()).toBe(true);
  }

  const toolsResults = await Promise.all(
    clients.map((client) => client.listTools())
  );

  for (const tools of toolsResults) {
    expect(tools.length).toBeGreaterThan(0);
  }

  await closeAllClients(clients);
});
```

### Example: Performance Test

```typescript
test('E2E-016: Tool call completes within 2 seconds', async () => {
  const client = new McpClient(mcpServiceUrl);
  await client.connect();

  const thresholdMs = 2000;
  const startTime = performance.now();

  const result = await client.callTool('append', {
    documentId: 'perf-test-doc',
    content: 'Performance test',
  });

  const duration = performance.now() - startTime;

  expect(result.isError).toBeFalsy();
  expect(getResultText(result)).toContain('success');
  expect(duration).toBeLessThanOrEqual(thresholdMs);

  await client.close();
});
```

### Example: Health Check (HTTP)

For non-MCP endpoints, use Playwright's request API:

```typescript
test('E2E-015: Health check', async ({ request }) => {
  const response = await request.get(`${mcpServiceUrl}/health`);

  expect(response.status()).toBe(200);

  const healthData = await response.json();
  expect(healthData.status).toBe('healthy');
  expect(healthData.dependencies).toBeDefined();
});
```

## Running Tests

### Local Development

```bash
# Run all E2E tests (starts Docker services automatically)
make test-e2e

# Run with specific environment
E2E_ENV=dev make test-e2e
```

### Environment Configuration

Tests support multiple environments via `E2E_ENV`:

| Environment | MCP Service URL |
|-------------|-----------------|
| `local` (default) | `http://localhost:8081` |
| `dev` | `https://mcp.dev.ondatra-ai.xyz` |

Configure environments in `tests/config/environments.ts`.

### Docker Test Pipeline

The `make test-e2e` command:
1. Starts backend, mcp-service, and frontend containers
2. Waits for health checks to pass
3. Builds and runs the Playwright test container
4. Cleans up containers after tests complete

## Test Scenarios Reference

| ID | Description | Type |
|----|-------------|------|
| E2E-011 | Connect and list tools | Connection |
| E2E-012 | Call replace_all tool | Tool execution |
| E2E-013 | Call append tool | Tool execution |
| E2E-014 | Concurrent client connections | Load |
| E2E-015 | Health check endpoint | HTTP |
| E2E-016 | Tool call within 2s SLA | Performance |

## Best Practices

1. **Always close connections**: Use `client.close()` in every test
2. **Check isError**: Validate `result.isError` is falsy before checking content
3. **Use getResultText**: Extract text content consistently with the helper
4. **Isolate tests**: Each test should create its own client instance
5. **Use meaningful document IDs**: Help identify test data in logs
