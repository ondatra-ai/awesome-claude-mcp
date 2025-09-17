<!-- Powered by BMAD™ Core -->

# PR Merge Task

## Purpose
Merge the current branch's pull request (squash recommended), update the main branch, and clean up local/remote branches.

## Requirements
- GitHub CLI `gh` installed and authenticated (for automated merge)
- Current branch must have an open PR

## Sequential Task Execution

### 1) Identify Branches
- Current branch: `git branch --show-current`
- Main branch: default to `main`; if `master` exists locally, use `master`.

### 2) Merge PR (squash, delete remote branch)
- Run: `gh pr merge --squash --delete-branch`
- If `gh` is not available, perform manual steps in GitHub UI.

### 3) Update Main and Clean Up
- `git checkout <MAIN_BRANCH>`
- `git pull origin <MAIN_BRANCH>`
- If merge succeeded, delete local feature branch: `git branch -D <CURRENT_BRANCH>`

### 4) Confirm
- Echo success message and ensure working tree is clean.

## Notes
- Ensure CI checks have passed before merging.
- Use protected branch rules if configured.
