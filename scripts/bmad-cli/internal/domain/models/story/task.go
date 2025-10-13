package story

type Task struct {
	Name               string   `json:"name"                yaml:"name"`
	AcceptanceCriteria []string `json:"acceptance_criteria" yaml:"acceptance_criteria"`
	Subtasks           []string `json:"subtasks"            yaml:"subtasks"`
	Status             string   `json:"status"              yaml:"status"`
}
