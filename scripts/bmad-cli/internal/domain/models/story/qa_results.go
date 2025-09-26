package story

type QAResults struct {
	ReviewDate    string     `yaml:"review_date" json:"review_date"`
	ReviewedBy    string     `yaml:"reviewed_by" json:"reviewed_by"`
	Assessment    Assessment `yaml:"assessment" json:"assessment"`
	GateStatus    string     `yaml:"gate_status" json:"gate_status"`
	GateReference string     `yaml:"gate_reference" json:"gate_reference"`
}
