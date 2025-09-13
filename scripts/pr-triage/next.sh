#!/usr/bin/env bash
set -euo pipefail

TMP_DIR="./tmp"
JSON="$TMP_DIR/PR_CONV.json"

if [[ ! -f "$JSON" ]]; then
  echo "No cached conversations found. Run scripts/pr-triage/fetch.sh first." >&2
  exit 1
fi

DEFAULT_BRANCH=$(gh repo view --json defaultBranchRef -q .defaultBranchRef.name)
git fetch origin "$DEFAULT_BRANCH" --quiet >/dev/null 2>&1 || true
git diff --name-only "origin/$DEFAULT_BRANCH"...HEAD > "$TMP_DIR/CHANGED_FILES.txt"

# Select first unresolved thread with at least one non‑outdated comment
jq -r 'first(.[] | select(.isResolved==false and any(.comments[]; .outdated!=true)))' "$JSON" > "$TMP_DIR/CURRENT_THREAD.json"

if [[ ! -s "$TMP_DIR/CURRENT_THREAD.json" ]] || [[ "$(jq -r 'keys|length' "$TMP_DIR/CURRENT_THREAD.json")" == "0" ]]; then
  echo "Nothing to triage right now."
  exit 0
fi

THREAD_ID=$(jq -r '.id' "$TMP_DIR/CURRENT_THREAD.json")
FILE_PATH=$(jq -r '.comments | map(select(.outdated!=true)) | .[0].file // ""' "$TMP_DIR/CURRENT_THREAD.json")
LINE_NO=$(jq -r '.comments | map(select(.outdated!=true)) | .[0].line // 0' "$TMP_DIR/CURRENT_THREAD.json")
LINK=$(jq -r '.comments | map(select(.outdated!=true)) | .[0].url // ""' "$TMP_DIR/CURRENT_THREAD.json")
COMMENT_FULL=$(jq -r '.comments | map(select(.outdated!=true)) | .[0].body // ""' "$TMP_DIR/CURRENT_THREAD.json")

# Preferred option heuristic (based on scope)
PREFERRED="Create ticket"
if [[ -n "$FILE_PATH" ]] && grep -qxF "$FILE_PATH" "$TMP_DIR/CHANGED_FILES.txt" 2>/dev/null; then
  PREFERRED="Proceed fix"
elif [[ "$FILE_PATH" == infrastructure/* || "$FILE_PATH" == services/* ]]; then
  PREFERRED="Proceed fix"
fi

# Output in the standardized format
echo "Thread: $THREAD_ID"
echo "Link: $LINK"
echo "Location: ${FILE_PATH:-unknown}:${LINE_NO}"
echo "Comment:"
printf "%s\n" "$COMMENT_FULL"
echo "Proposed Fix: Implement the reviewer’s suggestion in a minimal, scoped change aligned with architecture and coding standards; validate via terraform validate/plan or unit/E2E as applicable."
echo "Risk Analysis: Targeted change, limited blast radius; verify in plan/tests before merge."
echo "Risk: 5/10"
echo "Decision: $PREFERRED"
