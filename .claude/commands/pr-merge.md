# pr-merge

Safely merge the current branch's pull request using GitHub CLI with admin rights protection.

## Description

Merge the current branch's PR after validating merge requirements and ensuring proper access controls. Prevents admin bypass and enforces standard merge policies.

## Implementation

```bash
# Step 1: Validate environment and PR status
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed. Install with: brew install gh"
    exit 1
fi

if ! gh auth status &> /dev/null; then
    echo "❌ GitHub CLI is not authenticated. Run: gh auth login"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ]; then
    echo "❌ Cannot merge from main/master branch"
    exit 1
fi

# Check if PR exists
if ! gh pr view &> /dev/null; then
    echo "❌ No pull request found for branch '$CURRENT_BRANCH'"
    echo "Create a PR first with: gh pr create"
    exit 1
fi

# Step 2: Admin rights verification (CRITICAL)
USER_PERMISSION=$(gh api repos/:owner/:repo/collaborators/$(gh api user --jq '.login')/permission --jq '.permission')

if [ "$USER_PERMISSION" = "admin" ]; then
    echo "❌ CRITICAL: Admin rights detected"
    echo "This command explicitly forbids admin privilege usage"
    echo "Admin bypass of merge policies is not allowed"
    exit 1
fi

# Step 3: Check merge status
PR_STATUS=$(gh pr view --json mergeable,mergeStateStatus)
MERGEABLE=$(echo "$PR_STATUS" | jq -r '.mergeable')

if [ "$MERGEABLE" != "MERGEABLE" ]; then
    echo "❌ PR is not ready to merge"
    gh pr status
    exit 1
fi

# Step 4: Show status and merge
echo "✓ PR is ready to merge (admin bypass disabled)"
gh pr view
echo ""
echo "Proceeding with merge..."

# Execute merge (NO ADMIN FLAGS)
gh pr merge --squash --delete-branch

# Step 5: Local cleanup
echo "Updating local repository..."
MAIN_BRANCH="main"
if git show-ref --verify --quiet refs/heads/master; then
    MAIN_BRANCH="master"
fi

git checkout $MAIN_BRANCH
git pull origin $MAIN_BRANCH

echo "✅ Merge completed successfully!"
```

## Security

**CRITICAL**: This command explicitly forbids admin privilege usage:
- Detects admin permissions and exits with error
- Uses standard `gh pr merge` without admin flags
- Enforces all branch protection rules
- No bypass of required status checks

## Requirements

- GitHub CLI (`gh`) installed and authenticated
- Current branch must have an open, mergeable PR
- Standard merge permissions (admin rights blocked)
- All status checks passing and required approvals met