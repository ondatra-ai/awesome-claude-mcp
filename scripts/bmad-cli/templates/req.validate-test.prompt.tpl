## Validate Generated Test Against Architecture

### Scenario
**ID**: `{{.ScenarioID}}`
**Description**: {{.Description}}
**Level**: {{.Level}}
**Service**: {{.Service}}

### Generated Test File
Read the test file at: `tests/{{.Level}}/{{.Service}}.spec.ts`

### Architecture Definition (architecture.yaml)

```yaml
{{.ArchitectureContent}}
```

### Instructions

1. **Read** the generated test file at the path above
2. **Compare** every import, environment variable, helper, fixture, and service URL used in the test against the architecture.yaml content above
3. **Identify** any references in the test that are NOT defined or accounted for in architecture.yaml
4. **For each issue**, propose 2-3 concrete options for how to update architecture.yaml

### Output

Output the validation result as YAML between FILE_START/FILE_END markers:

=== FILE_START: {{.ResultPath}} ===
(YAML validation result here)
=== FILE_END: {{.ResultPath}} ===
