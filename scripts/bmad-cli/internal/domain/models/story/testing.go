package story

type Testing struct {
	TestLocation string            `json:"test_location" yaml:"test_location"`
	Frameworks   []string          `json:"frameworks"    yaml:"frameworks"`
	Requirements []string          `json:"requirements"  yaml:"requirements"`
	Coverage     map[string]string `json:"coverage"      yaml:"coverage"`
}
