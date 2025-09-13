#!/usr/bin/env bash
set -euo pipefail

TMP_DIR="./tmp"
JSON="$TMP_DIR/PR_CONV.json"

if [[ ! -f "$JSON" ]]; then
  echo "No cached conversations found. Run scripts/pr-triage/fetch.sh first." >&2
  exit 1
fi

IDS=$(jq -r 'map(select((.comments|length)>0 and all(.comments[]; .outdated==true) and (.isResolved==false))) | .[].id' "$JSON")
COUNT=0
if [[ -n "$IDS" ]]; then
  while read -r tid; do
    [[ -z "$tid" ]] && continue
    go run ./scripts/resolve-pr-conversation/main.go "$tid" "Auto-resolving: thread is fully outdated." >/dev/null || true
    COUNT=$((COUNT+1))
    sleep 0.1
  done <<< "$IDS"
fi

PR_NUMBER=$(gh pr view --json number -q .number)
go run ./scripts/list-pr-conversations/main.go "$PR_NUMBER" > "$JSON"

echo "Resolved $COUNT fully outdated threads."
