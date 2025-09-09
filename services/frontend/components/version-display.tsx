import type { IVersionDisplayProps } from '@/interfaces/version-display';

export function VersionDisplay({
  version,
  loadingState,
}: IVersionDisplayProps): JSX.Element {
  return (
    <div className="text-center">
      <div className="inline-flex items-center rounded-full border px-3 py-1 text-sm">
        <span className="text-muted-foreground">Backend Version:</span>
        <span data-testid="backend-version" className="ml-2">
          {loadingState.isLoading ? (
            <span className="text-muted-foreground">Loading...</span>
          ) : loadingState.error ? (
            <span className="text-destructive">
              Error: {loadingState.error}
            </span>
          ) : (
            <span className="font-mono font-semibold">{version?.version}</span>
          )}
        </span>
      </div>
    </div>
  );
}
