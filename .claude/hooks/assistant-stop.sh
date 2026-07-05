#!/bin/bash
# Stop hook: walk the main transcript from the persisted UUID cursor,
# append captured entries (assistant text + AskUserQuestion Q&A) to
# this session's history file, then advance the cursor.
#
# Idempotent — replayed Stop events (AskUserQuestion pauses, forced
# continuation, retries) contribute zero new records because the
# cursor moves past every entry we've already seen.
#
# Always emits a one-line diagnostic on hookSpecificOutput.systemMessage.
set -u

if [ -z "${CLAUDE_PROJECT_DIR:-}" ]; then
  echo "history-hooks: CLAUDE_PROJECT_DIR not set, skipping" >&2
  exit 0
fi

# CLAUDE_HISTORY_ROLE=0 → no-op.
[ "${CLAUDE_HISTORY_ROLE:-}" = "0" ] && exit 0

source "${CLAUDE_PROJECT_DIR}/.claude/hooks/lib.sh"

DIAG=""
step() { DIAG+="$1"$'\n'; }
emit_and_exit() {
  DIAG="$DIAG" python3 - <<'PY'
import json, os
msg = "[history hook]\n" + os.environ["DIAG"].rstrip()
print(json.dumps({"systemMessage": msg}))
PY
  exit 0
}

payload=$(cat)
session_id=$(printf '%s' "$payload" | extract_field session_id)
transcript=$(printf '%s' "$payload" | extract_field transcript_path)
stop_hook_active=$(printf '%s' "$payload" | python3 -c \
  "import json,sys; print(json.load(sys.stdin).get('stop_hook_active', False))" 2>/dev/null)

if [ -z "$session_id" ]; then
  step "  ERROR: no session_id"
  emit_and_exit
fi

# Forced-continue state: main loop already fired us once. Nothing new
# to add — bail before writing to avoid churn.
if [ "$stop_hook_active" = "True" ]; then
  step "  stop_hook_active — skip"
  emit_and_exit
fi

sf=$(state_file_for "$session_id")
if [ ! -f "$sf" ]; then
  step "  no state for session $session_id — skipped"
  emit_and_exit
fi

name=$(read_state_filename "$sf")
if [ -z "$name" ]; then
  step "  state has no filename — skipped"
  emit_and_exit
fi
HISTORY_FILE="${HISTORY_DIR}/${name}"

if [ -z "$transcript" ] || [ ! -f "$transcript" ]; then
  step "  ERROR: no transcript_path in Stop payload (got: '$transcript')"
  emit_and_exit
fi

# Wait for the transcript to stop growing so a partial flush doesn't
# cause silent tail loss.
wait_transcript_stable "$transcript"

cursor=$(read_state_uuid "$sf")
sha=$(git_sha)
result=$(dump_from_cursor "$transcript" "$HISTORY_FILE" "$cursor" "$sha" 2>/dev/null || true)
new_cursor=${result%%$'\t'*}
count=${result##*$'\t'}
case "$count" in ''|*[!0-9]*) count=0 ;; esac

# Persist advanced cursor even if count=0 — we still saw new entries
# we don't need to re-scan next time.
if [ -n "$new_cursor" ] && [ "$new_cursor" != "$cursor" ]; then
  write_state "$sf" "$name" "$new_cursor"
fi

if [ "$count" -eq 0 ]; then
  step "  nothing new to log (cursor: ${cursor:-<start>})"
else
  step "  Saved $count entries → $name (cursor → ${new_cursor:0:8}…)"
fi

emit_and_exit
