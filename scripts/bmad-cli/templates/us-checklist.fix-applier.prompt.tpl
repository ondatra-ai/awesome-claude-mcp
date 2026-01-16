## Current Story

**ID:** {{.Story.ID}}
**Title:** {{.Story.Title}}
**As a:** {{.Story.AsA}}
**I want:** {{.Story.IWant}}
**So that:** {{.Story.SoThat}}

### Current Acceptance Criteria
{{range .Story.AcceptanceCriteria}}
**{{.ID}}:**
{{.Description}}
{{end}}

---

## Fix Prompt to Apply

{{.FixPrompt}}

---

## Instructions

Apply the changes described in the "Fix Prompt to Apply" section above.

Output the complete updated acceptance_criteria array using FILE_START/FILE_END markers:

=== FILE_START: {{.ResultPath}} ===
(YAML array of acceptance criteria here)
=== FILE_END: {{.ResultPath}} ===
