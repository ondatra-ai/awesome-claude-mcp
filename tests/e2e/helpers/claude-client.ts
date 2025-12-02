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

    // Simple approach: get all text from body and parse it
    const response = await this.page.evaluate(() => {
      // Get the full text content from body
      const fullText = document.body.innerText;

      // Split by newlines and filter
      const lines = fullText.split('\n')
        .map(l => l.trim())
        .filter(l => l.length > 0);

      // Find response lines (exclude UI elements)
      const excludePatterns = [
        /^How can I help/,
        /^Thinking about/,
        /^(Opus|Sonnet|Haiku|Claude)\s/,
        /double-check responses/,
        /^Reply with/,
        /^Retry$/,
        /^\+$/,
        /^Send message$/,
        /^Remember this code/,
        /^What code did/,
        /^\d+ step/,
        /^Added memory/,
        /^\d+ result/,
        /^The memory has been/,
        /^Chats$/,
        /^Projects$/,
        /^Artifacts$/,
        /^Recents$/,
        /^Hide$/,
      ];

      const responseLines = lines.filter(line => {
        for (const pattern of excludePatterns) {
          if (pattern.test(line)) return false;
        }
        return line.length >= 3;
      });

      // Find the actual response - usually after thinking sections
      // Look for lines that could be Claude's direct response
      const candidates = responseLines.filter(line =>
        !line.includes('I should') &&
        !line.includes('successfully stored') &&
        line.length < 500
      );

      return candidates[candidates.length - 1] || '';
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
