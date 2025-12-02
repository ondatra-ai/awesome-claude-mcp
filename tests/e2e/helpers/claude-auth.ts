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
   */
  async applyAuthState(context: BrowserContext): Promise<boolean> {
    const state = this.parseAuthState();
    if (!state) {
      return false;
    }

    const { cookies, origins } = state as {
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
      origins: Array<{
        origin: string;
        localStorage: Array<{ name: string; value: string }>;
      }>;
    };

    if (cookies?.length > 0) {
      await context.addCookies(cookies);
    }

    if (origins?.length > 0) {
      const page = await context.newPage();
      for (const origin of origins) {
        if (origin.origin.includes('claude.ai')) {
          await page.goto(origin.origin, { waitUntil: 'domcontentloaded' });
          for (const item of origin.localStorage) {
            await page.evaluate(
              ({ key, value }) => localStorage.setItem(key, value),
              { key: item.name, value: item.value }
            );
          }
        }
      }
      await page.close();
    }

    return true;
  }

  /**
   * Validate if the current session is authenticated
   */
  async validateSession(page: Page): Promise<boolean> {
    await page.goto(`${CLAUDE_URL}/new`, { waitUntil: 'networkidle', timeout: 30000 });

    const url = page.url();
    if (url.includes('/login') || url.includes('challenge_redirect')) {
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
   * Main authentication flow
   */
  async authenticate(context: BrowserContext): Promise<IAuthResult> {
    let isNewLogin = false;
    const hasStoredState = await this.applyAuthState(context);
    const page = await context.newPage();

    try {
      let isValid = hasStoredState && (await this.validateSession(page));

      if (!isValid) {
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
