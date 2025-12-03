import { chromium, Browser, BrowserContext, Page } from 'playwright';

export interface IClaudeDesktopConfig {
  cdpUrl?: string;
  timeout?: number;
}

export class ClaudeDesktopClient {
  private browser: Browser | null = null;
  private context: BrowserContext | null = null;
  private page: Page | null = null;
  private config: IClaudeDesktopConfig;

  constructor(config: IClaudeDesktopConfig = {}) {
    this.config = {
      cdpUrl: config.cdpUrl || 'http://localhost:9222',
      timeout: config.timeout || 30000,
    };
  }

  async connect(): Promise<boolean> {
    try {
      this.browser = await chromium.connectOverCDP(this.config.cdpUrl!);
      const contexts = this.browser.contexts();

      if (contexts.length === 0) {
        throw new Error('No browser contexts found');
      }

      this.context = contexts[0];
      const pages = this.context.pages();

      // Find the main Claude.ai page
      this.page = pages.find(p => p.url().includes('claude.ai')) || pages[0];

      if (!this.page) {
        throw new Error('No Claude page found');
      }

      return true;
    } catch (error) {
      console.error('Failed to connect to Claude Desktop:', error);
      return false;
    }
  }

  async sendMessage(message: string): Promise<void> {
    if (!this.page) {
      throw new Error('Not connected to Claude Desktop');
    }

    // Find the input field and type the message
    const inputSelector = '[contenteditable="true"]';
    await this.page.waitForSelector(inputSelector, {
      timeout: this.config.timeout
    });

    await this.page.click(inputSelector);
    await this.page.fill(inputSelector, message);

    // Press Enter to send
    await this.page.keyboard.press('Enter');
  }

  async waitForResponse(timeout?: number): Promise<string[]> {
    if (!this.page) {
      throw new Error('Not connected to Claude Desktop');
    }

    const waitTime = timeout || this.config.timeout!;
    const startTime = Date.now();

    // Wait for loading spinner to disappear (Claude is done responding)
    // The spinner has aria-label containing "Loading" or is an SVG animation
    while (Date.now() - startTime < waitTime) {
      const isLoading = await this.page.evaluate(() => {
        // Check for loading indicators
        const spinner = document.querySelector('[aria-label*="Loading"], [aria-label*="loading"]');
        const svgSpinner = document.querySelector('svg[class*="animate"]');
        const stopButton = document.querySelector('button[aria-label*="Stop"]');
        return !!(spinner || svgSpinner || stopButton);
      });

      if (!isLoading) {
        // Double-check after a brief pause
        await this.page.waitForTimeout(1000);
        const stillLoading = await this.page.evaluate(() => {
          const spinner = document.querySelector('[aria-label*="Loading"], [aria-label*="loading"]');
          const svgSpinner = document.querySelector('svg[class*="animate"]');
          const stopButton = document.querySelector('button[aria-label*="Stop"]');
          return !!(spinner || svgSpinner || stopButton);
        });
        if (!stillLoading) {
          break;
        }
      }

      await this.page.waitForTimeout(1000);
    }

    // Get the response content
    const response = await this.getLastResponse();
    return response ? [response] : [];
  }

  async getLastResponse(): Promise<string> {
    if (!this.page) {
      throw new Error('Not connected to Claude Desktop');
    }

    // Find Claude's response by looking for the message content area
    const response = await this.page.evaluate(() => {
      // Try to find response containers - Claude Desktop uses specific structure
      // Look for paragraphs and lists that contain the actual response
      const mainContent = document.querySelector('[class*="prose"]') ||
                         document.querySelector('main') ||
                         document.body;

      // Get all paragraphs and list items from the response area
      const responseElements = mainContent.querySelectorAll('p, li, h1, h2, h3, strong');
      const responseTexts: string[] = [];

      // UI elements to exclude
      const excludeTexts = [
        'How can I help you today?',
        'Claude can make mistakes',
        'double-check responses',
        'Send message',
        'Retry',
        'Copy',
        'Edit',
        'Give positive feedback',
        'Give negative feedback',
        'Chats',
        'Projects',
        'Artifacts',
        'Recents',
        'Hide',
        'Open sidebar',
        'Scroll to bottom',
      ];

      responseElements.forEach(el => {
        const text = el.textContent?.trim();
        if (text && text.length > 0) {
          // Skip UI elements
          const isUI = excludeTexts.some(exclude =>
            text === exclude || text.startsWith(exclude)
          );
          if (!isUI) {
            responseTexts.push(text);
          }
        }
      });

      // Join all response text
      return responseTexts.join('\n');
    });

    return response;
  }

  async newChat(): Promise<void> {
    if (!this.page) {
      throw new Error('Not connected to Claude Desktop');
    }

    // Navigate directly to /new instead of clicking
    await this.page.goto('https://claude.ai/new', {
      waitUntil: 'domcontentloaded',
      timeout: 10000
    });

    // Wait for the page to be ready
    await this.page.waitForSelector('[contenteditable="true"]', {
      timeout: this.config.timeout
    });
  }

  async disconnect(): Promise<void> {
    // Don't close the browser - just disconnect from CDP
    // This keeps Claude Desktop running
    if (this.browser) {
      await this.browser.close();
      this.browser = null;
      this.context = null;
      this.page = null;
    }
  }

  isConnected(): boolean {
    return this.page !== null;
  }
}
