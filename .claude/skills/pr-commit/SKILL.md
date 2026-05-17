---
name: pr-commit
description: Run quality gates, update memory, commit and push, and update the PR. Use when the user says "commit", "push", "commit and push", or similar.
---

# PR Commit

1. Run `./.claude/skills/pr-commit/gates.sh` (lint + test + pre-commit).
2. Invoke the `update-memory` skill.
3. Run `./.claude/skills/pr-commit/commit.sh`.
4. Invoke the `pr-update` skill.
