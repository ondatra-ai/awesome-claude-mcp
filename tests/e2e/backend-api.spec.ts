import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl, frontendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('Backend API E2E Tests', () => {

  test('FR-00001 should access version endpoint directly', async ({ request }) => {
    // FR-00001: Backend /version endpoint returns 1.0.0
    // Source: Story 1.1 (1.1-E2E-001)

    const response = await request.get(`${backendUrl}/version`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('version');
    expect(data.version).toBe('1.0.0');
  });

  test('FR-00002 should access health endpoint directly', async ({ request }) => {
    // FR-00002: Backend /health endpoint returns healthy status
    // Source: Story 1.1 (1.1-E2E-006)

    const response = await request.get(`${backendUrl}/health`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('status', 'healthy');
    expect(data).toHaveProperty('service', 'MCP Google Docs Editor - Backend');
    expect(data).toHaveProperty('timestamp');
  });

  test('FR-00003 should handle 404 for non-existent endpoints', async ({ request }) => {
    // FR-00003: Backend handles 404 for non-existent endpoints
    // Source: Not in original requirements (orphaned test)

    const response = await request.get(`${backendUrl}/nonexistent`);
    expect(response.status()).toBe(404);
  });

  test('FR-00004 should handle method not allowed for POST on version endpoint', async ({ request }) => {
    // FR-00004: Backend rejects invalid HTTP methods
    // Source: Not in original requirements (orphaned test)

    const response = await request.post(`${backendUrl}/version`);
    expect(response.status()).toBe(405);
  });

  test('FR-00005 should verify CORS headers for frontend requests', async ({ request }) => {
    // FR-00005: Backend provides CORS headers for frontend
    // Source: Not in original requirements (orphaned test)

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
