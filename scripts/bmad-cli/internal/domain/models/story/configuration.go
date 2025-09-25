package story

type Configuration struct {
	EnvironmentVariables map[string]string `yaml:"environment_variables" json:"environment_variables"`
}
