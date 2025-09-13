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

2) Read conversations to `tmp/CONV.json`.
   - Run: `go run scripts/list-pr-conversations/main.go tmp/CONV.json [PR_NUMBER]`
   - MUST write the full conversations JSON array to `tmp/CONV.json`.
   - MUST NOT print JSON to standard output.

3) Create `tmp/CONV_ID.txt`.
   - Initialize file if missing: empty by default.
   - Clear it if it exists.
   - Purpose: store processed conversation IDs (one per line).

4) Identify the next conversation to process.
   - Run: `go run scripts/pr-triage/next-pr-conversation.go tmp/CONV.json tmp/CONV_ID.txt tmp/CONV_CURRENT.json`

5) Read `tmp/CONV_CURRENT.json`.
   - If it contains `{ "id": "No More Converations" }`, stop.
   - Otherwise, proceed with heuristic analysis for that conversation.
6) Append ID to `tmp/CONV_ID.txt`.
   - Append processed to `tmp/CONV_ID.txt`.
7) Repeat actions 4 → 6 until sentinel appears
   - Ensure any required thread reply was posted (per outcome in Heuristic analysis)
   - Verify the conversation ID was appended to `tmp/CONV_ID.txt` (no duplicates)
   - Run: `go run scripts/pr-triage/next-pr-conversation.go tmp/CONV.json tmp/CONV_ID.txt tmp/CONV_CURRENT.json`
   - If `tmp/CONV_CURRENT.json` contains `{ "id": "No More Converations" }`, stop
   - Otherwise, continue with item 5

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

4) If risk < 5, implement without human-in-the-loop
   - DO NOT COMMIT CHANGES
   - Conditions: strictly within PR scope, passes checklist alignment, and all tests pass locally.
   - Print BOTH blocks in this order:
     1) Heuristic Checklist Result (see below)
     2) Decision block (see below)
   - Apply the change, run tests, update the conversation with a concise summary of what changed and why.
   - Otherwise (risk ≥ 5):
     - Print BOTH blocks in this order:
       1) Heuristic Checklist Result (see below)
       2) Decision block (see below)
     - Do not apply/commit; wait for human approval.

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

### Decision Block (MUST PRINT)
Print this decision package after the checklist for every conversation.
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
- A conversation is considered processed only if:
  - The Heuristic Checklist Result block was printed, and
  - The Decision Block was printed, and
  - Any required thread reply was posted (when applicable).
- Only then append the conversation ID to `tmp/CONV_ID.txt`.
