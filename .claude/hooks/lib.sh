#!/bin/bash
# Shared helpers for the conversation-history hooks.
#
# Design (differs from the original linkedin-ai port):
#   1. Per-session state (improvement #1): state lives at
#      tmp/hooks-state/<session_id>.json, so two Claude sessions in the
#      same repo do not stomp on each other. Sub-agents have their own
#      state file: tmp/hooks-state/<session_id>.subagent-<agent_id>.json.
#
#   2. UUID cursor (improvement #2): state carries the uuid of the last
#      JSONL entry we processed. The Stop walk resumes from just past
#      that cursor, so repeated Stop events (AskUserQuestion pauses,
#      resumes, forced-continue) never double-write. Boundary detection
#      by "last real user prompt" is gone.
#
#   3. Stability wait (improvement #3): before parsing the transcript we
#      wait for its mtime+size to stop changing, so a partial flush
#      cannot cause silent tail loss.
#
#   Sidechain / meta rows are filtered inside the walker (improvement #5
#   groundwork for sub-agent isolation).
#
# All paths are relative to ${CLAUDE_PROJECT_DIR}.

if [ -z "${CLAUDE_PROJECT_DIR:-}" ]; then
  echo "history-hooks: CLAUDE_PROJECT_DIR not set, skipping" >&2
  exit 0
fi

HISTORY_DIR="${CLAUDE_PROJECT_DIR}/tmp/history"
STATE_DIR="${CLAUDE_PROJECT_DIR}/tmp/hooks-state"

mkdir -p "$HISTORY_DIR" "$STATE_DIR"

# ---------- payload helpers ----------

# Extract a top-level string field from a JSON payload on stdin.
# Usage:  echo "$payload" | extract_field session_id
extract_field() {
  local field="$1"
  python3 -c "import json,sys; print(json.load(sys.stdin).get('$field',''), end='')" 2>/dev/null
}

# ---------- state files ----------

state_file_for() {
  printf '%s/%s.json' "$STATE_DIR" "$1"
}

subagent_state_file_for() {
  printf '%s/%s.subagent-%s.json' "$STATE_DIR" "$1" "$2"
}

# Read fields from a state file. Empty output means "no state yet".
read_state_filename() {
  local sf="$1"
  [ -f "$sf" ] || return 0
  python3 - "$sf" <<'PY' 2>/dev/null
import json, sys
try:
    with open(sys.argv[1]) as f:
        print(json.load(f).get("filename", ""), end="")
except Exception:
    pass
PY
}

read_state_uuid() {
  local sf="$1"
  [ -f "$sf" ] || return 0
  python3 - "$sf" <<'PY' 2>/dev/null
import json, sys
try:
    with open(sys.argv[1]) as f:
        print(json.load(f).get("last_uuid", ""), end="")
except Exception:
    pass
PY
}

# Atomically rewrite a state file with (filename, last_uuid).
write_state() {
  local sf="$1" filename="$2" last_uuid="$3"
  python3 - "$sf" "$filename" "$last_uuid" <<'PY'
import json, os, sys
path, filename, last_uuid = sys.argv[1], sys.argv[2], sys.argv[3]
tmp = path + ".tmp." + str(os.getpid())
with open(tmp, "w") as f:
    json.dump({"filename": filename, "last_uuid": last_uuid}, f)
os.replace(tmp, path)
PY
}

# ---------- filename generation ----------

slugify() {
  printf '%s' "$1" \
    | head -c 120 \
    | tr '[:upper:]' '[:lower:]' \
    | tr -c 'a-z0-9' '-' \
    | sed -E 's/-+/-/g; s/^-//; s/-$//' \
    | cut -c1-40 \
    | sed -E 's/-$//'
}

# Resolve the role segment for a history filename, in priority order:
#   1. CLAUDE_HISTORY_ROLE explicit (non-empty, non-"0")
#   2. CLAUDE_CODE_ENTRYPOINT=sdk-cli  → "sdk"   (auto-detected `claude -p`)
#   3. "main"
# CLAUDE_HISTORY_ROLE="0" is handled at each hook's top and never
# reaches here.
resolve_role() {
  local r="${CLAUDE_HISTORY_ROLE:-}"
  if [ -n "$r" ] && [ "$r" != "0" ]; then
    printf '%s' "$r"
    return
  fi
  if [ "${CLAUDE_CODE_ENTRYPOINT:-}" = "sdk-cli" ]; then
    printf '%s' "sdk"
    return
  fi
  printf '%s' "main"
}

# Open a fresh history file for a session (main transcript). Naming:
#   tmp/history/<ts>-<role>-<slug>.md
# The role is picked by resolve_role() — see comment there.
start_history_file() {
  local session_id="$1" first_prompt="$2"
  local ts slug role base name n=0
  ts=$(date -u +"%Y%m%d-%H%M%S")
  slug=$(slugify "$first_prompt")
  [ -z "$slug" ] && slug="msg"
  role=$(slugify "$(resolve_role)")
  [ -z "$role" ] && role="main"
  base="${ts}-${role}-${slug}"
  name="${base}.md"
  while ! (set -o noclobber; : > "${HISTORY_DIR}/${name}") 2>/dev/null; do
    n=$((n + 1))
    name="${base}-${n}.md"
    [ "$n" -gt 20 ] && return 1
  done
  write_state "$(state_file_for "$session_id")" "$name" ""
  printf '%s' "$name"
}

# Read the first user-type entry from a sub-agent transcript and echo
#   "<ts-compact>\t<prompt-text>"
# ts-compact is YYYYMMDD-HHMMSS derived from the entry's timestamp
# field, or empty if unavailable. prompt-text is at most 200 chars.
_subagent_first_user_entry() {
  local transcript="$1"
  python3 - "$transcript" <<'PY' 2>/dev/null
import json, re, sys
path = sys.argv[1]
try:
    with open(path) as f:
        for line in f:
            line = line.strip()
            if not line: continue
            try:
                e = json.loads(line)
            except Exception:
                continue
            if e.get("type") != "user":
                continue
            if e.get("isMeta"):
                continue
            content = (e.get("message") or {}).get("content")
            if content is None:
                content = e.get("content")
            text = ""
            if isinstance(content, str):
                text = content
            elif isinstance(content, list):
                parts = []
                for b in content:
                    if isinstance(b, dict) and b.get("type") == "text":
                        t = (b.get("text") or "").strip()
                        if t: parts.append(t)
                text = "\n".join(parts)
            if not text.strip():
                continue
            ts_raw = (e.get("timestamp") or "").strip()
            m = re.match(r'^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})', ts_raw)
            ts_compact = "".join(m.groups()[:3]) + "-" + "".join(m.groups()[3:]) if m else ""
            print(ts_compact + "\t" + text[:200].replace("\t", " ").replace("\n", " "), end="")
            sys.exit(0)
except FileNotFoundError:
    pass
PY
}

# Open a fresh history file for a sub-agent at
#   tmp/history/<sub-ts>-<agent-type>-<prompt-slug>.md
# sub-ts and prompt-slug come from the sub-agent's own transcript
# (first user entry). Falls back to now_ts + short agent_id.
start_subagent_history_file() {
  local session_id="$1" agent_id="$2" agent_type="$3" transcript="$4"
  local ts slug short atype base name n=0 first
  first=$(_subagent_first_user_entry "$transcript")
  ts="${first%%$'\t'*}"
  slug=$(slugify "${first#*$'\t'}")
  [ -z "$ts" ] && ts=$(date -u +"%Y%m%d-%H%M%S")
  short=$(printf '%s' "$agent_id" | tr -c 'a-zA-Z0-9' '-' | cut -c1-8)
  [ -z "$slug" ] && slug="$short"
  atype=$(slugify "$agent_type")
  [ -z "$atype" ] && atype="subagent"
  base="${ts}-${atype}-${slug}"
  name="${base}.md"
  while ! (set -o noclobber; : > "${HISTORY_DIR}/${name}") 2>/dev/null; do
    n=$((n + 1))
    name="${base}-${n}.md"
    [ "$n" -gt 20 ] && return 1
  done
  write_state "$(subagent_state_file_for "$session_id" "$agent_id")" \
              "$name" ""
  printf '%s' "$name"
}

# Return the current git HEAD short SHA, or "-" if the repo is not a
# git checkout / git is unavailable.
git_sha() {
  git -C "$CLAUDE_PROJECT_DIR" rev-parse --short HEAD 2>/dev/null || echo "-"
}

# Return current UTC timestamp in ISO8601 (seconds precision).
now_ts() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

# Append a captured message to the history file recorded in $sf.
# Layout per message:
#
#   ## <heading>
#
#   _<timestamp> · <sha>_
#
#   <body>
#
append_to_history_by_state() {
  local sf="$1" heading="$2" timestamp="$3" sha="$4" body="$5"
  local name
  name=$(read_state_filename "$sf")
  [ -z "$name" ] && return 1
  {
    printf '## %s\n\n' "$heading"
    printf '_%s · %s_\n\n' "$timestamp" "$sha"
    printf '%s\n\n' "$body"
  } >> "${HISTORY_DIR}/${name}"
}

# ---------- transcript I/O ----------

# Wait until $transcript's mtime+size are stable across 2 consecutive
# ~100ms polls (or ~1.5s total). Guards against partial-flush tail loss.
wait_transcript_stable() {
  local transcript="$1"
  python3 - "$transcript" <<'PY' >/dev/null 2>&1
import os, sys, time
path = sys.argv[1]
prev = None
stable = 0
for _ in range(15):
    try:
        st = os.stat(path)
        cur = (st.st_mtime, st.st_size)
    except FileNotFoundError:
        cur = None
    if cur is not None and cur == prev:
        stable += 1
        if stable >= 2:
            break
    else:
        stable = 0
    prev = cur
    time.sleep(0.1)
PY
}

# Walk the transcript JSONL from just past $cursor_uuid, append captured
# entries to $history_file, and echo "<new_last_uuid>\t<count>".
#
# What we capture (per non-sidechain, non-meta entry):
#   - assistant text          → "## claude"
#   - AskUserQuestion tool_use → "## claude (asked)"
#   - AskUserQuestion tool_result → "## user (answered)"
#
# We DO NOT capture user text prompts here — prompt-submit.sh already
# logs those the moment they arrive. Skipping avoids duplication.
#
# The cursor is advanced to the last transcript entry we saw, even for
# entries that produced no records, so we never re-scan them.
dump_from_cursor() {
  local transcript="$1" history_file="$2" cursor_uuid="$3" sha="$4"
  python3 - "$transcript" "$history_file" "$cursor_uuid" "$sha" <<'PY'
import json, re, sys

transcript_path, history_path, cursor_uuid, sha = sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4]

def norm_ts(t):
    # "2026-07-04T16:32:02.911Z" → "2026-07-04T16:32:02Z"
    if not t:
        return "-"
    t = t.strip()
    m = re.match(r'^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.\d+)?Z?$', t)
    return (m.group(1) + "Z") if m else t

entries = []
try:
    with open(transcript_path) as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            try:
                entries.append(json.loads(line))
            except Exception:
                # Partial tail; caller waited for stability so this is rare.
                pass
except FileNotFoundError:
    sys.stdout.write("\t0")
    sys.exit(0)

def content_of(e):
    c = (e.get("message") or {}).get("content")
    if c is None:
        c = e.get("content")
    return c or []

# Locate the cursor. Empty cursor → start from the beginning.
start_idx = 0
if cursor_uuid:
    found = False
    for i, e in enumerate(entries):
        if e.get("uuid") == cursor_uuid:
            start_idx = i + 1
            found = True
            break
    if not found:
        # Cursor not found (transcript replaced / cleared). Restart.
        start_idx = 0

# Build the AskUserQuestion tool_use id → questions lookup across the
# replay window so we can reconstruct answers when their tool_result
# lands.
auq_inputs = {}
for e in entries[start_idx:]:
    if e.get("isSidechain") or e.get("isMeta"):
        continue
    c = content_of(e)
    if not isinstance(c, list):
        continue
    for b in c:
        if (isinstance(b, dict)
                and b.get("type") == "tool_use"
                and b.get("name") == "AskUserQuestion"):
            auq_inputs[b.get("id")] = (b.get("input") or {}).get("questions") or []

def format_questions(questions):
    out = []
    for q in questions:
        qt = (q.get("question") or "").strip()
        if not qt:
            continue
        out.append("- " + qt)
        for opt in (q.get("options") or []):
            lbl = (opt.get("label") or "").strip()
            dsc = (opt.get("description") or "").strip()
            if lbl and dsc:
                out.append("  - " + lbl + " — " + dsc)
            elif lbl:
                out.append("  - " + lbl)
    return "\n".join(out)

def format_answers(questions, result_text):
    # Observed result shape:
    #   'Your questions have been answered: "Q1"="A1", "Q2"="A2". ...'
    # Regex is best-effort; on miss we surface the raw text (truncated)
    # instead of a silent "<skipped>" so format drift is visible.
    out = []
    for q in questions:
        qt = (q.get("question") or "").strip()
        if not qt:
            continue
        needle = re.escape(qt)
        m = re.search(r'"' + needle + r'"\s*=\s*"((?:[^"\\]|\\.)*)"',
                      result_text or "")
        if m:
            try:
                ans = m.group(1).encode().decode("unicode_escape")
            except Exception:
                ans = m.group(1)
        else:
            raw = (result_text or "").strip().replace("\n", " ")[:200]
            ans = "<unparsed — raw: " + raw + ">"
        out.append("- " + qt + "\n  → " + ans)
    return "\n".join(out)

records = []
last_uuid = cursor_uuid

for e in entries[start_idx:]:
    if e.get("isSidechain") or e.get("isMeta"):
        # Skip sub-agent rows folded into the parent transcript, and
        # compact/resume metadata rows. Still advance the cursor.
        if e.get("uuid"):
            last_uuid = e["uuid"]
        continue

    etype = e.get("type")
    uuid = e.get("uuid") or ""
    ts = norm_ts(e.get("timestamp"))
    c = content_of(e)

    if etype == "assistant":
        if isinstance(c, str):
            c = [{"type": "text", "text": c}]
        if isinstance(c, list):
            for b in c:
                if not isinstance(b, dict):
                    continue
                if b.get("type") == "text":
                    text = (b.get("text") or "").strip()
                    if text:
                        records.append(("claude", ts, text))
                elif (b.get("type") == "tool_use"
                      and b.get("name") == "AskUserQuestion"):
                    questions = (b.get("input") or {}).get("questions") or []
                    body = format_questions(questions)
                    if body:
                        records.append(("claude (asked)", ts, body))
    elif etype == "user":
        # Skip pure user text — prompt-submit.sh logged it.
        # Only fold in AskUserQuestion tool_result answers here.
        if isinstance(c, list):
            for b in c:
                if not (isinstance(b, dict) and b.get("type") == "tool_result"):
                    continue
                tu_id = b.get("tool_use_id")
                if tu_id not in auq_inputs:
                    continue
                rc = b.get("content")
                if isinstance(rc, list):
                    result_text = "".join(
                        (x.get("text") or "") for x in rc
                        if isinstance(x, dict) and x.get("type") == "text"
                    )
                else:
                    result_text = rc or ""
                body = format_answers(auq_inputs[tu_id], result_text)
                if body:
                    records.append(("user (answered)", ts, body))

    if uuid:
        last_uuid = uuid

with open(history_path, "a") as f:
    for heading, ts, body in records:
        f.write("## " + heading + "\n\n_" + ts + " · " + sha + "_\n\n" + body + "\n\n")

sys.stdout.write((last_uuid or "") + "\t" + str(len(records)))
PY
}

# Same as dump_from_cursor but with no user-prompt suppression — used
# for sub-agent transcripts, which are complete little conversations
# in their own right.
dump_subagent_from_cursor() {
  local transcript="$1" history_file="$2" cursor_uuid="$3" sha="$4"
  python3 - "$transcript" "$history_file" "$cursor_uuid" "$sha" <<'PY'
import json, re, sys

transcript_path, history_path, cursor_uuid, sha = sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4]

def norm_ts(t):
    if not t:
        return "-"
    t = t.strip()
    m = re.match(r'^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.\d+)?Z?$', t)
    return (m.group(1) + "Z") if m else t

entries = []
try:
    with open(transcript_path) as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            try:
                entries.append(json.loads(line))
            except Exception:
                pass
except FileNotFoundError:
    sys.stdout.write("\t0")
    sys.exit(0)

def content_of(e):
    c = (e.get("message") or {}).get("content")
    if c is None:
        c = e.get("content")
    return c or []

start_idx = 0
if cursor_uuid:
    found = False
    for i, e in enumerate(entries):
        if e.get("uuid") == cursor_uuid:
            start_idx = i + 1
            found = True
            break
    if not found:
        start_idx = 0

records = []
last_uuid = cursor_uuid

for e in entries[start_idx:]:
    if e.get("isMeta"):
        if e.get("uuid"):
            last_uuid = e["uuid"]
        continue
    etype = e.get("type")
    uuid = e.get("uuid") or ""
    ts = norm_ts(e.get("timestamp"))
    c = content_of(e)

    if etype == "user":
        if isinstance(c, str):
            text = c.strip()
            if text:
                records.append(("user", ts, text))
        elif isinstance(c, list):
            texts = []
            for b in c:
                if isinstance(b, dict) and b.get("type") == "text":
                    t = (b.get("text") or "").strip()
                    if t:
                        texts.append(t)
            if texts:
                records.append(("user", ts, "\n\n".join(texts)))
    elif etype == "assistant":
        if isinstance(c, str):
            c = [{"type": "text", "text": c}]
        if isinstance(c, list):
            for b in c:
                if isinstance(b, dict) and b.get("type") == "text":
                    text = (b.get("text") or "").strip()
                    if text:
                        records.append(("claude", ts, text))

    if uuid:
        last_uuid = uuid

with open(history_path, "a") as f:
    for heading, ts, body in records:
        f.write("## " + heading + "\n\n_" + ts + " · " + sha + "_\n\n" + body + "\n\n")

sys.stdout.write((last_uuid or "") + "\t" + str(len(records)))
PY
}
