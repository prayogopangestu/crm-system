package domain

import "time"

type User struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"-"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Name           string    `json:"name,omitempty"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           string    `json:"role"`
	Status         string    `json:"status,omitempty"`
	AvatarURL      string    `json:"avatarUrl"`
	Initials       string    `json:"initials,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type RegisterInput struct {
	Name        string
	CompanyName string
	Email       string
	Password    string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateProfileInput struct {
	FirstName string
	LastName  string
	Email     string
}

type InviteInput struct {
	Name  string
	Email string
	Role  string
}

type Invitation struct {
	ID             string
	OrganizationID string
	UserID         string
	Email          string
	Role           string
	TokenHash      string
	ExpiresAt      time.Time
	AcceptedAt     *time.Time
}

type InviteResult struct {
	User      User   `json:"data"`
	InviteURL string `json:"inviteUrl"`
}
