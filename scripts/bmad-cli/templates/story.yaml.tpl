story:
  id: "{{.ID}}"
  title: "{{.Title | trim}}"
  as_a: "{{.AsA | trim}}"
  i_want: "{{.IWant | trim}}"
  so_that: "{{.SoThat | trim}}"
  status: "{{.Status | default "draft"}}"
  acceptance_criteria:{{range .AcceptanceCriteria}}
    - id: {{.ID}}
      description: "{{.Description}}"{{end}}

tasks:{{range .Tasks}}
  - name: "{{.Name}}"
    acceptance_criteria:{{range .AcceptanceCriteria}}
      - "{{.}}"{{end}}
    subtasks:{{range .Subtasks}}
      - "{{.}}"{{end}}
    status: "{{.Status}}"
{{end}}
dev_notes:
{{.DevNotes | toYaml | nindent 2}}

testing:
  test_location: "{{.Testing.TestLocation}}"
  frameworks:{{range .Testing.Frameworks}}
    - "{{.}}"{{end}}
  requirements:{{range .Testing.Requirements}}
    - "{{.}}"{{end}}
  coverage:{{range $key, $value := .Testing.Coverage}}
    {{$key}}: "{{$value}}"{{end}}

scenarios:
  test_scenarios:{{range .Scenarios.TestScenarios}}
    - id: "{{.ID}}"
      acceptance_criteria: [{{range $i, $ac := .AcceptanceCriteria}}{{if $i}}, {{end}}"{{$ac}}"{{end}}]
      steps:{{range .Steps}}{{if .Given}}
        - given: "{{.Given}}"{{end}}{{if .When}}
        - when: "{{.When}}"{{end}}{{if .Then}}
        - then: "{{.Then}}"{{end}}{{if .And}}
        - and: "{{.And}}"{{end}}{{if .But}}
        - but: "{{.But}}"{{end}}{{end}}{{if .ScenarioOutline}}
      scenario_outline: true
      examples:{{range .Examples}}
        - {{range $key, $val := .}}{{$key}}: {{$val}}
          {{end}}{{end}}{{end}}
      level: "{{.Level}}"
      priority: "{{.Priority}}"{{if .MitigatesRisks}}
      mitigates_risks: [{{range $i, $risk := .MitigatesRisks}}{{if $i}}, {{end}}"{{$risk}}"{{end}}]{{end}}{{end}}

change_log:{{range .ChangeLog}}
  - date: "{{.Date}}"
    version: "{{.Version}}"
    description: "{{.Description | trim}}"
    author: "{{.Author | trim}}"{{end}}
{{if .QAResults}}
qa_results:
  review_date: "{{.QAResults.ReviewDate}}"
  reviewed_by: "{{.QAResults.ReviewedBy}}"

  assessment:
    summary: "{{.QAResults.Assessment.Summary}}"

    strengths:{{range .QAResults.Assessment.Strengths}}
      - "{{.}}"{{end}}

    improvements:{{range .QAResults.Assessment.Improvements}}
      - "{{.}}"{{end}}

    risk_level: "{{.QAResults.Assessment.RiskLevel}}"
    risk_reason: "{{.QAResults.Assessment.RiskReason}}"
    testability_score: {{.QAResults.Assessment.TestabilityScore}}
    testability_max: {{.QAResults.Assessment.TestabilityMax}}
    testability_notes: "{{.QAResults.Assessment.TestabilityNotes}}"
    implementation_readiness: {{.QAResults.Assessment.ImplementationReadiness}}
    implementation_readiness_max: {{.QAResults.Assessment.ImplementationReadinessMax}}

  gate_status: "{{.QAResults.GateStatus}}"
  gate_reference: "{{.QAResults.GateReference}}"
{{end}}
dev_agent_record:
  agent_model_used: {{.DevAgentRecord.AgentModelUsed | default "null"}}
  debug_log_references: {{with .DevAgentRecord.DebugLogReferences}}{{range .}}
    - "{{.}}"{{end}}{{else}}[]{{end}}
  completion_notes: {{with .DevAgentRecord.CompletionNotes}}{{range .}}
    - "{{.}}"{{end}}{{else}}[]{{end}}
  file_list: {{with .DevAgentRecord.FileList}}{{range .}}
    - "{{.}}"{{end}}{{else}}[]{{end}}
