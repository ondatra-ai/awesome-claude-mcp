package story

type DevAgentRecord struct {
	AgentModelUsed     *string  `yaml:"agent_model_used" json:"agent_model_used"`
	DebugLogReferences []string `yaml:"debug_log_references" json:"debug_log_references"`
	CompletionNotes    []string `yaml:"completion_notes" json:"completion_notes"`
	FileList           []string `yaml:"file_list" json:"file_list"`
}
