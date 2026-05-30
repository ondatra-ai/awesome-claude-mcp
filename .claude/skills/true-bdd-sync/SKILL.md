---
name: true-bdd-sync
description: One-way mirror of bdd-cli from this monorepo into the standalone ondatra-ai/true-bdd repo, opening a sync PR. Use whenever the user wants to "sync true-bdd", "publish bdd-cli", "push bdd-cli to the standalone repo", or update the external true-bdd mirror after monorepo changes land on main.
---

# true-bdd Sync

One-shot mirror from the monorepo to the standalone `ondatra-ai/true-bdd` repo. The script clones (or refreshes) `true-bdd`, branches off its `main`, rsyncs four trees, auto-patches one path inside `tests/bdd/runner/runner.go`, asks `claude -p` to author the commit message + PR body (focused on the engine changes carried, not the mechanical file moves), commits, pushes, and opens (or updates) a sync PR. The script prints the PR URL on completion.

```bash
./.claude/skills/true-bdd-sync/sync.sh
```

## What gets mirrored

| Monorepo | true-bdd |
|---|---|
| `scripts/bdd-cli/src/**` | `src/**` |
| `scripts/bdd-cli/templates/**` | `templates/**` |
| `scripts/bdd-cli/tests/**` | `tests/**` |
| `bdd-cli/checklists/*.yaml` (excludes `*.tmp`) | `bdd-cli/checklists/*.yaml` |
| `scripts/bdd-cli/README.md` | `README.md` |

`rsync --delete` removes files inside the four mirrored trees that no longer exist in the monorepo. The README is mirrored as a single file (no `--delete`).

## What is not touched

- Project-specific config: `bdd-cli/{architecture,terms,acceptance-criteria-and-splitting,bdd-cli}.yaml`, all `*-schema.yaml`. These describe the MCP product, not the engine.
- `true-bdd`'s `.gitignore`, `LICENSE`, and anything outside the four mapped trees.

## Auto-patch

After mirroring, `sed` rewrites the one path inside `tests/bdd/runner/runner.go`'s `repoLayer()` from `"scripts/bdd-cli/templates"` to `"templates"` so the BDD harness resolves against `true-bdd`'s flat layout.

## Rules

- Working clone lives at `./tmp/sync/true-bdd/` and is reused across runs.
- If nothing changed since the last sync, the script exits 0 without committing or pushing.
- Branch selection: if any open sync PR (head ref starting with `chore/sync-from-monorepo`) already exists on the standalone repo, the new sync commit is appended to *that* PR's branch — one PR per review cycle, not one per monorepo SHA. If no open sync PR exists, the script branches off `main` as `chore/sync-from-monorepo` and opens a fresh PR.
- The script outputs the PR URL — report that back to the user.
