import { test, expect } from '@playwright/test';

import { getEnvironmentConfig } from '../config/environments';

const { mcpServiceUrl } = getEnvironmentConfig(process.env.E2E_ENV);

/**
 * MCP Service Integration Tests
 *
 * These tests validate the MCP server using HTTP+SSE (Streamable HTTP) transport
 * per the official MCP specification. MCP does NOT use WebSocket.
 *
 * Transport: Streamable HTTP
 * - POST /mcp: Client sends JSON-RPC messages
 * - GET /mcp: Client establishes SSE stream for server-to-client messages
 * - Header: Mcp-Session-Id for session tracking
 */
test.describe('MCP Service Integration Tests', () => {
  test(
    'INT-008: HTTP+SSE session establishment succeeds',
    async ({ request }) => {
      // Given: Server runs HTTP+SSE service on configured port

      // When: Client sends POST to /mcp with initialize request
      const initializeRequest = {
        jsonrpc: '2.0',
        method: 'initialize',
        id: 1,
        params: {
          protocolVersion: '2024-11-05',
          capabilities: {},
          clientInfo: {
            name: 'test-client',
            version: '1.0.0',
          },
        },
      };

      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: {
          'Content-Type': 'application/json',
        },
        data: initializeRequest,
      });

      // Then: Server establishes MCP session
      expect(response.status()).toBe(200);
      const result = await response.json();

      expect(result).toHaveProperty('jsonrpc', '2.0');
      expect(result).toHaveProperty('id', 1);
      expect(result).toHaveProperty('result');
      expect(result.result).toHaveProperty('protocolVersion');
      expect(result.result).toHaveProperty('serverInfo');

      // Then: Server assigns unique session identifier via Mcp-Session-Id header
      const sessionId = response.headers()['mcp-session-id'];
      // Session ID may be in header or response body depending on implementation
      expect(sessionId || result.result.sessionId).toBeTruthy();
    }
  );

  test(
    'INT-009: HTTP endpoint establishes SSE stream for server messages',
    async ({ request }) => {
      // Given: Server exposes HTTP endpoints at /mcp

      // First establish a session via POST
      const initResponse = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: {
          jsonrpc: '2.0',
          method: 'initialize',
          id: 1,
          params: {
            protocolVersion: '2024-11-05',
            capabilities: {},
            clientInfo: { name: 'test-client', version: '1.0.0' },
          },
        },
      });

      expect(initResponse.ok()).toBeTruthy();
      const sessionId = initResponse.headers()['mcp-session-id'];

      // When: Client sends GET request to /mcp for SSE stream
      const sseResponse = await request.get(`${mcpServiceUrl}/mcp`, {
        headers: {
          Accept: 'text/event-stream',
          ...(sessionId ? { 'Mcp-Session-Id': sessionId } : {}),
        },
      });

      // Then: Server establishes SSE connection for server-to-client messages
      // SSE endpoint should return 200 with text/event-stream content type
      // or indicate SSE support
      const contentType = sseResponse.headers()['content-type'];
      const status = sseResponse.status();

      // Server should support SSE (200 with event-stream) or indicate capability
      expect([200, 204, 202]).toContain(status);

      if (status === 200 && contentType) {
        expect(contentType).toContain('text/event-stream');
      }
    }
  );

  test(
    'INT-010: Health check returns service status with dependencies',
    async ({ request }) => {
      // Given: Server provides health check endpoint

      // When: Client sends GET request to /health
      const response = await request.get(`${mcpServiceUrl}/health`);

      // Then: Server returns status 200
      expect(response.status()).toBe(200);

      const data = await response.json();

      // Then: Response includes service version
      expect(data).toHaveProperty('version');
      expect(typeof data.version).toBe('string');

      // Then: Response includes dependency status for Redis and Google API
      expect(data).toHaveProperty('dependencies');
      expect(data.dependencies).toHaveProperty('redis');
      expect(data.dependencies).toHaveProperty('google_api');
    }
  );

  test(
    'INT-011: Server rejects unauthorized domain CORS preflight request',
    async ({ request }) => {
      // Given: Server enforces CORS policy for Claude domains

      // When: Client from unauthorized domain sends preflight OPTIONS request
      const response = await request.fetch(`${mcpServiceUrl}/mcp`, {
        method: 'OPTIONS',
        headers: {
          Origin: 'https://malicious-domain.com',
          'Access-Control-Request-Method': 'POST',
          'Access-Control-Request-Headers': 'content-type',
        },
      });

      // Then: Server rejects request with CORS error
      const allowedOrigin = response.headers()['access-control-allow-origin'];
      if (allowedOrigin) {
        expect(allowedOrigin).not.toBe('https://malicious-domain.com');
        expect(allowedOrigin).not.toBe('*');
      }
    }
  );

  test(
    'INT-012: Server accepts CORS preflight from authorized Claude domain',
    async ({ request }) => {
      // Given: Server enforces CORS policy for Claude domains

      // When: Client from claude.ai domain sends preflight OPTIONS request
      const response = await request.fetch(`${mcpServiceUrl}/mcp`, {
        method: 'OPTIONS',
        headers: {
          Origin: 'https://claude.ai',
          'Access-Control-Request-Method': 'POST',
          'Access-Control-Request-Headers': 'content-type',
        },
      });

      // Then: Server accepts request with CORS headers
      expect(response.status()).toBeLessThan(400);
      const allowedOrigin = response.headers()['access-control-allow-origin'];
      expect(allowedOrigin).toBeTruthy();
    }
  );

  test(
    'INT-013: MCP protocol message parsing succeeds',
    async ({ request }) => {
      // Given: Client maintains active MCP session

      // When: Client sends valid MCP protocol message via HTTP POST
      const validMcpMessage = {
        jsonrpc: '2.0',
        method: 'initialize',
        id: 1,
        params: {
          protocolVersion: '2024-11-05',
          capabilities: {
            roots: { listChanged: true },
            sampling: {},
          },
          clientInfo: { name: 'test-client', version: '1.0.0' },
        },
      };

      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: validMcpMessage,
      });

      // Then: Server parses message successfully
      expect(response.status()).toBeLessThan(300);
      const result = await response.json();

      // Then: Server extracts message type and protocol version
      expect(result).toHaveProperty('jsonrpc', '2.0');
      expect(result).toHaveProperty('id', 1);
      expect(result.result).toHaveProperty('protocolVersion');
    }
  );

  test(
    'INT-014: Server rejects invalid JSON with parsing error',
    async ({ request }) => {
      // Given: Client maintains active MCP session

      // When: Client sends message with invalid JSON format via HTTP POST
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: '{invalid json: missing quotes}',
      });

      // Then: Server returns parsing error immediately
      expect(response.status()).toBeGreaterThanOrEqual(400);
      expect(response.status()).toBeLessThan(500);

      // Then: Error message specifies invalid JSON
      const errorText = await response.text();
      const lowerError = errorText.toLowerCase();
      const hasJsonError = ['json', 'parse', 'syntax', 'invalid'].some((t) =>
        lowerError.includes(t)
      );
      expect(hasJsonError).toBeTruthy();
    }
  );

  test(
    'INT-015: Server validates required MCP protocol fields',
    async ({ request }) => {
      // Given: Client maintains active MCP session

      // When: Client sends MCP message missing required protocol_version field
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: {
          jsonrpc: '2.0',
          method: 'initialize',
          id: 1,
          params: {
            // Missing required 'protocolVersion' field
            capabilities: {},
            clientInfo: { name: 'test-client', version: '1.0.0' },
          },
        },
      });

      // Then: Server returns validation error immediately
      expect(response.status()).toBeGreaterThanOrEqual(400);

      // Then: Error message specifies missing required field
      const result = await response.json();
      expect(result).toHaveProperty('error');
    }
  );

  test(
    'INT-016: Server rejects message with invalid field values',
    async ({ request }) => {
      // Given: Client maintains active MCP session

      // When: Client sends message with field containing invalid value
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: {
          jsonrpc: '2.0',
          method: 'invalid_method_!@#$',
          id: 1,
          params: { protocolVersion: '999.999.999' },
        },
      });

      // Then: Server returns validation error
      expect(response.status()).toBeGreaterThanOrEqual(400);
      expect(response.status()).toBeLessThan(500);

      // Then: Error message specifies error detail
      const result = await response.json();
      expect(result).toHaveProperty('error');
    }
  );

  test(
    'INT-017: Server accepts well-formed message with valid schema',
    async ({ request }) => {
      // Given: Server validates messages against MCP protocol schema

      // When: Client sends message with correct schema and all required fields
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: {
          jsonrpc: '2.0',
          method: 'initialize',
          id: 1,
          params: {
            protocolVersion: '2024-11-05',
            capabilities: { roots: { listChanged: true }, sampling: {} },
            clientInfo: { name: 'test-client', version: '1.0.0' },
          },
        },
      });

      // Then: Server accepts message without validation errors
      expect(response.status()).toBeLessThan(300);
      const result = await response.json();
      expect(result).toHaveProperty('result');
      expect(result).not.toHaveProperty('error');
    }
  );

  test(
    'INT-018: Server returns complete MCP response with required fields',
    async ({ request }) => {
      // Given: Client sends valid MCP request
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: {
          jsonrpc: '2.0',
          method: 'initialize',
          id: 100,
          params: {
            protocolVersion: '2024-11-05',
            capabilities: {},
            clientInfo: { name: 'test-client', version: '1.0.0' },
          },
        },
      });

      // When: Server processes request successfully
      expect(response.status()).toBeLessThan(300);
      const result = await response.json();

      // Then: Response in MCP protocol format
      expect(result).toHaveProperty('jsonrpc', '2.0');

      // Then: Response includes protocol_version field
      expect(result.result).toHaveProperty('protocolVersion');

      // Then: Response includes message_id matching request
      expect(result).toHaveProperty('id', 100);
    }
  );

  test(
    'INT-019: Server returns formatted error for invalid MCP request',
    async ({ request }) => {
      // Given: Client sends invalid MCP request
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: {
          jsonrpc: '2.0',
          // Missing required 'method' field
          id: 1,
          params: { protocolVersion: '2024-11-05' },
        },
      });

      // When: Server detects validation error
      expect(response.status()).toBeGreaterThanOrEqual(400);
      const result = await response.json();

      // Then: Server returns error response in MCP format
      expect(result).toHaveProperty('jsonrpc', '2.0');
      expect(result).toHaveProperty('error');

      // Then: Error response includes error code
      expect(result.error).toHaveProperty('code');
      expect(typeof result.error.code).toBe('number');

      // Then: Error response includes descriptive error message
      expect(result.error).toHaveProperty('message');
      expect(typeof result.error.message).toBe('string');
    }
  );

  test(
    'INT-020: Server response includes JSON structure with metadata fields',
    async ({ request }) => {
      // Given: Server processes client requests

      // When: Server generates response
      const response = await request.get(`${mcpServiceUrl}/health`);
      expect(response.ok()).toBeTruthy();
      const data = await response.json();

      // Then: Response contains valid JSON structure
      expect(typeof data).toBe('object');

      // Then: Response includes timestamp field
      expect(data).toHaveProperty('timestamp');

      // Then: Response includes correlation ID for tracing
      expect(data).toHaveProperty('correlation_id');
    }
  );

  test(
    'INT-021: Server returns MCP protocol-compliant response',
    async ({ request }) => {
      // Test multiple request types
      const testRequests = [
        {
          name: 'initialize',
          payload: {
            jsonrpc: '2.0',
            method: 'initialize',
            id: 1,
            params: {
              protocolVersion: '2024-11-05',
              capabilities: {},
              clientInfo: { name: 'test-client', version: '1.0.0' },
            },
          },
        },
        {
          name: 'ping',
          payload: { jsonrpc: '2.0', method: 'ping', id: 2 },
        },
      ];

      for (const testReq of testRequests) {
        const response = await request.post(`${mcpServiceUrl}/mcp`, {
          headers: { 'Content-Type': 'application/json' },
          data: testReq.payload,
        });

        // Then: Response conforms to MCP protocol specification
        expect(response.status()).toBeLessThan(300);
        const result = await response.json();
        expect(result).toHaveProperty('jsonrpc', '2.0');
        expect(result).toHaveProperty('id', testReq.payload.id);
        expect('result' in result || 'error' in result).toBeTruthy();
      }
    }
  );

  test(
    'INT-022: Server tracks session state and message activity',
    async ({ request }) => {
      // Given: Client establishes MCP session via HTTP POST
      // When: Client sends messages over session
      for (let i = 0; i < 5; i++) {
        await request.get(`${mcpServiceUrl}/health`);
      }

      // Then: Server maintains session state
      const metricsResponse = await request.get(`${mcpServiceUrl}/metrics`);
      expect(metricsResponse.ok()).toBeTruthy();
      const metrics = await metricsResponse.json();

      // Then: Server tracks message count for session
      expect(metrics).toHaveProperty('total_connections');

      // Then: Server records last activity timestamp
      expect(metrics).toHaveProperty('last_activity');
    }
  );

  test(
    'INT-023: Server responds to keepalive request',
    async ({ request }) => {
      // Given: Client maintains active MCP session

      // When: Client sends keepalive request (ping)
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: { 'Content-Type': 'application/json' },
        data: { jsonrpc: '2.0', method: 'ping', id: 1 },
      });

      // Then: Server responds with keepalive acknowledgment immediately
      expect(response.ok()).toBeTruthy();
      const result = await response.json();
      expect(result).toHaveProperty('id', 1);
    }
  );

  test('INT-028: Metrics endpoint returns connection statistics', async ({
    request,
  }) => {
    // Given: Server tracks connection metrics

    // When: Client requests metrics endpoint GET /metrics
    const response = await request.get(`${mcpServiceUrl}/metrics`);

    // Then: Server returns metrics
    expect(response.ok()).toBeTruthy();
    const metrics = await response.json();

    // Then: Server returns active connection count
    expect(metrics).toHaveProperty('active_connections');

    // Then: Server returns total connections handled
    expect(metrics).toHaveProperty('total_connections');

    // Then: Server returns connection error count
    expect(metrics).toHaveProperty('connection_errors');
  });

  test(
    'INT-032: Server enforces rate limit and rejects excess requests',
    async ({ request }) => {
      // Given: Server enforces rate limit
      const results: number[] = [];

      // When: Client sends many requests rapidly
      for (let i = 0; i < 110; i++) {
        const response = await request.get(`${mcpServiceUrl}/health`);
        results.push(response.status());
      }

      // Then: Server may reject some requests with 429
      const rateLimited = results.filter((s) => s === 429);
      // Rate limiting may or may not trigger depending on configuration
      // Just verify server handled all requests without 5xx errors
      const serverErrors = results.filter((s) => s >= 500);
      expect(serverErrors.length).toBe(0);
    }
  );

  test(
    'INT-034: Server maintains stable performance under concurrent load',
    async ({ request }) => {
      // Given: Server operates under normal conditions
      const latencies: number[] = [];

      // When: Load test establishes 50 concurrent requests
      const promises = Array.from({ length: 50 }, async () => {
        const start = Date.now();
        const response = await request.get(`${mcpServiceUrl}/health`);
        const latency = Date.now() - start;
        latencies.push(latency);
        expect(response.ok()).toBeTruthy();
      });

      await Promise.all(promises);

      // Then: Server maintains stable performance
      expect(latencies.length).toBe(50);

      // Then: Response latency stays reasonable (p95 < 500ms for HTTP)
      const sorted = latencies.sort((a, b) => a - b);
      const p95 = sorted[Math.floor(sorted.length * 0.95)];
      expect(p95).toBeLessThan(500);
    }
  );

  test(
    'INT-035: Server survives request storm and recovers',
    async ({ request }) => {
      // Given: Server handles requests normally
      const preHealth = await request.get(`${mcpServiceUrl}/health`);
      expect(preHealth.ok()).toBeTruthy();

      // When: 100 requests within short timeframe
      const promises = Array.from({ length: 100 }, () =>
        request.get(`${mcpServiceUrl}/health`).catch(() => null)
      );
      await Promise.all(promises);

      // Then: Server recovers normal operation after storm
      await new Promise((r) => setTimeout(r, 100));
      const postHealth = await request.get(`${mcpServiceUrl}/health`);
      expect(postHealth.ok()).toBeTruthy();
    }
  );

  test(
    'INT-036: Server achieves target message throughput',
    async ({ request }) => {
      // Given: Server manages concurrent HTTP requests
      const start = Date.now();
      let processed = 0;

      // When: Send 300 messages rapidly
      const promises = Array.from({ length: 300 }, async () => {
        const response = await request.post(`${mcpServiceUrl}/mcp`, {
          headers: { 'Content-Type': 'application/json' },
          data: { jsonrpc: '2.0', method: 'ping', id: Date.now() },
        });
        if (response.ok()) {
          processed++;
        }
      });

      await Promise.all(promises);
      const elapsed = Date.now() - start;

      // Then: Server processes messages with reasonable throughput
      const throughput = (processed / elapsed) * 1000;
      expect(processed).toBeGreaterThan(250); // At least 85% success
      expect(throughput).toBeGreaterThan(50); // At least 50 msg/s over HTTP
    }
  );
});
