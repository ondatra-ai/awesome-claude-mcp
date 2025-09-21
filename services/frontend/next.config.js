/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  images: {
    domains: ['docs.google.com'],
  },
};

module.exports = nextConfig;
