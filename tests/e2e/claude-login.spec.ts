import * as path from 'path';

import dotenv from 'dotenv';
import { test, expect } from '@playwright/test';

import { IAuthResult } from './helpers/claude-auth-interfaces';
import { ClaudeAIClient } from './helpers/claude-ai-client';

const ENV_FILE_PATH = path.join(process.cwd(), '.env.test');
dotenv.config({ path: ENV_FILE_PATH });

/**
 * Claude.ai Login E2E Tests
 *
 * These tests are excluded from default test runs due to external dependencies.
 * Run explicitly with: npx playwright test e2e/claude-login.spec.ts
 *
 * Required environment variables:
 *   - MAILOSAUR_API_KEY
 *   - MAILOSAUR_SERVER_ID
 *   - CLAUDE_EMAIL
 *   - CLAUDE_AUTH_STATE (optional - cached auth state)
 */

const MAILOSAUR_API_KEY = process.env.MAILOSAUR_API_KEY;
const MAILOSAUR_SERVER_ID = process.env.MAILOSAUR_SERVER_ID;
const CLAUDE_EMAIL = process.env.CLAUDE_EMAIL;
const CLAUDE_AUTH_STATE = process.env.CLAUDE_AUTH_STATE;

// Fail fast if required env vars are missing
if (!MAILOSAUR_API_KEY || !MAILOSAUR_SERVER_ID || !CLAUDE_EMAIL) {
  throw new Error(
    'Missing required env vars: MAILOSAUR_API_KEY, MAILOSAUR_SERVER_ID, CLAUDE_EMAIL. ' +
      'Copy tests/.env.test.example to tests/.env.test and fill in values.'
  );
}

test.describe('Claude.ai Authentication with ClaudeAIClient', () => {

  let client: ClaudeAIClient;
  let authResult: IAuthResult;

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext();
    client = new ClaudeAIClient(context, {
      mailosaurApiKey: MAILOSAUR_API_KEY as string,
      mailosaurServerId: MAILOSAUR_SERVER_ID as string,
      claudeEmail: CLAUDE_EMAIL as string,
      authState: CLAUDE_AUTH_STATE,
      envFilePath: ENV_FILE_PATH,
    });

    authResult = await client.initialize();

    if (authResult.success) {
      const newState = await client.getAuthState();
      console.log(
        authResult.isNewLogin ? 'Fresh login via Mailosaur' : 'Reused cached auth state'
      );
      console.log(`Auth state: ${newState.length} chars`);
    }
  });

  test.afterAll(async () => {
    await client?.close();
  });

  test('should authenticate successfully', async () => {
    expect(authResult.success).toBe(true);
    expect(client.isReady()).toBe(true);
  });

  test('should be able to send a prompt', async () => {
    test.skip(!authResult.success, 'Authentication failed');

    const responses = await client.prompt('Reply with exactly: "TEST_SUCCESS_123"', 60000);

    expect(responses.length).toBeGreaterThan(0);
    expect(responses.join(' ')).toContain('TEST_SUCCESS');
  });

  test('should maintain conversation context', async () => {
    test.skip(!authResult.success, 'Authentication failed');

    await client.newChat();
    await client.prompt('Remember: The secret code is ALPHA-7. Reply with "Remembered."');

    const responses = await client.prompt('What is the secret code I told you?');

    expect(responses.length).toBeGreaterThan(0);
    expect(responses.join(' ')).toContain('ALPHA-7');
  });

  test('should handle new chat creation', async () => {
    test.skip(!authResult.success, 'Authentication failed');

    await client.prompt('Hello');
    await client.newChat();

    const conversation = await client.getConversation();
    expect(conversation.length).toBe(0);
  });
});

/**
 * UI Element Verification Tests (no auth required)
 */
test.describe('Claude.ai Login Page UI', () => {
  test('should show Google OAuth option', async ({ page }) => {
    await page.goto('https://claude.ai/login');
    await expect(page.getByRole('button', { name: /continue with google/i })).toBeVisible({
      timeout: 15000,
    });
  });

  test('should show SSO option', async ({ page }) => {
    await page.goto('https://claude.ai/login');
    await expect(page.getByRole('button', { name: /continue with sso/i })).toBeVisible({
      timeout: 15000,
    });
  });

  test('should show email input', async ({ page }) => {
    await page.goto('https://claude.ai/login');
    await expect(page.getByRole('textbox', { name: 'Email' })).toBeVisible({ timeout: 15000 });
  });

  test('should show continue with email button', async ({ page }) => {
    await page.goto('https://claude.ai/login');
    await expect(page.getByRole('button', { name: 'Continue with email' })).toBeVisible({
      timeout: 15000,
    });
  });
});

/**
 * Auth State Management Tests
 */
test.describe('Auth State Management', () => {
  test('should save auth state after successful login', async ({ browser }) => {
    const context = await browser.newContext();
    const client = new ClaudeAIClient(context, {
      mailosaurApiKey: MAILOSAUR_API_KEY as string,
      mailosaurServerId: MAILOSAUR_SERVER_ID as string,
      claudeEmail: CLAUDE_EMAIL as string,
      authState: CLAUDE_AUTH_STATE,
      envFilePath: ENV_FILE_PATH,
    });

    const result = await client.initialize();
    expect(result.success).toBe(true);

    const authState = await client.getAuthState();
    expect(authState.length).toBeGreaterThan(100);

    const decoded = Buffer.from(authState, 'base64').toString('utf-8');
    const parsed = JSON.parse(decoded);
    expect(parsed).toHaveProperty('cookies');
    expect(parsed).toHaveProperty('origins');

    await client.close();
    await context.close();
  });
});
