# Epics YAML Documentation

This directory contains machine-readable YAML representations of all project epics translated from the original PRD markdown files.

## File Structure

- `epic-01-foundation.yaml` - Foundation & Infrastructure (✅ COMPLETE)
- `epic-02-devops.yaml` - DevOps & Monitoring Infrastructure
- `epic-03-mcp-server.yaml` - MCP Server Setup
- `epic-04-oauth.yaml` - OAuth Authentication
- `epic-05-replace-all.yaml` - Replace All Operation
- `epic-06-append.yaml` - Append Operation
- `epic-07-prepend.yaml` - Prepend Operation
- `epic-08-replace-match.yaml` - Replace Match Operation
- `epic-09-insert-before.yaml` - Insert Before Operation
- `epic-10-insert-after.yaml` - Insert After Operation
- `epics-schema.yaml` - Yamale validation schema
- `README.md` - This documentation

## YAML Structure

Each epic YAML file follows this consistent structure:

```yaml
epic:
  id: <number>                    # Epic identifier (1-10)
  name: <string>                  # Epic name
  status: <enum>                  # COMPLETE|PLANNED|IN_PROGRESS|BLOCKED|CANCELLED
  goal: <string>                  # Epic objective
  completion_summary: <string>    # Optional: Summary for completed epics
  context: <string>               # Required: Additional context (empty string if none)

stories:
  - id: <string>                  # Story identifier (e.g., "1.1", "2.3")
    title: <string>               # Story title
    as_a: <string>                # User persona ("As a...")
    i_want: <string>              # User need ("I want...")
    so_that: <string>             # Business value ("So that...")
    status: <enum>                # Story status
    acceptance_criteria:          # List of acceptance criteria objects
      - id: <string>              # Criteria ID (e.g., "AC-1")
        description: <string>     # Criteria description
    notes: <string>               # Optional: Additional notes

dependencies:                     # List of epic dependencies
  - <string>

success_criteria:                 # List of success metrics
  - <string>

technical_notes:                  # List of technical considerations
  - <string>

```

## Status Values

### Epic Status
- `COMPLETE` - Epic fully implemented and deployed
- `PLANNED` - Epic scheduled for future development
- `IN_PROGRESS` - Epic currently being worked on
- `BLOCKED` - Epic blocked by dependencies or issues
- `CANCELLED` - Epic cancelled or deprioritized

### Story Status
Uses the same values as epic status for consistency.

## Validation

Use the `epics-schema.yaml` file to validate epic YAML files with [Yamale](https://github.com/23andMe/Yamale):

```bash
pip install yamale
yamale -s epics-schema.yaml epic-01-foundation.yaml
```

## Usage

These YAML files can be used for:

1. **Project Management Integration** - Import into tools like Jira, Linear, or GitHub Projects
2. **Progress Tracking** - Monitor epic and story completion status
3. **Dependency Analysis** - Understand epic relationships and blockers
4. **Automation** - Generate reports, dashboards, and metrics
5. **Documentation Generation** - Create formatted documentation from structured data
6. **Validation** - Ensure consistent epic structure and completeness

## Source Material

All YAML files are translated from the original PRD markdown files located in `docs/prd/`:

- `docs/prd/epic-list.md` - Epic overview
- `docs/prd/epic-1-foundation-infrastructure.md` through `docs/prd/epic-10-insert-after-operation.md`

## Development Workflow

The epic development follows this logical progression:

1. **Foundation** (Epic 1) - Infrastructure and deployment foundation ✅
2. **DevOps** (Epic 2) - Monitoring and operational tooling
3. **MCP Server** (Epic 3) - Protocol server implementation
4. **Authentication** (Epic 4) - OAuth integration
5. **Document Operations** (Epics 5-10) - Google Docs editing capabilities
   - Replace All → Append → Prepend → Replace Match → Insert Before → Insert After

Epic 10 includes MVP completion validation for the October 15, 2025 launch target.

## Maintenance

When updating epics:

1. Update the corresponding YAML file
2. Validate against the schema
3. Update the source markdown file if necessary
4. Maintain consistency between YAML and markdown representations
