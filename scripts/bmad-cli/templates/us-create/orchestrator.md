<!-- Powered by BMAD™ Core -->

# Story Creation Orchestrator

## Purpose

This orchestrator coordinates the step-by-step creation of user stories from epic definitions, transforming them into comprehensive, actionable story files ready for developer implementation. It serves as the main entry point for the `bmad-cli sm us-create` command.

## Usage

```bash
bmad-cli sm us-create <epic.story> [options]
```

Example:
```bash
bmad-cli sm us-create 3.1
```

## Orchestration Flow

The orchestrator executes the following steps sequentially, with each step building context for the next:

### Step 0: Load Configuration
**File:** `step-0-load-config.md`
**Purpose:** Load and validate core configuration, project structure, and workflow settings
**Context Output:** Configuration object with paths and settings

### Step 1: Identify Story
**File:** `step-1-identify-story.md`
**Purpose:** Locate the target epic file and extract the specific story definition
**Context Input:** Configuration object
**Context Output:** Story metadata, epic context, story requirements

### Step 2: Gather Requirements
**File:** `step-2-gather-requirements.md`
**Purpose:** Extract detailed story requirements and review previous story context
**Context Input:** Story metadata, epic context
**Context Output:** Enriched story requirements, previous story insights

### Step 3: Gather Architecture
**File:** `step-3-gather-architecture.md`
**Purpose:** Read relevant architecture documents based on story type and extract technical context
**Context Input:** Story requirements, configuration
**Context Output:** Technical specifications, architecture constraints, implementation details

### Step 4: Verify Structure
**File:** `step-4-verify-structure.md`
**Purpose:** Cross-reference story requirements with project structure and identify conflicts
**Context Input:** Story requirements, technical specifications
**Context Output:** Structure validation, conflict notes

### Step 5: Populate Template
**File:** `step-5-populate-template.md`
**Purpose:** Create the complete story file using the story template with all gathered context
**Context Input:** All previous context
**Context Output:** Complete story file

### Step 6: Review and Complete
**File:** `step-6-review-completion.md`
**Purpose:** Final validation, checklist execution, and completion summary
**Context Input:** Complete story file
**Context Output:** Validation results, completion summary

## Context Passing

Each step receives context from previous steps and adds its own context for subsequent steps. Context is structured as YAML data blocks that can be easily parsed and extended.

Example context structure:
```yaml
context:
  config:
    devStoryLocation: "scripts/bmad-cli/templates"
    prdSharded: true
    # ... other config
  story:
    epicNum: 3
    storyNum: 1
    title: "MCP Server Implementation"
    # ... story details
  technical:
    dataModels: []
    apiSpecs: []
    # ... technical details
```

## Error Handling

If any step fails:
1. Log the error with step context
2. Provide recovery suggestions
3. Allow user to restart from failed step or previous step
4. Preserve partial context for debugging

## Integration with BMAD Core

This orchestrator integrates with the BMAD Core system by:
- Using the SM agent persona defined in `.bmad-core/agents/sm.md`
- Following the task structure from `.bmad-core/tasks/create-next-story.md`
- Leveraging templates from `.bmad-core/templates/story-tmpl.yaml`
- Executing checklists from `.bmad-core/checklists/story-draft-checklist.md`
