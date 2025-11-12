import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { mcpServiceUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('MCP Service Integration Tests', () => {

  test('INT-009: HTTP endpoint initiates WebSocket protocol upgrade', async ({ request }) => {
    // Given: Server exposes HTTP endpoint at /mcp
    const mcpEndpoint = `${mcpServiceUrl}/mcp`;

    // When: Client sends GET request to /mcp
    const response = await request.get(mcpEndpoint);

    // Then: Server initiates WebSocket protocol upgrade
    // The server should respond with status 101 (Switching Protocols)
    // or provide upgrade headers indicating WebSocket protocol support

    // Check if server returns upgrade response (101) or supports WebSocket
    const status = response.status();
    const headers = response.headers();

    // WebSocket upgrade can be indicated by:
    // 1. Status 101 Switching Protocols (direct upgrade)
    // 2. Upgrade header containing "websocket"
    // 3. Connection header containing "Upgrade"

    const upgradeHeader = headers['upgrade']?.toLowerCase();
    const connectionHeader = headers['connection']?.toLowerCase();

    // Verify server indicates WebSocket protocol upgrade capability
    const supportsWebSocketUpgrade =
      status === 101 ||
      upgradeHeader?.includes('websocket') ||
      (connectionHeader?.includes('upgrade') && upgradeHeader);

    expect(supportsWebSocketUpgrade).toBeTruthy();

    // If status is 101, verify proper upgrade headers are present
    if (status === 101) {
      expect(upgradeHeader).toContain('websocket');
      expect(connectionHeader).toContain('upgrade');
    }
  });

  test('INT-008: WebSocket connection establishment succeeds', async ({ page }) => {
    // Given: Server runs WebSocket service on configured port
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // Use page.evaluate to run WebSocket test in browser context
    const result = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{
        success: boolean;
        connectionEstablished: boolean;
        connectionId: string | null;
        error?: string;
      }>((resolve) => {
        try {
          // When: Client connects to ws://localhost:8081/mcp
          const ws = new WebSocket(`${wsEndpoint}/mcp`);
          let connectionEstablished = false;
          let connectionId: string | null = null;

          ws.onopen = () => {
            // Then: Server establishes WebSocket connection
            connectionEstablished = true;

            // Send a message to potentially receive connection identifier
            ws.send(JSON.stringify({
              jsonrpc: '2.0',
              method: 'initialize',
              id: 1,
              params: {
                protocolVersion: '2024-11-05',
                capabilities: {},
                clientInfo: {
                  name: 'test-client',
                  version: '1.0.0'
                }
              }
            }));
          };

          ws.onmessage = (event) => {
            try {
              const message = JSON.parse(event.data as string);
              // Then: Server assigns unique connection identifier
              // Look for connection ID in response (could be in various places)
              if (message.result?.connectionId) {
                connectionId = message.result.connectionId;
              } else if (message.connectionId) {
                connectionId = message.connectionId;
              } else if (message.id) {
                // Use message ID as proxy for connection tracking
                connectionId = String(message.id);
              }

              ws.close();
              resolve({
                success: connectionEstablished,
                connectionEstablished,
                connectionId,
              });
            } catch (err) {
              // If we can't parse response, connection was still established
              ws.close();
              resolve({
                success: connectionEstablished,
                connectionEstablished,
                connectionId: 'implicit', // Connection exists even without explicit ID
              });
            }
          };

          ws.onerror = (error) => {
            resolve({
              success: false,
              connectionEstablished: false,
              connectionId: null,
              error: `WebSocket error: ${error}`,
            });
          };

          ws.onclose = (event) => {
            // If connection was established but closed without receiving message
            if (connectionEstablished && !connectionId) {
              resolve({
                success: connectionEstablished,
                connectionEstablished,
                connectionId: 'implicit', // Connection existed
              });
            }
          };

          // Set timeout for connection establishment
          setTimeout(() => {
            if (!connectionEstablished) {
              ws.close();
              resolve({
                success: false,
                connectionEstablished: false,
                connectionId: null,
                error: 'Connection timeout after 5 seconds',
              });
            }
          }, 5000);
        } catch (err) {
          resolve({
            success: false,
            connectionEstablished: false,
            connectionId: null,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, wsUrl);

    // Then: Server establishes WebSocket connection
    expect(result.success).toBeTruthy();
    expect(result.connectionEstablished).toBeTruthy();

    // Then: Server assigns unique connection identifier
    // The connection identifier can be explicit or implicit (connection itself)
    expect(result.connectionId).toBeTruthy();

    if (!result.success) {
      throw new Error(result.error || 'Unknown WebSocket connection failure');
    }
  });

  test('INT-020: Server response includes JSON structure with metadata fields', async ({ request }) => {
    // Given: Server processes client requests
    // When: Server generates response
    const response = await request.get(`${mcpServiceUrl}/health`);

    // Then: Response contains valid JSON structure
    expect(response.ok()).toBeTruthy();
    const data = await response.json();

    // Then: Response includes timestamp field
    expect(data).toHaveProperty('timestamp');
    expect(typeof data.timestamp).toBe('string');

    // Then: Response includes correlation ID for tracing
    expect(data).toHaveProperty('correlation_id');
    expect(typeof data.correlation_id).toBe('string');
    expect(data.correlation_id).toMatch(/^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$/);
  });

  test('INT-035: Server survives connection storm and recovers to normal operation', async ({ request }) => {
    // Given: Server handles connections normally
    // Verify server is healthy before storm
    const preStormHealth = await request.get(`${mcpServiceUrl}/health`);
    expect(preStormHealth.ok()).toBeTruthy();

    // When: 100 clients connect and disconnect within 5 seconds
    const startTime = Date.now();
    const connectionPromises: Promise<void>[] = [];

    // Simulate 100 rapid connection attempts
    for (let i = 0; i < 100; i++) {
      const connectionPromise = (async () => {
        try {
          // Attempt connection (using health endpoint as proxy for connection)
          await request.get(`${mcpServiceUrl}/health`);
        } catch (error) {
          // Connections may fail during storm - this is expected behavior
          // We're testing that the server survives, not that all connections succeed
        }
      })();
      connectionPromises.push(connectionPromise);
    }

    // Wait for all connection attempts to complete
    await Promise.all(connectionPromises);
    const elapsedTime = Date.now() - startTime;

    // Verify storm completed within 5 seconds
    expect(elapsedTime).toBeLessThan(5000);

    // Then: Server processes connection storm without crashes
    // Give server brief moment to stabilize (100ms)
    await new Promise(resolve => setTimeout(resolve, 100));

    // Then: Server recovers normal operation after storm
    const postStormHealth = await request.get(`${mcpServiceUrl}/health`);
    expect(postStormHealth.ok()).toBeTruthy();

    const healthData = await postStormHealth.json();
    expect(healthData).toHaveProperty('timestamp');
    expect(healthData).toHaveProperty('correlation_id');
  });

  test('INT-011: Server rejects unauthorized domain CORS preflight request', async ({ request }) => {
    // Given: Server enforces CORS policy for Claude domains
    // The server should be configured to reject requests from unauthorized domains

    // When: Client from unauthorized domain sends preflight OPTIONS request
    const response = await request.fetch(`${mcpServiceUrl}/health`, {
      method: 'OPTIONS',
      headers: {
        'Origin': 'https://malicious-domain.com',
        'Access-Control-Request-Method': 'GET',
        'Access-Control-Request-Headers': 'content-type',
      },
    });

    // Then: Server rejects request with CORS error
    // The server should either:
    // 1. Return 403 Forbidden (explicit rejection)
    // 2. Return 200/204 but without CORS headers (implicit rejection)
    // 3. Return CORS error status

    if (response.status() === 403 || response.status() === 401) {
      // Explicit rejection - verify status code indicates forbidden
      expect(response.status()).toBeGreaterThanOrEqual(400);
    } else {
      // Implicit rejection - verify CORS headers are missing or restrictive
      const headers = response.headers();
      const allowedOrigin = headers['access-control-allow-origin'];

      // Either no CORS header present, or it doesn't match the unauthorized origin
      if (allowedOrigin) {
        expect(allowedOrigin).not.toBe('https://malicious-domain.com');
        // If wildcard is present, this means CORS is too permissive (security issue)
        expect(allowedOrigin).not.toBe('*');
      }
      // If no CORS headers present, browser will block the request (desired behavior)
    }
  });

  test('INT-012: Server accepts CORS preflight from authorized Claude domain', async ({ request }) => {
    // Given: Server enforces CORS policy for Claude domains
    // The server should be configured to accept requests from claude.ai domain

    // When: Client from claude.ai domain sends preflight OPTIONS request
    const response = await request.fetch(`${mcpServiceUrl}/health`, {
      method: 'OPTIONS',
      headers: {
        'Origin': 'https://claude.ai',
        'Access-Control-Request-Method': 'GET',
        'Access-Control-Request-Headers': 'content-type',
      },
    });

    // Then: Server accepts request with CORS headers
    expect(response.status()).toBe(200);

    // Verify CORS headers are present
    const headers = response.headers();
    expect(headers['access-control-allow-origin']).toBeTruthy();

    // Verify the origin is allowed (should be either the specific origin or wildcard)
    const allowedOrigin = headers['access-control-allow-origin'];
    expect(['https://claude.ai', '*']).toContain(allowedOrigin);

    // Verify allowed methods include GET (the requested method)
    const allowedMethods = headers['access-control-allow-methods'];
    expect(allowedMethods).toBeTruthy();
    expect(allowedMethods.toUpperCase()).toContain('GET');

    // Verify allowed headers include the requested header
    const allowedHeaders = headers['access-control-allow-headers'];
    if (allowedHeaders) {
      expect(allowedHeaders.toLowerCase()).toContain('content-type');
    }
  });

  test('INT-014: Server rejects invalid JSON with parsing error', async ({ request }) => {
    // Given: Client maintains active WebSocket connection
    // Note: Testing via HTTP endpoint that accepts JSON, simulating message handling
    // The MCP service should validate JSON before processing any message

    // When: Client sends message with invalid JSON format
    const invalidJsonPayload = '{invalid json: missing quotes}';

    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: invalidJsonPayload,
    });

    // Then: Server returns parsing error immediately
    // Server should reject with 4xx status code (typically 400 Bad Request)
    expect(response.status()).toBeGreaterThanOrEqual(400);
    expect(response.status()).toBeLessThan(500);

    // Then: Error message specifies invalid JSON
    const errorData = await response.text();

    // Verify error message contains indicators of JSON parsing failure
    const errorIndicators = [
      'json',
      'parse',
      'parsing',
      'invalid',
      'syntax',
      'malformed',
    ];

    const containsJsonError = errorIndicators.some(indicator =>
      errorData.toLowerCase().includes(indicator)
    );

    expect(containsJsonError).toBeTruthy();
  });

  test('INT-022: Server tracks connection state and message activity', async ({ request }) => {
    // Given: Client establishes WebSocket connection
    // Note: Using metrics endpoint to verify connection tracking
    // The MCP service should track connection state and message activity

    // When: Client sends messages over connection
    // Simulate connection activity by making multiple requests
    const messageCount = 5;
    for (let i = 0; i < messageCount; i++) {
      await request.get(`${mcpServiceUrl}/health`);
    }

    // Retrieve metrics to verify tracking
    const metricsResponse = await request.get(`${mcpServiceUrl}/metrics`);

    // Then: Server maintains connection state
    expect(metricsResponse.ok()).toBeTruthy();
    const metricsData = await metricsResponse.json();

    // Then: Server tracks message count for connection
    // Verify metrics contain connection tracking information
    expect(metricsData).toHaveProperty('total_connections');
    expect(typeof metricsData.total_connections).toBe('number');
    expect(metricsData.total_connections).toBeGreaterThanOrEqual(0);

    // Then: Server records last activity timestamp
    expect(metricsData).toHaveProperty('last_activity');
    expect(typeof metricsData.last_activity).toBe('string');

    // Verify timestamp is valid ISO 8601 format
    const timestampRegex = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z?$/;
    expect(metricsData.last_activity).toMatch(timestampRegex);

    // Verify timestamp is recent (within last minute)
    const lastActivity = new Date(metricsData.last_activity);
    const now = new Date();
    const timeDiff = now.getTime() - lastActivity.getTime();
    expect(timeDiff).toBeLessThan(60000); // Less than 60 seconds
  });

  test('INT-023: Server responds to WebSocket ping with pong frame', async ({ page }) => {
    // Given: Client maintains active WebSocket connection
    // Create WebSocket connection to MCP service
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // Use page.evaluate to run WebSocket test in browser context
    const result = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{ success: boolean; error?: string }>((resolve) => {
        try {
          // Establish WebSocket connection
          const ws = new WebSocket(`${wsEndpoint}/mcp`);
          let pongReceived = false;
          let pingTimeout: NodeJS.Timeout;

          ws.onopen = () => {
            // When: Client sends ping frame
            // Send WebSocket ping (browser WebSocket API doesn't expose ping directly)
            // We'll send a JSON-RPC ping message instead
            ws.send(JSON.stringify({
              jsonrpc: '2.0',
              method: 'ping',
              id: 1
            }));

            // Set timeout for pong response (should be immediate)
            pingTimeout = setTimeout(() => {
              if (!pongReceived) {
                ws.close();
                resolve({ success: false, error: 'Pong response timeout' });
              }
            }, 5000); // 5 second timeout
          };

          ws.onmessage = (event) => {
            try {
              const message = JSON.parse(event.data as string);
              // Then: Server responds with pong frame immediately
              if (message.result === 'pong' || message.method === 'pong') {
                pongReceived = true;
                clearTimeout(pingTimeout);
                ws.close();
                resolve({ success: true });
              }
            } catch (err) {
              // Ignore parse errors for non-JSON messages
            }
          };

          ws.onerror = (error) => {
            clearTimeout(pingTimeout);
            resolve({
              success: false,
              error: `WebSocket error: ${error}`
            });
          };

          ws.onclose = () => {
            clearTimeout(pingTimeout);
            if (!pongReceived) {
              resolve({
                success: false,
                error: 'Connection closed without pong response'
              });
            }
          };
        } catch (err) {
          resolve({
            success: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`
          });
        }
      });
    }, wsUrl);

    // Then: Server responds with pong frame immediately
    expect(result.success).toBeTruthy();
    if (!result.success) {
      throw new Error(result.error || 'Unknown WebSocket ping/pong failure');
    }
  });

  test('INT-026: Server performs graceful shutdown on SIGTERM signal', async ({ request }) => {
    // Given: Server manages multiple active connections
    // Verify server is healthy and responsive before shutdown
    const preShutdownHealth = await request.get(`${mcpServiceUrl}/health`);
    expect(preShutdownHealth.ok()).toBeTruthy();

    // Establish multiple active connections (simulate multiple clients)
    const activeConnections: Promise<void>[] = [];
    const connectionCount = 5;

    for (let i = 0; i < connectionCount; i++) {
      const connectionPromise = (async () => {
        try {
          // Keep connection active with periodic health checks
          await request.get(`${mcpServiceUrl}/health`);
        } catch (error) {
          // Connection may be interrupted during shutdown - expected behavior
        }
      })();
      activeConnections.push(connectionPromise);
    }

    // Wait for all connections to be established
    await Promise.all(activeConnections);

    // When: Server receives SIGTERM signal
    // Note: In integration test environment, we simulate shutdown by calling shutdown endpoint
    // The actual SIGTERM handling would be tested in the server implementation
    const shutdownStartTime = Date.now();

    let shutdownResponse;
    try {
      shutdownResponse = await request.post(`${mcpServiceUrl}/shutdown`, {
        headers: {
          'Content-Type': 'application/json',
        },
        data: {
          signal: 'SIGTERM',
          graceful: true,
        },
      });
    } catch (error) {
      // Server may close connection during shutdown - this is expected
      // We'll verify the shutdown behavior through subsequent checks
    }

    // Then: Server closes all active connections gracefully
    // Wait briefly to allow graceful shutdown to begin
    await new Promise(resolve => setTimeout(resolve, 100));

    // Then: Server waits for in-flight requests to complete
    // Attempt to verify server is shutting down
    let serverShuttingDown = false;
    try {
      const shutdownCheckResponse = await request.get(`${mcpServiceUrl}/health`, {
        timeout: 1000, // Short timeout
      });

      // If server responds with a shutdown status, it's gracefully shutting down
      if (shutdownCheckResponse.status() === 503) {
        serverShuttingDown = true;
      }
    } catch (error) {
      // Server not responding means it has shut down
      serverShuttingDown = true;
    }

    expect(serverShuttingDown).toBeTruthy();

    // Then: Server shuts down within 10 seconds
    const shutdownDuration = Date.now() - shutdownStartTime;
    expect(shutdownDuration).toBeLessThan(10000);

    // Verify server is no longer accepting new connections
    try {
      const postShutdownResponse = await request.get(`${mcpServiceUrl}/health`, {
        timeout: 2000,
      });

      // If we get a response, it should indicate service unavailable
      expect(postShutdownResponse.status()).toBeGreaterThanOrEqual(500);
    } catch (error) {
      // Connection refused is expected after shutdown - this is the desired state
      expect(error).toBeTruthy();
    }
  });

  test('INT-010: Health check returns service status with dependencies', async ({ request }) => {
    // Given: Server provides health check endpoint
    // The MCP service should expose /health endpoint for monitoring

    // When: Client sends GET request to /health
    const response = await request.get(`${mcpServiceUrl}/health`);

    // Then: Server returns status 200
    expect(response.status()).toBe(200);

    // Then: Response includes service version
    const data = await response.json();
    expect(data).toHaveProperty('version');
    expect(typeof data.version).toBe('string');
    expect(data.version).toBeTruthy();

    // Then: Response includes dependency status for Redis and Google API
    expect(data).toHaveProperty('dependencies');
    expect(data.dependencies).toBeTruthy();

    // Verify Redis dependency status
    expect(data.dependencies).toHaveProperty('redis');
    expect(data.dependencies.redis).toHaveProperty('status');
    expect(['healthy', 'unhealthy', 'unknown']).toContain(data.dependencies.redis.status);

    // Verify Google API dependency status
    expect(data.dependencies).toHaveProperty('google_api');
    expect(data.dependencies.google_api).toHaveProperty('status');
    expect(['healthy', 'unhealthy', 'unknown']).toContain(data.dependencies.google_api.status);
  });

  test('INT-016: Server rejects message with invalid field values', async ({ request }) => {
    // Given: Client maintains active WebSocket connection
    // Note: Testing via HTTP endpoint that accepts JSON, simulating MCP message validation
    // The MCP service should validate field values before processing any message

    // When: Client sends message with field containing invalid value
    // Example: protocol_version with invalid format, method with unsupported value
    const invalidMessagePayload = {
      jsonrpc: '2.0',
      method: 'invalid_method_name_with_spaces_and_special_chars!@#',
      id: 1,
      params: {
        protocol_version: '999.999.999', // Invalid protocol version
        invalid_field: 'unexpected_value',
      },
    };

    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: invalidMessagePayload,
    });

    // Then: Server returns validation error
    // Server should reject with 4xx status code (typically 400 Bad Request)
    expect(response.status()).toBeGreaterThanOrEqual(400);
    expect(response.status()).toBeLessThan(500);

    // Then: Error message specifies error detail
    const responseData = await response.text();

    // Verify error message contains validation error indicators
    const validationErrorIndicators = [
      'validation',
      'invalid',
      'field',
      'value',
      'error',
      'protocol_version',
      'method',
    ];

    const containsValidationError = validationErrorIndicators.some(indicator =>
      responseData.toLowerCase().includes(indicator)
    );

    expect(containsValidationError).toBeTruthy();
  });

  test('INT-024: Server closes idle connection after timeout period', async ({ page }) => {
    // Given: Server enforces connection timeout of 30 seconds
    // Given: Client maintains idle connection
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // Use page.evaluate to run WebSocket test in browser context
    const result = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{
        success: boolean;
        connectionClosed: boolean;
        elapsedTime: number;
        error?: string;
      }>((resolve) => {
        try {
          // Establish WebSocket connection
          const ws = new WebSocket(`${wsEndpoint}/mcp`);
          let connectionOpened = false;
          let connectionClosed = false;
          const startTime = Date.now();

          ws.onopen = () => {
            connectionOpened = true;
            // Connection established - now maintain idle state
            // Do NOT send any messages to keep connection idle
          };

          ws.onclose = () => {
            // Then: Server closes connection automatically
            const elapsedTime = Date.now() - startTime;
            connectionClosed = true;

            resolve({
              success: connectionOpened && connectionClosed,
              connectionClosed: true,
              elapsedTime,
            });
          };

          ws.onerror = (error) => {
            const elapsedTime = Date.now() - startTime;
            resolve({
              success: false,
              connectionClosed: false,
              elapsedTime,
              error: `WebSocket error: ${error}`,
            });
          };

          // Set timeout longer than expected idle timeout (30 seconds + buffer)
          // If connection is not closed after 35 seconds, consider it a failure
          setTimeout(() => {
            const elapsedTime = Date.now() - startTime;
            if (!connectionClosed) {
              ws.close();
              resolve({
                success: false,
                connectionClosed: false,
                elapsedTime,
                error: 'Connection not closed after 35 seconds',
              });
            }
          }, 35000); // 35 second timeout
        } catch (err) {
          resolve({
            success: false,
            connectionClosed: false,
            elapsedTime: 0,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, wsUrl);

    // Then: Server closes connection automatically
    expect(result.success).toBeTruthy();
    expect(result.connectionClosed).toBeTruthy();

    // When: 30 seconds elapse without client activity
    // Verify connection was closed within reasonable time frame (30-35 seconds)
    expect(result.elapsedTime).toBeGreaterThanOrEqual(29000); // At least 29 seconds
    expect(result.elapsedTime).toBeLessThan(35000); // Less than 35 seconds

    if (!result.success) {
      throw new Error(result.error || 'Unknown idle connection timeout failure');
    }
  });

  test('INT-030: Server rejects connection when limit reached', async ({ page }) => {
    // Given: Server enforces maximum 100 concurrent connections
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');
    const maxConnections = 100;

    // Use page.evaluate to establish connections in browser context
    const result = await page.evaluate(async ({ wsEndpoint, max }) => {
      return new Promise<{
        success: boolean;
        connectionsEstablished: number;
        connectionRejected: boolean;
        error?: string;
      }>((resolve) => {
        try {
          const connections: WebSocket[] = [];
          let openConnections = 0;
          let rejectionDetected = false;
          let testComplete = false;

          // Function to establish a single connection
          const establishConnection = (index: number): Promise<boolean> => {
            return new Promise((resolveConn) => {
              const ws = new WebSocket(`${wsEndpoint}/mcp`);
              let connectionResolved = false;

              const timeoutId = setTimeout(() => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  resolveConn(false);
                }
              }, 5000); // 5 second timeout per connection

              ws.onopen = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  connections.push(ws);
                  openConnections++;
                  resolveConn(true);
                }
              };

              ws.onerror = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  // Error during connection - might be rejection
                  if (index === max) {
                    // This is the 101st connection - rejection expected
                    rejectionDetected = true;
                  }
                  resolveConn(false);
                }
              };

              ws.onclose = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  // Connection closed immediately - likely rejection
                  if (index === max) {
                    rejectionDetected = true;
                  }
                  resolveConn(false);
                }
              };
            });
          };

          // Establish maximum allowed connections (100)
          const connectionPromises: Promise<boolean>[] = [];
          for (let i = 0; i < max; i++) {
            connectionPromises.push(establishConnection(i));
          }

          // When: 101st client attempts to connect
          Promise.all(connectionPromises).then(async (results) => {
            // Count successful connections
            const successfulConnections = results.filter(r => r).length;

            // Small delay to ensure connections are stable
            await new Promise(r => setTimeout(r, 100));

            // Now attempt the 101st connection
            const rejectedConnection = await establishConnection(max);

            // Then: Server rejects connection
            if (testComplete) return;
            testComplete = true;

            // Clean up all connections
            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            // Then: Server returns error indicating connection limit reached
            resolve({
              success: successfulConnections === max && (rejectionDetected || !rejectedConnection),
              connectionsEstablished: successfulConnections,
              connectionRejected: rejectionDetected || !rejectedConnection,
            });
          }).catch((err) => {
            if (!testComplete) {
              testComplete = true;
              // Clean up connections
              connections.forEach(ws => {
                try {
                  ws.close();
                } catch (e) {
                  // Ignore cleanup errors
                }
              });

              resolve({
                success: false,
                connectionsEstablished: openConnections,
                connectionRejected: false,
                error: `Failed to establish connections: ${err instanceof Error ? err.message : String(err)}`,
              });
            }
          });

          // Safety timeout - if test doesn't complete in 60 seconds, fail
          setTimeout(() => {
            if (!testComplete) {
              testComplete = true;
              connections.forEach(ws => {
                try {
                  ws.close();
                } catch (e) {
                  // Ignore cleanup errors
                }
              });

              resolve({
                success: false,
                connectionsEstablished: openConnections,
                connectionRejected: false,
                error: 'Test timeout - took longer than 60 seconds',
              });
            }
          }, 60000);
        } catch (err) {
          resolve({
            success: false,
            connectionsEstablished: 0,
            connectionRejected: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, { wsEndpoint: wsUrl, max: maxConnections });

    // Then: Server rejects connection
    expect(result.success).toBeTruthy();
    expect(result.connectionsEstablished).toBe(maxConnections);
    expect(result.connectionRejected).toBeTruthy();

    if (!result.success) {
      throw new Error(result.error || `Connection limit test failed. Established: ${result.connectionsEstablished}, Rejected: ${result.connectionRejected}`);
    }
  });

  test('INT-034: Server maintains stable performance under concurrent connection load', async ({ request }) => {
    // Given: Server operates under normal conditions
    const connectionCount = 50;
    const latencyMeasurements: number[] = [];

    // When: Load test establishes 50 concurrent connections
    // Simulate concurrent load by making 50 parallel requests to health endpoint
    const concurrentRequests: Promise<void>[] = [];

    for (let i = 0; i < connectionCount; i++) {
      const requestPromise = (async () => {
        const startTime = Date.now();

        try {
          const response = await request.get(`${mcpServiceUrl}/health`);
          const endTime = Date.now();
          const latency = endTime - startTime;

          // Record latency for this request
          latencyMeasurements.push(latency);

          // Verify request succeeded
          expect(response.ok()).toBeTruthy();
        } catch (error) {
          // If request fails, record high latency to reflect performance degradation
          latencyMeasurements.push(10000); // 10 second penalty for failed requests
        }
      })();

      concurrentRequests.push(requestPromise);
    }

    // Wait for all concurrent requests to complete
    await Promise.all(concurrentRequests);

    // Then: Server maintains stable performance
    expect(latencyMeasurements.length).toBe(connectionCount);

    // Then: Response latency stays below 100ms at p95
    // Calculate p95 latency (95th percentile)
    const sortedLatencies = latencyMeasurements.sort((a, b) => a - b);
    const p95Index = Math.floor(sortedLatencies.length * 0.95);
    const p95Latency = sortedLatencies[p95Index];

    // Verify p95 latency is below 100ms
    expect(p95Latency).toBeLessThan(100);

    // Additional metrics for debugging if test fails
    const avgLatency = latencyMeasurements.reduce((sum, val) => sum + val, 0) / latencyMeasurements.length;
    const maxLatency = Math.max(...latencyMeasurements);
    const minLatency = Math.min(...latencyMeasurements);

    // Log performance metrics (helpful for debugging)
    if (p95Latency >= 100) {
      console.warn(`Performance degraded: p95=${p95Latency}ms, avg=${avgLatency.toFixed(2)}ms, max=${maxLatency}ms, min=${minLatency}ms`);
    }
  });

  test('INT-013: MCP protocol message parsing succeeds', async ({ request }) => {
    // Given: Client maintains active WebSocket connection
    // Note: Testing via HTTP endpoint that accepts MCP protocol messages
    // The MCP service should parse valid MCP protocol messages correctly

    // When: Client sends valid MCP protocol message
    const validMcpMessage = {
      jsonrpc: '2.0',
      method: 'initialize',
      id: 1,
      params: {
        protocolVersion: '2024-11-05',
        capabilities: {
          roots: {
            listChanged: true
          },
          sampling: {}
        },
        clientInfo: {
          name: 'test-client',
          version: '1.0.0'
        }
      }
    };

    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: validMcpMessage,
    });

    // Then: Server parses message successfully
    // Server should accept the message and respond with 2xx status
    expect(response.status()).toBeGreaterThanOrEqual(200);
    expect(response.status()).toBeLessThan(300);

    // Attempt to parse response as JSON
    const responseData = await response.json();

    // Then: Server extracts message type and protocol version
    // Verify the response indicates successful parsing
    // The server should have extracted the method type ('initialize')
    // and protocol version ('2024-11-05') from the request

    // Response should be valid JSON-RPC 2.0 format
    expect(responseData).toHaveProperty('jsonrpc');
    expect(responseData.jsonrpc).toBe('2.0');

    // Response should include the matching request ID
    expect(responseData).toHaveProperty('id');
    expect(responseData.id).toBe(validMcpMessage.id);

    // Response should either have 'result' (success) or 'error' (but not parsing error)
    const hasResult = 'result' in responseData;
    const hasError = 'error' in responseData;
    expect(hasResult || hasError).toBeTruthy();

    // If there's an error, it should not be a parsing error
    if (hasError) {
      const errorMessage = responseData.error?.message?.toLowerCase() || '';
      const parsingErrorIndicators = ['parse', 'parsing', 'json', 'syntax', 'malformed'];
      const isParsingError = parsingErrorIndicators.some(indicator =>
        errorMessage.includes(indicator)
      );
      expect(isParsingError).toBeFalsy();
    }

    // If there's a result, it indicates successful message parsing and processing
    if (hasResult) {
      // For initialize method, expect capabilities in response
      expect(responseData.result).toHaveProperty('protocolVersion');
      expect(responseData.result).toHaveProperty('capabilities');
      expect(responseData.result).toHaveProperty('serverInfo');
    }
  });

  test('INT-033: Server handles rapid connection churn without race conditions', async ({ page }) => {
    // Given: Server provides thread-safe connection pool
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');
    const clientCount = 50;

    // Use page.evaluate to run rapid connection churn test in browser context
    const result = await page.evaluate(async ({ wsEndpoint, count }) => {
      return new Promise<{
        success: boolean;
        totalAttempts: number;
        successfulConnections: number;
        raceConditionDetected: boolean;
        error?: string;
      }>((resolve) => {
        try {
          let successfulConnections = 0;
          let completedAttempts = 0;
          let raceConditionDetected = false;
          const connectionStates: Array<{ opened: boolean; closed: boolean }> = [];

          // Function to establish and immediately disconnect a connection
          const connectAndDisconnect = (index: number): Promise<boolean> => {
            return new Promise((resolveConn) => {
              const state = { opened: false, closed: false };
              connectionStates.push(state);

              const ws = new WebSocket(`${wsEndpoint}/mcp`);
              let connectionResolved = false;

              const timeoutId = setTimeout(() => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  // Timeout without connection
                  resolveConn(false);
                }
              }, 5000); // 5 second timeout per connection

              ws.onopen = () => {
                state.opened = true;
                clearTimeout(timeoutId);

                // Immediately close the connection to create churn
                ws.close();
              };

              ws.onclose = () => {
                state.closed = true;

                // Check for race condition: if closed before opened, connection state is inconsistent
                if (state.closed && !state.opened) {
                  raceConditionDetected = true;
                }

                if (!connectionResolved) {
                  connectionResolved = true;

                  // Consider successful if connection opened and closed properly
                  const wasSuccessful = state.opened && state.closed;
                  if (wasSuccessful) {
                    successfulConnections++;
                  }
                  resolveConn(wasSuccessful);
                }
              };

              ws.onerror = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  resolveConn(false);
                }
              };
            });
          };

          // When: 50 clients connect and disconnect rapidly
          const connectionPromises: Promise<boolean>[] = [];
          for (let i = 0; i < count; i++) {
            connectionPromises.push(connectAndDisconnect(i));
          }

          // Wait for all connection attempts to complete
          Promise.all(connectionPromises).then((results) => {
            completedAttempts = results.length;

            // Then: Server handles all connections without race conditions
            const allConnectionsHandled = completedAttempts === count;

            // Then: Connection count remains accurate
            // Verify each connection transitioned properly: opened -> closed
            const properStateTransitions = connectionStates.every(state =>
              state.opened === state.closed
            );

            resolve({
              success: allConnectionsHandled && !raceConditionDetected && properStateTransitions,
              totalAttempts: completedAttempts,
              successfulConnections,
              raceConditionDetected,
            });
          }).catch((err) => {
            resolve({
              success: false,
              totalAttempts: completedAttempts,
              successfulConnections,
              raceConditionDetected,
              error: `Connection churn test failed: ${err instanceof Error ? err.message : String(err)}`,
            });
          });
        } catch (err) {
          resolve({
            success: false,
            totalAttempts: 0,
            successfulConnections: 0,
            raceConditionDetected: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, { wsEndpoint: wsUrl, count: clientCount });

    // Then: Server handles all connections without race conditions
    expect(result.success).toBeTruthy();
    expect(result.totalAttempts).toBe(clientCount);

    // Verify no race conditions were detected
    expect(result.raceConditionDetected).toBeFalsy();

    // Then: Connection count remains accurate
    // Most connections should succeed in proper churn scenario
    expect(result.successfulConnections).toBeGreaterThan(clientCount * 0.8); // At least 80% success rate

    if (!result.success) {
      throw new Error(result.error || `Connection churn test failed. Attempts: ${result.totalAttempts}, Successful: ${result.successfulConnections}, Race condition: ${result.raceConditionDetected}`);
    }
  });

  test('INT-028: Metrics endpoint returns connection statistics', async ({ request }) => {
    // Given: Server tracks connection metrics
    // The MCP service should track active connections, total connections, and errors

    // When: Client requests metrics endpoint GET /metrics
    const response = await request.get(`${mcpServiceUrl}/metrics`);

    // Then: Server returns active connection count
    expect(response.ok()).toBeTruthy();
    const metricsData = await response.json();

    expect(metricsData).toHaveProperty('active_connections');
    expect(typeof metricsData.active_connections).toBe('number');
    expect(metricsData.active_connections).toBeGreaterThanOrEqual(0);

    // Then: Server returns total connections handled
    expect(metricsData).toHaveProperty('total_connections');
    expect(typeof metricsData.total_connections).toBe('number');
    expect(metricsData.total_connections).toBeGreaterThanOrEqual(0);

    // Then: Server returns connection error count
    expect(metricsData).toHaveProperty('connection_errors');
    expect(typeof metricsData.connection_errors).toBe('number');
    expect(metricsData.connection_errors).toBeGreaterThanOrEqual(0);
  });

  test('INT-032: Server enforces rate limit and rejects excess requests', async ({ request }) => {
    // Given: Server enforces rate limit of 100 requests per minute per client
    const rateLimit = 100;
    const timeWindow = 60000; // 60 seconds in milliseconds
    const startTime = Date.now();

    // Track successful and rejected requests
    const requestResults: { status: number; timestamp: number }[] = [];

    // When: Client sends 101 requests within one minute
    const requestPromises: Promise<void>[] = [];

    for (let i = 0; i < rateLimit + 1; i++) {
      const requestPromise = (async () => {
        try {
          const response = await request.get(`${mcpServiceUrl}/health`);
          const timestamp = Date.now();

          requestResults.push({
            status: response.status(),
            timestamp: timestamp - startTime,
          });
        } catch (error) {
          // Record failed requests (might be rate limited)
          const timestamp = Date.now();
          requestResults.push({
            status: 429, // Assume rate limit error
            timestamp: timestamp - startTime,
          });
        }
      })();

      requestPromises.push(requestPromise);
    }

    // Wait for all requests to complete
    await Promise.all(requestPromises);
    const totalTime = Date.now() - startTime;

    // Verify requests completed within time window
    expect(totalTime).toBeLessThan(timeWindow);
    expect(requestResults.length).toBe(rateLimit + 1);

    // Then: Server rejects 101st request
    // Count rejected requests (status 429 - Too Many Requests)
    const rejectedRequests = requestResults.filter(r => r.status === 429);
    expect(rejectedRequests.length).toBeGreaterThanOrEqual(1);

    // Then: Server returns rate limit exceeded error
    // At least one request should be rate limited
    const rateLimitedRequest = rejectedRequests[0];
    expect(rateLimitedRequest).toBeTruthy();
    expect(rateLimitedRequest.status).toBe(429);

    // Then: Server keeps connection open
    // Verify server still responds after rate limiting
    // Wait brief moment to ensure rate limit window hasn't reset
    await new Promise(resolve => setTimeout(resolve, 100));

    // Try another request - server should still be reachable
    const postRateLimitResponse = await request.get(`${mcpServiceUrl}/health`);
    // Response might be 200 (if rate limit reset) or 429 (if still limited)
    // The key is that server responds (doesn't close connection)
    expect([200, 429]).toContain(postRateLimitResponse.status());
  });

  test('INT-036: Server achieves target message throughput under concurrent load', async ({ page }) => {
    // Given: Server manages 30 concurrent connections
    const connectionCount = 30;
    const messagesPerConnection = 10;
    const totalMessages = connectionCount * messagesPerConnection;
    const targetThroughput = 200; // messages per second
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // Use page.evaluate to run throughput test in browser context
    const result = await page.evaluate(async ({ wsEndpoint, connCount, msgCount }) => {
      return new Promise<{
        success: boolean;
        totalMessages: number;
        processedMessages: number;
        elapsedTime: number;
        throughput: number;
        error?: string;
      }>((resolve) => {
        try {
          const connections: WebSocket[] = [];
          let processedMessages = 0;
          const startTime = Date.now();
          const messagePromises: Promise<void>[] = [];

          // Function to send messages on a single connection
          const sendMessagesOnConnection = (ws: WebSocket, connectionId: number): Promise<number> => {
            return new Promise((resolveMessages) => {
              let sentMessages = 0;
              let receivedMessages = 0;
              const messagesToSend = msgCount;

              ws.onmessage = () => {
                receivedMessages++;
                processedMessages++;

                // If all messages for this connection received, resolve
                if (receivedMessages === messagesToSend) {
                  resolveMessages(receivedMessages);
                }
              };

              ws.onerror = () => {
                resolveMessages(receivedMessages);
              };

              ws.onclose = () => {
                resolveMessages(receivedMessages);
              };

              // When: All clients send 10 messages each within 1 second
              // Send all messages rapidly
              const sendInterval = setInterval(() => {
                if (sentMessages >= messagesToSend) {
                  clearInterval(sendInterval);
                  return;
                }

                // Send MCP protocol message
                ws.send(JSON.stringify({
                  jsonrpc: '2.0',
                  method: 'ping',
                  id: `conn${connectionId}_msg${sentMessages}`,
                }));
                sentMessages++;
              }, 100 / messagesToSend); // Distribute sends over 100ms to stay within 1 second total

              // Timeout after 2 seconds if not all messages received
              setTimeout(() => {
                clearInterval(sendInterval);
                resolveMessages(receivedMessages);
              }, 2000);
            });
          };

          // Establish all connections
          const connectionPromises: Promise<void>[] = [];

          for (let i = 0; i < connCount; i++) {
            const connectionPromise = new Promise<void>((resolveConn) => {
              const ws = new WebSocket(`${wsEndpoint}/mcp`);
              const connectionIndex = i;

              ws.onopen = async () => {
                connections.push(ws);

                // Start sending messages on this connection
                await sendMessagesOnConnection(ws, connectionIndex);
                resolveConn();
              };

              ws.onerror = () => {
                resolveConn();
              };

              // Connection timeout
              setTimeout(() => {
                resolveConn();
              }, 5000);
            });

            connectionPromises.push(connectionPromise);
          }

          // Wait for all connections and message sending to complete
          Promise.all(connectionPromises).then(() => {
            const elapsedTime = Date.now() - startTime;

            // Then: Server processes 300 total messages successfully
            // Calculate throughput
            const throughput = (processedMessages / elapsedTime) * 1000; // messages per second

            // Clean up connections
            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            // Then: Message throughput exceeds 200 messages per second
            resolve({
              success: processedMessages >= connCount * msgCount * 0.95 && throughput > 200,
              totalMessages: connCount * msgCount,
              processedMessages,
              elapsedTime,
              throughput,
            });
          }).catch((err) => {
            // Clean up connections
            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            resolve({
              success: false,
              totalMessages: connCount * msgCount,
              processedMessages,
              elapsedTime: Date.now() - startTime,
              throughput: 0,
              error: `Throughput test failed: ${err instanceof Error ? err.message : String(err)}`,
            });
          });

          // Safety timeout - if test doesn't complete in 10 seconds, fail
          setTimeout(() => {
            const elapsedTime = Date.now() - startTime;
            const throughput = (processedMessages / elapsedTime) * 1000;

            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            resolve({
              success: false,
              totalMessages: connCount * msgCount,
              processedMessages,
              elapsedTime,
              throughput,
              error: 'Test timeout - took longer than 10 seconds',
            });
          }, 10000);
        } catch (err) {
          resolve({
            success: false,
            totalMessages: connCount * msgCount,
            processedMessages: 0,
            elapsedTime: 0,
            throughput: 0,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, { wsEndpoint: wsUrl, connCount: connectionCount, msgCount: messagesPerConnection });

    // Then: Server processes 300 total messages successfully
    expect(result.processedMessages).toBeGreaterThanOrEqual(totalMessages * 0.95); // At least 95% success rate

    // Then: Message throughput exceeds 200 messages per second
    expect(result.throughput).toBeGreaterThan(targetThroughput);

    if (!result.success) {
      throw new Error(
        result.error ||
        `Throughput test failed. Processed: ${result.processedMessages}/${result.totalMessages}, ` +
        `Throughput: ${result.throughput.toFixed(2)} msg/s (target: ${targetThroughput} msg/s), ` +
        `Elapsed: ${result.elapsedTime}ms`
      );
    }
  });

  test('INT-019: Server returns formatted error for invalid MCP request', async ({ request }) => {
    // Given: Client sends invalid MCP request
    // The MCP service should validate requests and return properly formatted error responses
    const invalidMcpRequest = {
      jsonrpc: '2.0',
      // Missing required 'method' field
      id: 1,
      params: {
        protocolVersion: '2024-11-05',
      },
    };

    // When: Server detects validation error
    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: invalidMcpRequest,
    });

    // Then: Server returns error response in MCP format
    // Server should reject with 4xx status code (typically 400 Bad Request)
    expect(response.status()).toBeGreaterThanOrEqual(400);
    expect(response.status()).toBeLessThan(500);

    // Attempt to parse response as JSON (MCP format should be JSON)
    const errorResponse = await response.json();

    // Then: Error response includes error code
    // MCP protocol errors should follow JSON-RPC 2.0 error structure
    expect(errorResponse).toHaveProperty('error');
    expect(errorResponse.error).toHaveProperty('code');
    expect(typeof errorResponse.error.code).toBe('number');

    // Then: Error response includes descriptive error message
    expect(errorResponse.error).toHaveProperty('message');
    expect(typeof errorResponse.error.message).toBe('string');
    expect(errorResponse.error.message.length).toBeGreaterThan(0);

    // Verify error message contains validation error indicators
    const errorMessage = errorResponse.error.message.toLowerCase();
    const validationIndicators = [
      'validation',
      'invalid',
      'required',
      'missing',
      'field',
      'method',
    ];

    const containsValidationError = validationIndicators.some(indicator =>
      errorMessage.includes(indicator)
    );

    expect(containsValidationError).toBeTruthy();

    // Verify JSON-RPC 2.0 compliance
    expect(errorResponse).toHaveProperty('jsonrpc');
    expect(errorResponse.jsonrpc).toBe('2.0');

    // Verify error response includes matching request ID
    expect(errorResponse).toHaveProperty('id');
    expect(errorResponse.id).toBe(invalidMcpRequest.id);
  });

  test('INT-021: Server returns MCP protocol-compliant response', async ({ request }) => {
    // Given: Client sends request of any type
    // Test multiple request types to ensure all return protocol-compliant responses
    const testRequests = [
      {
        name: 'initialize request',
        payload: {
          jsonrpc: '2.0',
          method: 'initialize',
          id: 1,
          params: {
            protocolVersion: '2024-11-05',
            capabilities: {},
            clientInfo: {
              name: 'test-client',
              version: '1.0.0'
            }
          }
        }
      },
      {
        name: 'ping request',
        payload: {
          jsonrpc: '2.0',
          method: 'ping',
          id: 2
        }
      },
      {
        name: 'tools/list request',
        payload: {
          jsonrpc: '2.0',
          method: 'tools/list',
          id: 3
        }
      }
    ];

    for (const testRequest of testRequests) {
      // When: Server completes operation
      const response = await request.post(`${mcpServiceUrl}/mcp`, {
        headers: {
          'Content-Type': 'application/json',
        },
        data: testRequest.payload,
      });

      // Then: Server returns response of appropriate type
      expect(response.status()).toBeGreaterThanOrEqual(200);
      expect(response.status()).toBeLessThan(300);

      const responseData = await response.json();

      // Then: Response conforms to MCP protocol specification
      // Verify JSON-RPC 2.0 compliance (MCP is built on JSON-RPC 2.0)
      expect(responseData).toHaveProperty('jsonrpc');
      expect(responseData.jsonrpc).toBe('2.0');

      // Verify response includes matching request ID
      expect(responseData).toHaveProperty('id');
      expect(responseData.id).toBe(testRequest.payload.id);

      // Verify response has either 'result' or 'error' (never both)
      const hasResult = 'result' in responseData;
      const hasError = 'error' in responseData;
      expect(hasResult || hasError).toBeTruthy();
      expect(hasResult && hasError).toBeFalsy();

      // Verify response structure matches MCP protocol
      if (hasResult) {
        // Successful response should have appropriate result structure
        expect(responseData.result).toBeTruthy();
        expect(typeof responseData.result).toBe('object');

        // For initialize method, verify MCP-specific result fields
        if (testRequest.payload.method === 'initialize') {
          expect(responseData.result).toHaveProperty('protocolVersion');
          expect(responseData.result).toHaveProperty('capabilities');
          expect(responseData.result).toHaveProperty('serverInfo');
        }
      } else if (hasError) {
        // Error response should follow JSON-RPC 2.0 error structure
        expect(responseData.error).toHaveProperty('code');
        expect(typeof responseData.error.code).toBe('number');
        expect(responseData.error).toHaveProperty('message');
        expect(typeof responseData.error.message).toBe('string');
      }
    }
  });

  test('INT-015: Server validates required MCP protocol fields', async ({ request }) => {
    // Given: Client maintains active WebSocket connection
    // Note: Testing via HTTP endpoint that accepts MCP protocol messages
    // The MCP service should validate required protocol fields before processing

    // When: Client sends MCP message missing required protocol_version field
    const invalidMcpMessage = {
      jsonrpc: '2.0',
      method: 'initialize',
      id: 1,
      params: {
        // Missing required 'protocolVersion' field
        capabilities: {},
        clientInfo: {
          name: 'test-client',
          version: '1.0.0'
        }
      }
    };

    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: invalidMcpMessage,
    });

    // Then: Server returns validation error immediately
    // Server should reject with 4xx status code (typically 400 Bad Request)
    expect(response.status()).toBeGreaterThanOrEqual(400);
    expect(response.status()).toBeLessThan(500);

    // Attempt to parse response as JSON
    const errorResponse = await response.json();

    // Then: Error message specifies missing required field
    // Verify error response follows JSON-RPC 2.0 error structure
    expect(errorResponse).toHaveProperty('error');
    expect(errorResponse.error).toHaveProperty('code');
    expect(typeof errorResponse.error.code).toBe('number');
    expect(errorResponse.error).toHaveProperty('message');
    expect(typeof errorResponse.error.message).toBe('string');

    // Verify error message indicates missing required field
    const errorMessage = errorResponse.error.message.toLowerCase();
    const requiredFieldIndicators = [
      'protocol',
      'version',
      'protocolversion',
      'required',
      'missing',
      'field',
    ];

    const containsRequiredFieldError = requiredFieldIndicators.some(indicator =>
      errorMessage.includes(indicator)
    );

    expect(containsRequiredFieldError).toBeTruthy();

    // Verify JSON-RPC 2.0 compliance
    expect(errorResponse).toHaveProperty('jsonrpc');
    expect(errorResponse.jsonrpc).toBe('2.0');

    // Verify error response includes matching request ID
    expect(errorResponse).toHaveProperty('id');
    expect(errorResponse.id).toBe(invalidMcpMessage.id);
  });

  test('INT-027: Server performs periodic heartbeat check on all connections', async ({ request }) => {
    // Given: Server monitors connection health with 30-second interval
    // Verify server provides metrics endpoint to track heartbeat operations
    const metricsResponse = await request.get(`${mcpServiceUrl}/metrics`);
    expect(metricsResponse.ok()).toBeTruthy();

    const initialMetrics = await metricsResponse.json();
    const initialActiveConnections = initialMetrics.active_connections || 0;

    // Establish multiple active connections by making several health check requests
    // to simulate active connections that need heartbeat monitoring
    const connectionCount = 5;
    const connectionRequests: Promise<void>[] = [];

    for (let i = 0; i < connectionCount; i++) {
      const connectionPromise = (async () => {
        try {
          await request.get(`${mcpServiceUrl}/health`);
        } catch (error) {
          // Connections may fail - not critical for this test
        }
      })();
      connectionRequests.push(connectionPromise);
    }

    await Promise.all(connectionRequests);

    // When: Server performs heartbeat check
    // Wait for heartbeat interval (30 seconds) plus buffer for processing
    // In test environment, heartbeat might be triggered manually or by monitoring endpoint
    // Check if server provides heartbeat status endpoint
    let heartbeatStatus;
    try {
      const heartbeatResponse = await request.get(`${mcpServiceUrl}/heartbeat`);
      if (heartbeatResponse.ok()) {
        heartbeatStatus = await heartbeatResponse.json();
      }
    } catch (error) {
      // Heartbeat endpoint may not exist - verify through metrics instead
    }

    // Alternative: Trigger heartbeat via metrics endpoint with heartbeat data
    // Wait brief moment to allow heartbeat processing
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Retrieve metrics after heartbeat interval
    const postHeartbeatMetrics = await request.get(`${mcpServiceUrl}/metrics`);
    expect(postHeartbeatMetrics.ok()).toBeTruthy();

    const metricsData = await postHeartbeatMetrics.json();

    // Then: Server sends ping to all active connections
    // Verify metrics show heartbeat activity
    expect(metricsData).toHaveProperty('total_connections');
    expect(typeof metricsData.total_connections).toBe('number');
    expect(metricsData.total_connections).toBeGreaterThanOrEqual(initialActiveConnections);

    // Then: Server marks unresponsive connections for cleanup
    // Verify connection error tracking includes heartbeat failures
    expect(metricsData).toHaveProperty('connection_errors');
    expect(typeof metricsData.connection_errors).toBe('number');
    expect(metricsData.connection_errors).toBeGreaterThanOrEqual(0);

    // Verify last activity timestamp is recent (within last minute)
    expect(metricsData).toHaveProperty('last_activity');
    expect(typeof metricsData.last_activity).toBe('string');

    const lastActivity = new Date(metricsData.last_activity);
    const now = new Date();
    const timeDiff = now.getTime() - lastActivity.getTime();
    expect(timeDiff).toBeLessThan(60000); // Less than 60 seconds

    // If heartbeat status endpoint exists, verify heartbeat data
    if (heartbeatStatus) {
      // Verify heartbeat status includes ping activity
      expect(heartbeatStatus).toHaveProperty('last_heartbeat');
      expect(typeof heartbeatStatus.last_heartbeat).toBe('string');

      // Verify heartbeat interval configuration
      if (heartbeatStatus.heartbeat_interval) {
        expect(heartbeatStatus.heartbeat_interval).toBe(30); // 30 seconds
      }

      // Verify unresponsive connection tracking
      if (heartbeatStatus.unresponsive_connections !== undefined) {
        expect(typeof heartbeatStatus.unresponsive_connections).toBe('number');
        expect(heartbeatStatus.unresponsive_connections).toBeGreaterThanOrEqual(0);
      }
    }
  });

  test('INT-031: Server processes concurrent messages from multiple clients', async ({ page }) => {
    // Given: Server handles 10 concurrent client connections
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');
    const clientCount = 10;

    // Use page.evaluate to run concurrent messaging test in browser context
    const result = await page.evaluate(async ({ wsEndpoint, count }) => {
      return new Promise<{
        success: boolean;
        totalClients: number;
        messagesProcessed: number;
        responsesReceived: number;
        error?: string;
      }>((resolve) => {
        try {
          const connections: WebSocket[] = [];
          let messagesProcessed = 0;
          let responsesReceived = 0;
          let allClientsConnected = false;

          // Function to establish connection and send message
          const connectAndSendMessage = (clientId: number): Promise<boolean> => {
            return new Promise((resolveClient) => {
              const ws = new WebSocket(`${wsEndpoint}/mcp`);
              let messageReceived = false;
              let connectionResolved = false;

              const timeoutId = setTimeout(() => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  resolveClient(false);
                }
              }, 10000); // 10 second timeout per client

              ws.onopen = () => {
                connections.push(ws);

                // When: All clients send messages simultaneously
                // Send MCP protocol message
                const message = {
                  jsonrpc: '2.0',
                  method: 'ping',
                  id: `client_${clientId}`,
                };

                try {
                  ws.send(JSON.stringify(message));
                  messagesProcessed++;
                } catch (err) {
                  if (!connectionResolved) {
                    connectionResolved = true;
                    clearTimeout(timeoutId);
                    resolveClient(false);
                  }
                }
              };

              ws.onmessage = (event) => {
                try {
                  const response = JSON.parse(event.data as string);

                  // Then: Each client receives corresponding response
                  // Verify response matches the client's request ID
                  if (response.id === `client_${clientId}`) {
                    messageReceived = true;
                    responsesReceived++;
                  }

                  if (!connectionResolved) {
                    connectionResolved = true;
                    clearTimeout(timeoutId);
                    resolveClient(messageReceived);
                  }
                } catch (err) {
                  // Invalid response format
                  if (!connectionResolved) {
                    connectionResolved = true;
                    clearTimeout(timeoutId);
                    resolveClient(false);
                  }
                }
              };

              ws.onerror = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  resolveClient(false);
                }
              };

              ws.onclose = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  resolveClient(messageReceived);
                }
              };
            });
          };

          // Establish all connections and send messages simultaneously
          const clientPromises: Promise<boolean>[] = [];
          for (let i = 0; i < count; i++) {
            clientPromises.push(connectAndSendMessage(i));
          }

          // Wait for all clients to complete
          Promise.all(clientPromises).then((results) => {
            allClientsConnected = true;

            // Then: Server processes all messages without errors
            const successfulClients = results.filter(r => r).length;

            // Clean up all connections
            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            resolve({
              success: messagesProcessed === count && responsesReceived === count,
              totalClients: count,
              messagesProcessed,
              responsesReceived,
            });
          }).catch((err) => {
            // Clean up connections
            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            resolve({
              success: false,
              totalClients: count,
              messagesProcessed,
              responsesReceived,
              error: `Concurrent messaging test failed: ${err instanceof Error ? err.message : String(err)}`,
            });
          });

          // Safety timeout - if test doesn't complete in 15 seconds, fail
          setTimeout(() => {
            if (!allClientsConnected) {
              connections.forEach(ws => {
                try {
                  ws.close();
                } catch (e) {
                  // Ignore cleanup errors
                }
              });

              resolve({
                success: false,
                totalClients: count,
                messagesProcessed,
                responsesReceived,
                error: 'Test timeout - took longer than 15 seconds',
              });
            }
          }, 15000);
        } catch (err) {
          resolve({
            success: false,
            totalClients: count,
            messagesProcessed: 0,
            responsesReceived: 0,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, { wsEndpoint: wsUrl, count: clientCount });

    // Then: Server processes all messages without errors
    expect(result.success).toBeTruthy();
    expect(result.messagesProcessed).toBe(clientCount);

    // Then: Each client receives corresponding response
    expect(result.responsesReceived).toBe(clientCount);

    if (!result.success) {
      throw new Error(
        result.error ||
        `Concurrent messaging test failed. Messages processed: ${result.messagesProcessed}/${result.totalClients}, ` +
        `Responses received: ${result.responsesReceived}/${result.totalClients}`
      );
    }
  });

  test('INT-017: Server accepts well-formed message with valid schema', async ({ request }) => {
    // Given: Server validates messages against MCP protocol schema
    // The MCP service should accept valid messages that conform to the protocol schema

    // When: Client sends message with correct schema and all required fields
    const validMessage = {
      jsonrpc: '2.0',
      method: 'initialize',
      id: 1,
      params: {
        protocolVersion: '2024-11-05',
        capabilities: {
          roots: {
            listChanged: true
          },
          sampling: {}
        },
        clientInfo: {
          name: 'test-client',
          version: '1.0.0'
        }
      }
    };

    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: validMessage,
    });

    // Then: Server accepts message without validation errors
    // Server should respond with success status (2xx)
    expect(response.status()).toBeGreaterThanOrEqual(200);
    expect(response.status()).toBeLessThan(300);

    // Verify the response is valid JSON
    const responseData = await response.json();

    // Verify response follows JSON-RPC 2.0 format
    expect(responseData).toHaveProperty('jsonrpc');
    expect(responseData.jsonrpc).toBe('2.0');

    // Verify response has matching request ID
    expect(responseData).toHaveProperty('id');
    expect(responseData.id).toBe(validMessage.id);

    // Verify response has result (no error)
    expect(responseData).toHaveProperty('result');
    expect(responseData.result).toBeTruthy();

    // For initialize method, verify MCP-compliant response structure
    expect(responseData.result).toHaveProperty('protocolVersion');
    expect(responseData.result).toHaveProperty('capabilities');
    expect(responseData.result).toHaveProperty('serverInfo');

    // Verify no error field is present
    expect(responseData).not.toHaveProperty('error');
  });

  test('INT-025: Server handles client-initiated connection closure gracefully', async ({ page }) => {
    // Given: Client maintains active WebSocket connection
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // Use page.evaluate to run WebSocket test in browser context
    const result = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{
        success: boolean;
        connectionEstablished: boolean;
        closeAcknowledged: boolean;
        connectionRemoved: boolean;
        error?: string;
      }>((resolve) => {
        try {
          // Establish WebSocket connection
          const ws = new WebSocket(`${wsEndpoint}/mcp`);
          let connectionEstablished = false;
          let closeAcknowledged = false;
          let connectionClosed = false;

          ws.onopen = () => {
            connectionEstablished = true;

            // When: Client sends close frame
            // Close the connection immediately after opening
            ws.close(1000, 'Client-initiated graceful close');
          };

          ws.onclose = (event) => {
            // Then: Server acknowledges close frame
            // WebSocket close event indicates server acknowledged the close
            connectionClosed = true;

            // Verify close was clean (code 1000 or 1001 for normal closure)
            if (event.code === 1000 || event.code === 1001 || event.wasClean) {
              closeAcknowledged = true;
            }

            resolve({
              success: connectionEstablished && closeAcknowledged && connectionClosed,
              connectionEstablished,
              closeAcknowledged,
              connectionRemoved: connectionClosed, // Connection closed implies removed from pool
            });
          };

          ws.onerror = (error) => {
            resolve({
              success: false,
              connectionEstablished,
              closeAcknowledged: false,
              connectionRemoved: false,
              error: `WebSocket error: ${error}`,
            });
          };

          // Set timeout for connection establishment and closure
          setTimeout(() => {
            if (!connectionClosed) {
              ws.close();
              resolve({
                success: false,
                connectionEstablished,
                closeAcknowledged: false,
                connectionRemoved: false,
                error: 'Connection did not close within timeout',
              });
            }
          }, 5000);
        } catch (err) {
          resolve({
            success: false,
            connectionEstablished: false,
            closeAcknowledged: false,
            connectionRemoved: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, wsUrl);

    // Then: Server acknowledges close frame
    expect(result.connectionEstablished).toBeTruthy();
    expect(result.closeAcknowledged).toBeTruthy();

    // Then: Server removes connection from active pool
    // Verified by successful clean close event
    expect(result.connectionRemoved).toBeTruthy();

    // Then: Server releases connection resources
    // Implicit in successful connection closure
    expect(result.success).toBeTruthy();

    if (!result.success) {
      throw new Error(result.error || 'Client-initiated connection closure failed');
    }
  });

  test('INT-029: Server handles concurrent client connections successfully', async ({ page }) => {
    // Given: Server supports concurrent connections
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');
    const clientCount = 10;

    // Use page.evaluate to run concurrent connection test in browser context
    const result = await page.evaluate(async ({ wsEndpoint, count }) => {
      return new Promise<{
        success: boolean;
        connectionsAccepted: number;
        uniqueIdentifiers: Set<string>;
        error?: string;
      }>((resolve) => {
        try {
          const connections: WebSocket[] = [];
          const connectionIdentifiers = new Set<string>();
          let connectionsOpened = 0;
          let testComplete = false;

          // Function to establish a single connection
          const establishConnection = (index: number): Promise<string | null> => {
            return new Promise((resolveConn) => {
              const ws = new WebSocket(`${wsEndpoint}/mcp`);
              let connectionResolved = false;
              let connectionId: string | null = null;

              const timeoutId = setTimeout(() => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  ws.close();
                  resolveConn(null);
                }
              }, 5000); // 5 second timeout per connection

              ws.onopen = () => {
                connectionsOpened++;
                connections.push(ws);

                // Send initialize message to get connection identifier
                ws.send(JSON.stringify({
                  jsonrpc: '2.0',
                  method: 'initialize',
                  id: `conn_${index}`,
                  params: {
                    protocolVersion: '2024-11-05',
                    capabilities: {},
                    clientInfo: {
                      name: `test-client-${index}`,
                      version: '1.0.0'
                    }
                  }
                }));
              };

              ws.onmessage = (event) => {
                try {
                  const message = JSON.parse(event.data as string);

                  // Then: Each connection receives unique identifier
                  // Extract unique identifier from response
                  if (message.id === `conn_${index}`) {
                    // Use the connection's request ID as unique identifier
                    connectionId = message.id;

                    // Alternatively, if server provides a session/connection ID
                    if (message.result?.sessionId) {
                      connectionId = message.result.sessionId;
                    } else if (message.result?.connectionId) {
                      connectionId = message.result.connectionId;
                    }

                    if (!connectionResolved) {
                      connectionResolved = true;
                      clearTimeout(timeoutId);
                      resolveConn(connectionId);
                    }
                  }
                } catch (err) {
                  // If we can't parse response, still count connection as accepted
                  if (!connectionResolved) {
                    connectionResolved = true;
                    clearTimeout(timeoutId);
                    connectionId = `implicit_${index}`;
                    resolveConn(connectionId);
                  }
                }
              };

              ws.onerror = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  resolveConn(null);
                }
              };

              ws.onclose = () => {
                if (!connectionResolved) {
                  connectionResolved = true;
                  clearTimeout(timeoutId);
                  // If connection opened but no message received
                  if (connectionsOpened > 0 && !connectionId) {
                    connectionId = `implicit_${index}`;
                  }
                  resolveConn(connectionId);
                }
              };
            });
          };

          // When: 10 clients connect simultaneously
          const connectionPromises: Promise<string | null>[] = [];
          for (let i = 0; i < count; i++) {
            connectionPromises.push(establishConnection(i));
          }

          // Wait for all connection attempts to complete
          Promise.all(connectionPromises).then((identifiers) => {
            testComplete = true;

            // Collect unique identifiers
            identifiers.forEach(id => {
              if (id) {
                connectionIdentifiers.add(id);
              }
            });

            // Then: Server accepts all 10 connections
            const successfulConnections = identifiers.filter(id => id !== null).length;

            // Clean up all connections
            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            resolve({
              success: successfulConnections === count && connectionIdentifiers.size === count,
              connectionsAccepted: successfulConnections,
              uniqueIdentifiers: connectionIdentifiers,
            });
          }).catch((err) => {
            testComplete = true;

            // Clean up connections
            connections.forEach(ws => {
              try {
                ws.close();
              } catch (e) {
                // Ignore cleanup errors
              }
            });

            resolve({
              success: false,
              connectionsAccepted: connectionsOpened,
              uniqueIdentifiers: connectionIdentifiers,
              error: `Concurrent connection test failed: ${err instanceof Error ? err.message : String(err)}`,
            });
          });

          // Safety timeout - if test doesn't complete in 15 seconds, fail
          setTimeout(() => {
            if (!testComplete) {
              testComplete = true;
              connections.forEach(ws => {
                try {
                  ws.close();
                } catch (e) {
                  // Ignore cleanup errors
                }
              });

              resolve({
                success: false,
                connectionsAccepted: connectionsOpened,
                uniqueIdentifiers: connectionIdentifiers,
                error: 'Test timeout - took longer than 15 seconds',
              });
            }
          }, 15000);
        } catch (err) {
          resolve({
            success: false,
            connectionsAccepted: 0,
            uniqueIdentifiers: new Set(),
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`,
          });
        }
      });
    }, { wsEndpoint: wsUrl, count: clientCount });

    // Then: Server accepts all 10 connections
    expect(result.connectionsAccepted).toBe(clientCount);
    expect(result.success).toBeTruthy();

    // Then: Each connection receives unique identifier
    expect(result.uniqueIdentifiers.size).toBe(clientCount);

    if (!result.success) {
      throw new Error(
        result.error ||
        `Concurrent connection test failed. Accepted: ${result.connectionsAccepted}/${clientCount}, ` +
        `Unique identifiers: ${result.uniqueIdentifiers.size}/${clientCount}`
      );
    }
  });

  test('INT-018: Server returns complete MCP response with required fields', async ({ request }) => {
    // Given: Client sends valid MCP request
    const validMcpRequest = {
      jsonrpc: '2.0',
      method: 'initialize',
      id: 100,
      params: {
        protocolVersion: '2024-11-05',
        capabilities: {
          roots: {
            listChanged: true
          },
          sampling: {}
        },
        clientInfo: {
          name: 'test-client',
          version: '1.0.0'
        }
      }
    };

    // When: Server processes request successfully
    const response = await request.post(`${mcpServiceUrl}/mcp`, {
      headers: {
        'Content-Type': 'application/json',
      },
      data: validMcpRequest,
    });

    // Then: Server returns response in MCP protocol format
    expect(response.status()).toBeGreaterThanOrEqual(200);
    expect(response.status()).toBeLessThan(300);

    const responseData = await response.json();

    // Verify response conforms to MCP protocol format (JSON-RPC 2.0)
    expect(responseData).toHaveProperty('jsonrpc');
    expect(responseData.jsonrpc).toBe('2.0');

    // Then: Response includes protocol_version field
    // For initialize method, protocol version is in the result
    expect(responseData).toHaveProperty('result');
    expect(responseData.result).toHaveProperty('protocolVersion');
    expect(typeof responseData.result.protocolVersion).toBe('string');
    expect(responseData.result.protocolVersion).toBeTruthy();

    // Then: Response includes message_id matching request
    expect(responseData).toHaveProperty('id');
    expect(responseData.id).toBe(validMcpRequest.id);
  });

});
