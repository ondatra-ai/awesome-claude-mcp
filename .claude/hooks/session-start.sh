#!/bin/bash
# SessionStart hook.
#   source=startup|clear → drop this session's state (main + any
#     sub-agent state) so the next prompt opens a new history file.
#   source=resume|compact → keep state; capture continues into the
#     same history file with the existing cursor.
set -u

if [ -z "${CLAUDE_PROJECT_DIR:-}" ]; then
  echo "history-hooks: CLAUDE_PROJECT_DIR not set, skipping" >&2
  exit 0
fi

# CLAUDE_HISTORY_ROLE=0 → hooks no-op (used by tooling that spawns
# `claude -p` and doesn't want a history file cluttering tmp/history).
[ "${CLAUDE_HISTORY_ROLE:-}" = "0" ] && exit 0

source "${CLAUDE_PROJECT_DIR}/.claude/hooks/lib.sh"

payload=$(cat)
session_id=$(printf '%s' "$payload" | extract_field session_id)
source_kind=$(printf '%s' "$payload" | extract_field source)

[ -z "$session_id" ] && exit 0

case "$source_kind" in
  startup|clear)
    rm -f "$(state_file_for "$session_id")"
    # Best-effort cleanup of any lingering sub-agent state for this
    # session (glob is safe: session_id contains no shell metachars).
    rm -f "${STATE_DIR}/${session_id}.subagent-"*.json 2>/dev/null || true
    ;;
esac

exit 0
