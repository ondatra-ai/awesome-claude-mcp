#!/usr/bin/env npx ts-node
/**
 * Manual Login Script
 *
 * Opens a headed browser for manual Claude.ai authentication.
 * Use this when Cloudflare blocks automated login attempts.
 *
 * Usage: npm run auth:login
 */
import * as path from 'path';

import { manualLogin } from '../e2e/helpers/claude-auth';

const envFilePath = path.join(process.cwd(), '.env.test');

manualLogin(envFilePath)
  .then(() => {
    console.log('\n🎉 Manual login complete! You can now run the tests.');
    process.exit(0);
  })
  .catch((error: Error) => {
    console.error('\n❌ Login failed:', error.message);
    process.exit(1);
  });
