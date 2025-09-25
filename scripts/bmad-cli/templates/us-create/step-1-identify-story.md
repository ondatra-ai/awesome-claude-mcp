<!-- Powered by BMAD™ Core -->

# Step 1: Identify Story for Preparation

## Purpose

Locate the target epic file and extract the specific story definition. This step validates story existence, checks for conflicts with existing stories, and prepares the story metadata for subsequent processing.

## Input Context

```yaml
context:
  config:
    devStoryLocation: "scripts/bmad-cli/templates"
    prdSharded: true
    prdLocation: "docs/epics/jsons"
    # ... other config

  target:
    epicNum: 3
    storyNum: 1
    epicStoryId: "3.1"
```

## Process

### 1. Locate Epic Files

Based on `prdSharded` configuration:

#### If PRD is Sharded (`prdSharded: true`)
- **Search Pattern:** `{prdLocation}/epic-{epicNum:02d}-*.yaml`
- **Example:** `docs/epics/jsons/epic-03-mcp-server.yaml`
- **Validation:** Ensure exactly one epic file matches the pattern

#### If PRD is Monolithic (`prdSharded: false`)
- **Search Location:** Use monolithic PRD file
- **Extract:** Epic section matching `epicNum`

### 2. Load and Parse Epic File

- **Load:** Target epic YAML file
- **Extract Epic Metadata:**
  - Epic ID
  - Epic name
  - Epic status
  - Epic goal
  - Epic context

- **Extract Stories Array:**
  - All stories defined in the epic
  - Validate story structure

### 3. Locate Target Story

- **Find:** Story with `id` matching `target.storyNum` (e.g., "3.1")
- **Validation:**
  - Story exists in epic
  - Story has required fields: `title`, `as_a`, `i_want`, `so_that`, `acceptance_criteria`

- **Error Handling:** If story not found:
  ```
  ERROR: Story 3.1 not found in epic 3
  Available stories: [list of available story IDs]
  ```

### 4. Check for Existing Story Files

- **Search Pattern:** `{devStoryLocation}/{epicNum}.{storyNum}.*`
- **Check for:** Any existing story files for this epic.story combination

#### If Existing Story Found
- **Load:** Existing story file
- **Check Status:** Extract current status
- **Status Validation:**
  - If status is 'Done': WARN user about overwriting completed story
  - If status is not 'Done': WARN about incomplete story
  - **User Prompt:** "ALERT: Found existing story! File: {filename} Status: {status} Would you like to: 1) Override and recreate 2) Cancel story creation 3) Continue from existing (if applicable)"

#### If No Existing Story
- **Proceed:** Continue with story creation

### 5. Validate Story Completeness

Check that the target story has all required elements:
- **Required Fields:**
  - `id`: Story identifier
  - `title`: Story title
  - `as_a`: User role
  - `i_want`: User goal
  - `so_that`: User benefit
  - `acceptance_criteria`: Array of criteria with `id` and `description`

- **Optional Fields:**
  - `status`: Current story status
  - `priority`: Story priority
  - `dependencies`: Story dependencies

## Output Context

```yaml
context:
  config:
    # ... previous config

  target:
    # ... previous target info

  epic:
    id: 3
    name: "MCP Server Setup"
    status: "PLANNED"
    goal: "Create functional MCP protocol server with tool registration and bidirectional communication"
    context: ""
    filePath: "docs/epics/jsons/epic-03-mcp-server.yaml"

  story:
    id: "3.1"
    title: "MCP Server Implementation"
    as_a: "Developer/Maintainer"
    i_want: "to implement MCP protocol server"
    so_that: "Claude can communicate with the service"
    status: "PLANNED"
    acceptance_criteria:
      - id: "AC-1"
        description: "WebSocket server implemented"
      - id: "AC-2"
        description: "HTTP endpoint for MCP available"
      # ... additional criteria

  existing:
    storyExists: false
    filePath: null
    currentStatus: null

  validation:
    epicFound: true
    storyFound: true
    storyComplete: true
    errors: []
    warnings: []
```

## Error Handling

### Epic File Not Found
- **Action:** HALT execution
- **Message:** "Epic file not found for epic {epicNum}. Expected pattern: {search_pattern}"

### Story Not Found in Epic
- **Action:** HALT execution
- **Message:** List available stories in epic

### Incomplete Story Definition
- **Action:** HALT execution
- **Message:** List missing required fields

### Existing Story Conflict
- **Action:** Prompt user for resolution
- **Options:** Override, Cancel, or Continue from existing

## Success Criteria

- Epic file located and loaded successfully
- Target story found and validated in epic
- Existing story conflicts resolved (if any)
- Story metadata prepared for requirements gathering

## Next Step

On success, pass enriched context to `step-2-gather-requirements.md`
