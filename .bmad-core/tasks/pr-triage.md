<!-- Powered by BMAD™ Core -->

# PR Triage (Auto-resolve ➜ Fix-or-Ticket)

## Purpose
Read ALL review conversations on the current PR, automatically resolve those that are fully outdated, then for the remainder: decide relevance to the work in progress, and either FIX immediately (if relevant) or CREATE a GitHub issue (if not relevant). No report file is produced.

## Requirements
- GitHub CLI `gh` authenticated for this repo
- Go toolchain available (for helper scripts)
- jq installed
- Helper scripts present:
  - `scripts/get-pr-number/main.go`
  - `scripts/list-pr-conversations/main.go`
  - `scripts/resolve-pr-conversation/main.go` (for optional resolve step)
- Local tests runnable (make test-unit, make test-e2e)
- Output policy: Do NOT display raw shell commands; show concise results and decisions only.

## Inputs
- None (non-interactive end-to-end flow)

## Sequential Task Execution

Preferred execution
<<<<<<< Updated upstream
- Use the wrapper script `scripts/pr-triage/run.sh` to perform triage with clean, standardized output. It internally calls `scripts/list-pr-conversations/main.go` and `scripts/resolve-pr-conversation/main.go`, applies auto‑resolve, selects the next actionable thread, and prints a decision package without showing any shell commands.
=======
- Step 1 — Fetch only: `scripts/pr-triage/fetch.sh` (caches all threads to `./tmp/PR_CONVERSATIONS.json`).
- Step 2 — Resolve outdated (optional): `scripts/pr-triage/resolve-outdated.sh` (resolves fully outdated threads then refreshes the cache).
- Step 3 — Process next item: `scripts/pr-triage/next.sh` (prints decision package for the next actionable thread with clean output).
- All scripts avoid printing raw shell commands; they show results only.

All‑in‑one (optional)
- Alternatively, `scripts/pr-triage/run.sh` performs all steps in one go and prints the decision package.

Output format (see template):
  - Thread: <id>
  - Link: <url>
  - Location: <file:line>
  - Comment: full review comment content
  - Proposed Fix: <concise action aligned with standards>
  - Risk Analysis: <short note>
  - Risk: <0–10>
  - Decision: <Proceed fix | Create ticket>

>>>>>>> Stashed changes
Template
- Reference: `.bmad-core/templates/pr-triage-output-tmpl.md` for the exact structure and labels used in output.

### 1) Detect PR Number
- Get current PR number for this branch:
  - `PR=$(gh pr view --json number -q .number)`
  - If empty, use `go run scripts/get-pr-number/main.go` to display info and HALT.

### 2) Fetch ALL Conversations (JSON)
- Use `scripts/list-pr-conversations/main.go` to fetch all review threads for the PR and produce JSON (both resolved and unresolved).
- Store JSON in a temp file (e.g., `./tmp/PR_CONVERSATIONS.json`) for analysis.

### 3) Auto-Resolve Outdated (MANDATORY)
- Identify threads where ALL comments are `outdated==true` and the thread is unresolved.
- Use `scripts/resolve-pr-conversation/main.go` to resolve those threads by ID with a standard note (e.g., “Auto-resolving: thread is fully outdated.”) before analysis.
- Re-fetch conversations with `scripts/list-pr-conversations/main.go` so resolution states are current for the next steps.

### 4) Load Architecture Context (Developer Analysis)
- Read these files to understand intended design and standards before making decisions:
  - `docs/architecture/tech-stack.md`
  - `docs/architecture/coding-standards.md`
  - `docs/architecture/source-tree.md`
  - `docs/architecture.md`
  - `docs/frontend-architecture.md`

### 5) Determine PR Scope (What is “in progress now”)
- Infer scope from changed files vs default branch (e.g., via `git diff --name-only origin/<default>...HEAD`).
- Build a list of paths and top-level areas (e.g., `services/frontend/`, `services/backend/`, `infrastructure/`, `docs/`).

### 6) Classify Remaining Conversations (Relevant vs Not Relevant Now)
- For each non-outdated thread from JSON:
  - If the first comment’s `file` exists and is in `CHANGED` (exact path) → RELEVANT
  - Else if directory of `file` matches the areas touched by `CHANGED` (same component) → RELEVANT
  - Else if it references `docs/` only → NOT RELEVANT NOW (track separately)
  - Else if the code referenced no longer exists → OUTDATED (reply and resolve)
  - Otherwise → NOT RELEVANT NOW
- Always use intent: if the issue describes a pattern the PR modifies, treat as RELEVANT even if code moved.

### 7) Human-In-The-Loop Approval (Process One-By-One)
- Process relevant threads sequentially, one at a time. Do not batch.
- For the current thread, present a decision package with enough context to act confidently:
  - Thread: short summary + link + file:line
  - Comment excerpt: 1–3 lines (trimmed)
  - Code context: brief diff or description of implicated code
  - Proposed fix: concrete steps (what exactly to change and where)
  - Architecture alignment: cite relevant points from coding-standards/tech-stack/source-tree
  - Risk/Effort: very short estimate (e.g., Low/Medium/High; ~N LOC)
  - Validations: tests to run (unit/integration/e2e) and any linters
  - Rollback: how to revert if needed
- Provide a Preferred option based on scope, risk, and effort (heuristic: if in-scope and low/medium risk → Prefer "Proceed fix"; otherwise → Prefer "Create ticket").
- Ask: "Default: <Preferred option>. Do you want to proceed with the default?"
- If the user declines, allow an explicit alternative:
  - 1) Proceed fix — implement the proposed fix now, then run validations, reply on thread, resolve
  - 2) Create ticket — open a GitHub issue with context and link; reply with the issue link; resolve
  - 3) Custom — provide instructions; follow them and then reply/resolve accordingly
- If no explicit approval is given, do not modify code or create issues; leave the thread pending and stop triage.

### 8) Act: Fix or Ticket (No report file)
- For RELEVANT items:
  1. Implement a minimal fix consistent with architecture and coding standards.
  2. Run validations: `make test-unit` and, when applicable, `make test-e2e`.
  3. Reply in the thread summarizing the fix and status of validations.
  4. Do NOT commit automatically; leave changes staged/unstaged for an explicit `@dev *pr-commit` later.
  5. Resolve the thread if the fix addresses the concern.
- For NOT RELEVANT NOW items:
  1. Create an issue with context and the PR thread URL:
     - `gh issue create --title "Follow-up from PR review: <short summary>" \
        --body "See: <thread URL>\n\nContext: <excerpt>\n\nSuggested action: <what to do>" \
        --label pr-review,tech-debt`
  2. Reply in the PR thread with the issue link and rationale.
  3. Resolve the thread to keep the PR focused.
- For newly detected OUTDATED items during analysis (file removed, already fixed):
  - Reply and resolve as OUTDATED.

### 9) Output
- Final state should be: outdated threads resolved; relevant ones fixed (pending review/commit); non-relevant converted to issues and resolved.

### 10) Summary (console only)
- Print counts to console:
  - Auto-resolved outdated: <n>
  - Fixed (relevant): <n>
  - Ticketed (not relevant now): <n>

## Checklist
- Execute checklist `.bmad-core/checklists/pr-triage-checklist.md` and confirm PASS.

## Notes
- This flow prioritizes keeping the PR focused: outdated items are resolved, relevant ones are fixed in-place (pending explicit commit), and non-relevant are ticketed for follow-up.
1) Read conversations to `tmp/CONV.json`.
   - The `main.go` script must be rewritten to write JSON directly to this file instead of standard output.

2) Create `tmp/CONV_ID.txt`.
   - This file stores processed conversation IDs (one per line). It is empty by default.

3) Identify the next conversation to process.
   - Write a Go script that reads `tmp/CONV.json` and `tmp/CONV_ID.txt`, finds the first conversation ID present in JSON that is not listed in `tmp/CONV_ID.txt`, and writes the complete conversation object to `tmp/CONV_CURRENT.json`. If none remain, write `{ "id": "No More Converations" }` to `tmp/CONV_CURRENT.json`.

4) Read `tmp/CONV_CURRENT.json`.
   - If it contains `{ "id": "No More Converations" }`, stop. Otherwise, proceed with heuristic analysis for that conversation.
