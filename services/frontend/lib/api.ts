import type { IHealthResponse, IVersionResponse } from '../interfaces/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  async getVersion(): Promise<IVersionResponse> {
    const response = await fetch(`${this.baseURL}/version`, {
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
    const response = await fetch(`${this.baseURL}/health`, {
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
