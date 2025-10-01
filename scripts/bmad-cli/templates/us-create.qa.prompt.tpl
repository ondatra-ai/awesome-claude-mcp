<!-- Powered by BMAD™ Core -->

# Generate QA Assessment

## Purpose

Generate comprehensive QA assessment for user story {{.Story.ID}}.

## Instructions
Perform comprehensive story validation following BMAD™ Core methodology:

### Validation Areas (Execute in Order):
1. **Template Completeness**: Verify all required sections are present and no placeholders remain
2. **File Structure Validation**: Check file paths clarity and source tree relevance
3. **Acceptance Criteria Assessment**: Ensure all ACs are satisfied and testable
4. **Testing Instructions Review**: Validate test approach clarity and scenario coverage
5. **Security Considerations**: Assess security requirements and data protection needs
6. **Task Sequence Validation**: Check logical order, dependencies, and granularity
7. **Anti-Hallucination Verification**: Ensure all technical claims are traceable to source documents
8. **Implementation Readiness**: Verify story is self-contained and actionable

### Assessment Categories:
- **Critical Issues**: Must fix - story blocked
- **Should-Fix Issues**: Important quality improvements
- **Nice-to-Have**: Optional enhancements
- **Anti-Hallucination Findings**: Unverifiable claims or missing references

## Output format:
CRITICAL: Save text content to file: ./tmp/{{.Story.ID}}-qa-assessment.yaml. Follow EXACTLY the format below:

=== FILE_START: ./tmp/{{.Story.ID}}-qa-assessment.yaml ===
qa_results:
  review_date: "{{now | date "2006-01-02"}}"
  reviewed_by: "Quinn (Technical QA Architect)"

  validation_areas:
    template_completeness:
      score: 8
      issues: []
      notes: "Assessment of template compliance"

    file_structure:
      score: 8
      issues: []
      notes: "File path and source tree assessment"

    acceptance_criteria:
      score: 8
      issues: []
      notes: "AC coverage and testability assessment"

    testing_instructions:
      score: 8
      issues: []
      notes: "Test approach and scenario assessment"

    security_considerations:
      score: 8
      issues: []
      notes: "Security requirements assessment"

    task_sequence:
      score: 8
      issues: []
      notes: "Task order and dependency assessment"

    anti_hallucination:
      score: 8
      issues: []
      notes: "Source verification and fact checking"

    implementation_readiness:
      score: 8
      issues: []
      notes: "Self-contained context assessment"

  assessment:
    summary: "Comprehensive assessment summary following BMAD™ Core validation methodology"

    critical_issues: []
    should_fix_issues: []
    nice_to_have: []
    anti_hallucination_findings: []

    strengths:
      - "Major strength 1"
      - "Major strength 2"

    risk_level: "Low/Medium/High"
    risk_reason: "Detailed risk assessment explanation"

    testability_score: 8
    testability_max: 10
    testability_notes: "Comprehensive testability assessment"

    implementation_readiness_score: 8
    implementation_readiness_max: 10
    confidence_level: "High/Medium/Low"

  gate_status: "GO/NO-GO"
  gate_reference: "docs/qa/gates/{{.Story.ID}}.yml"
=== FILE_END: ./tmp/{{.Story.ID}}-qa-assessment.yaml ===

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
