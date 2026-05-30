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

# 2. Pick the sync branch. If any open sync PR already exists on the
#    standalone repo, reuse its head ref so the new commit appends to
#    that PR (this is how follow-up syncs land on top of an in-review
#    PR instead of fragmenting into one PR per monorepo SHA). Otherwise
#    use a stable branch name (no SHA suffix — the SHA lives in the
#    commit message for traceability, not in the branch name).
SHORT_SHA=$(git -C "$REPO" rev-parse --short HEAD)
EXISTING_PR_BRANCH=$(
  cd "$TARGET" && gh pr list \
    --state open --json headRefName,headRefOid \
    --search 'head:chore/sync-from-monorepo' \
    -q '[.[] | select(.headRefName | startswith("chore/sync-from-monorepo"))][0].headRefName // ""'
)
if [ -n "$EXISTING_PR_BRANCH" ]; then
  BRANCH="$EXISTING_PR_BRANCH"
else
  BRANCH="chore/sync-from-monorepo"
fi

if git -C "$TARGET" show-ref --verify --quiet "refs/remotes/origin/$BRANCH"; then
  git -C "$TARGET" checkout -B "$BRANCH" "origin/$BRANCH" >/dev/null 2>&1
  ALREADY_SYNCED=1
else
  git -C "$TARGET" checkout -B "$BRANCH" origin/main >/dev/null 2>&1
  ALREADY_SYNCED=0
fi

# 3. Mirror the four trees and the README.
rsync -a --delete "$REPO/scripts/bdd-cli/src/" "$TARGET/src/"
rsync -a --delete "$REPO/scripts/bdd-cli/templates/" "$TARGET/templates/"
rsync -a --delete "$REPO/scripts/bdd-cli/tests/" "$TARGET/tests/"
mkdir -p "$TARGET/bdd-cli/checklists"
rsync -a --delete --exclude='*.tmp' \
  "$REPO/bdd-cli/checklists/" "$TARGET/bdd-cli/checklists/"
rsync -a "$REPO/scripts/bdd-cli/README.md" "$TARGET/README.md"

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

# 5. Decide whether to commit. Two outcomes after mirror:
#    - Clean working tree: the existing sync branch already matches the
#      current monorepo (or true-bdd main does, on a fresh run). No new
#      commit; on a re-run we still refresh PR metadata below.
#    - Dirty working tree: append a new sync commit. On a fresh run this
#      opens the PR; on a re-run this stacks the new monorepo changes
#      on top of the in-review PR (which is exactly the goal — one PR
#      per review cycle, not one PR per monorepo SHA).
if [ -z "$(git -C "$TARGET" status --porcelain)" ]; then
  if [ "$ALREADY_SYNCED" = "1" ]; then
    echo "Sync branch '${BRANCH}' already matches monorepo @ ${SHORT_SHA}; refreshing PR metadata."
  else
    echo "No changes to sync. true-bdd main already matches monorepo @ ${SHORT_SHA}."
    exit 0
  fi
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

if [ -n "$(git -C "$TARGET" status --porcelain)" ]; then
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

# 7. Open or update the sync PR. If section 6 ran a commit, the title
#    and body files are populated; otherwise (clean re-run) fall back
#    to the latest commit on the branch so `gh pr edit` still has
#    something to set.
if [ ! -s "$TITLE_FILE" ] || [ ! -s "$BODY_FILE" ]; then
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
