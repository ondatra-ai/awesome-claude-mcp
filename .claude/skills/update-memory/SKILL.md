---
name: update-memory
description: >
  Reviews staged changes and updates CLAUDE.md to reflect any architectural,
  structural, or conceptual changes. Called automatically by pr-commit
  before committing. Can also be invoked manually via "update memory",
  "sync claude.md", or "update project docs".
---

# Update Project Memory (CLAUDE.md)

You review the current staged changes (or recent diff) and update CLAUDE.md to keep it accurate. This runs before every commit to ensure the project documentation stays in sync with the codebase.

## What to check

Read the staged diff (`git diff --cached`) and look for changes that affect any section of CLAUDE.md:

### Repository structure
- New or renamed files/folders in `.claude/skills/`, `docs/`, `services/`, `scripts/`
- New skill folders → add to the structure tree if documented
- Deleted skills or folders → remove from documentation
- New file patterns → document the pattern

### Key concepts
- New terms, frameworks, or methodologies introduced in skill definitions or reference files
- Updated definitions or conventions
- New acronyms used in skills or docs

### Development setup
- New make targets or changed build commands
- New environment variables or dependencies
- Changed test commands or CI configuration

### Conventions
- New naming patterns or coding conventions
- New rules about how code should be structured
- New workflow conventions

## How to update

1. Read the staged diff
2. Read current CLAUDE.md
3. For each section, check if the diff introduces something that should be reflected
4. If yes — edit CLAUDE.md with the minimal change needed (don't rewrite sections that aren't affected)
5. If no changes needed — do nothing, don't touch the file
6. Stage CLAUDE.md if it was modified (`git add CLAUDE.md`)

## Rules

- **Minimal changes only** — don't reorganize or rewrite sections that weren't affected by the diff
- **Don't add speculative content** — only document what's actually in the codebase now
- **Don't remove content** unless the corresponding code/files were deleted
- **Keep the same style** — match the existing formatting, tone, and level of detail in CLAUDE.md
- **No commit messages or changelogs** — CLAUDE.md describes the current state, not the history
