#!/usr/bin/env python3
"""Conversation history capture for Claude Code hooks.

One script, five subcommands: session-start, prompt-submit, stop,
subagent-stop, new-task.

State file: tmp/history/hook-state.json
    Single JSON object keyed by session_id (and session_id.sub.agent_id
    for sub-agent cursors). Every update is delete-and-create (atomic
    tmp+rename) — file is never modified in place.

History file: tmp/history/<UTC-ts>-<session8>-<slug>.md

Off switch: CLAUDE_HISTORY_ROLE=0 skips all logging.
Rollover: /new-task removes the session's main entry from hook-state.json.
"""

import json
import os
import re
import subprocess
import sys
import time
from pathlib import Path


REPO = Path(os.environ.get("CLAUDE_PROJECT_DIR", "")).resolve()
HISTORY_DIR = REPO / "tmp" / "history"
STATE_FILE = HISTORY_DIR / "hook-state.json"


def _slugify(s: str) -> str:
    s = s[:120].lower()
    s = re.sub(r"[^a-z0-9]+", "-", s).strip("-")[:40].rstrip("-")
    return s or "msg"


def _now() -> str:
    return time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())


def _norm_ts(t: str) -> str:
    m = re.match(r"^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})", (t or "").strip())
    return (m.group(1) + "Z") if m else (t or "-")


def _git_sha() -> str:
    try:
        r = subprocess.run(
            ["git", "-C", str(REPO), "rev-parse", "--short", "HEAD"],
            capture_output=True, text=True, timeout=2,
        )
        return (r.stdout.strip() or "-") if r.returncode == 0 else "-"
    except Exception:
        return "-"


def _load_all() -> dict:
    if not STATE_FILE.exists():
        return {}
    try:
        return json.loads(STATE_FILE.read_text()) or {}
    except Exception:
        return {}


def _save_all(data: dict) -> None:
    HISTORY_DIR.mkdir(parents=True, exist_ok=True)
    tmp = STATE_FILE.with_name(STATE_FILE.name + f".tmp.{os.getpid()}")
    tmp.write_text(json.dumps(data))
    tmp.replace(STATE_FILE)


def _load_state(key: str):
    entry = _load_all().get(key) or {}
    return entry.get("filename", ""), entry.get("last_uuid", "")


def _save_state(key: str, filename: str, last_uuid: str) -> None:
    data = _load_all()
    data[key] = {"filename": filename, "last_uuid": last_uuid}
    _save_all(data)


def _open_new_file(session_id: str, first_prompt: str) -> str:
    HISTORY_DIR.mkdir(parents=True, exist_ok=True)
    ts = time.strftime("%Y%m%d-%H%M%S", time.gmtime())
    name = f"{ts}-{session_id[:8]}-{_slugify(first_prompt)}.md"
    (HISTORY_DIR / name).touch(exist_ok=True)
    return name


def _append(filename: str, heading: str, ts: str, body: str) -> None:
    line = f"## {heading}\n\n_{ts} · {_git_sha()}_\n\n{body}\n\n"
    with (HISTORY_DIR / filename).open("a") as f:
        f.write(line)


def _content(entry: dict):
    c = (entry.get("message") or {}).get("content")
    return c if c is not None else (entry.get("content") or [])


def _read_jsonl(path: str):
    """Read JSONL; silently drop malformed lines (partial tail — next Stop catches)."""
    out = []
    try:
        for line in Path(path).read_text().splitlines():
            line = line.strip()
            if not line:
                continue
            try:
                out.append(json.loads(line))
            except Exception:
                pass
    except FileNotFoundError:
        pass
    return out


def _from_cursor(entries, cursor: str):
    if not cursor:
        yield from entries
        return
    seen = False
    for e in entries:
        if seen:
            yield e
        elif e.get("uuid") == cursor:
            seen = True
    if not seen:
        # Cursor not found (transcript replaced) → replay all.
        yield from entries


def _last_uuid(entries) -> str:
    for e in reversed(entries):
        if e.get("uuid"):
            return e["uuid"]
    return ""


def _walk_main(entries):
    for e in entries:
        if e.get("isSidechain") or e.get("isMeta"):
            continue
        if e.get("type") != "assistant":
            continue
        ts = _norm_ts(e.get("timestamp"))
        c = _content(e)
        if isinstance(c, str):
            c = [{"type": "text", "text": c}]
        for b in c or []:
            if not isinstance(b, dict):
                continue
            if b.get("type") == "text" and (b.get("text") or "").strip():
                yield "claude", ts, b["text"].strip()
            elif b.get("type") == "tool_use" and b.get("name") == "Task":
                inp = b.get("input") or {}
                at = (inp.get("subagent_type") or "subagent").strip()
                p = (inp.get("prompt") or "").strip()
                if p:
                    yield f"claude to {at}", ts, p


def _walk_subagent(entries, agent_type: str):
    for e in entries:
        if e.get("isMeta"):
            continue
        if e.get("type") != "assistant":
            continue
        ts = _norm_ts(e.get("timestamp"))
        c = _content(e)
        if isinstance(c, str):
            c = [{"type": "text", "text": c}]
        for b in c or []:
            if isinstance(b, dict) and b.get("type") == "text":
                t = (b.get("text") or "").strip()
                if t:
                    yield agent_type, ts, t


# ---------- event handlers ----------

def session_start(payload: dict) -> None:
    if payload.get("source") not in ("startup", "clear"):
        return
    sid = payload.get("session_id") or ""
    if not sid:
        return
    data = _load_all()
    prefix = f"{sid}.sub."
    for key in list(data.keys()):
        if key == sid or key.startswith(prefix):
            data.pop(key)
    _save_all(data)


def new_task(_payload: dict) -> None:
    """Remove the session's main entry so the next prompt opens a fresh
    history file. Sub-agent cursors under <sid>.sub.<aid> are preserved."""
    sid = os.environ.get("CLAUDE_CODE_SESSION_ID") or ""
    if not sid:
        return
    data = _load_all()
    if sid in data:
        data.pop(sid)
        _save_all(data)


def prompt_submit(payload: dict) -> None:
    sid = payload.get("session_id") or ""
    prompt = payload.get("prompt") or payload.get("user_message") or ""
    if not sid or not prompt.strip():
        return
    filename, _ = _load_state(sid)
    if not filename:
        filename = _open_new_file(sid, prompt)
        _save_state(sid, filename, "")
    _append(filename, "user", _now(), prompt)


def stop(payload: dict) -> None:
    if payload.get("stop_hook_active"):
        return
    sid = payload.get("session_id") or ""
    tp = payload.get("transcript_path") or ""
    if not sid or not tp:
        return
    filename, cursor = _load_state(sid)
    if not filename:
        return
    entries = _read_jsonl(tp)
    for heading, ts, body in _walk_main(list(_from_cursor(entries, cursor))):
        _append(filename, heading, ts, body)
    new_cursor = _last_uuid(entries) or cursor
    if new_cursor and new_cursor != cursor:
        _save_state(sid, filename, new_cursor)


def subagent_stop(payload: dict) -> None:
    if payload.get("stop_hook_active"):
        return
    sid = payload.get("session_id") or ""
    aid = payload.get("agent_id") or ""
    at = payload.get("agent_type") or "subagent"
    tp = payload.get("agent_transcript_path") or payload.get("transcript_path") or ""
    if not sid or not aid or not tp:
        return
    filename, _ = _load_state(sid)
    if not filename:
        return  # main state missing (task rollover in flight); cursor preserved
    sub_key = f"{sid}.sub.{aid}"
    _, cursor = _load_state(sub_key)
    entries = _read_jsonl(tp)
    for heading, ts, body in _walk_subagent(list(_from_cursor(entries, cursor)), at):
        _append(filename, heading, ts, body)
    new_cursor = _last_uuid(entries) or cursor
    if new_cursor and new_cursor != cursor:
        _save_state(sub_key, "", new_cursor)


HANDLERS = {
    "session-start": session_start,
    "prompt-submit": prompt_submit,
    "stop": stop,
    "subagent-stop": subagent_stop,
    "new-task": new_task,
}


def main() -> None:
    if os.environ.get("CLAUDE_HISTORY_ROLE") == "0":
        return
    if not os.environ.get("CLAUDE_PROJECT_DIR"):
        return
    if len(sys.argv) < 2 or sys.argv[1] not in HANDLERS:
        return
    # new-task reads session_id from env, not stdin.
    if sys.argv[1] == "new-task":
        HANDLERS["new-task"]({})
        return
    try:
        payload = json.load(sys.stdin)
    except Exception:
        return
    HANDLERS[sys.argv[1]](payload)


if __name__ == "__main__":
    main()
