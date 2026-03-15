## Current Test

**Scenario ID:** {{.SubjectID}}
**Description:** {{.Subject.Description}}

### Current Test Content
{{- if .Subject.ArchitectureContent }}
```typescript
{{.Subject.ArchitectureContent}}
```
{{- else }}
No existing test content.
{{- end }}

---

## Fix Prompt to Apply

{{.FixPrompt}}

---

## Instructions

Apply the changes described in the "Fix Prompt to Apply" section above.

Output the complete updated test file content using FILE_START/FILE_END markers:

=== FILE_START: {{.ResultPath}} ===
(Complete test file content here)
=== FILE_END: {{.ResultPath}} ===
