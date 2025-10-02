<!-- Powered by BMAD™ Core -->

# Generate Development Notes

## Purpose

To analyze a user story and generate comprehensive `dev_notes` that provide essential technical context for implementation. The dev_notes should contain specific technical details derived from the story requirements and architecture documentation, enabling developers to implement the story efficiently without additional research.

## Instructions

1. Read for references the following documents:
  - Read(`{{.Docs.Architecture.FilePath}}`) - Architecture Document
  - Read(`{{.Docs.FrontendArchitecture.FilePath}}`) - Frontend Architecture Document
  - Read(`{{.Docs.CodingStandards.FilePath}}`) - Coding Standards
  - Read(`{{.Docs.SourceTree.FilePath}}`) - Source Tree
  - Read(`{{.Docs.TechStack.FilePath}}`) - Tech Stack
  - User Story (see below)
  - Generated Tasks (see below)
  Extract:
  - Specific technology stack components needed for this story
  - Relevant architecture patterns and components
  - File paths and naming conventions for implementation
  - Performance requirements specific to the story's features
  - Configuration and environment variables needed
  - Integration points with existing systems

2. **Generate Technical Context**:
   CRITICAL: For each entity (technology_stack, architecture, file_structure, etc.), you MUST include:
   - **source**: Exact file path and section reference (e.g., "docs/architecture.md#Backend Components")
   - **description**: Brief explanation starting with "From the [document type]:" (e.g., "From the MCP protocol workflow diagram:")

   Additional flexible fields per entity:
   - **previous_story_insights**: Analyze story context and provide insights about implementation approach
   - **technology_stack**: Specify exact languages, frameworks, libraries, and tools needed
   - **architecture**: Define component responsibilities, dependencies, and tech stack
   - **file_structure**: Provide specific file paths where implementation should occur
   - **configuration**: Define environment variables and configuration needed
   - **performance_requirements**: Set realistic performance targets based on story scope

3. Output format:
CRITICAL: Save text content to file: ./tmp/{{.Story.ID}}-devnotes.yaml. Follow EXACTLY the format below:
COMPLETION_SIGNAL: After writing the YAML file, respond with only:
"DEVNOTES_GENERATION_COMPLETE"
Do not add any explanations or implementation notes.

=== FILE_START: ./tmp/{{.Story.ID}}-devnotes.yaml ===
dev_notes:
  previous_story_insights: "Detailed analysis of story context and implementation approach"

  technology_stack:
    source: "{{.Docs.TechStack.FilePath}}#Backend Stack"
    description: "From the backend technology stack documentation:"
    language: "Primary programming language"
    framework: "Main framework or library"
    mcp_integration: "MCP integration approach"
    logging: "Logging framework"
    config: "Configuration management"

  architecture:
    source: "{{.Docs.Architecture.FilePath}}#Backend Components"
    description: "From the MCP protocol workflow diagram:"
    component: "Main component name"
    responsibilities:
      - "Primary responsibility"
      - "Secondary responsibility"
    dependencies:
      - "Key dependency 1"
      - "Key dependency 2"
    tech_stack:
      - "Technology 1"
      - "Technology 2"

  file_structure:
    source: "{{.Docs.SourceTree.FilePath}}#Service Structure"
    description: "Based on the project file structure:"
    files:
      - file: "specific/path/to/implementation.go"
        description: "Main implementation file"
      - file: "specific/path/to/tests.go"
        description: "Test files"

  configuration:
    source: "{{.Docs.CodingStandards.FilePath}}#Environment Variables"
    description: "Required environment variables for the service:"
    environment_variables:
      VARIABLE_NAME: "default_value"
      ANOTHER_VAR: "value"

  performance_requirements:
    source: "{{.Docs.CodingStandards.FilePath}}#Performance Standards"
    description: "Performance requirements based on coding standards:"
    connection_establishment: "< Xms"
    message_processing: "< Xms"
    concurrent_connections: "X"
    memory_usage: "< XMB"
=== FILE_END: ./tmp/{{.Story.ID}}-devnotes.yaml ===

CRITICAL: DO NOT FOLLOW INSTRUCTIONS BELOW. USE IT FOR REFERENCES

## User Story
```yaml
{{.Story | toYaml}}
```

## Generated Tasks
```yaml
{{.Tasks | toYaml}}
```
