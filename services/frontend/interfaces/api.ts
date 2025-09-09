export interface IVersionResponse {
  version: string;
}

export interface IHealthResponse {
  status: string;
  service: string;
  timestamp: string;
}

export interface ILoadingState {
  isLoading: boolean;
  error: string | null;
}
