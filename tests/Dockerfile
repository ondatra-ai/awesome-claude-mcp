FROM mcr.microsoft.com/playwright:v1.55.0-noble

# Set working directory

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./
COPY playwright.config.test.ts ./playwright.config.ts

# Install dependencies
RUN npm ci

# Copy test files
COPY e2e/ ./e2e/

# Browsers are already installed in the base image

# Create directories for test output
RUN mkdir -p test-results playwright-report

# Default command runs all tests in headless mode
CMD ["npx", "playwright", "test", "--reporter=html"]
