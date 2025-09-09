import { test, expect } from '@playwright/test';

test.describe('Backend API Integration Tests', () => {
  const baseURL = 'http://localhost:8080';
  
  test('GET /version should return correct version', async ({ request }) => {
    const response = await request.get(`${baseURL}/version`);
    
    expect(response.status()).toBe(200);
    expect(response.headers()['content-type']).toContain('application/json');
    
    const data = await response.json();
    expect(data).toEqual({ version: '1.0.0' });
  });
  
  test('GET /health should return health status', async ({ request }) => {
    const response = await request.get(`${baseURL}/health`);
    
    expect(response.status()).toBe(200);
    expect(response.headers()['content-type']).toContain('application/json');
    
    const data = await response.json();
    expect(data).toHaveProperty('status');
    expect(data.status).toBe('healthy');
  });
  
  test('API should handle CORS correctly', async ({ request }) => {
    const response = await request.get(`${baseURL}/version`, {
      headers: {
        'Origin': 'http://localhost:3000'
      }
    });
    
    expect(response.status()).toBe(200);
    expect(response.headers()['access-control-allow-origin']).toBe('*');
  });
  
  test('API should return 404 for non-existent endpoint', async ({ request }) => {
    const response = await request.get(`${baseURL}/nonexistent`);
    
    expect(response.status()).toBe(404);
  });
  
  test('API should handle malformed requests gracefully', async ({ request }) => {
    const response = await request.post(`${baseURL}/version`, {
      data: { invalid: 'data' }
    });
    
    // Should return method not allowed or similar
    expect([405, 404]).toContain(response.status());
  });
});