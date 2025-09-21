import { createApiClient } from '@/lib/api';

// Mock fetch globally
const mockFetch = jest.fn();
global.fetch = mockFetch;

const apiClient = createApiClient('http://localhost:8080');

describe('ApiClient', () => {
  beforeEach(() => {
    mockFetch.mockReset();
  });

  describe('getVersion', () => {
    it('returns version data on successful response', async () => {
      const mockVersionResponse = { version: '1.0.0' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockVersionResponse,
      });

      const result = await apiClient.getVersion();

      expect(result).toEqual(mockVersionResponse);
      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/version', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
    });

    it('throws error on failed response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
      });

      await expect(apiClient.getVersion()).rejects.toThrow(
        'Failed to fetch version: 500 Internal Server Error'
      );
    });

    it('handles network errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(apiClient.getVersion()).rejects.toThrow('Network error');
    });
  });

  describe('getHealth', () => {
    it('returns health data on successful response', async () => {
      const mockHealthResponse = {
        status: 'healthy',
        timestamp: '2023-01-01T00:00:00Z',
      };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockHealthResponse,
      });

      const result = await apiClient.getHealth();

      expect(result).toEqual(mockHealthResponse);
      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/health', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
    });

    it('throws error on failed response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 503,
        statusText: 'Service Unavailable',
      });

      await expect(apiClient.getHealth()).rejects.toThrow(
        'Failed to fetch health: 503 Service Unavailable'
      );
    });

    it('handles network errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(apiClient.getHealth()).rejects.toThrow('Network error');
    });
  });

  describe('factory', () => {
    it('throws when base URL is missing', () => {
      expect(() => createApiClient('')).toThrow('Backend base URL is required to create ApiClient');
    });
  });
});
