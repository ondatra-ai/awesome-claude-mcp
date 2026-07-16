---
name: update-bdd-cli-readme
description: >
  Reviews staged changes and updates scripts/bdd-cli/README.md to keep it
  accurate — especially when CLI commands, flags, checklists, or the BDD
  fixture manifest change. Called automatically by pr-commit before
  committing. Can also be invoked manually via "update bdd-cli readme" or
  "sync the bdd-cli readme".
---

# Update bdd-cli README

You review the current staged changes (or recent diff) and update
`scripts/bdd-cli/README.md` to keep it accurate. This runs before every
commit so the README stays in sync with the tool it documents.

## Scope gate

Only act when the staged diff touches `scripts/bdd-cli/` or `bdd-cli/`
(the config/data directory at the repo root). If neither is touched — do
nothing, don't read further.

## What to check

Read the staged diff (`git diff --cached`) and compare against the README
sections:

### Commands (highest priority)

The `Status` table and `Usage` examples must match the real cobra command
tree in `scripts/bdd-cli/src/cmd/`:

- New, renamed, or removed subcommands → update the `Status` table, the
  `Usage` code block, and the Vision paragraph that names the subcommand
  suites
- New, renamed, or removed flags (`--fix`, `--requirements`,
  `--architecture`, …) → update the paragraph under the `Status` table
- Changed command semantics (what a command reads, writes, or refuses to
  touch) → update that command's row in the `Status` table
- A stub becoming working (or a command being gutted back to a stub) →
  update its `State` cell

### Checklists and configuration

- New or renamed checklist files in `bdd-cli/checklists/` → the
  `Configuration` section's naming-convention examples
- Changes to what `bdd-cli.yaml` or `architecture.yaml` configure → the
  `Configuration` section

### Testing harness

- Changes to the fixture manifest schema in
  `scripts/bdd-cli/tests/bdd/runner/fixture_manifest.go` (new keys like
  `prep`/`teardown`, changed assertion strategies) → the `Testing`
  section
- Changes to how the runner prepares the tmpdir, snapshots, or judges →
  the `Testing` section
- Changed test invocations (build tags, paths) → the `Testing` code block

### Install / prerequisites

- Go version bumps in `scripts/bdd-cli/go.mod`, new external binaries
  required → the `Install` section

## How to update

1. Read the staged diff
2. Read the current `scripts/bdd-cli/README.md`
3. For each section above, check if the diff makes it stale
4. If yes — edit the README with the minimal change needed
5. If no changes needed — do nothing, don't touch the file
6. Stage the README if it was modified
   (`git add scripts/bdd-cli/README.md`)

## Rules

- **Minimal changes only** — don't reorganize or rewrite sections that
  weren't affected by the diff
- **Leave the essay sections alone** — `Background`, `Vision`,
  `How it compares`, and `References` express direction, not code state;
  only touch them when a command they name by backtick changes
- **Both-repo links only** — the README is mirrored verbatim into the
  standalone `ondatra-ai/true-bdd` repo (see the `true-bdd-sync` skill).
  Relative links must resolve in both layouts, which means they may only
  target paths inside `scripts/bdd-cli/` (mirrored to the repo root) such
  as `templates/`. Never link into monorepo-only paths like `docs/` or
  `services/`; name them in prose instead
- **Don't add speculative content** — only document what's actually in
  the code now
- **Keep the same style** — match the existing formatting, tone, and
  level of detail
- **No commit messages or changelogs** — the README describes the current
  state, not the history
