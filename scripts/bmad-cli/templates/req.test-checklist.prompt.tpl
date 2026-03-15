<!-- Powered by BMAD™ Core -->

# Test Checklist Validation

## Purpose
Evaluate a generated Playwright test against its BDD scenario specification.

## Instructions
1. Read the BDD scenario and test file
2. Answer the validation question
3. ALWAYS explain your reasoning BEFORE the answer block
{{- if .Docs }}

## Reference Documentation
{{- range $key, $doc := .Docs }}
- Read(`{{ $doc.FilePath }}`) - {{ $key | title }}
{{- end }}
{{- end }}

## BDD Scenario

**Scenario ID:** {{.SubjectID}}
**Description:** {{.Subject.Description}}
**Level:** {{.Subject.Level}}
**Service:** {{.Subject.Service}}
**Test File Path:** {{.Subject.TestFilePath}}

### Steps
{{.Subject.FormatSteps}}

### Test File
{{- if .Subject.ArchitectureContent }}
```
{{.Subject.ArchitectureContent}}
```
{{- else }}
No test file content available yet.
{{- end }}

---

## Validation Question

{{.Question}}{{- if .Rationale }} SO THAT WE ENSURE {{.Rationale}}{{- end }}
{{- if .FixTemplate }}

## If Validation Fails
If your answer does NOT match the expected criteria, generate a fix_prompt using this template:

{{ .FixTemplate }}
{{- end }}

## Answer Format
Output your answer using this exact format:

=== FILE_START: {{.ResultPath}} ===
answer: <your answer here>
{{- if .FixTemplate }}
fix_prompt: |
  <if validation fails, provide fix guidance here using template above>
  <if validation passes, leave this field empty or omit it>
{{- end }}
=== FILE_END: {{.ResultPath}} ===
