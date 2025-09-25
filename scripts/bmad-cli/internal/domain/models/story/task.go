package story

type Task struct {
	Name               string   `yaml:"name" json:"name"`
	AcceptanceCriteria []int    `yaml:"acceptance_criteria" json:"acceptance_criteria"`
	Subtasks           []string `yaml:"subtasks" json:"subtasks"`
	Status             string   `yaml:"status" json:"status"`
}
