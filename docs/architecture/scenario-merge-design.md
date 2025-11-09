# Scenario Merge Design - Story to Requirements

**Date**: 2025-10-06
**Status**: Design Phase
**Author**: Discussion with @killev

## Overview

This document captures the design decisions for merging scenarios from user stories into `requirements.yml` during the `us implement` command.

---

## Workflow Summary

```
1. Story Created (e.g., 3.1) → Contains scenarios with story-centric IDs
2. Implementation Starts → Merge story scenarios into requirements.yml
   - New scenarios: Append with requirement-centric IDs
   - Existing scenarios: Amend/update (story takes priority)
3. Implement Tests → Write code for all scenarios
4. Implement Features → Make tests pass
```

---

## Key Design Decisions

### 1. Workflow Direction
**Decision**: Stories → Requirements
**Rationale**: Stories are created first, then scenarios are merged into requirements.yml during implementation.

### 2. ID Format Strategy
**Decision**: Requirement-centric IDs (`INT-00015-01`)
**Rationale**:
- Single scenario can be tied to multiple user stories
- Requirements need stable IDs that persist across story updates
- Story scenarios link back to requirement IDs via `requirement_id` field

**Structure**:
- **Story**: `3.1-INT-001` (story-centric, source)
- **Requirement**: `INT_00015_01` (requirement-centric, target)
- **Mapping**: Bidirectional via `requirement_id` and `story_scenario_id`

### 3. File Path in Stories
**Decision**: No `file_path` in story scenarios
**Rationale**:
- Stories describe **not implemented** scenarios
- Requirements describe **implemented** scenarios
- File paths are implementation details, determined during implementation

### 4. Merging Approach
**Decision**: Append or amend existing scenarios in requirements.yml
**Strategy**:
- New scenarios → Append with new FR-XXXXX
- Updated scenarios → Amend existing FR-XXXXX (story takes priority)
- Scenario changes → Reset `implementation_status` to `pending`

### 5. Requirements.yml Ownership
**Decision**: Auto-generated during user story implementation
**Process**:
1. First step: Merge scenarios from user story into requirements.yml
2. Second step: Implement tests (code all scenarios)
3. Third step: Implement features (make tests pass)

### 6. Conflict Resolution
**Decision**: Intellectual merge with story priority
**Rule**: The latest story has higher priority when conflicts occur
**Example**: If Story 3.3 updates a scenario from Story 3.1, Story 3.3 wins

---

## Data Structure Design

### User Story Scenarios (Source - Not Implemented)

**File**: `docs/stories/{story-id}.yaml`

```yaml
scenarios:
  test_scenarios:
    - id: "3.1-INT-001"              # Story-centric, sequential
      requirement_id: null            # Populated after merge → "FR-00015"
      description: "..."              # Human-readable summary
      service: "backend"             # backend/frontend/performance/integration
      level: "integration"            # unit/integration/e2e
      priority: "P0"                  # P0/P1/P2
      acceptance_criteria: ["AC-1"]   # Links to story ACs
      given: "..."                    # Gherkin
      when: "..."
      then: "..."
```

**Fields NOT included** (implementation details):
- ❌ `file_path` - Only in requirements
- ❌ `implementation_status` - Only in requirements

---

### Requirements Scenarios (Target - Implementation State)

**File**: `docs/requirements.yml`

```yaml
scenarios:
  INT-015:
    description: "WebSocket connection establishment"
    service: "backend"
    requirement: "Server accepts WebSocket connections"
    level: "integration"
    priority: "P0"
    acceptance_criteria: ["AC-1"]
    implementation_status:
      status: "implemented"  # pending/implemented/failed/outdated
      file_path: "tests/integration/mcp-server.test.ts"
    last_updated: "2025-10-08"
    user_stories:             # Track modifications across stories
      - story_id: "3.1"
        story_file: "docs/stories/3.1-mcp-server-implementation.yaml"
        scenario_id: "3.1-INT-001"
        merge_date: "2025-10-06"
      - story_id: "3.3"
        story_file: "docs/stories/3.3-connection-retry.yaml"
        scenario_id: "3.3-INT-002"
        merge_date: "2025-10-08"
    merged_steps:
      given:
        - "Server accepts WebSocket connections on port 8081"
      when:
        - "Client sends WebSocket connection request"
      then:
        - "Server establishes WebSocket connection"
        - and: "Server returns connection ID"
```

---

## Merge Scenarios

### Case 1: New Scenario (Not in requirements.yml)

**Input**: Story 3.1 has new scenario
```yaml
- id: "3.1-INT-001"
  description: "WebSocket connection establishment"
  service: "backend"
  level: "integration"
```

**Output**: Create new FR-XXXXX
```yaml
FR-00015:  # Next available FR number
  title: "WebSocket connection establishment"
  story_references:
    - story_id: "3.1"
      scenario_id: "3.1-INT-001"
      merged_date: "2025-10-06"
  scenarios:
    INT_00015_01:  # Next available scenario number
      story_scenario_id: "3.1-INT-001"
      implementation_status: "pending"
      # ... copy all fields from story
```

---

### Case 2: Update Existing Scenario

**Input**: Story 3.3 modifies scenario from Story 3.1
```yaml
# Story 3.3
- id: "3.3-INT-005"
  requirement_id: "FR-00015"  # Explicitly targets existing requirement
  description: "WebSocket connection with new security check"
```

**Output**: Update FR-00015
```yaml
FR-00015:
  story_references:
    - story_id: "3.1"
      scenario_id: "3.1-INT-001"
      merged_date: "2025-10-06"
    - story_id: "3.3"            # ADD new reference
      scenario_id: "3.3-INT-005"
      merged_date: "2025-10-08"
  scenarios:
    INT_00015_01:
      story_scenario_id: "3.3-INT-005"  # UPDATE to latest
      implementation_status: "pending"  # Reset (needs re-implementation)
      last_updated: "2025-10-08"
      updated_by_story: "3.3"
      change_history:            # Track changes
        - date: "2025-10-06"
          story_id: "3.1"
          scenario_id: "3.1-INT-001"
        - date: "2025-10-08"
          story_id: "3.3"
          scenario_id: "3.3-INT-005"
```

---

### Case 3: Multiple Stories → Same Requirement

**Input**: Story 3.1 and Story 3.5 both need WebSocket testing

**Output**: Multiple scenarios under one requirement
```yaml
FR-00015:
  title: "WebSocket connection establishment"
  story_references:
    - story_id: "3.1"
      scenario_id: "3.1-INT-001"
    - story_id: "3.5"
      scenario_id: "3.5-INT-003"
  scenarios:
    INT_00015_01:  # From Story 3.1
      story_scenario_id: "3.1-INT-001"
    INT_00015_02:  # From Story 3.5 (additional scenario)
      story_scenario_id: "3.5-INT-003"
```

---

## Conflict Resolution Rules

**When Story 3.3 updates a scenario from Story 3.1:**

1. **Description changed?** → Update in requirements, mark `implementation_status: "pending"`
2. **Steps changed?** → Update steps, reset status to `"pending"`
3. **Priority changed?** → Use latest story's priority
4. **Category changed?** → Use latest story's category
5. **File path exists?** → Keep file path (implementation detail persists)
6. **Implementation exists?** → Flag for review/re-implementation

**Merge Algorithm**:
```python
def merge_scenario(story_scenario, existing_requirement_scenario):
    if scenario_content_changed(story_scenario, existing_requirement_scenario):
        # Story takes priority
        existing_requirement_scenario.update_from(story_scenario)
        existing_requirement_scenario.implementation_status = "pending"
        existing_requirement_scenario.add_to_change_history()
    else:
        # Just add story reference
        existing_requirement_scenario.add_story_reference(story_scenario)
```

---

## Open Questions (Pending Decisions)

### A. Scenario Matching Strategy
**Question**: How to match Story Scenario → Requirement Scenario?

**Options**:
1. **Explicit `requirement_id`** - Story author specifies target FR-XXXXX
2. **Fuzzy matching** - Auto-match by description/steps similarity
3. **Interactive prompt** - CLI asks user to confirm matches

**Status**: ⏳ Pending decision

---

### B. Change History Granularity
**Question**: What level of detail to track in `change_history`?

**Options**:
1. **Lightweight** - Just story references
2. **Heavyweight** - Full diffs of changes
3. **Flags only** - `modified_by_story` field only

**Status**: ⏳ Pending decision

---

### C. Implementation Status States
**Question**: What states should `implementation_status` support?

**Proposed States**:
- `pending` - Not implemented
- `implemented` - Code exists and passes
- `failed` - Code exists but test fails
- `outdated` - Scenario updated, needs re-implementation
- `deprecated` - Scenario no longer needed

**Status**: ⏳ Pending decision

---

### D. Multiple Scenarios per Requirement
**Question**: Can one requirement (FR-XXXXX) have multiple test scenarios?

**Current Design**: Yes
```yaml
FR-00015:
  scenarios:
    INT_00015_01:  # Scenario 1
    INT_00015_02:  # Scenario 2
```

**Status**: ⏳ Confirm if correct

---

### E. File Path Strategy
**Question**: When merging new scenario, how to determine `file_path`?

**Options**:
1. **Convention-based** - `tests/{level}/{category}-{story}.test.ts`
2. **Prompt user** - Ask during merge
3. **Leave empty** - Fill during implementation

**Status**: ⏳ Pending decision

---

## Proposed Merge Command Flow

```bash
$ bmad-cli us implement 3.1

Step 1: Analyzing story scenarios...
  Found 10 scenarios in Story 3.1

Step 2: Matching scenarios to requirements.yml...
  - 3.1-INT-001: No match found → Will create FR-00015
  - 3.1-INT-002: Matched to FR-00008 (98% similar) → Will update
  - 3.1-E2E-001: No match found → Will create FR-00016

Step 3: Merging scenarios...
  ✅ Created FR-00015 (INT_00015_01) from 3.1-INT-001
  ✅ Updated FR-00008 (INT_00008_02) from 3.1-INT-002
     ⚠️  Implementation status reset to 'pending'
  ✅ Created FR-00016 (E2E_00016_01) from 3.1-E2E-001

Step 4: Updated requirements.yml
  - 2 new requirements created
  - 1 requirement updated
  - 8 scenarios marked as pending

Next: Run test generation for pending scenarios
```

---

## Next Steps

1. Finalize decisions on open questions (A-E)
2. Create exact YAML schema for both structures
3. Implement merge logic in `us implement` command
4. Add validation for scenario structure
5. Create migration script for existing scenarios

---

## Related Documents

- `docs/requirements.yml` - Current requirements structure
- `docs/stories/*.yaml` - User story format
- `docs/architecture/bdd-guidelines.md` - BDD scenario standards
- `docs/test-naming.md` - Test naming conventions

---

## User Answers Summary

| Question | Answer |
|----------|--------|
| Workflow direction? | Stories → Requirements |
| ID format? | Requirement-centric (`INT-00015-01`), but with story mapping |
| File path in stories? | No - implementation detail only |
| Merge approach? | Append or amend existing scenarios |
| Requirements ownership? | Auto-generated during implementation |
| Conflict resolution? | Intellectual merge, latest story wins |

---

**Note**: This is a living document. Update as design evolves and questions are resolved.

---

## Detailed Analysis & Recommendations

### Structural Comparison: Current State

#### User Story Structure (3.1-mcp-server-implementation.yaml)
```yaml
scenarios:
  test_scenarios:
    - id: "3.1-INT-001"
      acceptance_criteria: ["AC-1"]
      given: "..."
      when: "..."
      then: "..."
      level: "integration"
      priority: "P0"
```

#### Requirements Structure (requirements.yml)
```yaml
requirements:
  FR-00001:
    title: "..."
    story_id: "1.1-E2E-001"
    story_reference: "Story 1.1"
    service: "backend"
    scenarios:
      UT_00001_01:
        description: "..."
        file_path: "services/backend/cmd/main_test.go"
        steps:
          - given: "..."
          - when: "..."
          - then: "..."
```

---

### Field Mapping Analysis

#### Fields Present in Story but Missing in Requirements

| Field | Story | Requirements | Action |
|-------|-------|--------------|--------|
| `acceptance_criteria` | ✅ | ❌ | **ADD to requirements** - Critical for story context |
| `priority` | ✅ | ❌ | **ADD to requirements** - Helps with implementation ordering |
| `level` | ✅ | ❌ | **KEEP in both** - Test type classification |

#### Fields Present in Requirements but Missing in Story

| Field | Story | Requirements | Action |
|-------|-------|--------------|--------|
| `description` | ❌ | ✅ | **ADD to story** - Human-readable summary |
| `file_path` | ❌ | ✅ | **KEEP only in requirements** - Implementation detail |
| `title` | ❌ | ✅ | **ADD at requirement level** - Not needed per scenario |
| `story_reference` | ❌ | ✅ | **KEEP in requirements** - Backwards compatibility |
| `category` | ❌ | ✅ | **ADD to story** - Or use `level` only? |

---

### Recommended Field Structure

#### User Story Scenario (Complete)

```yaml
scenarios:
  test_scenarios:
    - id: "3.1-INT-001"                    # Story-centric ID
      requirement_id: null                 # Populated after merge
      description: "WebSocket connection establishment succeeds"  # NEW
      service: "backend"                  # NEW (or keep just level?)
      level: "integration"                 # KEEP
      priority: "P0"                       # KEEP
      acceptance_criteria: ["AC-1"]        # KEEP
      given: "Server is running and ready to accept WebSocket connections"
      when: "Client attempts to establish a WebSocket connection"
      then: "Server accepts connection and sends welcome message"
```

**Rationale for each field**:
- `id` - Story-centric for easy tracing during story work
- `requirement_id` - Links back to requirements after merge
- `description` - Human-readable summary (easier than reading Gherkin)
- `category` - Business domain classification
- `level` - Technical test type (unit/integration/e2e)
- `priority` - Implementation order guidance
- `acceptance_criteria` - Maps to story ACs for traceability
- `given/when/then` - BDD specification

---

#### Requirements Scenario (Complete)

```yaml
requirements:
  FR-00015:
    title: "MCP WebSocket server accepts connections"
    story_references:
      - story_id: "3.1"
        scenario_id: "3.1-INT-001"
        merged_date: "2025-10-06"
    service: "backend"
    priority: "P0"
    acceptance_criteria: ["AC-1"]        # NEW
    scenarios:
      INT_00015_01:
        story_scenario_id: "3.1-INT-001"  # Back-reference
        description: "WebSocket connection establishment succeeds"
        file_path: "tests/integration/mcp-server.test.ts"
        implementation_status: "pending"
        last_updated: "2025-10-06"
        updated_by_story: "3.1"
        change_history:
          - date: "2025-10-06"
            story_id: "3.1"
            scenario_id: "3.1-INT-001"
        steps:
          - given: "Server is running and ready to accept WebSocket connections"
          - when: "Client attempts to establish a WebSocket connection"
          - then: "Server accepts connection and sends welcome message"
```

**Rationale for each field**:
- `story_references` - Multiple stories can update same requirement
- `acceptance_criteria` - Inherited from latest story
- `story_scenario_id` - Direct back-reference to source
- `file_path` - Where test is implemented (implementation detail)
- `implementation_status` - Test implementation state
- `change_history` - Tracks who modified what and when
- `updated_by_story` - Quick reference to latest modifier

---

### Hierarchy Design

#### Current Hierarchy Mismatch

**Requirements**:
```
requirements → FR-XXXXX → scenarios → UT/IT/EE_XXXXX_XX
```

**Stories**:
```
scenarios → test_scenarios → [array of scenarios]
```

#### Recommended: Keep Separate Hierarchies

**Rationale**: Each serves a different purpose
- **Stories** - Flat list for quick scenario creation during story work
- **Requirements** - Grouped by functional requirement for implementation tracking

**Implementation**: Use mapping fields to link them together rather than forcing identical structures.

---

### Traceability & Linking

#### Current State
- ✅ Requirements → Story (via `story_id` and `story_reference`)
- ❌ Story → Requirements (no field)

#### Recommended: Bidirectional Links

**Story → Requirements**:
```yaml
# In story
- id: "3.1-INT-001"
  requirement_id: "FR-00015"  # Added during merge
```

**Requirements → Story**:
```yaml
# In requirements
FR-00015:
  story_references:
    - story_id: "3.1"
      scenario_id: "3.1-INT-001"
  scenarios:
    INT_00015_01:
      story_scenario_id: "3.1-INT-001"
```

**Benefit**: Full traceability in both directions

---

### ID Format Decision Tree

#### Scenario: Story 3.1 creates new test

**Story-Centric Approach**:
```
Story creates: 3.1-INT-001
Requirements generates: INT_00015_01 (derived)
```

**Pros**:
- Easy to trace during story work
- Natural grouping by story
- Simple sequential numbering

**Cons**:
- Requirements IDs are derived (potential conflicts)
- Hard to find requirement from story ID alone

#### Recommendation: Story-Centric with Explicit Mapping

```yaml
# Story (source)
- id: "3.1-INT-001"
  requirement_id: null  # Will be "FR-00015" after merge

# Requirements (target)
INT_00015_01:
  story_scenario_id: "3.1-INT-001"
```

**Why**: Story authors think in story terms, implementers think in requirement terms. Keep both happy with explicit mapping.

---

### Scenario Matching Strategies

#### Option 1: Explicit requirement_id (Recommended for Updates)

```yaml
# Story 3.3 explicitly updates existing requirement
- id: "3.3-INT-005"
  requirement_id: "FR-00015"  # Explicit target
  description: "WebSocket connection with security"
```

**When to use**: Story author knows they're updating existing functionality

**Pros**: No ambiguity, direct control
**Cons**: Author must know requirement ID

---

#### Option 2: Fuzzy Matching (Recommended for New Scenarios)

```python
def find_matching_requirement(story_scenario):
    matches = []

    # 1. Description similarity (Levenshtein distance)
    for req in requirements:
        similarity = calculate_similarity(
            story_scenario.description,
            req.scenarios.description
        )
        if similarity > 0.85:
            matches.append((req, similarity, "description"))

    # 2. Category + Level + Gherkin steps
    for req in requirements:
        if (req.category == story_scenario.category and
            req.level == story_scenario.level):
            step_similarity = calculate_step_similarity(
                story_scenario.steps,
                req.scenarios.steps
            )
            if step_similarity > 0.80:
                matches.append((req, step_similarity, "steps"))

    # 3. Acceptance criteria overlap
    for req in requirements:
        ac_overlap = len(set(story_scenario.acceptance_criteria) &
                        set(req.acceptance_criteria))
        if ac_overlap > 0:
            matches.append((req, ac_overlap/len(...), "ac"))

    return sorted(matches, key=lambda x: x[1], reverse=True)
```

**When to use**: New scenarios without explicit requirement_id

**Pros**: Automatic, discovers relationships
**Cons**: Can have false positives

---

#### Option 3: Interactive Prompt (Recommended for Ambiguous Cases)

```bash
$ bmad-cli us implement 3.1

Analyzing scenario 3.1-INT-005:
  Description: "WebSocket connection establishment with security"
  Category: backend
  Level: integration

Found potential matches in requirements.yml:

  [1] FR-00015 - "WebSocket connection establishment" (85% match)
      INT_00015_01 - tests/integration/mcp-server.test.ts
      Similarity: Description (85%), Steps (78%)

  [2] FR-00027 - "Secure WebSocket authentication" (72% match)
      INT_00027_01 - tests/integration/auth.test.ts
      Similarity: Description (72%), AC overlap (2/3)

  [3] Create new requirement (FR-00030)

Select option [1/2/3]: 1

✅ Updating FR-00015 with scenario from Story 3.1
   Status reset to 'pending' - re-implementation required
```

**When to use**: When fuzzy matching returns multiple candidates or low confidence

**Pros**: User control, prevents errors
**Cons**: Requires human interaction

---

#### Recommended Hybrid Approach

```python
def merge_scenario(story_scenario, requirements):
    # Step 1: Check for explicit requirement_id
    if story_scenario.requirement_id:
        return update_requirement(
            requirements[story_scenario.requirement_id],
            story_scenario
        )

    # Step 2: Fuzzy matching
    matches = find_matching_requirement(story_scenario)

    # Step 3: Decision based on confidence
    if not matches:
        # No match - create new
        return create_new_requirement(story_scenario)

    elif len(matches) == 1 and matches[0].confidence > 0.90:
        # Single high-confidence match - auto-update
        log_info(f"Auto-matched to {matches[0].id} (conf: {matches[0].confidence})")
        return update_requirement(matches[0], story_scenario)

    else:
        # Multiple matches or low confidence - prompt user
        return interactive_match(story_scenario, matches)
```

---

### Change History Granularity

#### Option 1: Lightweight (Recommended)

```yaml
change_history:
  - date: "2025-10-06"
    story_id: "3.1"
    scenario_id: "3.1-INT-001"
  - date: "2025-10-08"
    story_id: "3.3"
    scenario_id: "3.3-INT-005"
```

**Pros**: Simple, tracks "who changed what when"
**Cons**: No detail about what changed

**Use case**: Sufficient for most needs, can reconstruct from git history

---

#### Option 2: Heavyweight

```yaml
change_history:
  - date: "2025-10-06"
    story_id: "3.1"
    scenario_id: "3.1-INT-001"
    changes:
      description:
        old: null
        new: "WebSocket connection establishment"
      steps:
        old: []
        new: [{given: "...", when: "...", then: "..."}]
  - date: "2025-10-08"
    story_id: "3.3"
    scenario_id: "3.3-INT-005"
    changes:
      description:
        old: "WebSocket connection establishment"
        new: "WebSocket connection with security"
      steps:
        old: [{given: "...", when: "...", then: "..."}]
        new: [{given: "...", when: "...", then: "..."}]
```

**Pros**: Full audit trail, can see exact changes
**Cons**: Very verbose, requirements.yml becomes huge

**Use case**: Regulatory compliance, detailed auditing

---

#### Option 3: Flags Only (Minimal)

```yaml
updated_by_story: "3.3"
last_updated: "2025-10-08"
```

**Pros**: Minimal overhead
**Cons**: No history beyond latest change

**Use case**: Not recommended - loses traceability

---

#### Recommendation: Lightweight + Git Integration

```yaml
# In requirements.yml
change_history:
  - date: "2025-10-06"
    story_id: "3.1"
    scenario_id: "3.1-INT-001"
    commit: "a1b2c3d"  # Git commit SHA
  - date: "2025-10-08"
    story_id: "3.3"
    scenario_id: "3.3-INT-005"
    commit: "e4f5g6h"  # Git commit SHA
```

**Rationale**: Lightweight in requirements.yml, detailed diffs available via git

**Access pattern**:
```bash
# View what changed
git diff a1b2c3d e4f5g6h -- docs/requirements.yml
```

---

### Implementation Status States

#### Recommended State Machine

```
┌─────────┐
│ pending │ ─────────────────────┐
└─────────┘                      │
     │                           │
     │ tests implemented         │ scenario updated
     ▼                           │
┌────────────┐                   │
│implemented │ ──────────────────┤
└────────────┘                   │
     │                           │
     │ test starts failing       │
     ▼                           │
┌─────────┐                      │
│ failed  │ ──────────────────────┘
└─────────┘
     │
     │ marked as no longer needed
     ▼
┌────────────┐
│ deprecated │
└────────────┘
```

#### State Definitions

| State | Meaning | Next States |
|-------|---------|-------------|
| `pending` | Not yet implemented | → `implemented` (tests written)<br>→ `deprecated` (no longer needed) |
| `implemented` | Code exists and passes | → `failed` (test breaks)<br>→ `pending` (scenario updated)<br>→ `deprecated` |
| `failed` | Code exists but test fails | → `implemented` (bug fixed)<br>→ `pending` (scenario updated) |
| `deprecated` | No longer needed | (terminal state) |

#### Example Lifecycle

```yaml
# Day 1: Story 3.1 merged
INT_00015_01:
  implementation_status: "pending"

# Day 2: Tests written
INT_00015_01:
  implementation_status: "implemented"

# Day 5: Bug introduced, test fails
INT_00015_01:
  implementation_status: "failed"

# Day 6: Bug fixed
INT_00015_01:
  implementation_status: "implemented"

# Day 10: Story 3.3 updates scenario
INT_00015_01:
  implementation_status: "pending"  # Reset - needs re-implementation
```

---

### Multiple Scenarios per Requirement

#### Recommendation: Yes, Support Multiple

**Rationale**: One functional requirement can have multiple test scenarios covering different aspects

#### Example: FR-00015 with Multiple Scenarios

```yaml
FR-00015:
  title: "WebSocket connection management"
  scenarios:
    INT_00015_01:  # Happy path
      description: "Successful connection establishment"
      story_scenario_id: "3.1-INT-001"

    INT_00015_02:  # Error case
      description: "Connection rejected when server full"
      story_scenario_id: "3.1-INT-002"

    INT_00015_03:  # Concurrent case
      description: "Multiple simultaneous connections"
      story_scenario_id: "3.5-INT-003"
```

**When to split vs merge**:
- **Same requirement**: Different test scenarios for same functionality
- **Different requirements**: Different functionalities, even if similar

---

### File Path Strategy

#### Option 1: Convention-Based (Recommended)

```python
def generate_file_path(scenario):
    # Pattern: tests/{level}/{category}.test.{ext}
    level = scenario.level  # unit/integration/e2e
    category = scenario.category  # backend/frontend/performance

    # Determine extension
    if category == "backend":
        ext = "go"
    elif category == "frontend":
        ext = "ts"
    else:
        ext = "ts"  # default

    # Generate path
    return f"tests/{level}/{category}.test.{ext}"

# Examples:
# 3.1-INT-001 → tests/integration/backend.test.go
# 3.1-E2E-001 → tests/e2e/backend.test.ts
# 3.2-UNIT-001 → tests/unit/frontend.test.ts
```

**Pros**:
- Consistent structure
- No human input needed
- Easy to find tests

**Cons**:
- May not match existing file structure
- All tests for category in one file (could be huge)

---

#### Option 2: Prompt User (Interactive)

```bash
$ bmad-cli us implement 3.1

Merging scenario 3.1-INT-001...
  Description: "WebSocket connection establishment"
  Level: integration
  Category: backend

Suggested file path: tests/integration/backend.test.go
Enter file path [or press Enter to accept]: tests/integration/mcp-server.test.go

✅ Using tests/integration/mcp-server.test.go
```

**Pros**: Full user control
**Cons**: Requires interaction, inconsistent results

---

#### Option 3: Leave Empty (Defer to Implementation)

```yaml
INT_00015_01:
  file_path: null  # Set during test generation
```

**Pros**: Maximum flexibility
**Cons**: No guidance for implementers

---

#### Recommendation: Convention with Override

```yaml
# In bmad-cli.yml config
test_file_conventions:
  integration:
    backend: "tests/integration/{story}-mcp.test.go"
    frontend: "tests/integration/{story}-ui.test.ts"
  e2e:
    backend: "tests/e2e/{story}.spec.ts"
  unit:
    backend: "services/backend/**/*_test.go"  # Co-located
    frontend: "services/frontend/__tests__/**/*.test.tsx"

# Allow override per story
story_overrides:
  "3.1":
    integration: "tests/integration/mcp-server.test.go"
```

**Merge behavior**:
```python
def assign_file_path(scenario, config):
    # 1. Check story-level override
    if override := config.story_overrides.get(scenario.story_id):
        return override

    # 2. Use convention
    template = config.test_file_conventions[scenario.level][scenario.category]
    return template.format(story=scenario.story_id)
```

---

## Final Recommendations Summary

### For Implementation

1. **Scenario Matching**: Hybrid approach
   - Explicit `requirement_id` when updating
   - Fuzzy matching for new scenarios
   - Interactive prompt for ambiguous cases

2. **Change History**: Lightweight + Git
   - Store story references and commit SHAs
   - Detailed diffs available via git

3. **Implementation Status**: 4-state machine
   - `pending`, `implemented`, `failed`, `deprecated`

4. **Multiple Scenarios**: Yes, support multiple scenarios per requirement

5. **File Paths**: Convention-based with config overrides

### Schema Updates Needed

#### User Story Schema
```yaml
# ADD these fields:
- description: string (human-readable)
- service: string (backend/frontend/performance)
- requirement_id: string | null (populated after merge)
```

#### Requirements Schema
```yaml
# ADD these fields to requirement level:
- acceptance_criteria: string[] (from story)
- priority: string (P0/P1/P2)

# ADD these fields to scenario level:
- story_scenario_id: string (back-reference)
- change_history: array (track modifications)
- updated_by_story: string (latest modifier)
```

---

## Implementation Checklist

- [ ] Update user story YAML schema
- [ ] Update requirements YAML schema
- [ ] Implement fuzzy matching algorithm
- [ ] Implement interactive prompt for matching
- [ ] Implement merge logic with conflict resolution
- [ ] Add file path convention system
- [ ] Create migration script for existing scenarios
- [ ] Update BMAD CLI `us implement` command
- [ ] Add validation for both schemas
- [ ] Write tests for merge logic
- [ ] Document merge workflow for users
- [ ] Update BDD guidelines with new structure
