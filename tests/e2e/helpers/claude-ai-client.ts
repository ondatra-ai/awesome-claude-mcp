/**
 * Claude.ai Browser Client
 *
 * Provides programmatic access to Claude.ai via browser automation.
 */
import { BrowserContext, Page } from '@playwright/test';

import { IClaudeAuthConfig, IAuthResult } from './claude-auth-interfaces';
import { ClaudeAuth, CLAUDE_URL } from './claude-auth';

/**
 * Claude.ai Browser Client
 */
export class ClaudeAIClient {
  private context: BrowserContext;
  private page: Page | null = null;
  private auth: ClaudeAuth;
  private authenticated = false;

  constructor(context: BrowserContext, config: IClaudeAuthConfig) {
    this.context = context;
    this.auth = new ClaudeAuth(config);
  }

  async initialize(): Promise<IAuthResult> {
    const result = await this.auth.authenticate(this.context);
    this.authenticated = result.success;

    if (result.success) {
      this.page = await this.context.newPage();
      await this.page.goto(`${CLAUDE_URL}/new`, { waitUntil: 'networkidle' });
    }

    return result;
  }

  async newChat(): Promise<void> {
    if (!this.page || !this.authenticated) {
      throw new Error('Client not initialized');
    }
    await this.page.goto(`${CLAUDE_URL}/new`, { waitUntil: 'networkidle' });
  }

  async prompt(message: string, timeout = 120000): Promise<string[]> {
    if (!this.page || !this.authenticated) {
      throw new Error('Client not initialized');
    }

    const chatInput = this.page
      .locator(
        'div[contenteditable="true"][data-placeholder], ' +
          'textarea[placeholder*="Claude"], ' +
          'div[contenteditable="true"]'
      )
      .first();

    await chatInput.waitFor({ timeout: 10000 });
    await chatInput.click();
    await chatInput.fill(message);

    const sendButton = this.page
      .getByRole('button', { name: /send/i })
      .or(this.page.locator('button[aria-label*="Send"]'));

    await sendButton.waitFor({ state: 'visible', timeout: 5000 });
    await this.page.waitForTimeout(500);
    await sendButton.click();
    await this.page.waitForTimeout(1000);

    return this.waitForResponse(timeout);
  }

  private async waitForResponse(timeout: number): Promise<string[]> {
    if (!this.page) throw new Error('No page');

    const startTime = Date.now();
    const stopButton = this.page.locator('button[aria-label*="Stop"], button:has-text("Stop")');

    while (Date.now() - startTime < timeout) {
      const isStreaming = await stopButton.isVisible().catch(() => false);
      if (!isStreaming) {
        await this.page.waitForTimeout(500);
        if (!(await stopButton.isVisible().catch(() => false))) break;
      }
      await this.page.waitForTimeout(500);
    }

    return this.extractResponses();
  }

  private async extractResponses(): Promise<string[]> {
    if (!this.page) return [];

    return this.page.evaluate(() => {
      const selectors = [
        '[data-testid="assistant-message"]',
        '[data-message-author-role="assistant"]',
        '.font-claude-message',
        '[class*="assistant"]',
      ];

      for (const selector of selectors) {
        const elements = document.querySelectorAll(selector);
        const results = Array.from(elements)
          .map((el) => el.textContent?.trim())
          .filter((text): text is string => !!text);
        if (results.length > 0) return results;
      }

      const allMessages = document.querySelectorAll('[data-testid*="message"]');
      const lastMessage = allMessages[allMessages.length - 1];
      const text = lastMessage?.textContent?.trim();
      return text ? [text] : [];
    });
  }

  async getConversation(): Promise<Array<{ role: string; content: string }>> {
    if (!this.page) return [];

    return this.page.evaluate(() => {
      const elements = document.querySelectorAll(
        '[data-testid*="message"], [data-message-author-role]'
      );

      return Array.from(elements)
        .map((el) => ({
          role:
            el.getAttribute('data-message-author-role') ||
            (el.querySelector('[data-testid="user-message"]') ? 'user' : 'assistant'),
          content: el.textContent?.trim() || '',
        }))
        .filter((msg) => msg.content);
    });
  }

  async getAuthState(): Promise<string> {
    const state = await this.context.storageState();
    return Buffer.from(JSON.stringify(state)).toString('base64');
  }

  async close(): Promise<void> {
    if (this.page) {
      await this.page.close();
      this.page = null;
    }
    this.authenticated = false;
  }

  isReady(): boolean {
    return this.authenticated && this.page !== null;
  }
}
