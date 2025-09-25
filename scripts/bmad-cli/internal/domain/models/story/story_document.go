package story

type StoryDocument struct {
	Story              Story                 `yaml:"story" json:"story"`
	UserStory          UserStory             `yaml:"user_story" json:"user_story"`
	AcceptanceCriteria []AcceptanceCriterion `yaml:"acceptance_criteria" json:"acceptance_criteria"`
	Tasks              []Task                `yaml:"tasks" json:"tasks"`
	DevNotes           DevNotes              `yaml:"dev_notes" json:"dev_notes"`
	Testing            Testing               `yaml:"testing" json:"testing"`
	ChangeLog          []ChangeLogEntry      `yaml:"change_log" json:"change_log"`
	QAResults          *QAResults            `yaml:"qa_results,omitempty" json:"qa_results,omitempty"`
	DevAgentRecord     DevAgentRecord        `yaml:"dev_agent_record" json:"dev_agent_record"`
}
