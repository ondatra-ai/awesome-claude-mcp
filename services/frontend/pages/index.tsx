import type { GetStaticProps, InferGetStaticPropsType } from 'next';

import { HomePageClient } from '@/components/home-page-client';

export const getStaticProps: GetStaticProps<{
  backendBaseUrl: string;
}> = () => {
  const value = process.env.NEXT_PUBLIC_API_URL?.trim();

  if (!value) {
    throw new Error(
      'NEXT_PUBLIC_API_URL environment variable is required for rendering the homepage.'
    );
  }

  return {
    props: {
      backendBaseUrl: value.replace(/\/$/, ''),
    },
  };
};

export default function HomePage({
  backendBaseUrl,
}: InferGetStaticPropsType<typeof getStaticProps>): JSX.Element {
  return <HomePageClient backendBaseUrl={backendBaseUrl} />;
}
