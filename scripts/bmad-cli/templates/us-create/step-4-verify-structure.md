<!-- Powered by BMAD™ Core -->

# Step 4: Verify Project Structure Alignment

## Purpose

Cross-reference story requirements with project structure guidelines to ensure file paths, component locations, and module names align with defined standards. Identify and document any structural conflicts or deviations.

## Input Context

```yaml
context:
  requirements:
    story:
      classification: "Backend"
    acceptance_criteria: []
    # ... requirements details

  technical_context:
    file_structure:
      base_path: "services/mcp-service"
      files:
        - path: "internal/server/mcp.go"
          purpose: "MCP server implementation"
          source: "[Source: architecture/unified-project-structure.md#service-structure]"
    # ... other technical context
```

## Process

### 1. Load Project Structure Guidelines

#### 1.1 Primary Structure Document
- **Load:** `docs/architecture/unified-project-structure.md` (or equivalent from config)
- **Extract:**
  - Directory structure standards
  - File naming conventions
  - Module organization patterns
  - Service structure guidelines

#### 1.2 Technology-Specific Patterns
Based on story classification and tech stack:
- **Backend/Go Services:** Extract Go-specific file organization
- **Frontend/React:** Extract React component structure
- **Shared Components:** Extract shared module patterns
- **Configuration Files:** Extract config file placement

### 2. Validate File Structure Requirements

#### 2.1 Check Base Path Alignment
- **Story Files:** Validate paths from `technical_context.file_structure`
- **Naming Conventions:** Check file names follow project standards
- **Directory Structure:** Ensure directories match organizational guidelines

Example validations:
```yaml
validations:
  - path: "services/mcp-service/internal/server/mcp.go"
    expected: "services/{service-name}/internal/{domain}/{file}.go"
    status: "COMPLIANT"

  - path: "services/mcp-service/cmd/main.go"
    expected: "services/{service-name}/cmd/main.go"
    status: "COMPLIANT"

  - path: "services/mcp-service/internal/handlers/websocket.go"
    expected: "services/{service-name}/internal/{domain}/{handler}.go"
    status: "DEVIATION"
    reason: "handlers should be in server domain, not separate handlers directory"
```

#### 2.2 Check Module Organization
- **Domain Separation:** Ensure logical domains are properly separated
- **Dependency Direction:** Check that dependencies flow in correct direction
- **Interface Boundaries:** Validate that interfaces are placed correctly

### 3. Validate Component Placement

#### 3.1 For Backend Stories
- **Service Boundaries:** Check that services are properly scoped
- **Internal Package Structure:** Validate internal package organization
- **Database Layer:** Check data access layer placement
- **API Layer:** Validate API handler placement

#### 3.2 For Frontend Stories
- **Component Hierarchy:** Check React component organization
- **Hook Placement:** Validate custom hooks location
- **Utility Functions:** Check shared utility placement
- **Style Organization:** Validate CSS/styling file placement

#### 3.3 For Infrastructure Stories
- **Configuration Files:** Check config file placement
- **Script Organization:** Validate script and tool placement
- **Documentation:** Check documentation file organization

### 4. Cross-Reference with Existing Codebase

#### 4.1 Scan Existing Structure
- **Current Structure:** Analyze existing project structure
- **Established Patterns:** Identify patterns already in use
- **Consistency Check:** Compare proposed files with existing patterns

#### 4.2 Identify Conflicts
- **Naming Conflicts:** Check for naming collisions
- **Structure Deviations:** Identify where existing structure deviates from guidelines
- **Migration Requirements:** Note if existing files need to be moved/renamed

### 5. Document Structure Decisions

#### 5.1 Compliant Structures
Document paths that align with guidelines:
```yaml
compliant_structures:
  - path: "services/mcp-service/internal/server/mcp.go"
    justification: "Follows service/internal/domain pattern"
    guideline_ref: "[unified-project-structure.md#service-structure]"
```

#### 5.2 Required Deviations
Document justified deviations from guidelines:
```yaml
deviations:
  - path: "services/mcp-service/internal/protocol/parser.go"
    deviation: "Creates new 'protocol' domain not in standard guidelines"
    justification: "MCP protocol handling is sufficiently complex to warrant separate domain"
    approval_needed: true
```

#### 5.3 Structural Conflicts
Document conflicts that need resolution:
```yaml
conflicts:
  - issue: "Existing auth service uses different internal structure"
    impact: "May confuse developers switching between services"
    resolution_options:
      - "Refactor existing auth service to match guidelines"
      - "Document exception for legacy services"
      - "Create migration plan for structural consistency"
```

## Output Context

```yaml
context:
  # ... previous context sections

  structure_validation:
    guidelines_source: "[Source: architecture/unified-project-structure.md]"

    compliant_files:
      - path: "services/mcp-service/cmd/main.go"
        pattern: "services/{service}/cmd/main.go"
        confidence: "HIGH"

    deviation_files:
      - path: "services/mcp-service/internal/protocol/validator.go"
        deviation_type: "NEW_DOMAIN"
        justification: "MCP protocol complexity requires dedicated domain"
        approval_required: false

    conflicts:
      - type: "NAMING_CONFLICT"
        description: "mcp.go conflicts with existing mcp/ directory"
        severity: "MEDIUM"
        resolution: "Rename file to mcp_server.go"

    migration_requirements:
      - action: "CREATE_DIRECTORY"
        path: "services/mcp-service/internal/protocol"
        reason: "New protocol domain needed"

    recommendations:
      - "Consider abstracting protocol handling for reuse in other services"
      - "Document MCP-specific architectural decisions"
      - "Plan integration testing structure for WebSocket connections"
```

## Error Handling

### Structure Guidelines Not Found
- **Action:** WARN user, continue with basic validation
- **Message:** "Project structure guidelines not found. Using basic validation only."

### Severe Structure Conflicts
- **Action:** Note conflicts, continue (don't halt)
- **Message:** Add to conflicts list for user review

### Inconsistent Existing Structure
- **Action:** Document inconsistencies
- **Message:** Add recommendations for structural improvements

## Success Criteria

- Project structure guidelines reviewed and applied
- File paths validated against organizational standards
- Structural conflicts identified and documented
- Migration requirements noted
- Structure validation context prepared for template population

## Next Step

On success, pass enriched context to `step-5-populate-template.md`
