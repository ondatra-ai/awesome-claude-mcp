/**
 * Claude.ai Authentication Interfaces
 */

export interface IClaudeAuthConfig {
  mailosaurApiKey: string;
  mailosaurServerId: string;
  claudeEmail: string;
  authState?: string;
  envFilePath?: string; // Path to .env.test for persisting auth state
}

export interface IAuthResult {
  success: boolean;
  isNewLogin: boolean;
  authState?: string;
  error?: string;
}
