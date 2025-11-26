/**
 * MCP Client Helper for E2E Tests
 *
 * Uses @modelcontextprotocol/sdk for MCP server communication.
 * Tests MCP protocol compliance without LLM integration.
 */
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import {
  StreamableHTTPClientTransport,
} from '@modelcontextprotocol/sdk/client/streamableHttp.js';
import {
  SSEClientTransport,
} from '@modelcontextprotocol/sdk/client/sse.js';

export interface IMcpClientOptions {
  name?: string;
  version?: string;
}

export interface IMcpTool {
  name: string;
  description?: string;
  inputSchema?: Record<string, unknown>;
}

export interface IToolResult {
  content: Array<{ type: string; text?: string }>;
  isError?: boolean;
}

/**
 * MCP Client for E2E testing.
 */
export class McpClient {
  private client: Client;
  private baseUrl: string;
  private connected: boolean = false;

  constructor(baseUrl: string, options: IMcpClientOptions = {}) {
    this.baseUrl = baseUrl;
    this.client = new Client({
      name: options.name || 'e2e-test-client',
      version: options.version || '1.0.0',
    });
  }

  /**
   * Connect to the MCP server.
   */
  async connect(): Promise<void> {
    if (this.connected) {
      return;
    }

    const mcpUrl = new URL('/mcp', this.baseUrl);

    try {
      const transport = new StreamableHTTPClientTransport(mcpUrl);
      await this.client.connect(transport);
      this.connected = true;
    } catch {
      const sseTransport = new SSEClientTransport(mcpUrl);
      await this.client.connect(sseTransport);
      this.connected = true;
    }
  }

  /**
   * List available tools from the MCP server.
   */
  async listTools(): Promise<IMcpTool[]> {
    if (!this.connected) {
      await this.connect();
    }

    const result = await this.client.listTools();
    return result.tools as IMcpTool[];
  }

  /**
   * Call a tool on the MCP server.
   */
  async callTool(
    name: string,
    args: Record<string, unknown> = {}
  ): Promise<IToolResult> {
    if (!this.connected) {
      await this.connect();
    }

    const result = await this.client.callTool({
      name,
      arguments: args,
    });

    return result as IToolResult;
  }

  /**
   * Close the connection.
   */
  async close(): Promise<void> {
    if (this.connected) {
      await this.client.close();
      this.connected = false;
    }
  }

  /**
   * Check if connected.
   */
  isConnected(): boolean {
    return this.connected;
  }
}

/**
 * Create multiple MCP clients for concurrent testing.
 */
export async function createMultipleClients(
  baseUrl: string,
  count: number
): Promise<McpClient[]> {
  const clients: McpClient[] = [];

  for (let i = 0; i < count; i++) {
    const client = new McpClient(baseUrl, {
      name: `e2e-test-client-${i}`,
      version: '1.0.0',
    });
    await client.connect();
    clients.push(client);
  }

  return clients;
}

/**
 * Close all clients.
 */
export async function closeAllClients(clients: McpClient[]): Promise<void> {
  await Promise.all(clients.map((client) => client.close()));
}

/**
 * Extract text content from tool result.
 */
export function getResultText(result: IToolResult): string {
  if (result.content && Array.isArray(result.content)) {
    const texts = result.content
      .filter((c) => c.type === 'text' && c.text)
      .map((c) => c.text);
    return texts.join('\n');
  }
  return JSON.stringify(result);
}
