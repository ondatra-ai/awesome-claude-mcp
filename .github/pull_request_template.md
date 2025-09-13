<!-- Enforces structure and conventional commit style for PRs -->

### PR Title (must follow Conventional Commits)

Example: `feat(auth)!: add OAuth2 device flow`

- Type: one of `feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert`
- Optional scope: `(scope)`
- Optional breaking change: `!`
- Subject: concise, imperative, no trailing period

### Summary (bulleted, no blank lines between bullets)

-
-
-

### Related Issue

<!-- Use GitHub keyword when applicable: closes #123 -->
**Related Issue**: #

### Why

Explain why this change is needed.

### Checklist

- [ ] Title follows Conventional Commits
- [ ] Lint passes (Go + Frontend)
- [ ] Unit tests pass
- [ ] E2E tests pass (if impacted)
- [ ] Docs updated (if applicable)
