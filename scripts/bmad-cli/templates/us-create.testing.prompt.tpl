<!-- Powered by BMAD™ Core -->

# Create Testing Requirements

## Purpose

Generate testing requirements for user story {{.Story.ID}}.

## Instructions
Generate testing requirements with:
- test_location: where tests go
- frameworks: testing tools to use
- requirements: what to test
- coverage: percentage targets

## Output format:
CRITICAL: Save text content to file: ./tmp/{{.Story.ID}}-testing.yaml. Follow EXACTLY the format below:

=== FILE_START: ./tmp/{{.Story.ID}}-testing.yaml ===
testing:
  test_location: "services/mcp-service"
  frameworks:
    - "Go testing package"
    - "testify"
  requirements:
    - "Unit tests"
    - "Integration tests"
  coverage:
    business_logic: "80%"
    overall: "75%"
=== FILE_END: ./tmp/{{.Story.ID}}-testing.yaml ===

## User Story
```yaml
{{.Story | toYaml}}
```

{{if .ArchitectureDocs}}
{{.ArchitectureDocs.Architecture.Content}}

{{.ArchitectureDocs.FrontendArchitecture.Content}}

{{.ArchitectureDocs.CodingStandards.Content}}

{{.ArchitectureDocs.SourceTree.Content}}
{{end}}
