<!-- Powered by BMAD™ Core -->

# Step 2: Gather Story Requirements and Previous Story Context

## Purpose

Extract detailed story requirements from the epic and review previous story context to understand implementation history, decisions, and lessons learned that inform the current story preparation.

## Input Context

```yaml
context:
  config:
    devStoryLocation: "scripts/bmad-cli/templates"
    # ... other config

  target:
    epicNum: 3
    storyNum: 1
    epicStoryId: "3.1"

  epic:
    id: 3
    name: "MCP Server Setup"
    # ... epic metadata

  story:
    id: "3.1"
    title: "MCP Server Implementation"
    as_a: "Developer/Maintainer"
    i_want: "to implement MCP protocol server"
    so_that: "Claude can communicate with the service"
    acceptance_criteria:
      # ... criteria array
```

## Process

### 1. Extract Detailed Story Requirements

#### 1.1 Primary Requirements
From the current story definition:
- **User Story Statement:** Combine `as_a`, `i_want`, `so_that` into complete user story
- **Acceptance Criteria:** Extract all AC with IDs and descriptions
- **Story Scope:** Determine if this is Backend, Frontend, Full-Stack, or Infrastructure
- **Dependencies:** Check if story references dependencies on other stories/epics

#### 1.2 Epic Context Requirements
From the parent epic:
- **Epic Goal:** Overall epic objective that this story contributes to
- **Epic Context:** Additional context that affects all stories in epic
- **Success Criteria:** Epic-level success criteria that story must support
- **Technical Notes:** Epic-level technical constraints or requirements

#### 1.3 Story Classification
Classify the story type for architecture context gathering:
- **Backend/API Story:** Server-side logic, APIs, data processing
- **Frontend/UI Story:** User interface, client-side logic
- **Full-Stack Story:** Both backend and frontend components
- **Infrastructure Story:** DevOps, deployment, configuration
- **Integration Story:** Third-party integrations, external services

### 2. Review Previous Story Context

#### 2.1 Locate Previous Stories
- **Search Pattern:** `{devStoryLocation}/{epicNum}.*.story.*`
- **Sort:** By story number (e.g., 3.1, 3.2, 3.3...)
- **Identify:** Most recent completed story in same epic

#### 2.2 Extract Previous Story Insights
If previous story exists, extract from `dev_agent_record` section:

##### Completion Notes
- Implementation deviations from original plan
- Technical decisions made during development
- Challenges encountered and how they were resolved
- Performance or quality considerations discovered

##### Debug Log References
- Links to detailed debug logs
- Error patterns that were resolved
- Testing issues and solutions

##### File List
- Files that were created or modified
- File structure decisions
- Code organization patterns established

##### Technical Decisions
- Library or framework choices made
- Architecture patterns implemented
- Database schema changes
- API design decisions

#### 2.3 Extract Relevant Insights
Filter previous story insights for relevance to current story:
- **Shared Components:** Components or modules that current story might use
- **Technical Patterns:** Established patterns that current story should follow
- **Known Issues:** Problems to avoid or solutions to reuse
- **Architecture Evolution:** How the system has evolved that affects current story

### 3. Consolidate Requirements

#### 3.1 Create Comprehensive Requirements
Combine story requirements with previous story insights:

```yaml
requirements:
  story:
    statement: "As a Developer/Maintainer, I want to implement MCP protocol server so that Claude can communicate with the service"
    scope: "Backend/Infrastructure"
    priority: "High"

  acceptance_criteria:
    - id: "AC-1"
      description: "WebSocket server implemented"
      notes: "Previous story established WebSocket patterns in user-auth module"
    # ... additional criteria with context notes

  dependencies:
    epic_dependencies: []
    story_dependencies: []
    technical_dependencies: ["OAuth Manager", "Cache Manager"]

  constraints:
    performance: ["< 1 second connection establishment", "< 100ms message processing"]
    security: ["Token validation required", "CORS handling"]
    compatibility: ["Go 1.21+", "Fiber 2.x framework"]
```

#### 3.2 Identify Implementation Patterns
Based on previous stories and current requirements:
- **Code Structure Patterns:** File organization, naming conventions
- **Integration Patterns:** How to connect with existing components
- **Testing Patterns:** Test structure and coverage expectations
- **Error Handling Patterns:** Consistent error handling approaches

## Output Context

```yaml
context:
  # ... previous context sections

  requirements:
    story:
      statement: "As a Developer/Maintainer, I want to implement MCP protocol server so that Claude can communicate with the service"
      scope: "Backend/Infrastructure"
      classification: "Backend"
      priority: "High"

    acceptance_criteria:
      - id: "AC-1"
        description: "WebSocket server implemented"
        context_notes: []
      - id: "AC-2"
        description: "HTTP endpoint for MCP available"
        context_notes: []
      # ... additional criteria

    dependencies:
      epic_level: []
      story_level: []
      technical: ["OAuth Manager", "Cache Manager", "Network layer"]

    constraints:
      performance: ["< 1 second connection establishment", "< 100ms message processing"]
      security: ["Token validation required", "CORS handling"]
      technical: ["Go 1.21", "Fiber 2.x", "gorilla/websocket"]

  previous_insights:
    story_completed: null  # or previous story ID if exists
    completion_notes: []
    technical_decisions: []
    file_patterns: []
    lessons_learned: []

  patterns:
    code_structure: []
    integration: []
    testing: []
    error_handling: []
```

## Error Handling

### Missing Acceptance Criteria
- **Action:** WARN user
- **Message:** "Story {id} missing acceptance criteria. This may result in incomplete implementation."

### Unclear Story Scope
- **Action:** Attempt classification, note uncertainty
- **Message:** Add warning to context about scope uncertainty

### Previous Story Analysis Failed
- **Action:** Continue without previous context
- **Message:** Log warning about missing previous story insights

## Success Criteria

- Story requirements extracted and classified
- Previous story context reviewed (if available)
- Implementation patterns identified
- Dependencies and constraints documented
- Requirements context prepared for architecture gathering

## Next Step

On success, pass enriched context to `step-3-gather-architecture.md`
