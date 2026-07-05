#!/bin/bash
# SubagentStop hook (improvement #4).
#
# Sub-agents run in their own JSONL transcript (agent_transcript_path).
# We capture each one into a dedicated file at
#   tmp/history/<sub-ts>-<agent-type>-<prompt-slug>.md
# with its own per-subagent cursor so it doesn't tangle with the parent
# session's history. Naming keeps main and sub-agent files coexisting
# flat in the same directory: main files carry a "-main-" segment,
# sub-agent files carry their agent type instead.
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

sf=$(subagent_state_file_for "$session_id" "$agent_id")

# Lazily open the sub-agent history file on first Stop. We pass the
# transcript path so the helper can derive <ts> and <slug> from the
# sub-agent's first user entry.
if [ ! -f "$sf" ]; then
  start_subagent_history_file "$session_id" "$agent_id" "$agent_type" "$sub_transcript" >/dev/null \
    || { step "  ERROR: failed to open sub-agent history file"; emit_and_exit; }
fi

name=$(read_state_filename "$sf")
if [ -z "$name" ]; then
  step "  state has no filename — skipped"
  emit_and_exit
fi
HISTORY_FILE="${HISTORY_DIR}/${name}"

wait_transcript_stable "$sub_transcript"

cursor=$(read_state_uuid "$sf")
sha=$(git_sha)
result=$(dump_subagent_from_cursor "$sub_transcript" "$HISTORY_FILE" "$cursor" "$sha" 2>/dev/null || true)
new_cursor=${result%%$'\t'*}
count=${result##*$'\t'}
case "$count" in ''|*[!0-9]*) count=0 ;; esac

if [ -n "$new_cursor" ] && [ "$new_cursor" != "$cursor" ]; then
  write_state "$sf" "$name" "$new_cursor"
fi

if [ "$count" -eq 0 ]; then
  step "  nothing new to log (cursor: ${cursor:-<start>})"
else
  step "  Saved $count entries → $name (cursor → ${new_cursor:0:8}…)"
fi

emit_and_exit
