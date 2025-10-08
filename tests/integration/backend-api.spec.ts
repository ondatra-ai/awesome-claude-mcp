import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('Backend API Integration Tests', () => {

  test('1.1-INT-001: server responds with correct status', async ({ request }) => {
    const response = await request.get(`${backendUrl}/version`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('version');
    expect(data.version).toBe('1.0.0');
  });

  test('ORPHAN: version endpoint rejects POST method', async ({ request }) => {
    const response = await request.post(`${backendUrl}/version`);
    expect(response.status()).toBe(405);
  });

  test('ORPHAN: health endpoint returns healthy status', async ({ request }) => {
    const response = await request.get(`${backendUrl}/health`);
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('status', 'healthy');
    expect(data).toHaveProperty('service', 'MCP Google Docs Editor - Backend');
    expect(data).toHaveProperty('timestamp');
  });

  test('ORPHAN: health endpoint rejects DELETE method', async ({ request }) => {
    const response = await request.delete(`${backendUrl}/health`);
    expect(response.status()).toBe(405);
  });

  test('ORPHAN: non-existent endpoint returns 404', async ({ request }) => {
    const response = await request.get(`${backendUrl}/nonexistent`);
    expect(response.status()).toBe(404);
  });
});
