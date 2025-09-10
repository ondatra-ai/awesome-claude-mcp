import type { ILoadingState, IVersionResponse } from './api';

export interface IVersionDisplayProps {
  version: IVersionResponse | null;
  loadingState: ILoadingState;
}
