package story

import "bmad-cli/internal/infrastructure/docs"

type StoryDocument struct {
	Story            Story                  `json:"story"                yaml:"story"`
	Scenarios        Scenarios              `json:"scenarios"            yaml:"scenarios"`
	ChangeLog        []ChangeLogEntry       `json:"change_log"           yaml:"change_log"`
	QAResults        *QAResults             `json:"qa_results,omitempty" yaml:"qa_results,omitempty"`
	DevAgentRecord   DevAgentRecord         `json:"dev_agent_record"     yaml:"dev_agent_record"`
	ArchitectureDocs *docs.ArchitectureDocs `json:"-"                    yaml:"-"`
}
