package contact

import "time"

type Contact struct {
	ID              string    `json:"id"`
	OrganizationID  string    `json:"-"`
	OwnerID         string    `json:"-"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Company         string    `json:"company"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	LastContacted   string    `json:"lastContacted"`
	LastContactedAt time.Time `json:"lastContactedAt,omitempty"`
	Initials        string    `json:"initials"`
	AvatarURL       string    `json:"avatarUrl,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
}

type Input struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	AvatarURL string `json:"avatarUrl"`
}

type List struct {
	Data  []Contact `json:"data"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
}

type Page struct {
	Page  int
	Limit int
}
