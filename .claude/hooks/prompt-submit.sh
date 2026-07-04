#!/bin/bash
# UserPromptSubmit hook: append the user's prompt to the session's
# current history file. If this is the session's first prompt (no
# state file yet), derive the filename slug from this prompt and
# create the file.
set -u

if [ -z "${CLAUDE_PROJECT_DIR:-}" ]; then
  echo "history-hooks: CLAUDE_PROJECT_DIR not set, skipping" >&2
  exit 0
fi

source "${CLAUDE_PROJECT_DIR}/.claude/hooks/lib.sh"

payload=$(cat)
session_id=$(printf '%s' "$payload" | extract_field session_id)
prompt=$(printf '%s' "$payload" | python3 -c "
import json, sys
try:
    d = json.load(sys.stdin)
except Exception:
    sys.exit(0)
# Claude Code uses 'prompt'; tolerate 'user_message' too.
print(d.get('prompt') or d.get('user_message') or '', end='')
" 2>/dev/null)

[ -z "$session_id" ] && exit 0
[ -z "$prompt" ] && exit 0

sf=$(state_file_for "$session_id")
if [ ! -f "$sf" ]; then
  start_history_file "$session_id" "$prompt" >/dev/null || exit 0
fi

append_to_history_by_state "$sf" "user" "$prompt"
exit 0
