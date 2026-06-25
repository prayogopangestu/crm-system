package task

import "time"

type Task struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"-"`
	Title          string    `json:"title"`
	Company        string    `json:"company"`
	Time           string    `json:"time"`
	Date           string    `json:"date"`
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	Completed      bool      `json:"completed"`
	Notes          string    `json:"notes"`
	Priority       string    `json:"priority"`
	Assignee       string    `json:"assignee"`
	AssigneeID     string    `json:"assigneeId,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type Input struct {
	Title      string `json:"title"`
	Company    string `json:"company"`
	Time       string `json:"time"`
	Date       string `json:"date"`
	Type       string `json:"type"`
	Priority   string `json:"priority"`
	Assignee   string `json:"assignee"`
	AssigneeID string `json:"assigneeId"`
	Notes      string `json:"notes"`
	Completed  bool   `json:"completed"`
}
