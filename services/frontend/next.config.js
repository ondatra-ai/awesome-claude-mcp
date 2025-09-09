/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  env: {
    NEXT_PUBLIC_API_URL:
      process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  },
  images: {
    domains: ['docs.google.com'],
  },
};

module.exports = nextConfig;
