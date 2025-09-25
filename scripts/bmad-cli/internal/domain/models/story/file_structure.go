package story

type FileStructureItem map[string][]string

type FileStructure struct {
	ServicePath string                       `yaml:"service_path" json:"service_path"`
	Structure   map[string]FileStructureItem `yaml:"structure" json:"structure"`
}
