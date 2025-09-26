package story

type DevNotes struct {
	PreviousStoryInsights   string                  `yaml:"previous_story_insights" json:"previous_story_insights"`
	TechnologyStack         TechnologyStack         `yaml:"technology_stack" json:"technology_stack"`
	Architecture            Architecture            `yaml:"architecture" json:"architecture"`
	FileStructure           FileStructure           `yaml:"file_structure" json:"file_structure"`
	Configuration           Configuration           `yaml:"configuration" json:"configuration"`
	PerformanceRequirements PerformanceRequirements `yaml:"performance_requirements" json:"performance_requirements"`
}
