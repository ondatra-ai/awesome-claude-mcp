import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('Backend API Integration Tests', () => {

  test('INT-001: server responds with correct status', async ({ request }) => {
    const response = await request.get(`${backendUrl}/version`);
    expect(response.ok(), 'Response should be successful').toBeTruthy();
    expect(response.status(), 'Expected HTTP 200 status').toBe(200);

    const data = await response.json();
    expect(data, 'Response should have version property').toHaveProperty('version');
    expect(data.version, 'Version should be 1.0.0').toBe('1.0.0');
  });

  test('INT-002: version endpoint rejects invalid methods', async ({ request }) => {
    const response = await request.post(`${backendUrl}/version`);
    expect(response.status(), 'POST method should return 405 Method Not Allowed').toBe(405);
  });

  test('INT-003: health endpoint returns healthy status', async ({ request }) => {
    const response = await request.get(`${backendUrl}/health`);
    expect(response.ok(), 'Health endpoint should be accessible').toBeTruthy();
    expect(response.status(), 'Expected HTTP 200 status').toBe(200);

    const data = await response.json();
    expect(data, 'Response should have healthy status').toHaveProperty('status', 'healthy');
    expect(data, 'Response should have service name').toHaveProperty('service', 'MCP Google Docs Editor - Backend');
    expect(data, 'Response should have timestamp').toHaveProperty('timestamp');
  });

  test('INT-004: health endpoint rejects invalid methods', async ({ request }) => {
    const response = await request.delete(`${backendUrl}/health`);
    expect(response.status(), 'DELETE method should return 405 Method Not Allowed').toBe(405);
  });

  test('INT-005: non-existent endpoint returns 404', async ({ request }) => {
    const response = await request.get(`${backendUrl}/nonexistent`);
    expect(response.status(), 'Non-existent endpoint should return 404').toBe(404);
  });
});
