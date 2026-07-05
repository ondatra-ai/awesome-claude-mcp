#!/bin/bash
# SubagentStop hook.
#
# Sub-agents run in their own JSONL transcript (agent_transcript_path).
# Their captured entries are appended to the CURRENT TASK's history file
# — the same file the main session writes to — so a task's history is
# a single readable record of main + every sub-agent it spawned.
#
# The current task's filename is looked up from the main session's state
# file (tmp/hooks-state/<session_id>.json). Sub-agent state is a cursor
# only (last_uuid), no filename.
#
# Speaker attribution:
#   - Sub-agent's assistant text → "## <agent_type>"
#   - Sub-agent's own user-role entries are skipped — main's
#     "## claude to <agent_type>" already captured that side.
#
# Edge cases:
#   - Main state missing (user ran /new-task and hasn't typed the next
#     prompt yet) → skip this batch and log; cursor unchanged so nothing
#     is lost.
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
msg = "[subagent history hook]\n" + os.environ["DIAG"].rstrip()
print(json.dumps({"systemMessage": msg}))
PY
  exit 0
}

payload=$(cat)
session_id=$(printf '%s' "$payload" | extract_field session_id)
agent_id=$(printf '%s' "$payload" | extract_field agent_id)
agent_type=$(printf '%s' "$payload" | extract_field agent_type)
sub_transcript=$(printf '%s' "$payload" | extract_field agent_transcript_path)
# Fallback: some hook versions only ship the parent transcript_path.
if [ -z "$sub_transcript" ]; then
  sub_transcript=$(printf '%s' "$payload" | extract_field transcript_path)
fi
stop_hook_active=$(printf '%s' "$payload" | python3 -c \
  "import json,sys; print(json.load(sys.stdin).get('stop_hook_active', False))" 2>/dev/null)

if [ -z "$session_id" ] || [ -z "$agent_id" ]; then
  step "  ERROR: session_id or agent_id missing"
  emit_and_exit
fi
if [ "$stop_hook_active" = "True" ]; then
  step "  stop_hook_active — skip"
  emit_and_exit
fi
if [ -z "$sub_transcript" ] || [ ! -f "$sub_transcript" ]; then
  step "  ERROR: no agent_transcript_path (got: '$sub_transcript')"
  emit_and_exit
fi

# Look up the current task's history file from the main session's state.
main_sf=$(state_file_for "$session_id")
name=$(read_state_filename "$main_sf")
if [ -z "$name" ]; then
  step "  main state missing (task rollover in progress?) — skipped; cursor preserved"
  emit_and_exit
fi
HISTORY_FILE="${HISTORY_DIR}/${name}"

# Sub-agent state file holds the cursor only. Create empty state on
# first Stop for this sub-agent so subsequent Stops can resume from
# where we leave off.
sf=$(subagent_state_file_for "$session_id" "$agent_id")
if [ ! -f "$sf" ]; then
  write_state "$sf" "" ""
fi

wait_transcript_stable "$sub_transcript"

cursor=$(read_state_uuid "$sf")
sha=$(git_sha)
[ -z "$agent_type" ] && agent_type="subagent"
result=$(dump_subagent_from_cursor \
  "$sub_transcript" "$HISTORY_FILE" "$cursor" "$sha" "$agent_type" 2>/dev/null || true)
new_cursor=${result%%$'\t'*}
count=${result##*$'\t'}
case "$count" in ''|*[!0-9]*) count=0 ;; esac

if [ -n "$new_cursor" ] && [ "$new_cursor" != "$cursor" ]; then
  # Preserve empty filename in state; we don't own a file.
  write_state "$sf" "" "$new_cursor"
fi

if [ "$count" -eq 0 ]; then
  step "  nothing new to log (cursor: ${cursor:-<start>})"
else
  step "  Saved $count $agent_type entries → $name (cursor → ${new_cursor:0:8}…)"
fi

emit_and_exit
