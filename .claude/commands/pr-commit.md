# /pr-commit Command

**Purpose**: Execute a complete code quality validation pipeline before committing changes to ensure production-ready code standards.

**What it does**: Automatically runs linting, unit tests, e2e tests, coverage validation (>80% requirement), pre-commit hooks, and commits with proper formatting. All changes are immediately pushed to the remote repository without manual confirmation.

**Use when**: You have completed development work and are ready to commit changes that meet all quality standards and testing requirements.

## Analysis Process
1. Run golangci-lint and fix all issues:
   - `golangci-lint run --fix` (or `go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix`)
   - If any issues are found, fix the actual problems in the code - never disable linting rules, remove files, or suppress errors
   - Re-run until all issues are resolved

2. Run unit tests with coverage and fix any failures:
   - `make coverage-unit` (if Makefile exists) or `go test ./... -coverprofile=coverage.out`
   - If any tests fail, fix the actual issues in the code - never skip or disable tests
   - Re-run until all tests pass

3. Run e2e tests with coverage and fix any failures:
   - `make coverage-e2e` (if Makefile exists) or appropriate e2e test command
   - If any tests fail, fix the actual issues in the code - never skip or disable tests
   - Re-run until all tests pass
   - **REQUIREMENT: All changed files must have >80% test coverage**

4. Run pre-commit hooks and fix issues if any:
   - `pre-commit run --all-files`
   - If any hook fails, fix the reported issues and re-run until all hooks pass

5. Analyze changes:
   - `git --no-pager status` to see which files have changed
   - `git --no-pager diff` to see the actual changes in the code

## Commit Process
1. Review and stage your changes: `git add .`
2. Create proper commit message (format below)
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
Use one of these approaches:

### Method 1: Build commit message using a temporary file
```bash
mkdir -p ./tmp
{
  echo "Your commit title"
  echo ""
  echo "- First bullet point"
  echo "- Second bullet point"
} > ./tmp/commit-msg.txt

git commit -F ./tmp/commit-msg.txt && rm ./tmp/commit-msg.txt && git push origin HEAD
```

### Method 2: For simple commits
```bash
git commit -m "Your commit title" -m "- First bullet point" && git push
```

## Important Rules
- Always create temporary files in the `./tmp/` folder and clean them up
- **Never leave commits unpushed** - changes must be pushed to remote repository
- **Never ask for confirmation when pushing** - always push automatically after commit
- **Never use `git commit --no-verify`** - fix all pre-commit issues instead
- If push fails, resolve immediately before continuing with other tasks