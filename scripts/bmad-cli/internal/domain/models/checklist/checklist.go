package checklist

// Checklist represents the entire validation checklist YAML structure.
type Checklist struct {
	Version          string              `yaml:"version"`
	LastUpdated      string              `yaml:"last_updated"`
	DefaultDocs      []string            `yaml:"default_docs,omitempty"`
	Template         TemplateSection     `yaml:"template"`
	Invest           InvestSection       `yaml:"invest"`
	DependenciesRisk DependenciesSection `yaml:"dependencies_risks"`
	AcceptanceCrit   AcceptanceSection   `yaml:"acceptance_criteria"`
	AntiPatterns     []AntiPattern       `yaml:"anti_patterns"`
	DefinitionReady  ReadySection        `yaml:"definition_of_ready"`
	Splitting        SplittingSection    `yaml:"splitting"`
}

// TemplateSection represents the story template validation section.
type TemplateSection struct {
	Origin      string      `yaml:"origin,omitempty"`
	Source      string      `yaml:"source,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Criteria    []Criterion `yaml:"criteria"`
}

// InvestSection represents the INVEST criteria section.
type InvestSection struct {
	Origin      string      `yaml:"origin,omitempty"`
	Source      string      `yaml:"source,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Criteria    []Criterion `yaml:"criteria"`
}

// DependenciesSection represents the dependencies and risks section.
type DependenciesSection struct {
	Description       string       `yaml:"description,omitempty"`
	Source            string       `yaml:"source,omitempty"`
	ValidationPrompts []Prompt     `yaml:"validation_prompts"`
	RiskScoring       *RiskScoring `yaml:"risk_scoring,omitempty"`
}

// RiskScoring represents risk scoring configuration.
type RiskScoring struct {
	Format            string            `yaml:"format"`
	Thresholds        map[string]string `yaml:"thresholds"`
	ValidationPrompts []Prompt          `yaml:"validation_prompts"`
}

// AcceptanceSection represents acceptance criteria validation section.
type AcceptanceSection struct {
	Description       string   `yaml:"description,omitempty"`
	Source            string   `yaml:"source,omitempty"`
	ValidationPrompts []Prompt `yaml:"validation_prompts"`
}

// AntiPattern represents a single anti-pattern to detect.
type AntiPattern struct {
	ID                string   `yaml:"id"`
	Name              string   `yaml:"name"`
	ValidationPrompts []Prompt `yaml:"validation_prompts"`
}

// ReadySection represents the definition of ready section.
type ReadySection struct {
	Description string               `yaml:"description,omitempty"`
	Source      string               `yaml:"source,omitempty"`
	Checklist   map[string][]Prompt  `yaml:"checklist"`
	Scoring     *ReadySectionScoring `yaml:"scoring,omitempty"`
}

// ReadySectionScoring represents scoring configuration for ready section.
type ReadySectionScoring struct {
	Method string `yaml:"method"`
	Note   string `yaml:"note"`
}

// SplittingSection represents the story splitting guidance section.
type SplittingSection struct {
	Description     string                    `yaml:"description,omitempty"`
	Source          string                    `yaml:"source,omitempty"`
	WhenToSplit     []Prompt                  `yaml:"when_to_split"`
	SPIDRTechniques map[string]SPIDRTechnique `yaml:"spidr_techniques"`
}

// SPIDRTechnique represents a single SPIDR splitting technique.
type SPIDRTechnique struct {
	Description string   `yaml:"description"`
	Trigger     string   `yaml:"trigger"`
	Validation  []Prompt `yaml:"validation"`
}

// Criterion represents a validation criterion with prompts.
type Criterion struct {
	ID                string   `yaml:"id"`
	Name              string   `yaml:"name"`
	Description       string   `yaml:"description,omitempty"`
	Source            string   `yaml:"source,omitempty"`
	ValidationPrompts []Prompt `yaml:"validation_prompts"`
}
