import type { IHealthResponse, IVersionResponse } from '../interfaces/api';

const sanitizeUrl = (url: string): string => url.replace(/\/$/, '');

const resolveDefaultRemoteUrl = (hostname: string): string | null => {
  if (hostname === 'dev.ondatra-ai.xyz') {
    return 'https://api.dev.ondatra-ai.xyz';
  }

  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    return 'http://localhost:8080';
  }

  return null;
};

const resolveBaseUrl = (): string => {
  const envValue = process.env.NEXT_PUBLIC_API_URL;
  if (envValue && envValue.trim() !== '') {
    return sanitizeUrl(envValue.trim());
  }

  if (typeof window !== 'undefined') {
    const inferred = resolveDefaultRemoteUrl(window.location.hostname);
    if (inferred) {
      return inferred;
    }
  }

  return 'http://localhost:8080';
};

const withBaseUrl = (path: string): string => {
  const baseUrl = resolveBaseUrl();
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;

  if (/^https?:\/\//i.test(baseUrl)) {
    return `${baseUrl}${normalizedPath}`;
  }

  return `${baseUrl}${normalizedPath}`;
};

class ApiClient {
  async getVersion(): Promise<IVersionResponse> {
    const response = await fetch(withBaseUrl('/version'), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch version: ${response.statusText}`);
    }

    return response.json() as Promise<IVersionResponse>;
  }

  async getHealth(): Promise<IHealthResponse> {
    const response = await fetch(withBaseUrl('/health'), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch health: ${response.statusText}`);
    }

    return response.json() as Promise<IHealthResponse>;
  }
}

export const apiClient = new ApiClient();
