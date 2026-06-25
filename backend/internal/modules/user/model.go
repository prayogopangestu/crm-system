package user

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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateProfileInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type InviteInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
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
