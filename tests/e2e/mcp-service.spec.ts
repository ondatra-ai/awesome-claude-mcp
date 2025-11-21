import { test, expect } from '@playwright/test';

import { getEnvironmentConfig } from '../config/environments';

const { mcpServiceUrl } = getEnvironmentConfig(process.env.E2E_ENV);

/**
 * MCP Service E2E Tests
 *
 * These tests verify complete MCP workflows using HTTP+SSE transport
 * (Streamable HTTP per MCP specification).
 *
 * Transport Protocol:
 * - POST /mcp: Client sends JSON-RPC messages
 * - GET /mcp: Client establishes SSE stream for server-to-client messages
 * - Header: Mcp-Session-Id for session tracking
 *
 * Note: For tests requiring Claude API integration (LLM simulation),
 * ensure .env.test contains ANTHROPIC_API_KEY.
 */
test.describe('MCP Service E2E Tests', () => {
  /**
   * E2E-011: MCP client completes handshake with server successfully
   *
   * Tests the complete initialization handshake flow:
   * 1. Client sends initialize request via POST /mcp
   * 2. Server returns initialize response with capabilities
   * 3. Client sends initialized notification
   * 4. Session is ready for tool calls
   */
  test('E2E-011: MCP client completes handshake with server', async ({
    request,
  }) => {
    // Given: MCP server runs on configured endpoint
    // When: Client completes MCP handshake from connection to first message

    // Step 1: Send initialize request
    const initializeRequest = {
      jsonrpc: '2.0',
      method: 'initialize',
      id: 1,
      params: {
        protocolVersion: '2024-11-05',
        capabilities: {},
        clientInfo: {
          name: 'playwright-e2e-test',
          version: '1.0.0',
        },
      },
    };

    const initResponse = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: initializeRequest,
    });

    // Then: Client establishes HTTP+SSE session successfully
    expect(initResponse.status()).toBe(200);

    const initResult = await initResponse.json();

    // Then: Server returns initialize response in MCP format
    expect(initResult).toHaveProperty('jsonrpc', '2.0');
    expect(initResult).toHaveProperty('id', 1);
    expect(initResult).toHaveProperty('result');
    expect(initResult.result).toHaveProperty('protocolVersion');
    expect(initResult.result).toHaveProperty('serverInfo');
    expect(initResult.result.serverInfo).toHaveProperty('name');

    // Get session ID from header or response
    const sessionId =
      initResponse.headers()['mcp-session-id'] ||
      initResult.result.sessionId ||
      'default';

    // Step 2: Send initialized notification (completes handshake)
    const initializedNotification = {
      jsonrpc: '2.0',
      method: 'notifications/initialized',
    };

    const notifyResponse = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
        'Mcp-Session-Id': sessionId,
      },
      data: initializedNotification,
    });

    // Notifications may return 200 with empty body or 204
    expect([200, 204]).toContain(notifyResponse.status());

    // Step 3: Verify session is ready by listing tools
    const toolsListRequest = {
      jsonrpc: '2.0',
      method: 'tools/list',
      id: 2,
      params: {},
    };

    const toolsResponse = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
        'Mcp-Session-Id': sessionId,
      },
      data: toolsListRequest,
    });

    expect(toolsResponse.status()).toBe(200);

    const toolsResult = await toolsResponse.json();
    expect(toolsResult).toHaveProperty('jsonrpc', '2.0');
    expect(toolsResult).toHaveProperty('id', 2);
    expect(toolsResult).toHaveProperty('result');
    expect(toolsResult.result).toHaveProperty('tools');
    expect(Array.isArray(toolsResult.result.tools)).toBeTruthy();
  });

  /**
   * E2E-012: Server processes complete MCP request-response cycle
   *
   * Tests a complete request-response cycle including:
   * 1. Session initialization
   * 2. Request validation
   * 3. Request processing
   * 4. MCP-compliant response
   */
  test('E2E-012: Complete MCP request-response cycle', async ({ request }) => {
    // Given: Client maintains active MCP session
    // Initialize session first
    const initResponse = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: { 'Content-Type': 'application/json' },
      data: {
        jsonrpc: '2.0',
        method: 'initialize',
        id: 1,
        params: {
          protocolVersion: '2024-11-05',
          capabilities: {},
          clientInfo: { name: 'playwright-e2e-test', version: '1.0.0' },
        },
      },
    });

    expect(initResponse.status()).toBe(200);
    const initResult = await initResponse.json();
    const sessionId =
      initResponse.headers()['mcp-session-id'] ||
      initResult.result?.sessionId ||
      'default';

    // Send initialized notification
    await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
        'Mcp-Session-Id': sessionId,
      },
      data: { jsonrpc: '2.0', method: 'notifications/initialized' },
    });

    // When: Client executes complete request-response cycle
    const testRequest = {
      jsonrpc: '2.0',
      method: 'tools/list',
      id: 2,
      params: {},
    };

    const startTime = Date.now();
    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
        'Mcp-Session-Id': sessionId,
      },
      data: testRequest,
    });
    const responseTime = Date.now() - startTime;

    // Then: Server validates request against MCP schema
    expect(response.status()).toBe(200);

    const result = await response.json();

    // Then: Server processes request according to MCP specification
    expect(result).toHaveProperty('jsonrpc', '2.0');
    expect(result).toHaveProperty('id', 2);

    // Then: Server returns response conforming to MCP protocol
    expect(result.result !== undefined || result.error !== undefined).toBe(
      true
    );

    // Response should have required MCP fields
    if (result.result) {
      expect(result.result).toHaveProperty('tools');
    } else if (result.error) {
      expect(result.error).toHaveProperty('code');
      expect(result.error).toHaveProperty('message');
    }

    // Response should be reasonably fast (under 5 seconds)
    expect(responseTime).toBeLessThan(5000);
  });

  /**
   * E2E-013: Server handles invalid requests with fail-fast error response
   *
   * Tests error handling behavior:
   * 1. Invalid request detection
   * 2. MCP-formatted error response
   * 3. Session remains usable after error
   */
  test('E2E-013: Invalid request error handling', async ({ request }) => {
    // Given: Client maintains active MCP session
    const initResponse = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: { 'Content-Type': 'application/json' },
      data: {
        jsonrpc: '2.0',
        method: 'initialize',
        id: 1,
        params: {
          protocolVersion: '2024-11-05',
          capabilities: {},
          clientInfo: { name: 'playwright-e2e-test', version: '1.0.0' },
        },
      },
    });

    expect(initResponse.status()).toBe(200);
    const initResult = await initResponse.json();
    const sessionId =
      initResponse.headers()['mcp-session-id'] ||
      initResult.result?.sessionId ||
      'default';

    // Send initialized notification
    await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
        'Mcp-Session-Id': sessionId,
      },
      data: { jsonrpc: '2.0', method: 'notifications/initialized' },
    });

    // When: Client sends invalid request (missing method field)
    const invalidRequest = {
      jsonrpc: '2.0',
      // Missing 'method' field - should trigger validation error
      id: 2,
      params: {},
    };

    const startTime = Date.now();
    const errorResponse = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
        'Mcp-Session-Id': sessionId,
      },
      data: invalidRequest,
    });
    const errorResponseTime = Date.now() - startTime;

    // Then: Server detects error immediately per fail-fast principle
    // Error should be detected within 1 second
    expect(errorResponseTime).toBeLessThan(1000);

    // Server may return 200 with error body or 400 for invalid request
    expect([200, 400]).toContain(errorResponse.status());

    const errorResult = await errorResponse.json();

    // Then: Server returns MCP-formatted error response
    expect(errorResult).toHaveProperty('jsonrpc', '2.0');

    if (errorResult.error) {
      expect(errorResult.error).toHaveProperty('code');
      expect(errorResult.error).toHaveProperty('message');
      expect(typeof errorResult.error.code).toBe('number');
      expect(typeof errorResult.error.message).toBe('string');
    }

    // Then: Session remains usable for subsequent valid requests
    const validRequest = {
      jsonrpc: '2.0',
      method: 'tools/list',
      id: 3,
      params: {},
    };

    const subsequentResponse = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
        'Mcp-Session-Id': sessionId,
      },
      data: validRequest,
    });

    expect(subsequentResponse.status()).toBe(200);

    const subsequentResult = await subsequentResponse.json();
    expect(subsequentResult).toHaveProperty('jsonrpc', '2.0');
    expect(subsequentResult).toHaveProperty('id', 3);
    expect(subsequentResult).toHaveProperty('result');
  });

  /**
   * E2E-014: Server isolates concurrent client sessions
   *
   * Tests session isolation:
   * 1. Multiple concurrent sessions
   * 2. No response cross-contamination
   * 3. Each session maintains independent state
   */
  test('E2E-014: Concurrent session isolation', async ({ request }) => {
    // Given: Multiple clients maintain active sessions
    const NUM_CLIENTS = 3;
    const sessions: { sessionId: string; clientName: string }[] = [];

    // Initialize multiple sessions concurrently
    const initPromises = [];
    for (let i = 0; i < NUM_CLIENTS; i++) {
      initPromises.push(
        request.post(`${mcpServiceUrl}/mcp`, {
          headers: { 'Content-Type': 'application/json' },
          data: {
            jsonrpc: '2.0',
            method: 'initialize',
            id: 1,
            params: {
              protocolVersion: '2024-11-05',
              capabilities: {},
              clientInfo: {
                name: `e2e-test-client-${i}`,
                version: '1.0.0',
              },
            },
          },
        })
      );
    }

    const initResponses = await Promise.all(initPromises);

    // Then: All clients connect successfully
    for (let i = 0; i < NUM_CLIENTS; i++) {
      expect(initResponses[i].status()).toBe(200);
      const result = await initResponses[i].json();
      const sessionId =
        initResponses[i].headers()['mcp-session-id'] ||
        result.result?.sessionId ||
        `session-${i}`;
      sessions.push({
        sessionId,
        clientName: `e2e-test-client-${i}`,
      });
    }

    // Send initialized notifications
    await Promise.all(
      sessions.map((session) =>
        request.post(`${mcpServiceUrl}/mcp`, {
          headers: {
            'Content-Type': 'application/json',
            'Mcp-Session-Id': session.sessionId,
          },
          data: { jsonrpc: '2.0', method: 'notifications/initialized' },
        })
      )
    );

    // When: Clients execute concurrent operations with unique IDs
    const requestPromises = sessions.map((session, index) => {
      return request.post(`${mcpServiceUrl}/mcp`, {
        headers: {
          'Content-Type': 'application/json',
          'Mcp-Session-Id': session.sessionId,
        },
        data: {
          jsonrpc: '2.0',
          method: 'tools/list',
          id: 100 + index, // Unique request ID per client
          params: {},
        },
      });
    });

    const responses = await Promise.all(requestPromises);

    // Then: Each client receives correct response
    for (let i = 0; i < NUM_CLIENTS; i++) {
      expect(responses[i].status()).toBe(200);
      const result = await responses[i].json();

      // Then: Responses isolated - correct ID matches request
      expect(result).toHaveProperty('jsonrpc', '2.0');
      expect(result).toHaveProperty('id', 100 + i);
      expect(result).toHaveProperty('result');

      // Then: No cross-contamination - each session independent
      expect(result.result).toHaveProperty('tools');
    }
  });

  /**
   * E2E-015: Health check returns status with session and dependency metrics
   *
   * Tests health endpoint which uses standard HTTP (not MCP protocol):
   * 1. Returns healthy status
   * 2. Includes session pool metrics
   * 3. Includes dependency health (Redis, Google API)
   */
  test('E2E-015: Health check with metrics', async ({ request }) => {
    // Given: Server accepts health check requests
    // When: Monitoring system queries health endpoint
    const response = await request.get(`${mcpServiceUrl}/health`);

    // Then: Server reports healthy status
    expect(response.status()).toBe(200);

    const healthData = await response.json();
    expect(healthData).toHaveProperty('status');
    expect(healthData.status).toBe('healthy');

    // Then: Server includes connection/session pool metrics
    // Note: Server may use 'connections' or 'sessions' depending on implementation
    const poolMetrics = healthData.sessions || healthData.connections;
    expect(poolMetrics).toBeTruthy();
    expect(poolMetrics).toHaveProperty('active');
    expect(poolMetrics).toHaveProperty('total');
    expect(typeof poolMetrics.active).toBe('number');
    expect(typeof poolMetrics.total).toBe('number');

    // Then: Server includes dependency health for Redis and Google API
    expect(healthData).toHaveProperty('dependencies');
    expect(healthData.dependencies).toHaveProperty('redis');
    expect(healthData.dependencies).toHaveProperty('googleAPI');

    // Validate Redis dependency health structure
    expect(healthData.dependencies.redis).toHaveProperty('status');
    expect(['healthy', 'unhealthy', 'degraded']).toContain(
      healthData.dependencies.redis.status
    );

    // Validate Google API dependency health structure
    expect(healthData.dependencies.googleAPI).toHaveProperty('status');
    expect(['healthy', 'unhealthy', 'degraded']).toContain(
      healthData.dependencies.googleAPI.status
    );
  });
});
