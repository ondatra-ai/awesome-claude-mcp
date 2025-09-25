<!-- Powered by BMAD™ Core -->

# Step 5: Populate Story Template with Full Context

## Purpose

Create the complete story file using the story template, populated with all context gathered from previous steps. This produces a comprehensive, self-contained story ready for developer implementation.

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
    # ... epic metadata

  story:
    # ... story requirements

  requirements:
    # ... detailed requirements

  technical_context:
    # ... architecture details

  structure_validation:
    # ... structure verification
```

## Process

### 1. Load Story Template

#### 1.1 Locate Template File
- **Primary:** `.bmad-core/templates/story-tmpl.yaml`
- **Fallback:** Use embedded template if core template not found
- **Validation:** Ensure template contains all required sections

#### 1.2 Template Sections Expected
- `story`: Basic story information
- `tasks`: Implementation tasks and subtasks
- `dev_notes`: Technical context and implementation guidance
- `testing`: Testing requirements and coverage
- `change_log`: Change history tracking
- `qa_results`: Quality assurance section
- `dev_agent_record`: Development completion tracking

### 2. Populate Basic Story Information

#### 2.1 Story Metadata
```yaml
story:
  id: "3.1"
  title: "MCP Server Implementation"
  as_a: "Developer/Maintainer"
  i_want: "to implement MCP protocol server"
  so_that: "Claude can communicate with the service"
  status: "DRAFT"
  acceptance_criteria:
    - id: AC-1
      description: "WebSocket server implemented"
    - id: AC-2
      description: "HTTP endpoint for MCP available"
    # ... all criteria from epic
```

### 3. Populate Dev Notes Section

This is the **CRITICAL** section that provides all technical context for developers.

#### 3.1 Previous Story Insights
```yaml
dev_notes:
  previous_story_insights: |
    From Story 3.0 completion: Established OAuth integration patterns and
    WebSocket connection handling in user-auth service. Token validation
    middleware created and ready for reuse.
    [Source: Previous story 3.0 dev_agent_record]
```

#### 3.2 Technology Stack
```yaml
  technology_stack:
    language: "Go 1.21"
    framework: "Fiber 2.x"
    mcp_integration: "Mark3Labs MCP-Go library (latest)"
    logging: "zerolog 1.x"
    config: "viper 1.x"
    websocket: "gorilla/websocket"
    source: "[Source: architecture/tech-stack.md#backend-stack]"
```

#### 3.3 Architecture Context
```yaml
  architecture:
    component: "MCP Protocol Handler"
    responsibilities:
      - "Handle MCP protocol communication with Claude AI"
      - "WebSocket endpoint for bidirectional communication"
      - "Tool registration and discovery"
      - "Request/response message handling"
    dependencies:
      - "Network layer"
      - "Command Processor"
    tech_stack:
      - "Go stdlib net/http"
      - "gorilla/websocket"
    source: "[Source: architecture/backend-architecture.md#service-components]"
```

#### 3.4 Data Models (If Applicable)
```yaml
  data_models:
    - name: "MCPMessage"
      fields:
        id: "string (required)"
        method: "string (required, enum)"
        params: "object (optional)"
      validation: ["id non-empty", "method in allowed values"]
      source: "[Source: architecture/data-models.md#mcp-protocol]"
```

#### 3.5 API Specifications (If Applicable)
```yaml
  api_specifications:
    endpoints:
      - path: "/mcp"
        method: "WebSocket"
        description: "MCP protocol WebSocket endpoint"
        auth_required: true
        auth_type: "Bearer Token"
        cors_enabled: true
        source: "[Source: architecture/rest-api-spec.md#websocket-endpoints]"
```

#### 3.6 File Structure
```yaml
  file_structure:
    files:
      - "services/mcp-service/cmd/main.go # MCP server entry point with Mark3Labs library"
      - "services/mcp-service/internal/server/mcp.go # Mark3Labs MCP server setup and configuration"
      - "services/mcp-service/internal/server/tools.go # Tool registration with schema validation"
      - "services/mcp-service/internal/server/handlers.go # Strongly-typed tool handlers"
      - "services/mcp-service/internal/server/middleware.go # Recovery and capability middleware"
    source: "[Source: architecture/unified-project-structure.md#service-structure]"
```

#### 3.7 Technical Constraints
```yaml
  technical_constraints:
    performance:
      connection_establishment: "< 1 second"
      message_processing: "< 100ms per message"
      concurrent_connections: "10+ connections"
      memory_usage: "< 128MB"
    security:
      token_validation: "Required for all connections"
      cors_handling: "Proper CORS headers for Claude client"
    source: "[Source: architecture/backend-architecture.md#performance-requirements]"
```

#### 3.8 Configuration
```yaml
  configuration:
    environment_variables:
      MCP_PORT: "8081"
      ENVIRONMENT: "development"
      LOG_LEVEL: "info"
      LOG_FORMAT: "json"
    source: "[Source: architecture/config-management.md#service-config]"
```

#### 3.9 Integration Points
```yaml
  integration_points:
    oauth_manager: "For token validation"
    cache_manager: "For session management"
    structured_logging: "For all MCP operations"
    source: "[Source: architecture/backend-architecture.md#integration-patterns]"
```

### 4. Generate Implementation Tasks

#### 4.1 Task Generation Strategy
- **Base on:** Epic acceptance criteria + Technical context
- **Structure:** Sequential tasks with clear subtasks
- **Link:** Each task to specific acceptance criteria
- **Include:** Unit testing as explicit subtasks

#### 4.2 Example Task Structure
```yaml
tasks:
  - name: "Set up WebSocket server infrastructure"
    acceptance_criteria: [AC-1, AC-5, AC-6]
    subtasks:
      - "Configure WebSocket endpoint `/mcp` with protocol upgrade"
      - "Implement connection lifecycle management (connect, disconnect, error handling)"
      - "Add concurrent connection support with goroutines"
      - "Implement connection state tracking"
      - "Unit test: WebSocket connection establishment"
      - "Unit test: Concurrent connection handling"
    status: "pending"

  - name: "Implement HTTP endpoint for MCP protocol"
    acceptance_criteria: [AC-2]
    subtasks:
      - "Create HTTP handler that upgrades to WebSocket protocol"
      - "Add proper CORS handling for Claude client connections"
      - "Configure endpoint routing in main server"
      - "Unit test: HTTP to WebSocket upgrade"
      - "Unit test: CORS header validation"
    status: "pending"
```

### 5. Populate Testing Section

#### 5.1 Testing Framework and Strategy
```yaml
testing:
  test_location: "Alongside implementation files (*_test.go)"
  frameworks:
    - "Go standard testing package"
    - "testify for assertions"
    - "gomock for interface mocking"
  source: "[Source: architecture/testing-strategy.md#go-testing]"
```

#### 5.2 Test Requirements
```yaml
  requirements:
    - "Unit tests for WebSocket connection lifecycle"
    - "Integration tests for MCP protocol compliance"
    - "Mock Claude client interactions"
    - "Test concurrent connection handling"
    - "Validate message parsing and response formatting"
  coverage:
    business_logic: "85%"
    overall: "80%"
```

### 6. Initialize Tracking Sections

#### 6.1 Change Log
```yaml
change_log:
  - date: "2025-09-25"  # Current date
    version: "1.0"
    description: "Initial story creation from epic 3.1"
    author: "Bob (Scrum Master)"
```

#### 6.2 QA Results (Template)
```yaml
qa_results:
  review_date: null
  reviewed_by: null
  assessment:
    summary: null
    strengths: []
    improvements: []
    risk_level: null
    testability_score: null
  gate_status: "PENDING"
```

#### 6.3 Dev Agent Record (Template)
```yaml
dev_agent_record:
  agent_model_used: null
  debug_log_references: []
  completion_notes: []
  file_list: []
```

### 7. Create Story File

#### 7.1 Target File Path
- **Path:** `{devStoryLocation}/{epicNum}.{storyNum}.story.md`
- **Example:** `scripts/bmad-cli/templates/3.1.story.md`

#### 7.2 File Creation
- **Format:** YAML format with proper structure
- **Validation:** Ensure all sections populated correctly
- **Backup:** If file exists, create backup before overwriting

## Output Context

```yaml
context:
  # ... previous context sections

  story_file:
    path: "scripts/bmad-cli/templates/3.1.story.md"
    status: "CREATED"
    backup_created: false

  story_content:
    sections_populated:
      - "story"
      - "tasks"
      - "dev_notes"
      - "testing"
      - "change_log"
      - "qa_results"
      - "dev_agent_record"

    technical_details_count:
      data_models: 3
      api_specifications: 2
      file_paths: 6
      constraints: 8

  validation:
    all_sections_complete: true
    source_citations_included: true
    acceptance_criteria_mapped: true
    errors: []
```

## Error Handling

### Template Not Found
- **Action:** Use embedded fallback template
- **Message:** WARN user about missing core template

### Template Malformed
- **Action:** HALT execution
- **Message:** Report template validation errors

### File Creation Failed
- **Action:** HALT execution
- **Message:** Report file system errors with permissions guidance

### Missing Technical Context
- **Action:** Create story with gaps noted
- **Message:** List missing context areas in story dev_notes

## Success Criteria

- Story file created with complete technical context
- All architecture details included with source citations
- Implementation tasks generated and linked to acceptance criteria
- Testing requirements specified
- File ready for developer handoff

## Next Step

On success, pass final context to `step-6-review-completion.md`
