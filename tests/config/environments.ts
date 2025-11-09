export type EnvironmentName = 'local' | 'dev';

export interface EnvironmentConfig {
  name: EnvironmentName;
  frontendUrl: string;
  backendUrl: string;
  mcpServiceUrl: string;
  /** Optional human friendly description for logging. */
  description?: string;
}

const DEFAULT_ENVIRONMENT: EnvironmentName = 'local';

const environments: Record<EnvironmentName, EnvironmentConfig> = {
  local: {
    name: 'local',
    frontendUrl: 'http://localhost:3000',
    backendUrl: 'http://localhost:8080',
    mcpServiceUrl: 'http://localhost:8081',
    description: 'Local docker-compose stack',
  },
  dev: {
    name: 'dev',
    frontendUrl: 'https://dev.ondatra-ai.xyz',
    backendUrl: 'https://api.dev.ondatra-ai.xyz',
    mcpServiceUrl: 'https://mcp.dev.ondatra-ai.xyz',
    description: 'Remote dev environment (frontend/backend hosted on Ondatra)',
  },
};

export function getEnvironmentConfig(rawName?: string): EnvironmentConfig {
  const normalizedName = (rawName ?? DEFAULT_ENVIRONMENT).toLowerCase() as EnvironmentName;

  if (normalizedName in environments) {
    return environments[normalizedName as EnvironmentName];
  }

  console.warn(
    `Unknown E2E_ENV "${rawName}" supplied. Falling back to "${DEFAULT_ENVIRONMENT}". ` +
    'Update tests/config/environments.ts if you need to add a new target.'
  );

  return environments[DEFAULT_ENVIRONMENT];
}

export const environmentRegistry = environments;
