import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { mcpServiceUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('MCP Service E2E Tests', () => {

  test('E2E-011: MCP client completes handshake with server successfully', async ({ page }) => {
    // Given: MCP server runs on configured endpoint
    // Convert HTTP URL to WebSocket URL
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // When: Client completes MCP handshake from connection to first message
    const handshakeResult = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{
        success: boolean;
        connectionEstablished: boolean;
        initializeRequestSent: boolean;
        initializeResponseReceived: boolean;
        error?: string;
      }>((resolve) => {
        try {
          const ws = new WebSocket(`${wsEndpoint}/mcp`);
          let connectionEstablished = false;
          let initializeRequestSent = false;
          let initializeResponseReceived = false;
          let handshakeTimeout: NodeJS.Timeout;

          ws.onopen = () => {
            // Then: Client establishes WebSocket connection successfully
            connectionEstablished = true;

            // Then: Client sends initialize request in MCP format
            const initializeRequest = {
              jsonrpc: '2.0',
              method: 'initialize',
              id: 1,
              params: {
                protocolVersion: '2024-11-05',
                capabilities: {},
                clientInfo: {
                  name: 'playwright-e2e-test',
                  version: '1.0.0'
                }
              }
            };

            ws.send(JSON.stringify(initializeRequest));
            initializeRequestSent = true;

            // Set timeout for initialize response (should be quick)
            handshakeTimeout = setTimeout(() => {
              if (!initializeResponseReceived) {
                ws.close();
                resolve({
                  success: false,
                  connectionEstablished,
                  initializeRequestSent,
                  initializeResponseReceived,
                  error: 'Initialize response timeout'
                });
              }
            }, 5000); // 5 second timeout
          };

          ws.onmessage = (event) => {
            try {
              const message = JSON.parse(event.data as string);

              // Then: Server returns initialize response in MCP format
              if (message.id === 1 && message.result) {
                initializeResponseReceived = true;
                clearTimeout(handshakeTimeout);
                ws.close();
                resolve({
                  success: true,
                  connectionEstablished,
                  initializeRequestSent,
                  initializeResponseReceived
                });
              }
            } catch (err) {
              clearTimeout(handshakeTimeout);
              ws.close();
              resolve({
                success: false,
                connectionEstablished,
                initializeRequestSent,
                initializeResponseReceived,
                error: `Message parse error: ${err instanceof Error ? err.message : String(err)}`
              });
            }
          };

          ws.onerror = (error) => {
            clearTimeout(handshakeTimeout);
            resolve({
              success: false,
              connectionEstablished,
              initializeRequestSent,
              initializeResponseReceived,
              error: `WebSocket error: ${error}`
            });
          };

          ws.onclose = () => {
            clearTimeout(handshakeTimeout);
            if (!initializeResponseReceived) {
              resolve({
                success: false,
                connectionEstablished,
                initializeRequestSent,
                initializeResponseReceived,
                error: 'Connection closed without initialize response'
              });
            }
          };
        } catch (err) {
          resolve({
            success: false,
            connectionEstablished: false,
            initializeRequestSent: false,
            initializeResponseReceived: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`
          });
        }
      });
    }, wsUrl);

    // Then: Client establishes WebSocket connection successfully
    expect(handshakeResult.connectionEstablished).toBeTruthy();

    // Then: Client sends initialize request in MCP format
    expect(handshakeResult.initializeRequestSent).toBeTruthy();

    // Then: Server returns initialize response in MCP format
    expect(handshakeResult.initializeResponseReceived).toBeTruthy();

    // Overall handshake success
    expect(handshakeResult.success).toBeTruthy();
    if (!handshakeResult.success) {
      throw new Error(handshakeResult.error || 'Unknown MCP handshake failure');
    }
  });

  test('E2E-012: Server processes complete MCP request-response cycle with validation', async ({ page }) => {
    // Given: Client maintains active MCP connection
    // Convert HTTP URL to WebSocket URL
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // When: Client executes complete request-response cycle
    const requestResponseResult = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{
        success: boolean;
        requestValidated: boolean;
        requestProcessed: boolean;
        responseConformsToMCP: boolean;
        responseHasRequiredFields: boolean;
        error?: string;
      }>((resolve) => {
        try {
          const ws = new WebSocket(`${wsEndpoint}/mcp`);
          let requestValidated = false;
          let requestProcessed = false;
          let responseConformsToMCP = false;
          let responseHasRequiredFields = false;
          let requestTimeout: NodeJS.Timeout;

          ws.onopen = () => {
            // First complete handshake, then send a request
            const initializeRequest = {
              jsonrpc: '2.0',
              method: 'initialize',
              id: 1,
              params: {
                protocolVersion: '2024-11-05',
                capabilities: {},
                clientInfo: {
                  name: 'playwright-e2e-test',
                  version: '1.0.0'
                }
              }
            };

            ws.send(JSON.stringify(initializeRequest));

            // Set timeout for request processing
            requestTimeout = setTimeout(() => {
              ws.close();
              resolve({
                success: false,
                requestValidated,
                requestProcessed,
                responseConformsToMCP,
                responseHasRequiredFields,
                error: 'Request-response cycle timeout'
              });
            }, 10000); // 10 second timeout
          };

          let handshakeComplete = false;

          ws.onmessage = (event) => {
            try {
              const message = JSON.parse(event.data as string);

              // Handle handshake response first
              if (!handshakeComplete && message.id === 1 && message.result) {
                handshakeComplete = true;

                // Now send a complete MCP request
                const testRequest = {
                  jsonrpc: '2.0',
                  method: 'tools/list',
                  id: 2,
                  params: {}
                };

                ws.send(JSON.stringify(testRequest));
                return;
              }

              // Then: Server validates request against MCP schema
              // If we receive a response (not an error), request was validated
              if (message.id === 2) {
                requestValidated = true;

                // Then: Server processes request according to MCP specification
                // Check response has either result or error (MCP spec)
                if (message.result !== undefined || message.error !== undefined) {
                  requestProcessed = true;

                  // Then: Server returns response conforming to MCP protocol
                  // Check for required MCP fields: jsonrpc, id
                  if (message.jsonrpc === '2.0' && message.id === 2) {
                    responseConformsToMCP = true;

                    // Check for response structure (result or error object)
                    if (message.result || (message.error && message.error.code && message.error.message)) {
                      responseHasRequiredFields = true;
                    }
                  }
                }

                clearTimeout(requestTimeout);
                ws.close();
                resolve({
                  success: requestValidated && requestProcessed && responseConformsToMCP && responseHasRequiredFields,
                  requestValidated,
                  requestProcessed,
                  responseConformsToMCP,
                  responseHasRequiredFields
                });
              }
            } catch (err) {
              clearTimeout(requestTimeout);
              ws.close();
              resolve({
                success: false,
                requestValidated,
                requestProcessed,
                responseConformsToMCP,
                responseHasRequiredFields,
                error: `Message parse error: ${err instanceof Error ? err.message : String(err)}`
              });
            }
          };

          ws.onerror = (error) => {
            clearTimeout(requestTimeout);
            resolve({
              success: false,
              requestValidated,
              requestProcessed,
              responseConformsToMCP,
              responseHasRequiredFields,
              error: `WebSocket error: ${error}`
            });
          };

          ws.onclose = () => {
            clearTimeout(requestTimeout);
            if (!responseHasRequiredFields) {
              resolve({
                success: false,
                requestValidated,
                requestProcessed,
                responseConformsToMCP,
                responseHasRequiredFields,
                error: 'Connection closed before complete response'
              });
            }
          };
        } catch (err) {
          resolve({
            success: false,
            requestValidated: false,
            requestProcessed: false,
            responseConformsToMCP: false,
            responseHasRequiredFields: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`
          });
        }
      });
    }, wsUrl);

    // Then: Server validates request against MCP schema
    expect(requestResponseResult.requestValidated).toBeTruthy();

    // Then: Server processes request according to MCP specification
    expect(requestResponseResult.requestProcessed).toBeTruthy();

    // Then: Server returns response conforming to MCP protocol
    expect(requestResponseResult.responseConformsToMCP).toBeTruthy();
    expect(requestResponseResult.responseHasRequiredFields).toBeTruthy();

    // Overall request-response cycle success
    expect(requestResponseResult.success).toBeTruthy();
    if (!requestResponseResult.success) {
      throw new Error(requestResponseResult.error || 'Unknown MCP request-response cycle failure');
    }
  });

  test('E2E-013: Server handles invalid requests with fail-fast error response', async ({ page }) => {
    // Given: Client maintains active MCP connection
    // Convert HTTP URL to WebSocket URL
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // When: Client sends invalid request triggering error flow
    const errorHandlingResult = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{
        success: boolean;
        errorDetectedImmediately: boolean;
        errorResponseInMCPFormat: boolean;
        connectionRemainsOpen: boolean;
        subsequentRequestSucceeds: boolean;
        errorResponseTime?: number;
        error?: string;
      }>((resolve) => {
        try {
          const ws = new WebSocket(`${wsEndpoint}/mcp`);
          let errorDetectedImmediately = false;
          let errorResponseInMCPFormat = false;
          let connectionRemainsOpen = false;
          let subsequentRequestSucceeds = false;
          let errorResponseTime = 0;
          let errorRequestStartTime = 0;
          let requestTimeout: NodeJS.Timeout;

          ws.onopen = () => {
            // First complete handshake
            const initializeRequest = {
              jsonrpc: '2.0',
              method: 'initialize',
              id: 1,
              params: {
                protocolVersion: '2024-11-05',
                capabilities: {},
                clientInfo: {
                  name: 'playwright-e2e-test',
                  version: '1.0.0'
                }
              }
            };

            ws.send(JSON.stringify(initializeRequest));

            // Set timeout for overall test
            requestTimeout = setTimeout(() => {
              ws.close();
              resolve({
                success: false,
                errorDetectedImmediately,
                errorResponseInMCPFormat,
                connectionRemainsOpen,
                subsequentRequestSucceeds,
                errorResponseTime,
                error: 'Test timeout'
              });
            }, 15000); // 15 second timeout
          };

          let handshakeComplete = false;
          let errorReceived = false;

          ws.onmessage = (event) => {
            try {
              const message = JSON.parse(event.data as string);

              // Handle handshake response
              if (!handshakeComplete && message.id === 1 && message.result) {
                handshakeComplete = true;

                // Now send an INVALID request (missing required fields)
                const invalidRequest = {
                  jsonrpc: '2.0',
                  // Missing 'method' field - this should trigger validation error
                  id: 2,
                  params: {}
                };

                errorRequestStartTime = Date.now();
                ws.send(JSON.stringify(invalidRequest));
                return;
              }

              // Handle error response to invalid request
              if (!errorReceived && message.id === 2) {
                errorResponseTime = Date.now() - errorRequestStartTime;

                // Then: Server detects error immediately per fail-fast principle
                // Immediate = within 1000ms (1 second)
                if (errorResponseTime < 1000) {
                  errorDetectedImmediately = true;
                }

                // Then: Server returns MCP-formatted error response
                // Check for MCP error format: jsonrpc, id, error object with code and message
                if (message.jsonrpc === '2.0' &&
                    message.error &&
                    message.error.code !== undefined &&
                    message.error.message) {
                  errorResponseInMCPFormat = true;
                }

                errorReceived = true;

                // Then: Server keeps connection open for subsequent requests
                // Verify connection is still open by sending a valid request
                const validRequest = {
                  jsonrpc: '2.0',
                  method: 'tools/list',
                  id: 3,
                  params: {}
                };

                // Small delay to ensure connection state is stable
                setTimeout(() => {
                  try {
                    ws.send(JSON.stringify(validRequest));
                    connectionRemainsOpen = true;
                  } catch (err) {
                    connectionRemainsOpen = false;
                  }
                }, 100);

                return;
              }

              // Handle response to subsequent valid request
              if (errorReceived && message.id === 3) {
                // If we receive a response to the subsequent request, connection is definitely open
                subsequentRequestSucceeds = true;

                clearTimeout(requestTimeout);
                ws.close();
                resolve({
                  success: errorDetectedImmediately && errorResponseInMCPFormat && connectionRemainsOpen && subsequentRequestSucceeds,
                  errorDetectedImmediately,
                  errorResponseInMCPFormat,
                  connectionRemainsOpen,
                  subsequentRequestSucceeds,
                  errorResponseTime
                });
              }
            } catch (err) {
              clearTimeout(requestTimeout);
              ws.close();
              resolve({
                success: false,
                errorDetectedImmediately,
                errorResponseInMCPFormat,
                connectionRemainsOpen,
                subsequentRequestSucceeds,
                errorResponseTime,
                error: `Message parse error: ${err instanceof Error ? err.message : String(err)}`
              });
            }
          };

          ws.onerror = (error) => {
            clearTimeout(requestTimeout);
            resolve({
              success: false,
              errorDetectedImmediately,
              errorResponseInMCPFormat,
              connectionRemainsOpen,
              subsequentRequestSucceeds,
              errorResponseTime,
              error: `WebSocket error: ${error}`
            });
          };

          ws.onclose = () => {
            clearTimeout(requestTimeout);
            // If connection closed before test completed, this is a failure
            if (!subsequentRequestSucceeds) {
              resolve({
                success: false,
                errorDetectedImmediately,
                errorResponseInMCPFormat,
                connectionRemainsOpen,
                subsequentRequestSucceeds,
                errorResponseTime,
                error: 'Connection closed before test completion'
              });
            }
          };
        } catch (err) {
          resolve({
            success: false,
            errorDetectedImmediately: false,
            errorResponseInMCPFormat: false,
            connectionRemainsOpen: false,
            subsequentRequestSucceeds: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`
          });
        }
      });
    }, wsUrl);

    // Then: Server detects error immediately per fail-fast principle
    expect(errorHandlingResult.errorDetectedImmediately).toBeTruthy();
    if (errorHandlingResult.errorResponseTime) {
      expect(errorHandlingResult.errorResponseTime).toBeLessThan(1000);
    }

    // Then: Server returns MCP-formatted error response
    expect(errorHandlingResult.errorResponseInMCPFormat).toBeTruthy();

    // Then: Server keeps connection open for subsequent requests
    expect(errorHandlingResult.connectionRemainsOpen).toBeTruthy();
    expect(errorHandlingResult.subsequentRequestSucceeds).toBeTruthy();

    // Overall error handling success
    expect(errorHandlingResult.success).toBeTruthy();
    if (!errorHandlingResult.success) {
      throw new Error(errorHandlingResult.error || 'Unknown error handling failure');
    }
  });

  test('E2E-014: Server isolates concurrent client sessions without cross-contamination', async ({ page }) => {
    // Given: Multiple clients maintain active connections
    // Convert HTTP URL to WebSocket URL
    const wsUrl = mcpServiceUrl.replace('http://', 'ws://').replace('https://', 'wss://');

    // When: Clients execute concurrent operations
    const sessionIsolationResult = await page.evaluate(async (wsEndpoint) => {
      return new Promise<{
        success: boolean;
        allClientsConnected: boolean;
        allClientsInitialized: boolean;
        responsesIsolated: boolean;
        noResponseCrossContamination: boolean;
        connectionStateConsistent: boolean;
        error?: string;
      }>((resolve) => {
        try {
          const NUM_CLIENTS = 3;
          const clients: {
            ws: WebSocket;
            id: number;
            connected: boolean;
            initialized: boolean;
            receivedOwnResponse: boolean;
            receivedOthersResponse: boolean;
          }[] = [];

          let allClientsConnected = false;
          let allClientsInitialized = false;
          let responsesIsolated = false;
          let noResponseCrossContamination = false;
          let connectionStateConsistent = false;
          let testTimeout: NodeJS.Timeout;

          // Create multiple clients
          for (let i = 0; i < NUM_CLIENTS; i++) {
            const ws = new WebSocket(`${wsEndpoint}/mcp`);
            clients.push({
              ws,
              id: i,
              connected: false,
              initialized: false,
              receivedOwnResponse: false,
              receivedOthersResponse: false
            });

            const clientIndex = i;
            const client = clients[clientIndex];

            ws.onopen = () => {
              // Then: Server isolates each client session correctly
              client.connected = true;

              // Check if all clients are connected
              if (clients.every(c => c.connected)) {
                allClientsConnected = true;

                // Send initialize requests with unique client identifiers
                clients.forEach((c, idx) => {
                  const initRequest = {
                    jsonrpc: '2.0',
                    method: 'initialize',
                    id: idx + 1,
                    params: {
                      protocolVersion: '2024-11-05',
                      capabilities: {},
                      clientInfo: {
                        name: `e2e-test-client-${idx}`,
                        version: '1.0.0'
                      }
                    }
                  };
                  c.ws.send(JSON.stringify(initRequest));
                });
              }
            };

            ws.onmessage = (event) => {
              try {
                const message = JSON.parse(event.data as string);

                // Then: Responses reach correct clients without cross-contamination
                // Check if this is the initialize response for this client
                if (message.id === clientIndex + 1 && message.result) {
                  client.initialized = true;
                  client.receivedOwnResponse = true;

                  // Check if all clients are initialized
                  if (clients.every(c => c.initialized)) {
                    allClientsInitialized = true;

                    // Now send unique requests from each client
                    clients.forEach((c, idx) => {
                      const uniqueRequest = {
                        jsonrpc: '2.0',
                        method: 'tools/list',
                        id: 100 + idx, // Use unique IDs to track responses
                        params: {}
                      };
                      c.ws.send(JSON.stringify(uniqueRequest));
                    });
                  }
                }

                // Check if this is the tools/list response
                if (message.id >= 100 && message.id < 100 + NUM_CLIENTS) {
                  const expectedClientIndex = message.id - 100;

                  // This response should only be received by the corresponding client
                  if (expectedClientIndex === clientIndex) {
                    client.receivedOwnResponse = true;
                  } else {
                    // If a client receives another client's response, that's cross-contamination
                    client.receivedOthersResponse = true;
                  }

                  // Then: Server maintains consistent connection state per client
                  // Check if all clients received their own responses and none received others'
                  if (clients.every(c => c.receivedOwnResponse)) {
                    responsesIsolated = true;

                    // Verify no cross-contamination
                    if (clients.every(c => !c.receivedOthersResponse)) {
                      noResponseCrossContamination = true;
                    }

                    // Verify connection state consistency
                    if (clients.every(c => c.connected && c.initialized)) {
                      connectionStateConsistent = true;
                    }

                    // Close all connections and resolve
                    clearTimeout(testTimeout);
                    clients.forEach(c => c.ws.close());
                    resolve({
                      success: allClientsConnected && allClientsInitialized &&
                               responsesIsolated && noResponseCrossContamination &&
                               connectionStateConsistent,
                      allClientsConnected,
                      allClientsInitialized,
                      responsesIsolated,
                      noResponseCrossContamination,
                      connectionStateConsistent
                    });
                  }
                }
              } catch (err) {
                clearTimeout(testTimeout);
                clients.forEach(c => c.ws.close());
                resolve({
                  success: false,
                  allClientsConnected,
                  allClientsInitialized,
                  responsesIsolated,
                  noResponseCrossContamination,
                  connectionStateConsistent,
                  error: `Message parse error: ${err instanceof Error ? err.message : String(err)}`
                });
              }
            };

            ws.onerror = (error) => {
              clearTimeout(testTimeout);
              clients.forEach(c => c.ws.close());
              resolve({
                success: false,
                allClientsConnected,
                allClientsInitialized,
                responsesIsolated,
                noResponseCrossContamination,
                connectionStateConsistent,
                error: `WebSocket error on client ${clientIndex}: ${error}`
              });
            };
          }

          // Set overall timeout
          testTimeout = setTimeout(() => {
            clients.forEach(c => c.ws.close());
            resolve({
              success: false,
              allClientsConnected,
              allClientsInitialized,
              responsesIsolated,
              noResponseCrossContamination,
              connectionStateConsistent,
              error: 'Test timeout - not all operations completed'
            });
          }, 15000); // 15 second timeout

        } catch (err) {
          resolve({
            success: false,
            allClientsConnected: false,
            allClientsInitialized: false,
            responsesIsolated: false,
            noResponseCrossContamination: false,
            connectionStateConsistent: false,
            error: `Exception: ${err instanceof Error ? err.message : String(err)}`
          });
        }
      });
    }, wsUrl);

    // Then: Server isolates each client session correctly
    expect(sessionIsolationResult.allClientsConnected).toBeTruthy();
    expect(sessionIsolationResult.allClientsInitialized).toBeTruthy();

    // Then: Responses reach correct clients without cross-contamination
    expect(sessionIsolationResult.responsesIsolated).toBeTruthy();
    expect(sessionIsolationResult.noResponseCrossContamination).toBeTruthy();

    // Then: Server maintains consistent connection state per client
    expect(sessionIsolationResult.connectionStateConsistent).toBeTruthy();

    // Overall session isolation success
    expect(sessionIsolationResult.success).toBeTruthy();
    if (!sessionIsolationResult.success) {
      throw new Error(sessionIsolationResult.error || 'Unknown session isolation failure');
    }
  });

  test('E2E-015: Health check returns status with connection and dependency metrics', async ({ request }) => {
    // Given: Server accepts health check requests
    // Note: This is an E2E test but uses Request API because health check is HTTP endpoint, not WebSocket

    // When: Monitoring system queries health endpoint during active connections
    const response = await request.get(`${mcpServiceUrl}/health`);

    // Then: Server reports healthy status
    expect(response.status()).toBe(200);

    const healthData = await response.json();
    expect(healthData).toHaveProperty('status');
    expect(healthData.status).toBe('healthy');

    // Then: Server includes connection pool metrics
    expect(healthData).toHaveProperty('connections');
    expect(healthData.connections).toHaveProperty('active');
    expect(healthData.connections).toHaveProperty('total');
    expect(typeof healthData.connections.active).toBe('number');
    expect(typeof healthData.connections.total).toBe('number');

    // Then: Server includes dependency health for Redis and Google API
    expect(healthData).toHaveProperty('dependencies');
    expect(healthData.dependencies).toHaveProperty('redis');
    expect(healthData.dependencies).toHaveProperty('googleAPI');

    // Validate Redis dependency health structure
    expect(healthData.dependencies.redis).toHaveProperty('status');
    expect(['healthy', 'unhealthy', 'degraded']).toContain(healthData.dependencies.redis.status);

    // Validate Google API dependency health structure
    expect(healthData.dependencies.googleAPI).toHaveProperty('status');
    expect(['healthy', 'unhealthy', 'degraded']).toContain(healthData.dependencies.googleAPI.status);
  });

});
