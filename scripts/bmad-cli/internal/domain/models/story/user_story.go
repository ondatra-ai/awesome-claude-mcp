package story

type UserStory struct {
	AsA    string `yaml:"as_a" json:"as_a"`
	IWant  string `yaml:"i_want" json:"i_want"`
	SoThat string `yaml:"so_that" json:"so_that"`
}
