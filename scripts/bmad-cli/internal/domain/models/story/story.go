package story

type Story struct {
	ID                 string                `yaml:"id" json:"id"`
	Title              string                `yaml:"title" json:"title"`
	AsA                string                `yaml:"as_a" json:"as_a"`
	IWant              string                `yaml:"i_want" json:"i_want"`
	SoThat             string                `yaml:"so_that" json:"so_that"`
	Status             string                `yaml:"status" json:"status"`
	AcceptanceCriteria []AcceptanceCriterion `yaml:"acceptance_criteria" json:"acceptance_criteria"`
}
