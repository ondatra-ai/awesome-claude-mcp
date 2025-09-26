package story

type Architecture struct {
	Component        string   `yaml:"component" json:"component"`
	Responsibilities []string `yaml:"responsibilities" json:"responsibilities"`
	Dependencies     []string `yaml:"dependencies" json:"dependencies"`
	TechStack        []string `yaml:"tech_stack" json:"tech_stack"`
}
