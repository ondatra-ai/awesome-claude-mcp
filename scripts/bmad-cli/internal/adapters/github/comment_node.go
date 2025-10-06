package github

type commentNode struct {
	Path     *string `json:"path"`
	Line     *int    `json:"line"`
	Body     string  `json:"body"`
	Outdated bool    `json:"outdated"`
	URL      string  `json:"url"`
}
