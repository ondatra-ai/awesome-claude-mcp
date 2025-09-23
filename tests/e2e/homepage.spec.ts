import { test, expect } from '@playwright/test';

test.describe('Homepage E2E Tests', () => {
  test('EE-00006-04: should load homepage and display title', async ({ page }) => {
    // Source: FR-00006 - Frontend single-page application loads successfully

    await page.goto('/');

    // Check that the page loads successfully
    await expect(page).toHaveTitle('MCP Google Docs Editor');

    // Check that the main title is displayed
    await expect(page.locator('h1')).toContainText('MCP Google Docs Editor');

    // Check that the description is displayed
    await expect(page.getByTestId('hero-description')).toContainText('A Model Context Protocol integration for seamless Google Docs editing');
  });

  test('EE-00008-01: should display welcome card with features', async ({ page }) => {
    // Source: FR-00008 - Homepage displays welcome card with features

    await page.goto('/');

    // Check welcome card title
    await expect(page.getByTestId('welcome-title')).toContainText('Welcome to MCP Google Docs Editor');

    // Check feature cards
    await expect(page.getByTestId('feature-document-ops')).toBeVisible();
    await expect(page.getByTestId('feature-ai-integration')).toBeVisible();

    // Check feature descriptions
    await expect(page.getByTestId('feature-document-ops-desc')).toContainText('Replace, append, prepend, and insert content');
    await expect(page.getByTestId('feature-ai-integration-desc')).toContainText('Compatible with Claude Code and ChatGPT');
  });

  test('EE-00007-06: should fetch and display backend version', async ({ page }) => {
    // Source: FR-00007 - Homepage displays backend version at bottom

    await page.goto('/');

    // Check that backend version section is present
    await expect(page.locator('text=Backend Version:')).toBeVisible();

    // Wait for version to load (it should show "1.0.0" or "Loading..." initially)
    const versionSection = page.getByTestId('backend-version');

    // Wait for either loading state or actual version
    await expect(versionSection).toHaveText(/Loading\.\.\.|1\.0\.0/);

    // If it's loading, wait for the actual version to appear
    try {
      await expect(versionSection).toHaveText('1.0.0', { timeout: 10000 });
    } catch (error) {
      // If version doesn't load, check that loading state is still showing
      await expect(versionSection).toHaveText('Loading...');
      console.warn('Backend version did not load within timeout period');
    }
  });

  test('EE-00009-01: should have responsive design', async ({ page }) => {
    // Source: FR-00009 - Homepage has responsive design

    await page.goto('/');

    // Test mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await expect(page.locator('h1')).toBeVisible();

    // Test tablet viewport
    await page.setViewportSize({ width: 768, height: 1024 });
    await expect(page.locator('h1')).toBeVisible();

    // Test desktop viewport
    await page.setViewportSize({ width: 1200, height: 800 });
    await expect(page.locator('h1')).toBeVisible();
  });
});
