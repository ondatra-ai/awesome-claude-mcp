#!/usr/bin/env bash
set -euo pipefail

TMP_DIR="./tmp"
mkdir -p "$TMP_DIR"

PR_NUMBER=$(gh pr view --json number -q .number)
DEFAULT_BRANCH=$(gh repo view --json defaultBranchRef -q .defaultBranchRef.name)

# Fetch all review threads (resolved + unresolved)
go run ./scripts/list-pr-conversations/main.go "$PR_NUMBER" > "$TMP_DIR/PR_CONV.json"

# Auto‑resolve fully outdated threads
AUTO_IDS=$(jq -r 'map(select( all(.comments[]; .outdated==true) and (.isResolved==false) )) | .[].id' "$TMP_DIR/PR_CONV.json")
if [[ -n "$AUTO_IDS" ]]; then
  while read -r tid; do
    [[ -z "$tid" ]] && continue
    go run ./scripts/resolve-pr-conversation/main.go "$tid" "Auto-resolving: thread is fully outdated." >/dev/null || true
    sleep 0.1
  done <<< "$AUTO_IDS"
fi

# Refresh after auto‑resolve
go run ./scripts/list-pr-conversations/main.go "$PR_NUMBER" > "$TMP_DIR/PR_CONV.json"

# Compute changed files vs base (scope inference)
git fetch origin "$DEFAULT_BRANCH" --quiet >/dev/null 2>&1 || true
git diff --name-only "origin/$DEFAULT_BRANCH"...HEAD > "$TMP_DIR/CHANGED_FILES.txt"

# Pick first unresolved thread that still has at least one non‑outdated comment
jq -r '[.[] | select(.isResolved==false and any(.comments[]; .outdated!=true))][0] // empty' "$TMP_DIR/PR_CONV.json" > "$TMP_DIR/CURRENT_THREAD.json"

# If nothing to triage, exit with a clear, single‑line message
if [[ ! -s "$TMP_DIR/CURRENT_THREAD.json" ]] || [[ "$(jq -r 'keys|length' "$TMP_DIR/CURRENT_THREAD.json")" == "0" ]]; then
  echo "Nothing to triage right now."
  exit 0
fi

# Extract key fields
THREAD_ID=$(jq -r '.id' "$TMP_DIR/CURRENT_THREAD.json")
FILE_PATH=$(jq -r '.comments[0].file // ""' "$TMP_DIR/CURRENT_THREAD.json")
LINE_NO=$(jq -r '.comments[0].line // 0' "$TMP_DIR/CURRENT_THREAD.json")
AUTHOR=$(jq -r '.comments[0].author // ""' "$TMP_DIR/CURRENT_THREAD.json")
LINK=$(jq -r '.comments[0].url // ""' "$TMP_DIR/CURRENT_THREAD.json")
COMMENT_FULL=$(jq -r '.comments[0].body // ""' "$TMP_DIR/CURRENT_THREAD.json")

# No code preview required; we'll show the full review comment content instead.

# Preferred option heuristic (scope: changed files or known project areas)
PREFERRED="Create ticket"
if [[ -n "$FILE_PATH" ]] && grep -qxF "$FILE_PATH" "$TMP_DIR/CHANGED_FILES.txt" 2>/dev/null; then
  PREFERRED="Proceed fix"
elif [[ "$FILE_PATH" == infrastructure/* || "$FILE_PATH" == services/* ]]; then
  PREFERRED="Proceed fix"
fi

# Proposed fix (generic, avoids leaking internal heuristics)
PROPOSED_FIX="Implement the reviewer’s suggestion in a minimal, scoped change aligned with architecture and coding standards; validate via terraform validate/plan or unit/E2E as applicable."

# Risk analysis and score (0–10) — simple heuristic
RISK_ANALYSIS="Targeted change, limited blast radius; verify in plan/tests before merge."
RISK_SCORE=4
if [[ "$FILE_PATH" == infrastructure/terraform/modules/alb/* ]]; then RISK_SCORE=5; fi
if [[ "$FILE_PATH" == infrastructure/terraform/environments/*/main.tf ]]; then RISK_SCORE=3; fi

# Print a clean decision package only (no internal commands/logic)
echo "Thread: $THREAD_ID"
echo "Link: $LINK"
echo "Location: ${FILE_PATH:-unknown}:${LINE_NO}"
echo "Comment:"
printf "%s\n" "$COMMENT_FULL"
echo "Proposed Fix: $PROPOSED_FIX"
echo "Risk Analysis: $RISK_ANALYSIS"
echo "Risk: $RISK_SCORE/10"
echo "Decision: $PREFERRED"
