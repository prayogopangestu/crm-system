package domain

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

type TaskInput struct {
	Title      string `json:"title" validate:"required,min=3"`
	Company    string `json:"company" validate:"required,min=2"`
	Time       string `json:"time" validate:"required"`
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
	Type       string `json:"type" validate:"required,oneof=Meeting Call Proposal Other"`
	Priority   string `json:"priority" validate:"required,oneof=Tinggi Sedang Rendah"`
	Assignee   string `json:"assignee"`
	AssigneeID string `json:"assigneeId"`
	Notes      string `json:"notes"`
	Completed  bool   `json:"completed"`
}
