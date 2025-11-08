<!-- Powered by BMAD™ Core -->

# Generate QA Assessment

## Purpose

Generate comprehensive QA assessment for user story {{.Story.ID}}.

## Instructions
Analyze the story and provide:
- Assessment summary with strengths and improvements
- Risk level and testability scores
- Gate status (PASS/CONCERNS/FAIL/WAIVED)

## Output format:
CRITICAL: Save text content to file: {{.TmpDir}}/{{.Story.ID}}-qa-assessment.yaml. Follow EXACTLY the format below:

=== FILE_START: {{.TmpDir}}/{{.Story.ID}}-qa-assessment.yaml ===
qa_results:
  review_date: "2025-09-28"
  reviewed_by: "Quinn (Test Architect)"
  assessment:
    summary: "Brief assessment summary"
    strengths:
      - "Strength 1"
      - "Strength 2"
    improvements:
      - "Improvement 1"
      - "Improvement 2"
    risk_level: "Low/Medium/High"
    risk_reason: "Risk explanation"
    testability_score: 8
    testability_max: 10
    testability_notes: "Testability assessment"
    implementation_readiness: 9
    implementation_readiness_max: 10
  gate_status: "PASS"  # Valid: PASS, CONCERNS, FAIL, WAIVED
  gate_reference: "docs/qa/gates/{{.Story.ID}}.yml"
=== FILE_END: {{.TmpDir}}/{{.Story.ID}}-qa-assessment.yaml ===

## User Story
```yaml
{{.Story | toYaml}}
```

## Tasks
```yaml
{{.Tasks | toYaml}}
```

## Development Notes
```yaml
{{.DevNotes | toYaml}}
```
