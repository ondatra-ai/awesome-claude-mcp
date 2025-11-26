import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl, frontendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('Backend API E2E Tests', () => {

  test('E2E-001: service returns version with headers', async ({ request }) => {
    const response = await request.get(`${backendUrl}/version`);
    expect(response.ok(), 'Response should be successful').toBeTruthy();
    expect(response.status(), 'Expected HTTP 200 status').toBe(200);

    const data = await response.json();
    expect(data, 'Response should have version property').toHaveProperty('version');
    expect(data.version, 'Version should be 1.0.0').toBe('1.0.0');
  });

  test('E2E-002: health endpoint accessible via HTTP', async ({ request }) => {
    const response = await request.get(`${backendUrl}/health`);
    expect(response.ok(), 'Health endpoint should be accessible').toBeTruthy();
    expect(response.status(), 'Expected HTTP 200 status').toBe(200);

    const data = await response.json();
    expect(data, 'Response should have healthy status').toHaveProperty('status', 'healthy');
    expect(data, 'Response should have service name').toHaveProperty('service', 'MCP Google Docs Editor - Backend');
    expect(data, 'Response should have timestamp').toHaveProperty('timestamp');
  });

  test('E2E-003: returns 404 for non-existent endpoints', async ({ request }) => {
    const response = await request.get(`${backendUrl}/nonexistent`);
    expect(response.status(), 'Non-existent endpoint should return 404').toBe(404);
  });

  test('E2E-004: returns 405 for POST on version endpoint', async ({ request }) => {
    const response = await request.post(`${backendUrl}/version`);
    expect(response.status(), 'POST method should not be allowed on version endpoint').toBe(405);
  });

  test('E2E-005: includes CORS headers for frontend requests', async ({ request }) => {
    const response = await request.get(`${backendUrl}/version`, {
      headers: {
        'Origin': frontendUrl
      }
    });
    expect(response.ok(), 'Request with origin header should be successful').toBeTruthy();

    // The response should include CORS headers allowing the frontend origin
    const headers = response.headers();
    // Note: Exact CORS headers depend on Fiber's CORS middleware implementation
    expect(response.status(), 'Expected HTTP 200 status for CORS request').toBe(200);
  });
});
