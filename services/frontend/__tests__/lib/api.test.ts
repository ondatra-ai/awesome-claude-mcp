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
    it('UT-00007-01: should construct correct request for version', async () => {
      // UT-00007-01: API client constructs correct request
      // Source: FR-00007 - Homepage displays backend version at bottom
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

    it('ORPHAN: should handle API errors gracefully', async () => {
      // ORPHAN: Testing error handling not specified in requirements
      // Note: This could map to UT-00007-03 if error handling is considered part of version display

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
      });

      await expect(apiClient.getVersion()).rejects.toThrow(
        'Failed to fetch version: 500 Internal Server Error'
      );
    });

    it('ORPHAN: should handle network errors', async () => {
      // ORPHAN: Testing network error handling not specified in requirements
      // Note: This could map to UT-00007-03 if error handling is considered part of version display

      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(apiClient.getVersion()).rejects.toThrow('Network error');
    });
  });

  describe('getHealth', () => {
    it('ORPHAN: should return health data on successful response', async () => {
      // ORPHAN: Testing health endpoint functionality not specified in requirements
      // Reason: Health endpoint is tested at E2E level, this is additional unit coverage
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

    it('ORPHAN: should handle health API errors', async () => {
      // ORPHAN: Testing health endpoint error handling not specified in requirements
      // Reason: Health endpoint error handling not part of functional requirements

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 503,
        statusText: 'Service Unavailable',
      });

      await expect(apiClient.getHealth()).rejects.toThrow(
        'Failed to fetch health: 503 Service Unavailable'
      );
    });

    it('ORPHAN: should handle health network errors', async () => {
      // ORPHAN: Testing health endpoint network error handling not specified in requirements
      // Reason: Health endpoint error handling not part of functional requirements

      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(apiClient.getHealth()).rejects.toThrow('Network error');
    });
  });

  describe('factory', () => {
    it('ORPHAN: should validate base URL is required', () => {
      // ORPHAN: Testing API client factory validation not specified in requirements
      // Reason: Input validation for API client factory not part of functional requirements

      expect(() => createApiClient('')).toThrow('Backend base URL is required to create ApiClient');
    });
  });
});
