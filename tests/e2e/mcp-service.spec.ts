import { test, expect } from '@playwright/test';

import { getEnvironmentConfig } from '../config/environments';

import {
  McpClient,
  createMultipleClients,
  closeAllClients,
  getResultText,
} from './helpers/mcp-client';

const { mcpServiceUrl } = getEnvironmentConfig(process.env.E2E_ENV);

/**
 * MCP Service E2E Tests
 *
 * Tests MCP protocol compliance using @modelcontextprotocol/sdk.
 * Validates tool listing and execution via the MCP server.
 */
test.describe('MCP Service E2E Tests', () => {
  /**
   * E2E-011: MCP client connects and lists tools
   */
  test('E2E-011: Connect and list tools', async () => {
    const client = new McpClient(mcpServiceUrl);
    await client.connect();

    expect(client.isConnected()).toBe(true);

    const tools = await client.listTools();

    expect(tools.length).toBeGreaterThan(0);
    expect(tools[0]).toHaveProperty('name');

    await client.close();
  });

  /**
   * E2E-012: MCP client calls replace_all tool
   */
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

  /**
   * E2E-013: MCP client calls append tool
   */
  test('E2E-013: Call append tool', async () => {
    const client = new McpClient(mcpServiceUrl);
    await client.connect();

    const result = await client.callTool('append', {
      documentId: 'test-doc-456',
      content: 'This is a test paragraph.',
    });

    expect(result.isError).toBeFalsy();
    expect(getResultText(result)).toContain('success');

    await client.close();
  });

  /**
   * E2E-014: Multiple clients connect concurrently
   */
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

  /**
   * E2E-015: Health check returns healthy status
   */
  test('E2E-015: Health check', async ({ request }) => {
    const response = await request.get(`${mcpServiceUrl}/health`);

    expect(response.status()).toBe(200);

    const healthData = await response.json();
    expect(healthData.status).toBe('healthy');
    expect(healthData.dependencies).toBeDefined();
  });

  /**
   * E2E-016: Tool call completes within SLA
   */
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
});
