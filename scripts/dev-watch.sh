#!/bin/bash

# Development Environment with File Watching
# This script starts services with file watching for auto-reload on changes

set -e

echo "🚀 Starting MCP Google Docs Editor Development Environment with File Watching..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose down --remove-orphans || true

# Start services with build and watch for changes
echo "🔧 Building and starting services with file watching..."
echo "📁 File changes will automatically reload the services"
docker-compose up --build --watch

echo "✅ Development environment with file watching started!"
echo "📊 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "👀 Watching for file changes..."
echo "💡 Use Ctrl+C to stop all services"
