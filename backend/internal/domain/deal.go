package domain

import "time"

type Assignee struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
}

type Deal struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"-"`
	Title          string    `json:"title"`
	Company        string    `json:"company"`
	Value          int64     `json:"value"`
	Priority       string    `json:"priority"`
	Stage          string    `json:"stage"`
	LostReason     string    `json:"lostReason,omitempty"`
	Assignee       Assignee  `json:"assignee"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type DealInput struct {
	Title      string `json:"title" validate:"required,min=2"`
	Company    string `json:"company" validate:"required,min=2"`
	Value      int64  `json:"value" validate:"gte=0"`
	Priority   string `json:"priority" validate:"required,oneof=High Medium Low"`
	Stage      string `json:"stage" validate:"required"`
	AssigneeID string `json:"assigneeId"`
	LostReason string `json:"lostReason"`
}

type StageUpdateInput struct {
	Stage      string `json:"stage" validate:"required"`
	LostReason string `json:"lostReason"`
}
