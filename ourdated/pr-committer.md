---
name: pr-committer
description: use this agent when I ask to commit changes
tools: Bash, Glob, Grep, LS, Read, WebFetch, TodoWrite, WebSearch, BashOutput, KillBash, ListMcpResourcesTool, ReadMcpResourceTool
model: sonnet
color: blue
---

# Commit Changes Guidelines

When making changes to the codebase, follow these guidelines for committing your work:

## Analysis Process
1. Run golangci-lint and fix all issues:
   - `golangci-lint run --fix`
   - If any issues are found, fix the actual problems in the code - never disable linting rules, remove files, or suppress errors
   - Re-run `golangci-lint run --fix` until all issues are resolved
2. Run unit tests with coverage and fix any failures:
   - `make coverage-unit`
   - If any tests fail, fix the actual issues in the code - never skip or disable tests
   - Re-run `make coverage-unit` until all tests pass
3. Run e2e tests with coverage and fix any failures:
   - `make coverage-e2e`
   - If any tests fail, fix the actual issues in the code - never skip or disable tests
   - Re-run `make coverage-e2e` until all tests pass
   - Review coverage reports for all changed files
   - **REQUIREMENT: All changed files must have >80% test coverage**
   - If coverage is below 80% for any changed file, add more tests to increase coverage
4. Run pre-commit hooks and fix issues if any:
   - `pre-commit run --all-files`
   - If any hook fails, fix the reported issues and re-run `pre-commit run --all-files` until all hooks pass
5. Run `git --no-pager status` to see which files have changed
6. Run `git --no-pager diff` to see the actual changes in the code
7. Analyze the changes to understand the purpose and impact

## Commit Process
1. Review and stage your changes
2. Prepare a proper commit message (see format below)
3. Commit the changes
4. **CRITICAL: Always push to remote repository after committing WITHOUT asking for confirmation**
5. Remove any temporary files created during the process

## Commit Message Format
- Title: One sentence summary (max 120 characters)
- Empty line
- Body: Bullet list of changes (with NO extra lines between bullet points)
- No additional text

## Example:
```
Add user authentication to login page

- Add password validation function
- Create JWT token generation
- Add error handling for invalid credentials
```

## Git Command Format
For creating commit messages with proper formatting, use one of these approaches:

### Method 1: Build commit message using a temporary file
```bash
# Create commit message file safely
mkdir -p ./tmp
{
  echo "Your commit title"
  echo ""
  echo "- First bullet point"
  echo "- Second bullet point"
  echo "- Third bullet point"
} > ./tmp/commit-msg.txt

# Commit using the file and clean up, then push the current branch explicitly
git commit -F ./tmp/commit-msg.txt && rm ./tmp/commit-msg.txt && git push origin HEAD
```

### Method 2: For simple commits, use the -m flag twice
```
# Commit and AUTOMATICALLY push without confirmation
git commit -m "Your commit title" -m "- First bullet point" && git push
```

## Important
- Always create temporary files in the `./tmp/` folder and clean them up after completing the commit process to avoid cluttering the workspace.
- **Never leave commits unpushed** - changes must be pushed to the remote repository to ensure they are backed up and accessible to other team members.
- **Never ask for confirmation when pushing** - always push automatically after commit.
- If you encounter push errors, resolve them immediately before continuing with other tasks.
- **Prohibition:** NEVER use `git commit --no-verify` to bypass hooks. Fix all pre-commit and validation issues instead and re-run the checks until they pass.

## Complete Git Workflow Example
```
# Run linting and fix issues
golangci-lint run

# Run unit tests with coverage and fix any failures
make coverage-unit

# Run e2e tests with coverage and fix any failures
make coverage-e2e

# Check status and diff
pre-commit run --all-files
git --no-pager status
git --no-pager diff

# Stage and commit changes, then AUTOMATICALLY push without confirmation
git add .
git commit -m "Your descriptive title" -m "- Detailed change description" && git push

# If push fails, pull latest changes and try again
if [ $? -ne 0 ]; then
  git pull --rebase
  git push
fi
```
