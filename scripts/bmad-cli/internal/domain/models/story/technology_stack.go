package story

type TechnologyStack struct {
	Language       string `yaml:"language" json:"language"`
	Framework      string `yaml:"framework" json:"framework"`
	MCPIntegration string `yaml:"mcp_integration" json:"mcp_integration"`
	Logging        string `yaml:"logging" json:"logging"`
	Config         string `yaml:"config" json:"config"`
}
