import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    trace: 'on-first-retry',
  },

  projects: [
    {
      name: 'frontend',
      testDir: './frontend',
      use: { 
        ...devices['Desktop Chrome'],
        baseURL: 'http://localhost:3000',
      },
    },
    {
      name: 'backend-api',
      testDir: './backend',
      use: { 
        ...devices['Desktop Chrome'],
        baseURL: 'http://localhost:8080',
      },
    },
  ],

  webServer: [
    {
      command: 'cd ../../services/backend && go run ./cmd/main.go',
      port: 8080,
      reuseExistingServer: !process.env.CI,
    },
    {
      command: 'cd ../../services/frontend && npm run dev',
      port: 3000,
      reuseExistingServer: !process.env.CI,
    },
  ],
});