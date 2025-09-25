<!-- Powered by BMAD™ Core -->

# Step 6: Story Draft Completion and Review

## Purpose

Perform final validation of the created story file, execute quality checklists, and provide completion summary to the user. This step ensures the story is ready for developer handoff.

## Input Context

```yaml
context:
  story_file:
    path: "scripts/bmad-cli/templates/3.1.story.md"
    status: "CREATED"

  story_content:
    sections_populated: []
    technical_details_count: {}

  # ... all previous context
```

## Process

### 1. Story File Validation

#### 1.1 File Integrity Check
- **Verify:** Story file exists and is readable
- **Validate:** File format is valid YAML
- **Check:** All required sections are present and populated

#### 1.2 Content Completeness Review
- **Story Metadata:** Verify all basic story information is complete
- **Acceptance Criteria:** Ensure all ACs from epic are included
- **Dev Notes:** Validate technical context sections are populated
- **Tasks:** Check implementation tasks are generated and linked to ACs
- **Testing:** Verify testing requirements are specified

#### 1.3 Source Citation Validation
- **Architecture References:** Ensure all technical details include source citations
- **Format Check:** Validate citation format: `[Source: architecture/{file}.md#{section}]`
- **Coverage:** Check that major technical decisions reference architecture docs

### 2. Execute Quality Checklist

#### 2.1 Load Story Draft Checklist
- **Load:** `.bmad-core/checklists/story-draft-checklist.md`
- **Execute:** Run checklist validation against created story
- **Document:** Checklist results and any failing items

#### 2.2 Checklist Categories
Expected checklist items to validate:

##### Story Definition Quality
- [ ] User story follows "As a... I want... So that..." format
- [ ] Acceptance criteria are testable and measurable
- [ ] Story scope is appropriate for single development iteration
- [ ] Dependencies are clearly identified

##### Technical Context Completeness
- [ ] Technology stack requirements specified with versions
- [ ] File structure and naming conventions defined
- [ ] Integration points with existing systems documented
- [ ] Performance and security constraints included

##### Implementation Readiness
- [ ] Tasks break down acceptance criteria into actionable items
- [ ] File paths follow project structure guidelines
- [ ] Testing requirements specify frameworks and coverage
- [ ] Error handling and logging requirements included

##### Documentation Standards
- [ ] All technical details include architecture source references
- [ ] Previous story insights incorporated where relevant
- [ ] Configuration requirements documented
- [ ] Change log entry created

### 3. Identify Story Quality Issues

#### 3.1 Missing Context Analysis
Identify areas where context might be insufficient:
- **Architecture Gaps:** Areas where no specific guidance was found
- **Integration Unclear:** Unclear integration points with existing systems
- **Testing Incomplete:** Missing testing scenarios or frameworks

#### 3.2 Complexity Assessment
Evaluate story complexity and implementation risk:
- **High Complexity Indicators:**
  - Multiple new dependencies
  - Complex integration requirements
  - Performance-critical functionality
  - New architectural patterns

#### 3.3 Developer Readiness Score
Assess how ready the story is for developer handoff:
```yaml
readiness_assessment:
  technical_context: 9/10  # Strong architecture guidance
  implementation_clarity: 8/10  # Clear tasks and subtasks
  testing_guidance: 7/10  # Good coverage requirements
  integration_clarity: 6/10  # Some integration points unclear
  overall_score: 7.5/10
```

### 4. Generate Completion Summary

#### 4.1 Story Creation Summary
```yaml
completion_summary:
  story_created: "scripts/bmad-cli/templates/3.1.story.md"
  epic_source: "docs/epics/jsons/epic-03-mcp-server.yaml"
  status: "Draft"

  technical_context:
    architecture_docs_referenced: 8
    data_models_defined: 3
    api_endpoints_specified: 2
    file_paths_documented: 6
    testing_requirements: 5

  implementation_guidance:
    tasks_generated: 4
    subtasks_total: 16
    acceptance_criteria_mapped: 6
    unit_tests_specified: 8
```

#### 4.2 Key Technical Components
```yaml
  key_components:
    technology_stack: "Go 1.21, Fiber 2.x, Mark3Labs MCP-Go, gorilla/websocket"
    primary_files:
      - "services/mcp-service/internal/server/mcp.go"
      - "services/mcp-service/internal/server/handlers.go"
    integration_points: ["OAuth Manager", "Cache Manager", "Structured Logging"]
    performance_targets: ["<1s connection", "<100ms message processing"]
```

### 5. Document Deviations and Conflicts

#### 5.1 Architecture Deviations
```yaml
  deviations:
    - type: "Structure"
      description: "Created new 'protocol' domain for MCP-specific logic"
      justification: "Protocol complexity warrants dedicated domain"
      approval_required: false
```

#### 5.2 Epic vs Architecture Conflicts
```yaml
  conflicts:
    - area: "File Structure"
      epic_requirement: "WebSocket server in main service"
      architecture_guidance: "WebSocket handlers in separate domain"
      resolution: "Follow architecture guidance with protocol domain"
```

### 6. Provide Next Steps Guidance

#### 6.1 For Complex Stories
If story complexity is high or has significant gaps:
```yaml
recommendations:
  complexity_level: "Medium-High"
  suggested_actions:
    - "Review story with Product Owner for scope validation"
    - "Consider splitting into smaller stories if implementation > 1 week"
    - "Schedule architecture review for new protocol domain"
    - "Plan integration testing with existing OAuth/Cache systems"
```

#### 6.2 For Simple Stories
If story is straightforward and complete:
```yaml
recommendations:
  complexity_level: "Low-Medium"
  suggested_actions:
    - "Story ready for developer assignment"
    - "Schedule brief technical walkthrough if developer unfamiliar with MCP"
```

## Output Context

```yaml
final_context:
  story_file:
    path: "scripts/bmad-cli/templates/3.1.story.md"
    status: "DRAFT_COMPLETE"
    size_kb: 8.5
    last_modified: "2025-09-25T10:30:00Z"

  validation:
    checklist_passed: true
    checklist_score: "8/10"
    failing_items: []
    warnings: ["Integration testing scenarios could be more specific"]

  readiness:
    developer_ready: true
    complexity_level: "Medium"
    estimated_effort: "3-5 days"
    risk_level: "Low-Medium"

  recommendations:
    immediate_actions: []
    optional_reviews: ["Architecture review for protocol domain"]
    next_story_suggestions: ["3.2 Tool Registration"]
```

## User Summary Output

Provide a comprehensive summary to the user:

```
✅ Story 3.1 Created Successfully!

📁 File: scripts/bmad-cli/templates/3.1.story.md
📊 Status: DRAFT COMPLETE
🎯 Epic: MCP Server Setup (3)

📋 Content Summary:
• 6 Acceptance Criteria mapped to implementation tasks
• 8 Architecture documents referenced for technical context
• 16 Subtasks generated with testing requirements
• 3 Data models defined with validation rules
• 2 API endpoints specified with auth requirements

🏗️ Technical Stack:
• Go 1.21 with Fiber 2.x framework
• Mark3Labs MCP-Go library for protocol handling
• WebSocket with gorilla/websocket
• Integration with OAuth Manager and Cache Manager

✅ Quality Checklist: 8/10 passed
⚠️  Warnings: Integration testing scenarios could be more specific

🚀 Next Steps:
• Story ready for developer assignment
• Consider brief technical walkthrough for MCP protocol concepts
• Optional: Architecture review for new protocol domain

📈 Complexity: Medium (estimated 3-5 days)
🎭 Risk Level: Low-Medium

Would you like to:
1) Create the next story (3.2)
2) Review the generated story file
3) Run additional validation
```

## Error Handling

### Story File Validation Failed
- **Action:** Report specific validation errors
- **Message:** Provide corrective actions needed

### Checklist Execution Failed
- **Action:** Continue with manual review
- **Message:** WARN user about missing checklist validation

### Severe Quality Issues
- **Action:** Flag story as requiring review
- **Message:** List specific issues that need resolution

## Success Criteria

- Story file validated and complete
- Quality checklist executed successfully
- Completion summary provided to user
- Next steps clearly communicated
- Story ready for developer handoff or further review

## Final Step

This is the final step in the story creation orchestration. Return control to user with completion summary and recommendations.
