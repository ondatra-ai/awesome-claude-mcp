import type { IHealthResponse, IVersionResponse } from '../interfaces/api';

const normalizeBaseUrl = (rawBaseUrl: string): string =>
  rawBaseUrl.replace(/\/$/, '');

const buildUrl = (baseUrl: string, path: string): string => {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${baseUrl}${normalizedPath}`;
};

const handleError = (response: Response, context: string): never => {
  throw new Error(
    `Failed to fetch ${context}: ${response.status} ${response.statusText}`
  );
};

class ApiClient {
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async getVersion(): Promise<IVersionResponse> {
    const response = await fetch(buildUrl(this.baseUrl, '/version'), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      handleError(response, 'version');
    }

    return response.json() as Promise<IVersionResponse>;
  }

  async getHealth(): Promise<IHealthResponse> {
    const response = await fetch(buildUrl(this.baseUrl, '/health'), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      handleError(response, 'health');
    }

    return response.json() as Promise<IHealthResponse>;
  }
}

export const createApiClient = (rawBaseUrl: string): ApiClient => {
  if (!rawBaseUrl || rawBaseUrl.trim() === '') {
    throw new Error('Backend base URL is required to create ApiClient');
  }

  return new ApiClient(normalizeBaseUrl(rawBaseUrl.trim()));
};
