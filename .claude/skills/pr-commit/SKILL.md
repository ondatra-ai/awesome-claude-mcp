---
name: pr-commit
description: Run a full code quality validation pipeline (linting, unit tests, e2e tests, coverage check, pre-commit hooks) then commit and push changes. Use this skill whenever the user wants to commit code, is done with development work, says "commit this", "push my changes", "run checks and commit", or wants to validate and commit their work. This is more thorough than a simple git commit — it ensures production-ready quality before committing.
---

# PR Commit

Execute a complete code quality validation pipeline before committing changes. This ensures all code meets production standards before it reaches the remote repository.

## Validation Pipeline

Run these checks in order. If any step fails, fix the issue and re-run that step before proceeding. Never disable linting rules, skip tests, or use `--no-verify`.

### 1. Linting

Run linters and fix all issues in the actual code (not by suppressing warnings):

```bash
make lint-backend    # golint + go fmt
make lint-frontend   # ESLint --fix + Prettier --write
```

Re-run until clean.

### 2. Unit Tests

```bash
make test-unit
```

Fix failing tests by fixing the code, not by disabling tests. Re-run until all pass.

### 3. E2E Tests

```bash
make test-e2e
```

This runs E2E tests with Docker containers. Fix failures in the code. All changed files should have >80% test coverage.

### 4. Pre-commit Hooks

```bash
pre-commit run --all-files
```

If configured, run and fix any issues reported.

### 5. Review Changes

```bash
git --no-pager status
git --no-pager diff
```

Review what will be committed.

## Commit and Push

After all checks pass:

1. Stage changes: `git add .`
2. Create commit with a proper message (see format below)
3. Push immediately — do not ask for confirmation
4. Clean up any temp files in `./tmp/`

## Commit Message Format

- **Title**: One sentence summary, max 120 characters
- Empty line
- **Body**: Bullet list of changes (no blank lines between bullets)

**Example:**
```
Add user authentication to login page

- Add password validation function
- Create JWT token generation
- Add error handling for invalid credentials
```

Use a temp file for multi-line messages:
```bash
mkdir -p ./tmp
{
  echo "Your commit title"
  echo ""
  echo "- First change"
  echo "- Second change"
} > ./tmp/commit-msg.txt
git commit -F ./tmp/commit-msg.txt && rm ./tmp/commit-msg.txt && git push origin HEAD
```

## Rules

- Always push after committing — never leave commits unpushed
- Never use `git commit --no-verify`
- Use `./tmp/` for any temporary files and clean them up afterwards
- If push fails, resolve immediately
