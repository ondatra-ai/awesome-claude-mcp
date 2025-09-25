package story

type Story struct {
	ID     string `yaml:"id" json:"id"`
	Title  string `yaml:"title" json:"title"`
	Status string `yaml:"status" json:"status"`
}
