/**
 * Claude.ai Authentication Manager
 *
 * Manages Claude.ai session state for E2E tests:
 * - Reuses cached auth state from CLAUDE_AUTH_STATE env var
 * - Falls back to email login via Mailosaur when state is invalid
 * - Persists new auth state back to .env.test file
 */
import * as fs from 'fs';

import { BrowserContext, Page } from '@playwright/test';
import MailosaurClient from 'mailosaur';

import { IClaudeAuthConfig, IAuthResult } from './claude-auth-interfaces';

export const CLAUDE_URL = 'https://claude.ai';

/**
 * Claude Authentication Manager
 */
export class ClaudeAuth {
  private config: IClaudeAuthConfig;
  private mailosaur: MailosaurClient;

  constructor(config: IClaudeAuthConfig) {
    this.config = config;
    this.mailosaur = new MailosaurClient(config.mailosaurApiKey);
  }

  /**
   * Parse stored auth state from base64 string
   */
  private parseAuthState(): object | null {
    if (!this.config.authState) {
      return null;
    }
    const decoded = Buffer.from(this.config.authState, 'base64').toString('utf-8');
    return JSON.parse(decoded);
  }

  /**
   * Apply stored auth state to browser context
   * Only applies cookies - localStorage is not needed for authentication
   */
  async applyAuthState(context: BrowserContext): Promise<boolean> {
    const state = this.parseAuthState();
    if (!state) {
      return false;
    }

    const { cookies } = state as {
      cookies: Array<{
        name: string;
        value: string;
        domain: string;
        path: string;
        expires?: number;
        httpOnly?: boolean;
        secure?: boolean;
        sameSite?: 'Strict' | 'Lax' | 'None';
      }>;
    };

    if (cookies?.length > 0) {
      // Filter to only claude.ai cookies to avoid issues with other domains
      const claudeCookies = cookies.filter(
        (c) => c.domain.includes('claude.ai') || c.domain.includes('anthropic')
      );
      await context.addCookies(claudeCookies);
    }

    return true;
  }

  /**
   * Validate if the current session is authenticated
   */
  async validateSession(page: Page): Promise<boolean> {
    try {
      await page.goto(`${CLAUDE_URL}/new`, { waitUntil: 'domcontentloaded', timeout: 15000 });
    } catch {
      return false;
    }

    const url = page.url();
    if (url.includes('/login') || url.includes('challenge_redirect')) {
      return false;
    }

    // Wait a bit for any client-side redirects
    await page.waitForTimeout(2000);
    const finalUrl = page.url();

    if (finalUrl.includes('/login') || finalUrl.includes('challenge_redirect')) {
      return false;
    }

    const chatInput = page.locator(
      '[data-testid="chat-input"], textarea[placeholder*="Claude"], [contenteditable="true"]'
    );
    return chatInput.first().isVisible().catch(() => false);
  }

  /**
   * Perform email login via Mailosaur
   */
  async performEmailLogin(page: Page): Promise<void> {
    await page.goto(`${CLAUDE_URL}/login`, { waitUntil: 'networkidle' });

    const emailInput = page.getByRole('textbox', { name: 'Email' });
    await emailInput.waitFor({ timeout: 15000 });
    await emailInput.fill(this.config.claudeEmail);

    const continueButton = page.getByRole('button', { name: 'Continue with email' });
    await continueButton.click();
    await page.waitForTimeout(2000);

    const message = await this.mailosaur.messages.get(
      this.config.mailosaurServerId,
      { sentTo: this.config.claudeEmail },
      { timeout: 60000 }
    );

    const magicLink = message.html?.links?.find((link) =>
      link.href?.includes('claude.ai')
    );

    if (!magicLink?.href) {
      throw new Error('Magic link not found in email');
    }

    await page.goto(magicLink.href, { waitUntil: 'networkidle', timeout: 30000 });
    await page.waitForURL(/claude\.ai\/(new|chat)?/, { timeout: 30000 });
  }

  /**
   * Apply stealth evasions to avoid bot detection
   */
  private async applyStealthEvasions(context: BrowserContext): Promise<void> {
    await context.addInitScript(() => {
      // Hide webdriver
      Object.defineProperty(navigator, 'webdriver', { get: () => undefined });

      // Override plugins
      Object.defineProperty(navigator, 'plugins', {
        get: () => [1, 2, 3, 4, 5],
      });

      // Override languages
      Object.defineProperty(navigator, 'languages', {
        get: () => ['en-US', 'en'],
      });

      // Chrome detection
      (window as Window & { chrome?: object }).chrome = { runtime: {} };

      // Override permissions
      const originalQuery = window.navigator.permissions.query;
      window.navigator.permissions.query = (parameters: PermissionDescriptor) =>
        parameters.name === 'notifications'
          ? Promise.resolve({ state: 'denied' } as PermissionStatus)
          : originalQuery(parameters);
    });
  }

  /**
   * Main authentication flow
   */
  async authenticate(context: BrowserContext): Promise<IAuthResult> {
    let isNewLogin = false;

    // Apply stealth evasions before any navigation
    await this.applyStealthEvasions(context);

    const hasStoredState = await this.applyAuthState(context);
    const page = await context.newPage();

    try {
      const isValid = hasStoredState && (await this.validateSession(page));

      if (!isValid) {
        // Check if blocked by Cloudflare
        const currentUrl = page.url();
        if (currentUrl.includes('challenge_redirect')) {
          throw new Error('Cloudflare challenge detected - automated login blocked.');
        }

        await this.performEmailLogin(page);
        isNewLogin = true;
      }

      const state = await context.storageState();
      const authState = Buffer.from(JSON.stringify(state)).toString('base64');

      if (isNewLogin) {
        this.persistAuthState(authState);
      }

      return { success: true, isNewLogin, authState };
    } catch (error) {
      return {
        success: false,
        isNewLogin,
        error: error instanceof Error ? error.message : 'Authentication failed',
      };
    } finally {
      await page.close();
    }
  }

  /**
   * Persist auth state to .env.test file
   * Only called when envFilePath is configured
   */
  private persistAuthState(authState: string): void {
    const envFilePath = this.config.envFilePath;
    if (!envFilePath) {
      return; // Persistence not configured - skip by design
    }

    const content = fs.readFileSync(envFilePath, 'utf-8');
    const regex = /^CLAUDE_AUTH_STATE=.*$/m;

    if (!regex.test(content)) {
      throw new Error('CLAUDE_AUTH_STATE= line not found in .env.test file');
    }

    const updated = content.replace(regex, `CLAUDE_AUTH_STATE=${authState}`);
    fs.writeFileSync(envFilePath, updated, 'utf-8');
  }
}
