<!-- Powered by BMAD™ Core -->

# PR Conversations Read

## Purpose
Analyze review conversations on the current PR and generate a structured report of OUTDATED vs STILL RELEVANT items. Optionally, auto-resolve conversations that are clearly outdated.

## Requirements
- GitHub CLI `gh` authenticated for this repo
- Go toolchain available (for helper scripts)
- jq installed
- Helper scripts present:
  - `scripts/get-pr-number/main.go`
  - `scripts/list-pr-conversations/main.go`
  - `scripts/resolve-pr-conversation/main.go` (for optional resolve step)

## Inputs
- None

## Sequential Task Execution

### 1) Detect PR Number
- Get current PR number for this branch:
  - `PR=$(gh pr view --json number -q .number)`
  - If empty, use `go run scripts/get-pr-number/main.go` to display info and HALT.

### 2) Fetch ALL Conversations (JSON)
- Run: `go run ./scripts/list-pr-conversations/main.go "$PR" > ./tmp/PR_CONVERSATIONS.json`
- Validate JSON: `jq 'type=="array"' ./tmp/PR_CONVERSATIONS.json`
- Note: This list includes both resolved and unresolved threads.

### 3) Auto-Resolve Outdated (MANDATORY)
- Identify threads where ALL comments are `outdated==true`.
- Resolve each of those threads before analysis:
  - `IDS=$(jq -r 'map(select( all(.comments[]; .outdated==true) and (.isResolved==false) )) | .[].id' ./tmp/PR_CONVERSATIONS.json)`
  - For each: `go run ./scripts/resolve-pr-conversation/main.go "$id" "Auto-resolving: thread is fully outdated."`
- Optional (recommended): Re-fetch the conversations JSON to capture updated resolution states:
  - `go run ./scripts/list-pr-conversations/main.go "$PR" > ./tmp/PR_CONVERSATIONS.json`

### 4) Classify Conversations (Comprehensive Heuristics)
- Heuristics (intent-based):
  - If a comment references a file path that no longer exists in HEAD → OUTDATED
    - `FILE=$(jq -r '.comments[0].file // empty' <<<"$node"); [ -n "$FILE" ] && git ls-files --error-unmatch -- "$FILE" >/dev/null 2>&1 || mark OUTDATED`
  - Else if every comment is marked outdated → OUTDATED (already auto-resolved)
  - Else if path is under `docs/` → RELEVANT (Docs) with recommendation: defer or track separately
  - Else → RELEVANT
- This pass operates on the refreshed JSON from step 3.

### 5) Generate Report
- Use template `.bmad-core/templates/pr-conversations-report-tmpl.md` (structure mirrored below).
- Produce `./tmp/PR_CONVERSATIONS.md` with sections:
  - Auto-Resolved Outdated (fixed pre-analysis)
  - Still Relevant After Auto-Resolve (needs attention)
- Fill each item with: file:line (if present), conversation id, author, body, and a short one-line description.
- Command scaffold (example):
  ```bash
  export PR
  jq -r '
    def firstFile: (.[0].file // "unknown");
    def firstLine: (.[0].line // 0);
    def desc(s): (s | gsub("\n"; " ") | .[0:120]);

    . as $allRaw |
    # Identify auto-resolved (all comments outdated prior to analysis)
    ( $allRaw | map(select( all(.comments[]; .outdated==true) and (.isResolved==false) )) ) as $auto |
    # After auto-resolve and optional re-fetch
    ( $allRaw ) as $all |
    "# All Conversations for PR #" + env.PR + ":\n\n" +
    "- Auto-resolved (outdated): " + ($auto | length | tostring) + "\n" +
    "- Remaining (relevant candidates): " + (($all | map(select(any(.comments[]; .outdated!=true))) ) | length | tostring) + "\n\n" +

    "## ❌ Auto-Resolved Outdated:\n\n" +
    ( $auto
      | map( "### **" + ( (.comments|first|.file // "unknown") + ":" + ((.comments|first|.line // 0)|tostring) ) + "**\n" +
              "Id: " + .id + "\n" +
              "Author: " + ((.comments|first|.author) // "unknown") + "\n" +
              "Description: " + (desc((.comments|first|.body) // "")) + "\n----\n" +
              ((.comments|first|.body) // "") + "\n----\n" +
              "Status: OUTDATED: All comments were marked outdated and resolved automatically.\n\n" )
      | join("") ) +
    "\n## ✅ Still Relevant After Auto-Resolve:\n\n" +
    ( $all
      | map(select( any(.comments[]; .outdated!=true) ))
      | map( "### **" + ( (.comments|first|.file // "unknown") + ":" + ((.comments|first|.line // 0)|tostring) ) + "**\n" +
              "Id: " + .id + "\n" +
              "Author: " + ((.comments|first|.author) // "unknown") + "\n" +
              "Description: " + (desc((.comments|first|.body) // "")) + "\n----\n" +
              ((.comments|first|.body) // "") + "\n----\n" +
              "Status: RELEVANT: At least one comment not marked outdated.\n" +
              "Recommendation: Review intent and address accordingly.\n" +
              "Decision:\n\n" )
      | join("") )
  ' ./tmp/PR_CONVERSATIONS.json > ./tmp/PR_CONVERSATIONS.md
  ```

### 6) Output
- Print a short summary:
  - Count of OUTDATED vs RELEVANT
  - Location of the report file `./tmp/PR_CONVERSATIONS.md`

## Checklist
- Execute checklist `.bmad-core/checklists/pr-conversations-checklist.md` and confirm PASS.

## Notes
- This task uses a conservative classifier based on the `outdated` flags; a human pass should refine relevance.
