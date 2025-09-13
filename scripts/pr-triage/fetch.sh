#!/usr/bin/env bash
set -euo pipefail

TMP_DIR="./tmp"
mkdir -p "$TMP_DIR"

PR_NUMBER=$(gh pr view --json number -q .number)

go run ./scripts/list-pr-conversations/main.go "$PR_NUMBER" > "$TMP_DIR/PR_CONV.json"

TOTAL=$(jq 'length' "$TMP_DIR/PR_CONV.json")
UNRES=$(jq '[.[] | select(.isResolved==false)] | length' "$TMP_DIR/PR_CONV.json")

echo "Fetched $TOTAL threads; $UNRES unresolved."
