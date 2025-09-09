# pr-merge

Simple merge command to merge the current branch's pull request.

## Description

Merge the current branch's PR and clean up local branches.

## Implementation

```bash
# Get current branch
CURRENT_BRANCH=$(git branch --show-current)

# Execute merge with squash and delete remote branch
gh pr merge --squash --delete-branch

# Switch to main and clean up local branch
MAIN_BRANCH="main"
if git show-ref --verify --quiet refs/heads/master; then
    MAIN_BRANCH="master"
fi

git checkout $MAIN_BRANCH
git pull origin $MAIN_BRANCH

# Delete local branch if merge was successful
if [ $? -eq 0 ]; then
    git branch -D $CURRENT_BRANCH
    echo "✅ Merge completed successfully!"
    echo "✅ Local branch '$CURRENT_BRANCH' deleted"
else
    echo "❌ Failed to update main branch"
    exit 1
fi
```

## Requirements

- GitHub CLI (`gh`) installed and authenticated
- Current branch must have an open PR