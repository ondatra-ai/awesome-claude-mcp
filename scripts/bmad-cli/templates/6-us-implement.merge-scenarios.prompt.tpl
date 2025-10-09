<!-- Powered by BMAD™ Core -->

# Merge Test Scenario into Requirements Registry

## Your Task

You are merging a test scenario from user story {{.StoryNumber}} into `docs/requirements.yml`.

**CRITICAL**: Use fuzzy matching to understand scenario behavior and detect conflicts.

---

## Scenario to Merge

```yaml
id: "{{.ScenarioID}}"
level: "{{.Level}}"
priority: "{{.Priority}}"
acceptance_criteria: {{.AcceptanceCriteria}}
steps:
{{.Steps}}
```

---

## Merge Process

### Step 1: Read Current Requirements

```
Read {{.RequirementsFile}}
```

Understand the current state:
- What scenarios exist at this level ({{.Level}})?
- What is the highest ID number? (e.g., if INT-007 exists, next is INT-008)

---

### Step 2: Fuzzy Match Analysis

Analyze the scenario behavior using Given-When-Then logic:

**Context (Given)**: What state/preconditions are set?
**Trigger (When)**: What action occurs?
**Outcome (Then)**: What result is expected?

Find scenarios in {{.RequirementsFile}} that have:
- **Similar Context**: Same or similar Given conditions
- **Similar Trigger**: Same or similar When action
- **Similar/Conflicting Outcome**: Same or different Then expectations

**Matching Rules:**
- **EXACT MATCH** (95%+ similarity): Same context + trigger + outcome → UPDATE
- **CONFLICT** (same context/trigger, different outcome): → REPORT ERROR, ask user
- **NO MATCH** (< 80% similarity): → CREATE NEW

---

### Step 3: Merge Strategy Decision

Based on fuzzy match results:

#### Strategy A: CREATE NEW (No match found)
- Generate new flat ID: `{{.Level}}-XXX` (next sequential number)
- Example: If INT-007 exists → create INT-008
- Set `file_path: null` for new scenarios

#### Strategy B: UPDATE EXISTING (Exact match found)
- Preserve existing ID and `file_path`
- Update `merged_steps` with new steps
- Add to `user_stories` array
- Update `last_updated` to today's date (2025-10-09)

#### Strategy C: CONFLICT DETECTED (Same context/trigger, different outcome)
- **DO NOT MERGE**
- Report error with clear explanation
- Show both scenarios side-by-side
- Ask user to resolve manually

---

### Step 4: Generate Missing Fields

For NEW scenarios, generate:

1. **description**:
   - Analyze Given-When-Then steps
   - Create concise summary (< 80 chars)
   - Example: "WebSocket connection establishment succeeds"

2. **category**:
   - Extract from story context
   - Values: backend, frontend, performance, integration
   - For Story {{.StoryNumber}}: Use "backend" (MCP server story)

3. **requirement**:
   - **DO NOT include this field** (removed from format)

---

### Step 5: Build Merged Entry

For **NEW** scenarios:
```yaml
{{.Level}}-XXX:  # Next sequential ID
  description: "..." # Generated from steps
  category: "backend"  # From story context
  level: "{{.Level}}"
  priority: "{{.Priority}}"
  acceptance_criteria: {{.AcceptanceCriteria}}
  implementation_status:
    status: "pending"
    file_path: null
  last_updated: "2025-10-09"
  user_stories:
    - story_id: "{{.StoryNumber}}"
      story_file: "docs/stories/{{.StoryNumber}}-*.yaml"
      scenario_id: "{{.ScenarioID}}"
      merge_date: "2025-10-09"
  merged_steps:
    given: [...]
    when: [...]
    then: [...]
```

For **UPDATE** scenarios:
```yaml
{{.Level}}-XXX:  # Existing ID (preserve)
  description: "..." # Keep existing
  category: "..."  # Keep existing
  level: "{{.Level}}"
  priority: "{{.Priority}}"
  acceptance_criteria: {{.AcceptanceCriteria}}
  implementation_status:
    status: "pending"  # Reset to pending
    file_path: "..."  # PRESERVE existing path
  last_updated: "2025-10-09"  # Update date
  user_stories:
    - ... # Existing entries
    - story_id: "{{.StoryNumber}}"  # ADD new entry
      story_file: "docs/stories/{{.StoryNumber}}-*.yaml"
      scenario_id: "{{.ScenarioID}}"
      merge_date: "2025-10-09"
  merged_steps:  # Update with new steps
    given: [...]
    when: [...]
    then: [...]
```

---

### Step 6: Execute Merge

Use the Edit tool to update `{{.RequirementsFile}}`:

```
Edit {{.RequirementsFile}}
```

Add or update the scenario entry under the `scenarios:` section.

---

## Important Rules

1. **ID Format**: Flat IDs only (INT-008, E2E-011, etc.) - NO story prefix
2. **Sequential IDs**: Find highest existing number, increment by 1
3. **Preserve Paths**: For UPDATE, keep existing `file_path` value
4. **Reset Status**: Always set `implementation_status.status: "pending"` for updates
5. **No Requirement Field**: Do not include `requirement` field (removed from format)
6. **Today's Date**: Use 2025-10-09 for `last_updated` and `merge_date`

---

## Output Requirements

After completing the merge, provide a brief summary:

```
Merge Analysis:
- Strategy: [CREATE NEW | UPDATE | CONFLICT]
- Scenario ID: [Generated or Existing ID]
- Match Confidence: [Percentage if applicable]
- Actions Taken: [Description of edits made]
```

**If CONFLICT detected**, output:
```
CONFLICT DETECTED:
- Existing Scenario: [ID and description]
- New Scenario: {{.ScenarioID}}
- Conflict Reason: [Explain why they conflict]
- User Action Required: Resolve manually
```
