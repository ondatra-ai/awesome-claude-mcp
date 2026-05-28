import type { ReactNode } from 'react';

export const metadata = {
  title: 'Frontend',
  description: 'BDD fixture Next.js home page.',
};

export default function RootLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
