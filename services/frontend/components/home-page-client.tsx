'use client';

import { useEffect, useMemo, useState } from 'react';

import { VersionDisplay } from '@/components/version-display';
import { WelcomeCard } from '@/components/welcome-card';
import type { ILoadingState, IVersionResponse } from '@/interfaces/api';
import { createApiClient } from '@/lib/api';

export function HomePageClient({
  backendBaseUrl,
}: {
  backendBaseUrl: string;
}): JSX.Element {
  const apiClient = useMemo(
    () => createApiClient(backendBaseUrl),
    [backendBaseUrl]
  );

  const [version, setVersion] = useState<IVersionResponse | null>(null);
  const [loadingState, setLoadingState] = useState<ILoadingState>({
    isLoading: true,
    error: null,
  });

  useEffect(() => {
    const fetchVersion = async (): Promise<void> => {
      try {
        setLoadingState({ isLoading: true, error: null });
        const versionData = await apiClient.getVersion();
        setVersion(versionData);
        setLoadingState({ isLoading: false, error: null });
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : 'Failed to fetch version';
        setLoadingState({ isLoading: false, error: errorMessage });
      }
    };

    void fetchVersion();
  }, [apiClient]);

  return (
    <div className="container mx-auto max-w-4xl px-4 py-12">
      <div className="space-y-8">
        <div className="text-center space-y-4">
          <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl lg:text-7xl">
            MCP Google Docs Editor
          </h1>
          <p
            data-testid="hero-description"
            className="mx-auto max-w-[700px] text-lg text-muted-foreground sm:text-xl"
          >
            A Model Context Protocol integration for seamless Google Docs
            editing with Claude Code and ChatGPT
          </p>
        </div>

        <WelcomeCard />

        <VersionDisplay version={version} loadingState={loadingState} />
      </div>
    </div>
  );
}
