---
description: Start a new task — history rolls over on the next prompt
allowed-tools: Bash(python3:*)
---
!python3 "${CLAUDE_PROJECT_DIR}/.claude/hooks/history.py" new-task
