<!-- Powered by BMAD™ Core -->

# pr-triage

Stage and process PR review conversations deterministically using local files.

## Purpose

Provide a predictable, file‑based workflow to read all PR conversations once, track processed items, and surface the next conversation for analysis.

## Artifcts

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
   - Cleat it if it exists.
   - Purpose: store processed conversation IDs (one per line).

4) Identify the next conversation to process.
   - Run: `go run scripts/pr-triage/next-pr-conversation.go tmp/CONV.json tmp/CONV_ID.txt tmp/CONV_CURRENT.json`

5) Read `tmp/CONV_CURRENT.json`.
   - If it contains `{ "id": "No More Converations" }`, stop.
   - Otherwise, proceed with heuristic analysis for that conversation.
6) Goto item 4.

### Default Behavior: Auto‑Iterate (no prompts)

- The triage run iterates through conversations automatically without pausing between items.
- Stops only when the selector writes the sentinel `{ "id": "No More Converations" }` or on an error.
- Per item behavior:
  - Low risk (<5): apply locally (no commit), post a concise reply on the thread, resolve the thread, append ID to `tmp/CONV_ID.txt`, and continue.
  - Medium/High (≥5): generate a YAML recommendation using the template, save to `tmp/pr-triage-<THREAD_ID>.yaml`, print it to stdout exactly, post a non‑resolving reply linking the recommendation or summarizing it, append ID to `tmp/CONV_ID.txt`, and continue.

Optional interactive mode (opt‑in)

- To pause for confirmation after each item, set an environment flag or use a wrapper option (suggested):

```bash
# Example wrapper semantics (optional):
PR_TRIAGE_INTERACTIVE=1 bash scripts/pr-triage/run.sh
# or
scripts/pr-triage/run.sh --interactive
```
In interactive mode, the workflow pauses after printing any recommendation (risk ≥ 5) or after a low‑risk change, awaiting user input before proceeding.

## Heuristic analysis

1) Perform Heuristic checklist
   - Use `.bmad-core/checklists/triage-heuristic-checklist.md` against the single conversation in `tmp/CONV_CURRENT.json`.
   - Ensure each item is evaluated; record any conflicts or notes.
   - Out Perform Heuristic checklist results

2) Determine the best option to proceed
   - Choose one: implement now, request changes/clarification, defer and create an issue, or escalate (architect/PO/QA).
   - Base the choice on alignment, scope fit, and value/effort.

3) Determine the risk score (1–10)
   - 1–3: low (cosmetic, doc, trivial refactor, isolated change)
   - 4–6: medium (localized behavior, limited blast radius)
   - 7–10: high (architecture/security/performance implications or wide impact)
   - Record the numeric score and rationale.

4) If risk < 5, implement without human-in-the-loop
   - Conditions: strictly within PR scope, passes checklist alignment, and all tests pass locally.
   - DO NOT COMMIT CHANGES
   - Apply the change, run tests, update the conversation with a concise summary of what changed and why.
   - Otherwise (risk ≥ 5), produce a YAML recommendation using `bmad-core/templates/pr-triage-output-tmpl.yaml` and both:
     - Save it to `tmp/pr-triage-<THREAD_ID>.yaml` (exact template fields, no extra keys), and
     - Print the YAML to stdout exactly as written (no wrappers, banners, or formatting beyond the YAML itself).

### Console Output Requirements (risk ≥ 5)

- Always print the recommendation YAML to the console as-is so reviewers can read it without opening files.
- The printed content MUST match the file contents byte-for-byte (aside from terminal line endings).
- Suggested flow:

```bash
# Prepare recommendation file from template fields
THREAD_ID=$(jq -r '.id' tmp/CONV_CURRENT.json)
OUT="tmp/pr-triage-$THREAD_ID.yaml"
# …generate YAML into "$OUT" using the template fields…

# Print to console exactly (no headers or footers)
cat "$OUT"
```

Note: Do not auto-apply or commit medium/high-risk items; wait for explicit human approval after printing the recommendation.
