<!-- Powered by BMAD™ Core -->

# pr-triage

Stage and process PR review conversations deterministically using local files.

## Purpose

Provide a predictable, file‑based workflow to read all PR conversations once, track processed items, and surface the next conversation for analysis.

## Artifacts

- `scripts/list-pr-conversations/main.go`
  - Reads review threads for the current PR (or a provided PR number) and writes the full conversations JSON array to a file.
  - Usage: `go run scripts/list-pr-conversations/main.go tmp/CONV.json [PR_NUMBER]`
  - Guarantees: writes to `tmp/CONV.json` (or specified path), produces no JSON on stdout.

- `tmp/CONV.json`
  - Cached conversations JSON (array of conversation objects) used as the single source of truth for the triage run.

- `tmp/CONV_ID.txt`
  - Line‑separated list of processed conversation IDs.
  - Empty by default; updated after each conversation is processed.

- Next selection script (to be written)
  - Reads `tmp/CONV.json` and `tmp/CONV_ID.txt`.
  - Writes the first unprocessed conversation object to `tmp/CONV_CURRENT.json`.
  - If none remain, writes exactly `{ "id": "No More Converations" }` to `tmp/CONV_CURRENT.json`.
  - Suggested path/name: `scripts/pr-triage/next.go` (or equivalent).

- `tmp/CONV_CURRENT.json`
  - The single conversation object to analyze next, or the sentinel `{ "id": "No More Converations" }` when finished.

## Process

1) Initialization
   - Delete `tmp/CONV.json` if it exists
   - Delete `tmp/CONV_ID.txt` if it exists
   - Delete `tmp/CONV_CURRENT.json` if it exists
   - Create/clear `tmp/CONV_ID.txt` (empty file). Note: Do this ONCE per triage session. Do NOT clear inside the loop.

2) Read conversations to `tmp/CONV.json`.
   - Run: `go run scripts/list-pr-conversations/main.go tmp/CONV.json [PR_NUMBER]`
   - MUST write the full conversations JSON array to `tmp/CONV.json`.
   - MUST NOT print JSON to standard output.

3) Identify the next conversation to process.
   - Run: `go run scripts/pr-triage/next-pr-conversation.go tmp/CONV.json tmp/CONV_ID.txt tmp/CONV_CURRENT.json`

4) Read `tmp/CONV_CURRENT.json`.
   - If it contains `{ "id": "No More Converations" }`, stop.
   - Otherwise, proceed with heuristic analysis for that conversation.
5) Append ID to `tmp/CONV_ID.txt`.
   - Append processed to `tmp/CONV_ID.txt`.
6) Repeat actions 3 → 5 until sentinel appears
   - Ensure any required thread reply was posted (per outcome in Heuristic analysis)
   - Verify the conversation ID was appended to `tmp/CONV_ID.txt` (no duplicates)
   - Run: `go run scripts/pr-triage/next-pr-conversation.go tmp/CONV.json tmp/CONV_ID.txt tmp/CONV_CURRENT.json`
   - If `tmp/CONV_CURRENT.json` contains `{ "id": "No More Converations" }`, stop (do NOT print blocks for the sentinel)
   - Otherwise, continue with item 4
   - Append the conversation ID to `tmp/CONV_ID.txt` AFTER meeting the Output Gate conditions for that conversation

### Continuous Triage Loop (example)
```bash
# One-time initialization (outside loop)
rm -f tmp/CONV.json tmp/CONV_ID.txt tmp/CONV_CURRENT.json
mkdir -p tmp
go run scripts/list-pr-conversations/main.go tmp/CONV.json
: > tmp/CONV_ID.txt

# Loop: identify + analyze + append until sentinel
while true; do
  go run scripts/pr-triage/next-pr-conversation.go tmp/CONV.json tmp/CONV_ID.txt tmp/CONV_CURRENT.json
  if jq -e '.id == "No More Converations"' tmp/CONV_CURRENT.json >/dev/null; then
    echo "No more conversations."
    break
  fi

  # Heuristic analysis (print checklist)
  # Risk branch: Low (<5) → auto-apply + print Action Block; High (≥5) → print Approval Block
  # Post required thread reply per branch

  # Append processed ID (after Output Gate conditions)
  jq -r '.id' tmp/CONV_CURRENT.json >> tmp/CONV_ID.txt
done
```

## Heuristic analysis
1) Perform Heuristic checklist
   - Use `.bmad-core/checklists/triage-heuristic-checklist.md` against the single conversation in `tmp/CONV_CURRENT.json`.
   - MUST print the Heuristic Checklist Result block exactly once per conversation (see "Heuristic Checklist Output").

2) Determine the best option to proceed
   - Choose one: implement now, request changes/clarification, defer and create an issue, or escalate (architect/PO/QA).
   - Base the choice on alignment, scope fit, and value/effort.

3) Determine the risk score (1–10)
   - 1–3: low (cosmetic, doc, trivial refactor, isolated change)
   - 4–6: medium (localized behavior, limited blast radius)
   - 7–10: high (architecture/security/performance implications or wide impact)
   - Record the numeric score and rationale.
   - Thresholds (authoritative): Low risk = risk < 5. Medium/High risk = risk ≥ 5.

4) If risk < 5, implement without human-in-the-loop
   - DO NOT COMMIT CHANGES
   - Conditions: strictly within PR scope, passes checklist alignment, and all tests pass locally.
   - MUST print BOTH blocks in this order:
     1) Heuristic Checklist Result (see below)
     2) Low‑Risk Action Block (no prompt; see below)
   - Auto‑apply the change, run validations/tests, post a resolving reply to the thread.
   - Append the conversation ID to `tmp/CONV_ID.txt`.

5) If risk ≥ 5 (approval required)
   - MUST NOT apply or commit changes.
   - MUST print BOTH blocks in this order:
     1) Heuristic Checklist Result (see below)
     2) Medium/High‑Risk Approval Block (see below)
   - Post a non‑resolving reply summarizing the recommendation and await approval.
   - Append the conversation ID to `tmp/CONV_ID.txt`.

### Heuristic Checklist Output (MUST PRINT)
Print this block exactly once per conversation, before any action or decision output.

BEGIN_HEURISTIC
Heuristic Checklist Result
- Locate code: OK | ISSUE: …
- Read conversation intent: OK | ISSUE: …
- Already fixed in current code: Yes | No
- Standards alignment: OK | CONFLICTS: …
- Pros/cons analyzed: OK | NOTES: …
- Scope fit (this PR/story): In-scope | Out-of-scope
- Better/confirming solution: N/A | Brief summary
- Postpone criteria met: Yes | No (explain)
- Risk score (1–10): N — one‑sentence rationale
END_HEURISTIC

### Low‑Risk Action Block (risk < 5 — MUST PRINT, no prompt)
Print this action package after the checklist for every low‑risk conversation.
```txt
BEGIN_ACTION
Id: "{{thread_id}}"
Url: "{{thread_url}}"
Location: "{{file}}:{{line}}"
Summary: {{one‑line description of fix}}
Actions Taken: {{what changed at a high level}}
Tests/Checks: {{brief result}}
Resolution: Posted reply and resolved
END_ACTION
```

### Medium/High‑Risk Approval Block (risk ≥ 5 — MUST PRINT)
Print this approval package (with question) after the checklist for every medium/high‑risk conversation.
```txt
  Id: "{{thread_id}}"
  Url: "{{thread_url}}"
  Location: "{{file}}:{{line}}"
  Comment: {{comment_body}}
  Proposed Fix: {{proposed_fix}}
  Risk: "{{risk_score}}"

  Should I proceed with the {{proposed_fix}}?
  1. Yes
  2. No, do ... instead
```

### Output Gate
- Low risk (risk < 5) conversation is processed only if:
  - The Heuristic Checklist Result block was printed, and
  - The Low‑Risk Action Block was printed, and
  - The change was applied + validations run, and
  - A resolving reply was posted to the thread.
  - Then append the conversation ID to `tmp/CONV_ID.txt`.

- Medium/High risk (risk ≥ 5) conversation is processed only if:
  - The Heuristic Checklist Result block was printed, and
  - The Medium/High‑Risk Approval Block was printed, and
  - A non‑resolving reply was posted to the thread.
  - Then append the conversation ID to `tmp/CONV_ID.txt`.
