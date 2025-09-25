package story

type Testing struct {
	TestLocation string            `yaml:"test_location" json:"test_location"`
	Frameworks   []string          `yaml:"frameworks" json:"frameworks"`
	Requirements []string          `yaml:"requirements" json:"requirements"`
	Coverage     map[string]string `yaml:"coverage" json:"coverage"`
}
