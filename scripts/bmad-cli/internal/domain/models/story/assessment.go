package story

type Assessment struct {
	Summary                    string   `yaml:"summary" json:"summary"`
	Strengths                  []string `yaml:"strengths" json:"strengths"`
	Improvements               []string `yaml:"improvements" json:"improvements"`
	RiskLevel                  string   `yaml:"risk_level" json:"risk_level"`
	RiskReason                 string   `yaml:"risk_reason" json:"risk_reason"`
	TestabilityScore           int      `yaml:"testability_score" json:"testability_score"`
	TestabilityMax             int      `yaml:"testability_max" json:"testability_max"`
	TestabilityNotes           string   `yaml:"testability_notes" json:"testability_notes"`
	ImplementationReadiness    int      `yaml:"implementation_readiness" json:"implementation_readiness"`
	ImplementationReadinessMax int      `yaml:"implementation_readiness_max" json:"implementation_readiness_max"`
}
