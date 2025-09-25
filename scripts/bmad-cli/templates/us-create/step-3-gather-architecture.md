<!-- Powered by BMAD™ Core -->

# Step 3: Gather Architecture Context

## Purpose

Read relevant architecture documents based on story type and extract technical context, specifications, and constraints needed for story implementation. This step ensures all technical details come from authoritative architecture sources.

## Input Context

```yaml
context:
  config:
    architectureVersion: "v4"
    architectureSharded: true
    architectureShardedLocation: "docs/architecture"
    architectureFile: "docs/architecture.md"  # if monolithic

  requirements:
    story:
      classification: "Backend"  # Backend|Frontend|Full-Stack|Infrastructure
    # ... requirements details
```

## Process

### 1. Determine Architecture Reading Strategy

#### 1.1 Architecture Version and Structure
- **If `architectureVersion >= v4` AND `architectureSharded: true`:**
  - Read `{architectureShardedLocation}/index.md` first
  - Follow structured reading order based on story classification
- **Else (Monolithic Architecture):**
  - Use `architectureFile` for all sections
  - Extract relevant sections based on story classification

#### 1.2 Validate Architecture Accessibility
- **Check:** All required architecture files exist and are readable
- **Error Handling:** If architecture files missing, HALT with clear error message

### 2. Read Core Architecture Documents

#### 2.1 Universal Documents (Read for ALL Stories)
Read these documents regardless of story type:

1. **`tech-stack.md`** - Technology stack and versions
   - Programming languages and versions
   - Frameworks and libraries
   - Development tools and dependencies
   - Version constraints and compatibility

2. **`unified-project-structure.md`** - File organization
   - Directory structure standards
   - File naming conventions
   - Module organization patterns
   - Code organization guidelines

3. **`coding-standards.md`** - Code quality requirements
   - Code style guidelines
   - Documentation standards
   - Security coding practices
   - Performance considerations

4. **`testing-strategy.md`** - Testing requirements
   - Unit testing frameworks and patterns
   - Integration testing approaches
   - Test coverage requirements
   - Mocking and test data strategies

#### 2.2 Story-Type Specific Documents

Based on `requirements.story.classification`:

##### For Backend/API Stories
Additionally read:
- **`data-models.md`** - Data structures and validation
- **`database-schema.md`** - Database design and schemas
- **`backend-architecture.md`** - Server-side architecture patterns
- **`rest-api-spec.md`** - API design standards
- **`external-apis.md`** - Third-party integration patterns

##### For Frontend/UI Stories
Additionally read:
- **`frontend-architecture.md`** - Client-side architecture
- **`components.md`** - UI component specifications
- **`core-workflows.md`** - User interaction patterns
- **`data-models.md`** - Client-side data structures

##### For Full-Stack Stories
Read both Backend AND Frontend document sets above.

##### For Infrastructure Stories
Additionally read:
- **`deployment.md`** - Deployment and DevOps practices
- **`monitoring.md`** - Observability and logging
- **`security.md`** - Security architecture and practices

### 3. Extract Story-Specific Technical Details

For each architecture document read, extract ONLY information directly relevant to the current story.

#### 3.1 Extraction Guidelines
- **DO NOT invent** new libraries, patterns, or standards not in source documents
- **ONLY extract** information that applies to current story requirements
- **ALWAYS cite** source documents: `[Source: architecture/{filename}.md#{section}]`
- **Be specific** - extract exact technical details, not general descriptions

#### 3.2 Categories to Extract

##### Data Models and Schemas
- Specific data structures the story will use
- Validation rules and constraints
- Field definitions and types
- Relationships between models

##### API Specifications
- Endpoint definitions the story must implement or consume
- Request/response formats and schemas
- Authentication and authorization requirements
- Error response formats

##### Component Specifications (UI Stories)
- Component interface definitions
- Props and state management patterns
- Styling and theme requirements
- Event handling patterns

##### File Structure Requirements
- Exact file paths where new code should be created
- Naming conventions for new files/modules
- Directory organization requirements
- Import/export patterns

##### Testing Requirements
- Specific test types required for this story
- Test file locations and naming
- Coverage requirements
- Mock and fixture patterns

##### Technical Constraints
- Version requirements and compatibility
- Performance benchmarks and limits
- Security requirements and practices
- Error handling standards

#### 3.3 Source Documentation
Every technical detail MUST include its source reference:
- Format: `[Source: architecture/{filename}.md#{section}]`
- Example: `[Source: architecture/backend-architecture.md#websocket-handling]`

If information for a category is not found in architecture docs:
- State explicitly: `"No specific guidance found in architecture docs"`
- Do not invent or assume details

### 4. Organize Technical Context

Structure extracted information by implementation area:

```yaml
technical_context:
  technology_stack:
    language: "Go 1.21"
    framework: "Fiber 2.x"
    libraries: ["gorilla/websocket", "zerolog"]
    source: "[Source: architecture/tech-stack.md#backend-stack]"

  data_models:
    - name: "MCPMessage"
      fields: {"id": "string", "method": "string", "params": "object"}
      validation: ["id required", "method enum validation"]
      source: "[Source: architecture/data-models.md#mcp-protocol]"

  api_specifications:
    endpoints:
      - path: "/mcp"
        method: "WebSocket"
        description: "MCP protocol WebSocket endpoint"
        auth: "Bearer token required"
        source: "[Source: architecture/rest-api-spec.md#websocket-endpoints]"

  file_structure:
    base_path: "services/mcp-service"
    files:
      - path: "internal/server/mcp.go"
        purpose: "MCP server implementation"
        source: "[Source: architecture/unified-project-structure.md#service-structure]"

  testing_requirements:
    unit_coverage: "85%"
    frameworks: ["go test", "testify"]
    patterns: ["table-driven tests", "interface mocking"]
    source: "[Source: architecture/testing-strategy.md#go-testing]"

  constraints:
    performance: ["connection_time < 1s", "message_latency < 100ms"]
    security: ["token_validation_required", "cors_enabled"]
    source: "[Source: architecture/backend-architecture.md#performance-requirements]"
```

## Output Context

```yaml
context:
  # ... previous context sections

  technical_context:
    technology_stack:
      language: "Go 1.21"
      framework: "Fiber 2.x"
      mcp_integration: "Mark3Labs MCP-Go library (latest)"
      logging: "zerolog 1.x"
      config: "viper 1.x"
      source: "[Source: architecture/tech-stack.md#backend-stack]"

    data_models:
      # ... extracted data models with sources

    api_specifications:
      # ... extracted API specs with sources

    component_specifications:  # Only for Frontend/UI stories
      # ... extracted component specs with sources

    file_structure:
      # ... extracted file paths and conventions with sources

    testing_requirements:
      # ... extracted testing details with sources

    constraints:
      # ... extracted constraints with sources

  architecture_gaps:
    missing_specifications: []  # List areas where no guidance was found
    unclear_requirements: []    # List areas where guidance was ambiguous
    conflicts: []               # List any conflicts between documents
```

## Error Handling

### Architecture Files Not Found
- **Action:** HALT execution
- **Message:** List specific missing files and expected locations

### Conflicting Architecture Information
- **Action:** Note conflicts in context, continue
- **Message:** Add conflicts to `architecture_gaps.conflicts` for later resolution

### No Relevant Architecture Found
- **Action:** Continue with empty technical context
- **Message:** Note in `architecture_gaps.missing_specifications`

## Success Criteria

- Relevant architecture documents identified and read
- Story-specific technical details extracted with source citations
- Technical constraints and requirements documented
- File structure and implementation guidance prepared
- Architecture context ready for structure verification

## Next Step

On success, pass enriched context to `step-4-verify-structure.md`
