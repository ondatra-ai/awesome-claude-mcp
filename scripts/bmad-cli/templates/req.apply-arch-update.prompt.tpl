## Apply Architecture Update

### Current architecture.yaml

```yaml
{{.ArchitectureContent}}
```

### Issue Found

**Type**: {{.IssueType}}
**Name**: {{.IssueName}}
**Test File**: {{.TestFile}}

### Chosen Update

{{.ChosenOption}}

### Instructions

Apply the chosen update to the architecture.yaml content above. Output the complete updated file between FILE_START/FILE_END markers:

=== FILE_START: {{.ResultPath}} ===
(complete updated architecture.yaml here)
=== FILE_END: {{.ResultPath}} ===
