import { test, expect } from '@playwright/test';

import { ClaudeDesktopClient } from './helpers/claude-client';

/**
 * Claude Desktop E2E Tests via CDP
 *
 * Prerequisites:
 * 1. Launch Claude Desktop with remote debugging:
 *    pkill -x "Claude" && open -a Claude --args --remote-debugging-port=9222
 *
 * 2. Run tests:
 *    npx playwright test e2e/claude.spec.ts
 */

// Run tests serially - Claude Desktop is a single instance
test.describe.configure({ mode: 'serial' });

// Increase timeout for Claude Desktop tests (responses can take time)
test.setTimeout(120000);

test.describe('Claude Desktop via CDP', () => {
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

  test('should connect to Claude Desktop', async () => {
    expect(client.isConnected()).toBe(true);
  });

  test('should send a prompt and receive response', async () => {
    await client.newChat();
    await client.sendMessage('Reply with exactly: "CDP_TEST_SUCCESS"');

    const responses = await client.waitForResponse(60000);

    expect(responses.length).toBeGreaterThan(0);
    expect(responses.join(' ')).toContain('CDP_TEST_SUCCESS');
  });

  test('should maintain conversation context', async () => {
    await client.newChat();
    await client.sendMessage(
      'Remember this code: DELTA-42. Reply with "Code stored."'
    );
    await client.waitForResponse(30000);

    await client.sendMessage('What code did I tell you to remember?');
    const responses = await client.waitForResponse(30000);

    expect(responses.join(' ')).toContain('DELTA-42');
  });

  test('should take screenshot for debugging', async ({ }, testInfo) => {
    const screenshotPath = testInfo.outputPath('claude-desktop-test.png');
    await client.takeScreenshot(screenshotPath);
    // Verify screenshot was created
    const fs = await import('fs');
    expect(fs.existsSync(screenshotPath)).toBe(true);
  });
});

/**
 * MCP Integration Tests
 *
 * These tests verify MCP server functionality through Claude Desktop.
 * Requires MCP servers to be configured in Claude Desktop.
 */
test.describe('MCP Integration via Claude Desktop', () => {
  let client: ClaudeDesktopClient;

  test.beforeAll(async () => {
    client = new ClaudeDesktopClient({
      cdpUrl: 'http://localhost:9222',
      timeout: 60000,
    });

    const connected = await client.connect();
    test.skip(!connected, 'Claude Desktop not available');
  });

  test.afterAll(async () => {
    await client?.disconnect();
  });

  test.skip('should list available MCP tools', async () => {
    await client.newChat();
    await client.sendMessage('What MCP tools do you have available?');

    const responses = await client.waitForResponse(30000);
    expect(responses.length).toBeGreaterThan(0);
    // Add specific tool assertions based on your MCP server
  });

  test.skip('should execute MCP tool', async () => {
    await client.newChat();
    // Replace with your actual MCP tool invocation
    await client.sendMessage('Use the [your-tool] to [action]');

    const responses = await client.waitForResponse(60000);
    expect(responses.length).toBeGreaterThan(0);
    // Add specific result assertions
  });
});
