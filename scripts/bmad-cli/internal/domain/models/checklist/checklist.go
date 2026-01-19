package checklist

// Checklist represents the entire validation checklist YAML structure.
// Structure: stages[] -> sections[] -> validation_prompts[].
type Checklist struct {
	Version     string   `yaml:"version"`
	LastUpdated string   `yaml:"last_updated"`
	DefaultDocs []string `yaml:"default_docs,omitempty"`
	Stages      []Stage  `yaml:"stages"`
}

// Stage represents a validation stage in the pipeline.
// Stages are processed in order: Story Creation → Refinement → Architecture → Ready Gate.
type Stage struct {
	ID          string    `yaml:"id"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description,omitempty"`
	Sections    []Section `yaml:"sections"`
}

// Section represents a validation section within a stage.
type Section struct {
	ID                string   `yaml:"id"`
	Name              string   `yaml:"name"`
	Description       string   `yaml:"description,omitempty"`
	Source            string   `yaml:"source,omitempty"`
	ValidationPrompts []Prompt `yaml:"validation_prompts"`
}
