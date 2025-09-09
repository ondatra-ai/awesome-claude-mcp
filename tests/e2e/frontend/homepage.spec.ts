import { test, expect } from '@playwright/test';

test.describe('Homepage', () => {
  test('should load successfully', async ({ page }) => {
    await page.goto('http://localhost:3000');
    
    // Check title
    await expect(page).toHaveTitle(/MCP Google Docs Editor/);
    
    // Check main heading
    await expect(page.locator('h1')).toContainText('MCP Google Docs Editor');
    
    // Check version display
    await expect(page.locator('text=Backend version: 1.0.0')).toBeVisible();
  });
  
  test('should display backend version at bottom', async ({ page }) => {
    await page.goto('http://localhost:3000');
    
    // Wait for version to load
    await expect(page.locator('text=Backend version:')).toBeVisible();
    await expect(page.locator('text=1.0.0')).toBeVisible();
  });
});