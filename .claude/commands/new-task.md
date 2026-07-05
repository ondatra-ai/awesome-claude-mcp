---
description: Start a new task — history rolls over on the next prompt
allowed-tools: Bash(rm:*)
---
!rm -f "${CLAUDE_PROJECT_DIR}/tmp/hooks-state/${CLAUDE_CODE_SESSION_ID}.json"
