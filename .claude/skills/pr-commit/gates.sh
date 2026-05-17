#!/usr/bin/env bash
# Usage: gates.sh
# Run the awesome-claude-mcp quality pipeline before the commit step.
# Aborts on any failure. Mirrors the make targets the project uses to
# gate production-ready code.
set -euo pipefail

make lint-backend
make lint-frontend
make lint-scripts
make test-unit
make test-e2e
pre-commit run --all-files
