---
name: pr-commit
description: Run a full code quality validation pipeline (linting, unit tests, e2e tests, coverage check, pre-commit hooks) then commit and push changes. Use this skill whenever the user wants to commit code, is done with development work, says "commit this", "push my changes", "run checks and commit", or wants to validate and commit their work. This is more thorough than a simple git commit — it ensures production-ready quality before committing.
---

# PR Commit

Execute a complete code quality validation pipeline before committing changes. This ensures all code meets production standards before it reaches the remote repository.

## Steps

### 0. Ensure Working Branch

```bash
git rev-parse --abbrev-ref HEAD
```

If on `main`: create a new branch before proceeding. Derive the branch name from the staged changes (e.g., `feat/add-ebos-research`, `fix/update-skill`). Use kebab-case with a conventional prefix (`feat/`, `fix/`, `chore/`, `docs/`).

```bash
git checkout -b <branch-name>
```

If already on a feature branch: continue as-is.

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

### 5. Update Project Memory

Before staging and committing, invoke the `update-memory` skill to check if CLAUDE.md needs updating based on the changes. If it modifies CLAUDE.md, it will stage it automatically — the update will be included in this commit.

### 6. Review Changes

```bash
git --no-pager status
git --no-pager diff
```

Review what will be committed.

### 7. Stage and Commit

```bash
git add .
```

Create commit with a proper message (see format below).

### 8. Push

```bash
git push origin HEAD
```

Push immediately — do not ask for confirmation. If push fails, resolve immediately. Clean up any temp files in `./tmp/`.

### 9. Update PR

Invoke the `pr-update` skill to update the PR title and description to reflect all commits on the branch.

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
- Always update the PR description after pushing
- Never use `git commit --no-verify`
- Use `./tmp/` for any temporary files and clean them up afterwards
- If push fails, resolve immediately
