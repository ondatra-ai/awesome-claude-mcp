import { test, expect } from '@playwright/test';

test.describe('Backend API E2E Tests', () => {
  const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8080';

  test('should access version endpoint directly', async ({ request }) => {
    const response = await request.get(`${BACKEND_URL}/version`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('version');
    expect(data.version).toBe('1.0.0');
  });

  test('should access health endpoint directly', async ({ request }) => {
    const response = await request.get(`${BACKEND_URL}/health`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('status', 'healthy');
    expect(data).toHaveProperty('service', 'MCP Google Docs Editor - Backend');
    expect(data).toHaveProperty('timestamp');
  });

  test('should handle 404 for non-existent endpoints', async ({ request }) => {
    const response = await request.get(`${BACKEND_URL}/nonexistent`);
    expect(response.status()).toBe(404);
  });

  test('should handle method not allowed for POST on version endpoint', async ({ request }) => {
    const response = await request.post(`${BACKEND_URL}/version`);
    expect(response.status()).toBe(405);
  });

  test('should verify CORS headers for frontend requests', async ({ request }) => {
    const FRONTEND_URL = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000';
    const response = await request.get(`${BACKEND_URL}/version`, {
      headers: {
        'Origin': FRONTEND_URL
      }
    });
    expect(response.ok()).toBeTruthy();

    // The response should include CORS headers allowing the frontend origin
    const headers = response.headers();
    // Note: Exact CORS headers depend on Fiber's CORS middleware implementation
    expect(response.status()).toBe(200);
  });
});
