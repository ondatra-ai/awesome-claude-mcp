package models

type Comment struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	URL      string `json:"url"`
	Body     string `json:"body"`
	Outdated bool   `json:"outdated"`
}
