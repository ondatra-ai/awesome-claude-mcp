package models

type Thread struct {
	ID         string    `json:"id"`
	IsResolved bool      `json:"isResolved"`
	Comments   []Comment `json:"comments"`
}
