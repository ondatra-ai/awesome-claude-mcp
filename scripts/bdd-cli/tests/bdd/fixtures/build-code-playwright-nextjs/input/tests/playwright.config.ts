import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './integration',
  timeout: 30_000,
  fullyParallel: false,
  workers: 1,
  reporter: 'json',
  use: {
    baseURL: 'http://127.0.0.1:3000',
  },
  webServer: {
    command: 'npm run dev -- --hostname 127.0.0.1 --port 3000',
    cwd: '../services/frontend',
    url: 'http://127.0.0.1:3000',
    timeout: 120_000,
    reuseExistingServer: true,
  },
  projects: [
    {
      name: 'chromium',
      use: { browserName: 'chromium' },
    },
  ],
});
