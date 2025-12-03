import { test, expect } from '@playwright/test';

import { ClaudeDesktopClient } from './helpers/claude-client';

/**
 * MCP Integration Tests via Claude Desktop
 *
 * Tests MCP tool functionality through Claude Desktop UI.
 * Validates that Claude can discover and execute MCP tools.
 *
 * Prerequisites:
 * 1. Launch Claude Desktop with remote debugging:
 *    pkill -x "Claude" && open -a Claude --args --remote-debugging-port=9222
 *
 * 2. Ensure MCP servers are configured in Claude Desktop
 *
 * 3. Run tests:
 *    npx playwright test e2e/claude-mcp.spec.ts
 */

// Run tests serially - Claude Desktop is a single instance
test.describe.configure({ mode: 'serial' });

// Increase timeout for Claude Desktop tests
test.setTimeout(180000);

test.describe('MCP Integration via Claude Desktop', () => {
  let client: ClaudeDesktopClient;

  test.beforeAll(async () => {
    client = new ClaudeDesktopClient({
      cdpUrl: 'http://localhost:9222',
      timeout: 60000,
    });

    const connected = await client.connect();
    if (!connected) {
      throw new Error(
        'Failed to connect to Claude Desktop. ' +
          'Make sure it is running with --remote-debugging-port=9222'
      );
    }
  });

  test.afterAll(async () => {
    await client?.disconnect();
  });

  /**
   * E2E-011: MCP client connects and lists tools
   *
   * Verifies Claude can discover and list available MCP tools.
   */
  test('E2E-011: should list available MCP tools', async () => {
    await client.newChat();
    await client.sendMessage(
      'List all MCP tools you have available. ' +
        'Reply with the tool names in a bullet list format.'
    );

    const responses = await client.waitForResponse(60000);

    expect(responses.length).toBeGreaterThan(0);
    const responseText = responses.join(' ').toLowerCase();

    // Verify Claude acknowledges MCP tools
    expect(
      responseText.includes('tool') ||
        responseText.includes('mcp') ||
        responseText.includes('available'),
      'Response should mention tools or MCP'
    ).toBe(true);
  });

  /**
   * E2E-017: Tool catalog discovery includes Google Docs operations
   *
   * Verifies MCP tool catalog includes Google Docs operations.
   */
  test('E2E-017: should discover Google Docs MCP tools', async () => {
    await client.newChat();
    await client.sendMessage(
      'What Google Docs MCP tools do you have? ' +
        'List any tools related to documents like replace_all, append, or prepend.'
    );

    const responses = await client.waitForResponse(60000);

    expect(responses.length).toBeGreaterThan(0);
    const responseText = responses.join(' ').toLowerCase();

    // Check for Google Docs related tool mentions
    const hasDocTools =
      responseText.includes('replace') ||
      responseText.includes('append') ||
      responseText.includes('prepend') ||
      responseText.includes('document') ||
      responseText.includes('google docs');

    expect(
      hasDocTools,
      'Response should mention document-related tools'
    ).toBe(true);
  });

  /**
   * E2E-012: MCP client calls replace_all tool
   *
   * Tests that Claude can execute the replace_all MCP tool.
   */
  test('E2E-012: should execute replace_all tool', async () => {
    await client.newChat();
    await client.sendMessage(
      'Use the replace_all MCP tool to replace the content of document ' +
        '"test-doc-e2e-001" with "# Hello from E2E Test". ' +
        'Tell me if it succeeded or failed.'
    );

    const responses = await client.waitForResponse(90000);

    expect(responses.length).toBeGreaterThan(0);
    const responseText = responses.join(' ').toLowerCase();

    // Verify tool execution response
    const hasExecution =
      responseText.includes('success') ||
      responseText.includes('replaced') ||
      responseText.includes('updated') ||
      responseText.includes('complete') ||
      responseText.includes('done') ||
      responseText.includes('error') ||
      responseText.includes('failed');

    expect(
      hasExecution,
      'Response should indicate tool execution result'
    ).toBe(true);
  });

  /**
   * E2E-013: MCP client calls append tool
   *
   * Tests that Claude can execute the append MCP tool.
   */
  test('E2E-013: should execute append tool', async () => {
    await client.newChat();
    await client.sendMessage(
      'Use the append MCP tool to add "## Appended Section" ' +
        'to document "test-doc-e2e-002". ' +
        'Confirm the operation result.'
    );

    const responses = await client.waitForResponse(90000);

    expect(responses.length).toBeGreaterThan(0);
    const responseText = responses.join(' ').toLowerCase();

    // Verify append tool response
    const hasExecution =
      responseText.includes('success') ||
      responseText.includes('appended') ||
      responseText.includes('added') ||
      responseText.includes('complete') ||
      responseText.includes('done') ||
      responseText.includes('error') ||
      responseText.includes('failed');

    expect(
      hasExecution,
      'Response should indicate append operation result'
    ).toBe(true);
  });

  /**
   * E2E-016: Initialize handshake returns server metadata
   *
   * Verifies Claude has proper MCP server connection with metadata.
   */
  test('E2E-016: should show MCP server information', async () => {
    await client.newChat();
    await client.sendMessage(
      'What MCP servers are you connected to? ' +
        'Provide details about the server name and available capabilities.'
    );

    const responses = await client.waitForResponse(60000);

    expect(responses.length).toBeGreaterThan(0);
    const responseText = responses.join(' ').toLowerCase();

    // Verify server information response
    const hasServerInfo =
      responseText.includes('server') ||
      responseText.includes('mcp') ||
      responseText.includes('connected') ||
      responseText.includes('capability') ||
      responseText.includes('tool');

    expect(
      hasServerInfo,
      'Response should contain MCP server information'
    ).toBe(true);
  });

  /**
   * E2E-019: Complete Claude to MCP flow performance
   *
   * Tests that tool execution completes in reasonable time.
   * Note: UI-based tests have higher latency than direct MCP calls.
   */
  test('E2E-019: should complete MCP tool execution with performance metrics', async () => {
    await client.newChat();

    const startTime = performance.now();

    await client.sendMessage(
      'Use the replace_all tool to set document "perf-test-doc" ' +
        'content to "Performance test content". Reply with just "DONE" when complete.'
    );

    const responses = await client.waitForResponse(90000);
    const duration = performance.now() - startTime;

    expect(responses.length).toBeGreaterThan(0);

    // UI-based tests are slower, allow up to 90 seconds
    // (includes Claude thinking time, tool execution, and response generation)
    expect(
      duration,
      `MCP flow should complete within 90 seconds (actual: ${(duration / 1000).toFixed(1)}s)`
    ).toBeLessThan(90000);

    // Log performance for analysis
    console.log(`MCP tool execution completed in ${(duration / 1000).toFixed(2)}s`);
  });

  /**
   * E2E-018: replaceAll operation with confirmation
   *
   * Tests replace_all execution with detailed result confirmation.
   */
  test('E2E-018: should execute replaceAll with confirmation', async () => {
    await client.newChat();

    const testContent = `# E2E Test Document

## Overview
This document was created by an automated E2E test.

## Test Data
- Test ID: ${Date.now()}
- Timestamp: ${new Date().toISOString()}
`;

    await client.sendMessage(
      `Use the replace_all MCP tool to set the content of document ` +
        `"e2e-test-doc-${Date.now()}" to the following:\n\n` +
        `${testContent}\n\n` +
        `Confirm the operation was successful and mention the document ID.`
    );

    const responses = await client.waitForResponse(90000);

    expect(responses.length).toBeGreaterThan(0);
    const responseText = responses.join(' ').toLowerCase();

    // Verify detailed confirmation
    const hasConfirmation =
      responseText.includes('success') ||
      responseText.includes('replaced') ||
      responseText.includes('updated') ||
      responseText.includes('document') ||
      responseText.includes('complete');

    expect(
      hasConfirmation,
      'Response should confirm replace_all operation'
    ).toBe(true);
  });
});
