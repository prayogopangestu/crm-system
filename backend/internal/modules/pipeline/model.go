package pipeline

import "time"

type Stage struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"-"`
	Key            string    `json:"key"`
	Name           string    `json:"name"`
	Color          string    `json:"color"`
	Position       int       `json:"position"`
	IsSystem       bool      `json:"isSystem,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

type Input struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
