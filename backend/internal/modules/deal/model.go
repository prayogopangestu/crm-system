package deal

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

type Input struct {
	Title      string `json:"title"`
	Company    string `json:"company"`
	Value      int64  `json:"value"`
	Priority   string `json:"priority"`
	Stage      string `json:"stage"`
	AssigneeID string `json:"assigneeId"`
	LostReason string `json:"lostReason"`
}

type StageInput struct {
	Stage      string `json:"stage"`
	LostReason string `json:"lostReason"`
}
