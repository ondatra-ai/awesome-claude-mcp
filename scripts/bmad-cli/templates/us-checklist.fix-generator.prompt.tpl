# Generate Fix Prompt for User Story

## Reference Documentation
- Read(`docs/architecture/bdd-guidelines.md`) - BDD Guidelines

## Original User Story

**Story ID:** {{.Story.ID}}
**Title:** {{.Story.Title}}

**As a** {{.Story.AsA}}
**I want** {{.Story.IWant}}
**So that** {{.Story.SoThat}}

## Current Acceptance Criteria

{{- range $i, $ac := .Story.AcceptanceCriteria }}
{{ add $i 1 }}. **{{ $ac.ID }}:** {{ $ac.Description }}
{{- end }}

## Validation Failure

The following check failed and needs to be fixed:

### Failed Check: {{ .FailedCheck.SectionPath }}

**Question:** {{ .FailedCheck.Question }}
**Expected:** {{ .FailedCheck.ExpectedAnswer }}
**Actual:** {{ .FailedCheck.ActualAnswer }}

**Suggested Fix:**
{{ .FailedCheck.FixPrompt }}

## Task

Generate a complete fix prompt that:
1. Lists the original story and ALL acceptance criteria
2. Clearly identifies which ACs need to change (by ID/number)
3. Provides the complete rewritten ACs (not just fragments)
4. Includes any NEW ACs that should be added (e.g., edge cases)
5. Shows the final complete AC list after fixes

## Output Format

Output using this exact format:

=== FILE_START: {{.ResultPath}} ===
# Fix Prompt for Story {{.Story.ID}}: {{.Story.Title}}

## Instructions
Apply the following changes to the acceptance criteria for this story.

## Original Acceptance Criteria
{{- range $i, $ac := .Story.AcceptanceCriteria }}
{{ add $i 1 }}. {{ $ac.ID }}: {{ $ac.Description }}
{{- end }}

## Required Changes

<For each AC that needs to change, show:>

### Change #N: [AC-ID]
**Before:** <original description>
**Issue:** <what's wrong with it>
**After:**
```gherkin
Scenario: <scenario name>
  Given <context>
  When <action>
  Then <observable outcome>
```

<Also add any NEW ACs needed for edge cases or missing coverage>

## Complete Fixed Acceptance Criteria

<List ALL ACs after applying changes, ready to copy-paste into story file>

1. AC-1:
   Scenario: <name>
     Given <context>
     When <action>
     Then <outcome>

2. AC-2:
   ...

=== FILE_END: {{.ResultPath}} ===
