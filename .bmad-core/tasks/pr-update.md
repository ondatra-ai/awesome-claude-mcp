<!-- Powered by BMAD™ Core -->

# PR Update Task

## Purpose
Create or update a pull request with a clean, conventional format and related-issue linkage, then perform post-PR housekeeping.

## Requirements
- Git configured with upstream remote
- Optional: GitHub CLI `gh` authenticated (recommended)
- Temporary files written under `./tmp/`

## Inputs (elicit)
- Type: one of [feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert]
- Scope (optional): short scope like `auth`, `docs`, `mcp`, `frontend`
- Breaking Change (yes/no)
- Subject: concise, imperative, no trailing period (≤100 chars)
- PR Bullets: each change as a bullet line, no blank lines
- Optional: Related Issue Number (e.g., 37) — leave blank to auto-detect

## Sequential Task Execution

### 1) Analyze Changes
- Run: `git --no-pager status`
- Run: `git --no-pager diff --stat`
- Ensure working tree is committed (create commits before opening/updating PR).

### 2) Detect Related Issue (if not provided)
- Search repository for issue references in format `#<number>` or documented links.
- Prefer most recent explicit references in docs or story files.

### 3) Build Title and PR Body
- Compose PR Title using Conventional Commits:
  - Format: `type(scope)!: subject`
  - Include `(scope)` if provided; include `!` if breaking
  - Validate: `type` is allowed, `subject` has no trailing period and ≤100 chars
- Create `./tmp/pr-body.md` with content:
  - Bullet list of changes (no blank lines between them)
  - If related issue is known: add line `**Related Issue**: #<number>`
  - One paragraph of rationale after the bullets.

### 4) Create or Update PR (via GitHub CLI if available)
- If `gh` is available:
  - Try `gh pr view --json number` to detect existing PR for current branch.
  - If exists: `gh pr edit --title "<PR Title>" --body-file ./tmp/pr-body.md`
  - Else create: `gh pr create --fill=false --title "<PR Title>" --body-file ./tmp/pr-body.md`
- If `gh` is not available:
  - Push branch and open a browser link to create PR manually; include `./tmp/pr-body.md` contents.

### 5) Post-PR Actions (if issue exists and gh is available)
- Fetch issue title: `gh issue view <number> --json title -q .title`
- Ensure PR body contains `**Related Issue**: #<number>: <issue_title>`; update if needed.
- Ensure issue body references PR: append or update `**Related PR**: #<pr_number>: <PR Title>`.
- Apply labels (e.g., `has-pr`, `in-progress`) as appropriate.

### 6) Request Review (optional)
- Use `gh pr edit --add-reviewer <user>` as needed.

### 7) Cleanup
- Remove `./tmp/pr-body.md` after creation/edit completes.

## Notes
- Title must not be duplicated in the body; body begins directly with bullets.
- Use `--body-file` for reliability with special characters.
