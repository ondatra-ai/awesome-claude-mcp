package story

import "bmad-cli/internal/infrastructure/docs"

type StoryDocument struct {
	Story            Story                  `yaml:"story" json:"story"`
	Tasks            []Task                 `yaml:"tasks" json:"tasks"`
	DevNotes         DevNotes               `yaml:"dev_notes" json:"dev_notes"`
	Testing          Testing                `yaml:"testing" json:"testing"`
	Scenarios        Scenarios              `yaml:"scenarios" json:"scenarios"`
	ChangeLog        []ChangeLogEntry       `yaml:"change_log" json:"change_log"`
	QAResults        *QAResults             `yaml:"qa_results,omitempty" json:"qa_results,omitempty"`
	DevAgentRecord   DevAgentRecord         `yaml:"dev_agent_record" json:"dev_agent_record"`
	ArchitectureDocs *docs.ArchitectureDocs `yaml:"-" json:"-"` // Not serialized, used for generation
}
