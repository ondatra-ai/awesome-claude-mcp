#!/usr/bin/env python3
"""Conversation history capture for Claude Code hooks.

One script, three subcommands: prompt-submit, stop, new-task.

State file: tmp/history/hook-state
    JSON object keyed by session_id: {"filename": "...", "n": N}.
    `n` = number of transcript JSONL entries already processed.
    Every update is atomic tmp+rename — file is never modified in place.

History file: tmp/history/<UTC-ts>-<session8>-<slug>.md

Sub-agents: captured from the main transcript itself — Task tool_use
becomes "## claude to <agent_type>", and its tool_result becomes
"## <agent_type>". No SubagentStop hook required.

Off switch: CLAUDE_HISTORY_ROLE=0 skips all logging.
Rollover: /new-task removes the session's entry.
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
STATE_FILE = HISTORY_DIR / "hook-state"


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


def _load_state(sid: str):
    entry = _load_all().get(sid) or {}
    return entry.get("filename", ""), entry.get("n", 0)


def _save_state(sid: str, filename: str, n: int) -> None:
    data = _load_all()
    data[sid] = {"filename": filename, "n": n}
    _save_all(data)


def _open_new_file(sid: str, first_prompt: str) -> str:
    HISTORY_DIR.mkdir(parents=True, exist_ok=True)
    ts = time.strftime("%Y%m%d-%H%M%S", time.gmtime())
    name = f"{ts}-{sid[:8]}-{_slugify(first_prompt)}.md"
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
    """Read JSONL; drop malformed lines (partial tail — next Stop catches)."""
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


def _walk(entries, start_idx: int):
    """Yield (heading, ts, body) for entries[start_idx:].

    Two-pass: first build a full-history lookup of Task tool_use ids →
    subagent_type so that a tool_result appearing after start_idx (but
    whose tool_use was before) can still be attributed to the right
    agent.
    """
    task_types = {}
    for e in entries:
        if e.get("isSidechain") or e.get("isMeta"):
            continue
        if e.get("type") != "assistant":
            continue
        for b in _content(e) or []:
            if (isinstance(b, dict)
                    and b.get("type") == "tool_use"
                    and b.get("name") == "Task"):
                at = ((b.get("input") or {}).get("subagent_type") or "subagent").strip()
                task_types[b.get("id")] = at

    for e in entries[start_idx:]:
        if e.get("isSidechain") or e.get("isMeta"):
            continue
        ts = _norm_ts(e.get("timestamp"))
        etype = e.get("type")
        c = _content(e)
        if etype == "assistant":
            if isinstance(c, str):
                c = [{"type": "text", "text": c}]
            for b in c or []:
                if not isinstance(b, dict):
                    continue
                if b.get("type") == "text" and (b.get("text") or "").strip():
                    yield "claude", ts, b["text"].strip()
                elif b.get("type") == "tool_use" and b.get("name") == "Task":
                    at = ((b.get("input") or {}).get("subagent_type") or "subagent").strip()
                    p = ((b.get("input") or {}).get("prompt") or "").strip()
                    if p:
                        yield f"claude to {at}", ts, p
        elif etype == "user" and isinstance(c, list):
            for b in c:
                if not (isinstance(b, dict) and b.get("type") == "tool_result"):
                    continue
                at = task_types.get(b.get("tool_use_id"))
                if not at:
                    continue
                rc = b.get("content")
                if isinstance(rc, list):
                    text = "".join(
                        (x.get("text") or "") for x in rc
                        if isinstance(x, dict) and x.get("type") == "text"
                    )
                else:
                    text = rc or ""
                text = text.strip()
                if text:
                    yield at, ts, text


# ---------- event handlers ----------

def new_task(_payload: dict) -> None:
    """Drop the session's entry so the next prompt opens a fresh file."""
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
        _save_state(sid, filename, 0)
    _append(filename, "user", _now(), prompt)


def stop(payload: dict) -> None:
    if payload.get("stop_hook_active"):
        return
    sid = payload.get("session_id") or ""
    tp = payload.get("transcript_path") or ""
    if not sid or not tp:
        return
    filename, n = _load_state(sid)
    if not filename:
        return
    entries = _read_jsonl(tp)
    for heading, ts, body in _walk(entries, n):
        _append(filename, heading, ts, body)
    if len(entries) != n:
        _save_state(sid, filename, len(entries))


HANDLERS = {
    "prompt-submit": prompt_submit,
    "stop": stop,
    "new-task": new_task,
}


def main() -> None:
    if os.environ.get("CLAUDE_HISTORY_ROLE") == "0":
        return
    if not os.environ.get("CLAUDE_PROJECT_DIR"):
        return
    if len(sys.argv) < 2 or sys.argv[1] not in HANDLERS:
        return
    if sys.argv[1] == "new-task":  # reads session_id from env, not stdin
        HANDLERS["new-task"]({})
        return
    try:
        payload = json.load(sys.stdin)
    except Exception:
        return
    HANDLERS[sys.argv[1]](payload)


if __name__ == "__main__":
    main()
