import type { IHealthResponse, IVersionResponse } from '../interfaces/api';

const resolveBaseUrl = (): string => {
  const raw = process.env.NEXT_PUBLIC_API_URL?.trim();

  if (!raw) {
    throw new Error(
      'NEXT_PUBLIC_API_URL environment variable is required for API client'
    );
  }

  return raw.replace(/\/$/, '');
};

const buildUrl = (path: string): string => {
  const baseUrl = resolveBaseUrl();
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${baseUrl}${normalizedPath}`;
};

const handleError = (response: Response, context: string): never => {
  throw new Error(
    `Failed to fetch ${context}: ${response.status} ${response.statusText}`
  );
};

class ApiClient {
  async getVersion(): Promise<IVersionResponse> {
    const response = await fetch(buildUrl('/version'), {
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
    const response = await fetch(buildUrl('/health'), {
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

export const apiClient = new ApiClient();
