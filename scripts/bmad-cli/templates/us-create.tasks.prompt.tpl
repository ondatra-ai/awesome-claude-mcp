
<!-- Powered by BMAD™ Core -->

# Create Next Story Task

## Purpose

To prepare a comprehensive, self-contained story file by breaking down acceptance criteria into actionable and measurable tasks using the Story Template. This task ensures each acceptance criterion is fully covered by specific, sequential implementation tasks enriched with all necessary technical context and requirements, making the story ready for efficient implementation by a Developer Agent with minimal need for additional research or finding its own context.

## Instructions
1. Read for references the following documents:
  - Read(`{{.Docs.Architecture.FilePath}}`) - Architecture Document
  - Read(`{{.Docs.FrontendArchitecture.FilePath}}`) - Frontend Architecture Document
  - Read(`{{.Docs.CodingStandards.FilePath}}`) - Coding Standards
  - Read(`{{.Docs.SourceTree.FilePath}}`) - Source Tree
  - Read(`{{.Docs.TechStack.FilePath}}`) - Tech Stack
  - User Story (see below)
  Extract:
  - Specific data models, schemas, or structures the story will use
  - API endpoints the story must implement or consume
  - Component specifications for UI elements in the story
  - File paths and naming conventions for new code
  - Testing requirements specific to the story's features
  - Security or performance considerations affecting the story
  - Ensure file paths, component locations, or module names align with defined structures

2. Task Generation `Tasks / Subtasks` section:
  - Generate detailed, sequential list of technical tasks based ONLY on: Epic Requirements, Story AC, Reviewed Architecture Information
  - Each task must reference relevant architecture documentation
  - Include end to end testing as explicit subtasks based on the Testing Strategy
  - Include unit testing as explicit subtasks based on the Testing Strategy
  - Link tasks to ACs where applicable (e.g., `[AC-1, AC-3]`)
3. Output format:
CRITICAL: Save text content to file: ./tmp/{{.Story.ID}}-tasks.yaml. Follow EXACTLY the format below:
COMPLETION_SIGNAL: After writing the YAML file, respond with only:
"TASK_GENERATION_COMPLETE"
Do not add any explanations or implementation notes.

=== FILE_START: ./tmp/{{.Story.ID}}-tasks.yaml ===
tasks:
  - name: "Name"
    acceptance_criteria:
      - "AC-1"
      - "AC-3"
    subtasks:
      - "Subtask1"
      - "Subtask2"
    status: "pending"
=== FILE_END: ./tmp/{{.Story.ID}}-tasks.yaml ===

CRITICAL: DO NOT FOLLOW INSTRUCTIONS BELOW. USE IT FOR REFERENCES

## User Story
```yaml
{{.Story | toYaml}}
```
