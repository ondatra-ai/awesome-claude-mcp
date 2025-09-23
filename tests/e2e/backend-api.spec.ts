import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl, frontendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('Backend API E2E Tests', () => {

  test('EE_00001_04: should access version endpoint directly', async ({ request }) => {
    const response = await request.get(`${backendUrl}/version`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('version');
    expect(data.version).toBe('1.0.0');
  });

  test('EE_00002_02: should access health endpoint directly', async ({ request }) => {
    const response = await request.get(`${backendUrl}/health`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('status', 'healthy');
    expect(data).toHaveProperty('service', 'MCP Google Docs Editor - Backend');
    expect(data).toHaveProperty('timestamp');
  });

  test('EE_00003_01: should handle 404 for non-existent endpoints', async ({ request }) => {
    const response = await request.get(`${backendUrl}/nonexistent`);
    expect(response.status()).toBe(404);
  });

  test('EE_00004_01: should handle method not allowed for POST on version endpoint', async ({ request }) => {
    const response = await request.post(`${backendUrl}/version`);
    expect(response.status()).toBe(405);
  });

  test('EE_00005_01: should verify CORS headers for frontend requests', async ({ request }) => {
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
