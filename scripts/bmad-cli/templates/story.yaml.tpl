story:
  id: "{{.ID}}"
  title: "{{.Title}}"
  as_a: "{{.AsA}}"
  i_want: "{{.IWant}}"
  so_that: "{{.SoThat}}"
  status: "{{.Status}}"
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

change_log:{{range .ChangeLog}}
  - date: "{{.Date}}"
    version: "{{.Version}}"
    description: "{{.Description}}"
    author: "{{.Author}}"{{end}}
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
  agent_model_used: {{if .DevAgentRecord.AgentModelUsed}}{{.DevAgentRecord.AgentModelUsed}}{{else}}null{{end}}
  debug_log_references: {{if .DevAgentRecord.DebugLogReferences}}{{range .DevAgentRecord.DebugLogReferences}}
    - "{{.}}"{{end}}{{else}}[]{{end}}
  completion_notes: {{if .DevAgentRecord.CompletionNotes}}{{range .DevAgentRecord.CompletionNotes}}
    - "{{.}}"{{end}}{{else}}[]{{end}}
  file_list: {{if .DevAgentRecord.FileList}}{{range .DevAgentRecord.FileList}}
    - "{{.}}"{{end}}{{else}}[]{{end}}
