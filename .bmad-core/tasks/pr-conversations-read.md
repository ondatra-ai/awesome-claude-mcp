<!-- Powered by BMAD™ Core -->

# PR Conversations Read (Auto-resolve ➜ Fix-or-Ticket)

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

## Inputs
- None (non-interactive end-to-end flow)

## Sequential Task Execution

### 1) Detect PR Number
- Get current PR number for this branch:
  - `PR=$(gh pr view --json number -q .number)`
  - If empty, use `go run scripts/get-pr-number/main.go` to display info and HALT.

### 2) Fetch ALL Conversations (JSON)
- Run: `go run ./scripts/list-pr-conversations/main.go "$PR" > ./tmp/PR_CONVERSATIONS.json`
- Validate JSON: `jq 'type=="array"' ./tmp/PR_CONVERSATIONS.json`
- Note: This list includes both resolved and unresolved threads.

### 3) Auto-Resolve Outdated (MANDATORY)
- Identify threads where ALL comments are `outdated==true`.
- Resolve each of those threads before analysis:
  - `IDS=$(jq -r 'map(select( all(.comments[]; .outdated==true) and (.isResolved==false) )) | .[].id' ./tmp/PR_CONVERSATIONS.json)`
  - For each: `go run ./scripts/resolve-pr-conversation/main.go "$id" "Auto-resolving: thread is fully outdated."`
- Optional (recommended): Re-fetch the conversations JSON to capture updated resolution states:
  - `go run ./scripts/list-pr-conversations/main.go "$PR" > ./tmp/PR_CONVERSATIONS.json`

### 4) Load Architecture Context (Developer Analysis)
- Read these files to understand intended design and standards before making decisions:
  - `docs/architecture/tech-stack.md`
  - `docs/architecture/coding-standards.md`
  - `docs/architecture/source-tree.md`
  - `docs/architecture.md`
  - `docs/frontend-architecture.md`

### 5) Determine PR Scope (What is “in progress now”)
- Infer scope from changed files vs default branch:
  - `BASE=$(gh repo view --json defaultBranchRef -q .defaultBranchRef.name)`
  - `CHANGED=$(git diff --name-only origin/$BASE...HEAD)`
- Build a list of paths and top-level areas (e.g., `services/frontend/`, `services/backend/`, `infrastructure/`, `docs/`).

### 6) Classify Remaining Conversations (Relevant vs Not Relevant Now)
- For each non-outdated thread from JSON:
  - If the first comment’s `file` exists and is in `CHANGED` (exact path) → RELEVANT
  - Else if directory of `file` matches the areas touched by `CHANGED` (same component) → RELEVANT
  - Else if it references `docs/` only → NOT RELEVANT NOW (track separately)
  - Else if the code referenced no longer exists → OUTDATED (reply and resolve)
  - Otherwise → NOT RELEVANT NOW
- Always use intent: if the issue describes a pattern the PR modifies, treat as RELEVANT even if code moved.

### 7) Human-In-The-Loop Approval (Per Relevant Thread)
- For each RELEVANT thread, pause and request approval before acting. Present a concise fix plan grounded in architecture docs and current changes:

  Proposal format (example to send to human):
  - Thread: <short summary> (<file:line>)
  - Proposed fix: <1–2 lines describing what and where>
  - Tests: <unit/integration/e2e to run>
  - Rationale: aligns with <coding-standards|tech-stack|source-tree>

  Options (reply with a number):
  1) Proceed: implement the proposed fix now (then run tests, reply, resolve)
  2) Create ticket: defer fix, file a GitHub issue and link it (reply, resolve)
  3) Not relevant now: explain briefly, resolve
  4) Defer: skip this thread for now (leave unresolved)
  5) Custom: provide instructions (I will follow them)

  - If no explicit approval is given, do not modify code or create issues; move to next thread.

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
- Execute checklist `.bmad-core/checklists/pr-conversations-checklist.md` and confirm PASS.

## Notes
- This flow prioritizes keeping the PR focused: outdated items are resolved, relevant ones are fixed in-place (pending explicit commit), and non-relevant are ticketed for follow-up.
