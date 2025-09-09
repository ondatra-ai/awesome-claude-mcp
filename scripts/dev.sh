#!/bin/bash

# Development Environment Starter
# This script starts all services for local development with Docker Compose

set -e

echo "🚀 Starting MCP Google Docs Editor Development Environment..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose down --remove-orphans || true

# Start services with build
echo "🔧 Building and starting services..."
docker-compose up --build

echo "✅ Development environment started!"
echo "📊 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "💡 Use Ctrl+C to stop all services"
