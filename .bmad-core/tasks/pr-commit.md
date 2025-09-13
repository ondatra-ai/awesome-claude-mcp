<!-- Powered by BMAD™ Core -->

# PR Commit Task

## Purpose
Execute a complete quality gate before committing and pushing changes: lint, unit tests, e2e tests, coverage check, pre-commit hooks, then create and push a well‑formatted commit.

## Prerequisites
- Active Git branch with staged or unstaged changes to commit
- Project dependencies installed (`make init` available and executed as needed)
- Docker running for E2E pipeline
- GitHub remote configured (origin) and authenticated SSH/HTTPS

## Workflow Rules
- Do not bypass failures (no --no-verify). Auto-fix issues and re-run until green.
- Do NOT attempt to install missing tools automatically. If any required tool is missing, HALT and ask the user to install it, then re-run.
- Maintain >80% coverage on changed files (enforce via local tooling if configured; otherwise verify reports).
- Push after successful commit without asking for confirmation.
- Create any temporary files under `./tmp/` and clean them up.
- Strict order: lint → unit → e2e → pre-commit → stage → compose message → commit → push. On any failure, attempt automated fixes and iterate until passing.
- Fix policy: For any failing step (lint/tests/hooks), attempt up to 5 fix iterations. First apply tool auto-fixes; if issues remain, modify SOURCE CODE to conform. Never modify configuration files (e.g., `.golangci.yml`, `.eslintrc.json`, `tsconfig.json`, `package.json`, CI configs). If still failing after 5 attempts, STOP and report remaining issues.
- E2E tests are authoritative: NEVER edit E2E test code to make tests pass. Fix the service/application code (frontend/backend/mcp/infra) to satisfy the tests.

## Inputs
- None (message is auto-generated; no prompts)

## Sequential Task Execution

### 0) Tooling Check (HALT if missing)
- Required commands must exist on PATH; if any are missing, STOP and install before proceeding:
  - `git`, `make`, `go`, `npm`, `docker`, `docker compose`, `pre-commit`
- Verify with `command -v <tool>`; if not found, instruct: install the tool and re-run this command.

### 1) Linting
- Run backend lint: `make lint-backend`
- Run frontend lint: `make lint-frontend`
- If issues are reported, iterate fixes up to 5 attempts:
  1. Apply auto-fixes and quick formatters
     - Backend: `gofmt -w`, `go vet ./...` (scoped to backend), optional `golangci-lint run --fix` if available
     - Frontend: `npx prettier --write` on changed files/directories
  2. Re-run linters; if still failing, MODIFY SOURCE CODE in implicated files to satisfy rules
  3. Re-run linters; repeat until clean or 5 attempts reached
  - Never edit linter configs or disable rules; conform code to current configuration

### 2) Tests
- Unit tests: `make test-unit`
- End-to-end tests: `make test-e2e`
- If failures occur, iterate fixes up to 5 attempts:
  - Inspect failing output; fix the SOURCE CODE to address defects (do not weaken test settings). Prefer minimal, targeted changes that satisfy test intent.
  - Re-run unit and e2e suites until green or attempts exhausted.
  - Do NOT change E2E tests; treat them as the contract. Only modify service/application code to meet test expectations.

### 3) Pre-commit Hooks
- Run: `pre-commit run --all-files`
- If `pre-commit` is not installed or no hooks configured, HALT and prompt to install/configure it, then re-run.
- If hooks report issues, iterate fixes up to 5 attempts:
  - Apply hook auto-fixes (formatting, import ordering, headers), then modify SOURCE CODE as needed
  - Never modify hook configuration; conform code to hooks
  - Re-run until clean or attempts exhausted

### 4) Review Changes
- Show status: `git --no-pager status`
- Optional: `git --no-pager diff` to inspect hunks.

### 5) Stage and Commit (Auto-generate Conventional Commit)
- Stage: `git add -A`
- Infer Conventional Commit fields:
  - TYPE (in order): from branch prefix `feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert`; else if only docs changes → `docs`; if test files touched → `test`; if infra/docker touched → `build`; if CI/.github/Makefile touched → `ci`; else `chore`.
  - SCOPE: `frontend` if any `services/frontend/` files; `backend` if any `services/backend/`; `mcp` if any `services/mcp-service/`; `infra` if any `infrastructure/`; `docs` if any `docs/`; if multiple areas, use `repo`.
  - SUBJECT: `update <scope> files` (e.g., `update frontend files`).
- Build commit message file at `./tmp/commit-msg.txt` using shell:
  - Title line: `TYPE(scope): SUBJECT` (scope omitted if `repo`).
  - Blank line, then a bullet per changed file with verb inferred from status (add/modify/delete/rename):

Example shell snippet to compose message:
```bash
mkdir -p ./tmp
COMMIT_FILE=./tmp/commit-msg.txt

# Determine branch and TYPE heuristic
BRANCH=$(git rev-parse --abbrev-ref HEAD)
case "$BRANCH" in
  feat/*) TYPE=feat;; fix/*) TYPE=fix;; docs/*) TYPE=docs;; style/*) TYPE=style;;
  refactor/*) TYPE=refactor;; perf/*) TYPE=perf;; test/*) TYPE=test;; build/*) TYPE=build;;
  ci/*) TYPE=ci;; chore/*) TYPE=chore;; revert/*) TYPE=revert;; *) TYPE="";;
esac

# Scope detection from staged changes
CHANGES=$(git status --porcelain)
SCOPE="repo"
echo "$CHANGES" | grep -E ' services/frontend/' >/dev/null && SCOPE="frontend"
echo "$CHANGES" | grep -E ' services/backend/' >/dev/null && SCOPE=${SCOPE%%|*}${SCOPE:+|}backend
echo "$CHANGES" | grep -E ' services/mcp-service/' >/dev/null && SCOPE=${SCOPE%%|*}${SCOPE:+|}mcp
echo "$CHANGES" | grep -E ' infrastructure/' >/dev/null && SCOPE=${SCOPE%%|*}${SCOPE:+|}infra
echo "$CHANGES" | grep -E ' docs/|\.md$' >/dev/null && SCOPE=${SCOPE%%|*}${SCOPE:+|}docs
echo "$SCOPE" | grep -q '|' && SCOPE="repo"

# TYPE fallback heuristics if not from branch
if [ -z "$TYPE" ]; then
  echo "$CHANGES" | grep -E ' docs/|\.md$' >/dev/null && TYPE=docs || true
fi
if [ -z "$TYPE" ]; then
  echo "$CHANGES" | grep -E '_test\.go|tests/|\.spec\.|\.test\.' >/dev/null && TYPE=test || true
fi
if [ -z "$TYPE" ]; then
  echo "$CHANGES" | grep -E ' infrastructure/|Dockerfile' >/dev/null && TYPE=build || true
fi
if [ -z "$TYPE" ]; then
  echo "$CHANGES" | grep -E ' \.github/|Makefile|scripts/ci/' >/dev/null && TYPE=ci || true
fi
[ -z "$TYPE" ] && TYPE=chore

# Compose title and bullets
SUBJECT="update ${SCOPE} files"
TITLE="$TYPE"
[ "$SCOPE" != "repo" ] && TITLE+="(${SCOPE})"
TITLE+=": ${SUBJECT}"

{
  echo "$TITLE"
  echo ""
  # One bullet per changed path
  echo "$CHANGES" | while read -r line; do
    [ -z "$line" ] && continue
    status=${line:0:2}
    path=${line:3}
    verb="update"
    case "$status" in
      A*) verb="add";;
      M*) verb="modify";;
      D*) verb="delete";;
      R*) verb="rename";;
    esac
    echo "- $verb $path"
  done
} > "$COMMIT_FILE"

git commit -F "$COMMIT_FILE"
```

### 6) Push
- Push current HEAD to remote tracking branch: `git push origin HEAD`
- If push fails, resolve immediately (auth, branch, rebase) and re-run push.

### 7) Cleanup
- Remove `./tmp/commit-msg.txt` if it exists.

## Quality Gate Checklist (PASS required)
- Tools present: git, make, go, npm, docker, docker compose, pre-commit
- Backend lint passes (format, vet, optional golangci)
- Frontend lint passes (ESLint --fix, Prettier)
- Unit tests pass (backend + frontend)
- E2E tests pass (Playwright via Docker Compose)
- Pre-commit hooks pass on all files
- Commit message follows Conventional Commits (type(scope): subject) and body bullet rules
- Commit pushed to origin HEAD

## Artifacts Produced
- Git commit on current branch with Conventional Commit title and bullet list
- Remote branch updated (origin/BRANCH)
- Temporary file removed: `./tmp/commit-msg.txt`

## Failure Handling & Reporting
- If a tool is missing: stop, report missing tool, ask user to install, then re-run
- If a step fails after 5 fix attempts: stop and summarize remaining issues with file paths, error excerpts, and suggested fixes
- Never edit configuration files or E2E tests to pass gates

## Notes
- If your repo uses conventional commits, ensure the Title conforms (e.g., feat:, fix:, chore:).
- If coverage gates are enforced in CI only, still verify local coverage reports before commit.
