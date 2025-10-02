<!-- Powered by BMAD™ Core -->

# Create Testing Requirements

## Purpose

Generate comprehensive testing requirements for user story {{.Story.ID}} that ensure complete validation of all acceptance criteria through unit, integration, and end-to-end tests.

## Instructions

1. Read for references the following documents:
  - Read(`{{.Docs.Architecture.FilePath}}`) - Architecture Document
  - Read(`{{.Docs.FrontendArchitecture.FilePath}}`) - Frontend Architecture Document
  - Read(`{{.Docs.CodingStandards.FilePath}}`) - Coding Standards
  - Read(`{{.Docs.SourceTree.FilePath}}`) - Source Tree
  - Read(`{{.Docs.TechStack.FilePath}}`) - Tech Stack
  - User Story (see below)
  - Generated Tasks (see below)
  - Generated DevNotes (see below)
  Extract:
  - Testing frameworks and tools specified in tech stack
  - Test coverage requirements from coding standards
  - Test location conventions from source tree
  - Performance benchmarks that need testing
  - Security requirements that need validation
  - Integration points that need testing

2. Generate Testing Requirements:
  - **test_location**: Specify where tests should be located based on source tree
  - **frameworks**: List testing tools from tech stack (unit, integration, E2E frameworks)
  - **requirements**: Create specific test scenarios for EACH acceptance criterion
  - **coverage**: Define percentage targets by component (align with coding standards)
  - Link test requirements to acceptance criteria (e.g., tests for AC-1, AC-2)
  - Include unit tests, integration tests, and E2E tests as appropriate
  - Specify load/performance tests if story has performance requirements
  - Include security tests if story has security considerations

3. Output format:
CRITICAL: Save text content to file: ./tmp/{{.Story.ID}}-testing.yaml. Follow EXACTLY the format below:
COMPLETION_SIGNAL: After writing the YAML file, respond with only:
"TESTING_GENERATION_COMPLETE"
Do not add any explanations or implementation notes.

=== FILE_START: ./tmp/{{.Story.ID}}-testing.yaml ===
testing:
  test_location: "services/backend"
  frameworks:
    - "Go testing package"
    - "testify"
    - "httptest"
  requirements:
    - "Unit test for feature X (AC-1)"
    - "Integration test for endpoint Y (AC-2)"
    - "E2E test for workflow Z (AC-1, AC-2)"
    - "Load test for concurrent operations (AC-3)"
  coverage:
    business_logic: "90%"
    http_handlers: "85%"
    integration: "80%"
    overall: "85%"
=== FILE_END: ./tmp/{{.Story.ID}}-testing.yaml ===

CRITICAL: DO NOT FOLLOW INSTRUCTIONS BELOW. USE IT FOR REFERENCES

## User Story
```yaml
{{.Story | toYaml}}
```

## Generated Tasks
```yaml
{{.Tasks | toYaml}}
```

## Generated DevNotes
```yaml
{{.DevNotes | toYaml}}
```
