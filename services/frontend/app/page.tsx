import { HomePageClient } from '@/components/home-page-client';

function resolveBackendBaseUrl(): string {
  const value = process.env.NEXT_PUBLIC_API_URL?.trim();

  if (!value) {
    throw new Error(
      'NEXT_PUBLIC_API_URL environment variable is required for rendering the homepage.'
    );
  }

  return value.replace(/\/$/, '');
}

export default function HomePage(): JSX.Element {
  const backendBaseUrl = resolveBackendBaseUrl();

  return <HomePageClient backendBaseUrl={backendBaseUrl} />;
}
