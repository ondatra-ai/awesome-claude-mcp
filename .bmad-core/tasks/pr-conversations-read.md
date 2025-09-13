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

## Inputs (optional)
- `autoResolveOutdated` (yes/no, default: no) — if yes, resolve all conversations where all comments are marked `outdated=true`.

## Sequential Task Execution

### 1) Detect PR Number
- Get current PR number for this branch:
  - `PR=$(gh pr view --json number -q .number)`
  - If empty, use `go run scripts/get-pr-number/main.go` to display info and HALT.

### 2) Fetch Conversations (JSON)
- Run: `go run ./scripts/list-pr-conversations/main.go "$PR" > ./tmp/PR_CONVERSATIONS.json`
- Validate JSON: `jq 'type=="array"' ./tmp/PR_CONVERSATIONS.json`

### 3) Classify Conversations
- Heuristics (intent-based):
  - If every comment in a thread has `outdated=true` → OUTDATED
  - If a comment references a file that no longer exists → OUTDATED
  - If the suggestion/issue intent still applies (pattern still present) → RELEVANT
  - Vague comments → default RELEVANT
- For this automated run, apply lightweight rules using the JSON only:
  - OUTDATED if all comments have `outdated==true`
  - Otherwise RELEVANT
- Note: Human review may upgrade/downgrade relevance later.

### 4) Generate Report
- Use template `.bmad-core/templates/pr-conversations-report-tmpl.md`.
- Produce `./tmp/PR_CONVERSATIONS.md` with sections:
  - OUTDATED (auto)
  - STILL RELEVANT (auto)
- Fill each item with: file:line (if present), conversation id, author, body, and a short one-line description.
- Command scaffold (example):
  ```bash
  jq -r '
    def firstFile: (.[0].file // "unknown");
    def firstLine: (.[0].line // 0);
    def desc(s): (s | gsub("\n"; " ") | .[0:120]);

    . as $all |
    "# All Conversations for PR #" + env.PR + ":\n\n" +
    "## ❌ OUTDATED (Fixed by previous changes):\n\n" +
    ( $all
      | map(select( all(.comments[]; .outdated==true) ))
      | map( "### **" + ( (.comments|first|.file // "unknown") + ":" + ((.comments|first|.line // 0)|tostring) ) + "**\n" +
              "Id: " + .id + "\n" +
              "Author: " + ((.comments|first|.author) // "unknown") + "\n" +
              "Description: " + (desc((.comments|first|.body) // "")) + "\n----\n" +
              ((.comments|first|.body) // "") + "\n----\n" +
              "Status: OUTDATED: All comments marked outdated.\n\n" )
      | join("") ) +
    "\n## ✅ STILL RELEVANT (Need to be fixed):\n\n" +
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

### 5) Optional: Auto-Resolve Outdated
- If `autoResolveOutdated == yes`:
  - Extract IDs: `IDS=$(jq -r 'map(select( all(.comments[]; .outdated==true) )) | .[].id' ./tmp/PR_CONVERSATIONS.json)`
  - For each ID: `go run ./scripts/resolve-pr-conversation/main.go "$id" "Auto-resolving: thread is outdated."`

### 6) Output
- Print a short summary:
  - Count of OUTDATED vs RELEVANT
  - Location of the report file `./tmp/PR_CONVERSATIONS.md`

## Checklist
- Execute checklist `.bmad-core/checklists/pr-conversations-checklist.md` and confirm PASS.

## Notes
- This task uses a conservative classifier based on the `outdated` flags; a human pass should refine relevance.
