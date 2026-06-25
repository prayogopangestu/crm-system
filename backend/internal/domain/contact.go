package domain

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

type ContactInput struct {
	Name      string `json:"name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Company   string `json:"company" validate:"required,min=2"`
	Role      string `json:"role"`
	Status    string `json:"status" validate:"required,oneof=Negosiasi Menang 'Prospek Awal' Proposal Kalah Kualifikasi"`
	AvatarURL string `json:"avatarUrl"`
}

type ContactList struct {
	Data  []Contact `json:"data"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
}
