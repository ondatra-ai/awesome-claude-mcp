#!/usr/bin/env bash
# Usage: sync.sh
# One-way mirror from this monorepo into the standalone
# ondatra-ai/true-bdd repo. Clones (or refreshes) true-bdd, rsyncs four
# trees, auto-patches the BDD runner's template path, commits, pushes,
# and opens or updates a sync PR. Prints the PR URL.
set -euo pipefail

if ! command -v claude >/dev/null 2>&1; then
  echo "claude CLI not found in PATH." >&2
  exit 1
fi

REPO=$(git rev-parse --show-toplevel)
WORK="$REPO/tmp/sync"
TARGET="$WORK/true-bdd"
TARGET_REMOTE="git@github.com:ondatra-ai/true-bdd.git"

mkdir -p "$WORK"

# 1. Clone or refresh true-bdd to a clean origin/main, with all
#    remote-tracking refs up to date. Force the full-tree fetch refspec
#    so prior single-branch / shallow checkouts still pick up every
#    remote branch (without this, sync-from-monorepo-* branches stay
#    invisible and the re-run guard misfires).
if [ -d "$TARGET/.git" ]; then
  git -C "$TARGET" config remote.origin.fetch \
    '+refs/heads/*:refs/remotes/origin/*'
  git -C "$TARGET" fetch --quiet --prune origin
  git -C "$TARGET" checkout --quiet main
  git -C "$TARGET" reset --quiet --hard origin/main
  git -C "$TARGET" clean -fdq
else
  rm -rf "$TARGET"
  git clone --quiet "$TARGET_REMOTE" "$TARGET"
fi

# 2. Branch off origin/main. If a sync for this monorepo SHA already
#    exists on the remote, reuse it instead of recreating (which would
#    produce a divergent commit and a push conflict on re-run).
SHORT_SHA=$(git -C "$REPO" rev-parse --short HEAD)
BRANCH="chore/sync-from-monorepo-${SHORT_SHA}"
if git -C "$TARGET" show-ref --verify --quiet "refs/remotes/origin/$BRANCH"; then
  git -C "$TARGET" checkout -B "$BRANCH" "origin/$BRANCH" >/dev/null 2>&1
  ALREADY_SYNCED=1
else
  git -C "$TARGET" checkout -B "$BRANCH" origin/main >/dev/null 2>&1
  ALREADY_SYNCED=0
fi

# 3. Mirror the four trees.
rsync -a --delete "$REPO/scripts/bdd-cli/src/" "$TARGET/src/"
rsync -a --delete "$REPO/scripts/bdd-cli/templates/" "$TARGET/templates/"
rsync -a --delete "$REPO/scripts/bdd-cli/tests/" "$TARGET/tests/"
mkdir -p "$TARGET/bdd-cli/checklists"
rsync -a --delete --exclude='*.tmp' \
  "$REPO/bdd-cli/checklists/" "$TARGET/bdd-cli/checklists/"

# 4. Patch tests/bdd/runner/runner.go: rewrite the one template path
#    that differs between the monorepo's nested layout and true-bdd's
#    flat layout.
RUNNER="$TARGET/tests/bdd/runner/runner.go"
if [ -f "$RUNNER" ]; then
  if sed --version >/dev/null 2>&1; then
    sed -i 's|"scripts/bdd-cli/templates"|"templates"|g' "$RUNNER"
  else
    sed -i '' 's|"scripts/bdd-cli/templates"|"templates"|g' "$RUNNER"
  fi
fi

# 5. Bail if nothing changed.
#    - On a fresh sync: clean status against origin/main → no monorepo
#      changes since the standalone repo's main.
#    - On a re-run (ALREADY_SYNCED=1): the working tree should be clean
#      because the existing sync branch already contains the diff. If
#      anything would change, that's a drift signal — abort loudly so
#      the operator can investigate.
if [ -z "$(git -C "$TARGET" status --porcelain)" ]; then
  if [ "$ALREADY_SYNCED" = "1" ]; then
    echo "Sync for monorepo @ ${SHORT_SHA} already pushed; skipping to PR step."
  else
    echo "No changes to sync. true-bdd main already matches monorepo @ ${SHORT_SHA}."
    exit 0
  fi
elif [ "$ALREADY_SYNCED" = "1" ]; then
  echo "ERROR: existing sync branch '${BRANCH}' diverges from a fresh mirror." >&2
  echo "       Either the monorepo or the standalone repo has been edited since the sync was pushed." >&2
  echo "       Resolve manually before re-running." >&2
  exit 1
fi

# 6. Stage, ask Claude to write the commit message + PR body, commit,
#    push — only on a fresh sync. A re-run skips straight to the PR
#    step and reuses the previously-generated message.
MONO_REMOTE=$(git -C "$REPO" remote get-url origin)
MONO_REPO_URL=$(echo "$MONO_REMOTE" \
  | sed -E 's|^git@github\.com:|https://github.com/|; s|\.git$||')

MSG_FILE="$WORK/sync-msg.txt"
TITLE_FILE="$WORK/sync-title.txt"
BODY_FILE="$WORK/sync-body.md"

if [ "$ALREADY_SYNCED" = "0" ]; then
  git -C "$TARGET" add -A

  PROMPT='Generate a commit message + PR body for a one-way SYNC commit that mirrors engine code from a monorepo into a standalone repo. This is a mechanical mirror, not a feature change — the noteworthy content is which engine changes from the source repo are now reaching the standalone consumers.

Format:
LINE 1: title (max 120 chars; "chore(sync): ..." prefix; name the most impactful engine change carried, not the file moves)
LINE 2: blank
LINE 3 onward: body. Start with a one-sentence framing of what landed (cite source-repo commits by their conventional-commit titles, not by SHA). Then a short bullet list of the noteworthy engine-level changes carried in this sync (skip mechanical renames, focus on new features / refactors / fixes). Then a final paragraph naming the source repo and SHA for traceability.

No markdown code fences, no Co-authored-by, no Generated-with trailers, no surrounding quotes, no explanation. Output only the message.'

  {
    echo "=== Sync mapping (mechanical) ==="
    echo "scripts/bdd-cli/src/        -> src/"
    echo "scripts/bdd-cli/templates/  -> templates/"
    echo "scripts/bdd-cli/tests/      -> tests/"
    echo "bdd-cli/checklists/*.yaml   -> bdd-cli/checklists/"
    echo "tests/bdd/runner/runner.go: repoLayer() template path rewritten for flat layout"
    echo ""
    echo "=== Source repo: ${MONO_REPO_URL}@${SHORT_SHA} ==="
    echo ""
    echo "=== Source-repo recent commits (last 15) ==="
    git -C "$REPO" log -15 --pretty=format:"%h %s"
    echo ""
    echo ""
    echo "=== Stat of changes this sync carries to the standalone repo ==="
    git -C "$TARGET" diff --cached --stat
  } | claude -p "$PROMPT" \
    | sed -e '/^```[a-zA-Z]*$/d' -e '/^```$/d' > "$MSG_FILE"

  if [ ! -s "$MSG_FILE" ]; then
    echo "Claude returned an empty sync message." >&2
    exit 1
  fi

  head -n 1 "$MSG_FILE" > "$TITLE_FILE"
  tail -n +3 "$MSG_FILE" > "$BODY_FILE"

  if [ ! -s "$TITLE_FILE" ] || [ ! -s "$BODY_FILE" ]; then
    echo "Parsed empty title or body from Claude output." >&2
    exit 1
  fi

  git -C "$TARGET" commit --quiet -F "$MSG_FILE"
  git -C "$TARGET" push --quiet -u origin "$BRANCH"
fi

# 7. Open or update the sync PR. On a re-run the message files won't
#    have been regenerated this turn — fall back to the existing PR
#    body if so (gh pr edit accepts a body-file; gh pr create requires
#    one).

#    On a fresh sync the title/body were generated above. On a re-run
#    they were not — read them from the existing commit on the branch.
if [ "$ALREADY_SYNCED" = "1" ]; then
  git -C "$TARGET" log -1 --pretty='%s' "$BRANCH" > "$TITLE_FILE"
  git -C "$TARGET" log -1 --pretty='%b' "$BRANCH" > "$BODY_FILE"
fi

PR_TITLE=$(cat "$TITLE_FILE")

(
  cd "$TARGET"
  # Only treat OPEN PRs as "reuse me" — closed/merged PRs that happen
  # to match the branch name are ignored (a fresh PR is created
  # instead, because closed PRs whose head SHA has changed can't be
  # reopened).
  OPEN_PR_NUMBER=$(gh pr list \
    --head "$BRANCH" --state open --json number -q '.[0].number // ""')

  if [ -n "$OPEN_PR_NUMBER" ]; then
    gh pr edit "$OPEN_PR_NUMBER" \
      --title "$PR_TITLE" --body-file "$BODY_FILE" >/dev/null
    gh pr view "$OPEN_PR_NUMBER" --json url -q .url
  else
    gh pr create \
      --base main --head "$BRANCH" \
      --title "$PR_TITLE" --body-file "$BODY_FILE" >/dev/null
    gh pr view "$BRANCH" --json url -q .url
  fi
)

rm -f "$MSG_FILE" "$TITLE_FILE" "$BODY_FILE"
