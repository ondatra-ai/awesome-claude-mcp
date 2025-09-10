# /pr-update Command

Create or update pull requests following conventional commit guidelines with GitHub issue integration.

## Analysis Process
1. Run `git --no-pager status` to see which files have changed
2. Run `git --no-pager diff` to compare current branch with origin default branch
3. Analyze changes to understand purpose and impact
4. Check project files for GitHub issue references

## GitHub Issue Integration
1. **Issue Reference Detection**: Look for issue references in format:
   - `**Issue Reference**: [#<number>](https://github.com/owner/repo/issues/<number>)` (in project documentation)
   - Search project documentation and relevant files for issue references

2. **Issue Link in PR Body**: If an issue reference is found:
   - Check if issue link is already mentioned in PR body
   - If not mentioned, add "**Related Issue**: #<number>: <issue_title>" after bullet list
   - If already mentioned, update existing issue link with current title

3. **Issue Update After PR**: After creating/updating PR, update the related issue:
   - Check if PR link is already mentioned in issue body
   - If not mentioned, update issue body to include: "**Related PR**: #<pr_number>: <pr_title>"
   - If already mentioned, update existing PR link with current title
   - Use appropriate labels if PR addresses/closes the issue

## Pull Request Types
- `fix`: Patches a bug (correlates with PATCH in Semantic Versioning)
- `feat`: Introduces a new feature (correlates with MINOR in Semantic Versioning)
- `BREAKING CHANGE`: Breaking API change (correlates with MAJOR in Semantic Versioning)
- Other types: `build:`, `chore:`, `ci:`, `docs:`, `style:`, `refactor:`, `perf:`, `test:`

## Pull Request Format
- **Title**: One sentence summary (max 120 characters)
- **Body**:
  - Bullet list of changes (NO extra lines between bullet points)
  - **Related Issue**: #<number>: <issue_title> (if issue reference found)
  - Brief explanation of why this PR is needed
- Do NOT repeat the title in the body

## Example with Issue Reference:
```
feat: Add user authentication to login page

- Add password validation function
- Create JWT token generation
- Add error handling for invalid credentials

**Related Issue**: #37: GitHub Reader Step Implementation

This feature is necessary to secure user accounts and prevent unauthorized access.
```

## Example without Issue Reference:
```
feat: Add user authentication to login page

- Add password validation function
- Create JWT token generation
- Add error handling for invalid credentials

This feature is necessary to secure user accounts and prevent unauthorized access.
```

## PR Description Formatting
When creating a PR:
- Title is not repeated in body
- All bullet points written without extra lines between them
- Body starts directly with bullet points (no introductory text)
- Check if issue reference already exists in PR body, update existing or add new
- Use `--body-file` instead of `--body` with GitHub CLI to avoid escape issues
- Create PR body file in `./tmp/` folder (e.g., `./tmp/pr-body.md`)

## Process Steps
1. Analyze changes between origin default branch and current branch
2. Check project files for issue references
3. Create PR title and bullet list description
4. Check if issue reference already exists in PR body, update or add
5. Create new PR or update existing PR tied to current branch
6. **Post-PR Actions**: If issue reference found:
   - Check if PR link already mentioned in issue body
   - Update issue body to include PR link with title if not mentioned
   - Update existing PR link if already mentioned
   - Apply appropriate labels (e.g., "in-progress", "has-pr")
   - Update issue status if appropriate
7. Request review if needed
