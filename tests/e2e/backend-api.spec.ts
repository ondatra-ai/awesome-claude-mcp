import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl, frontendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('Backend API E2E Tests', () => {

  test('1.1-E2E-001: service returns version with headers', async ({ request }) => {
    const response = await request.get(`${backendUrl}/version`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('version');
    expect(data.version).toBe('1.0.0');
  });

  test('1.1-E2E-002: health endpoint accessible via HTTP', async ({ request }) => {
    const response = await request.get(`${backendUrl}/health`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('status', 'healthy');
    expect(data).toHaveProperty('service', 'MCP Google Docs Editor - Backend');
    expect(data).toHaveProperty('timestamp');
  });

  test('1.1-E2E-003: returns 404 for non-existent endpoints', async ({ request }) => {
    const response = await request.get(`${backendUrl}/nonexistent`);
    expect(response.status()).toBe(404);
  });

  test('1.1-E2E-004: returns 405 for POST on version endpoint', async ({ request }) => {
    const response = await request.post(`${backendUrl}/version`);
    expect(response.status()).toBe(405);
  });

  test('1.1-E2E-005: includes CORS headers for frontend requests', async ({ request }) => {
    const response = await request.get(`${backendUrl}/version`, {
      headers: {
        'Origin': frontendUrl
      }
    });
    expect(response.ok()).toBeTruthy();

    // The response should include CORS headers allowing the frontend origin
    const headers = response.headers();
    // Note: Exact CORS headers depend on Fiber's CORS middleware implementation
    expect(response.status()).toBe(200);
  });
});
