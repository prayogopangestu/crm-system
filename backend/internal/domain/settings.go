package domain

import "time"

type PipelineStage struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"-"`
	Key            string    `json:"key"`
	Name           string    `json:"name"`
	Color          string    `json:"color"`
	Position       int       `json:"position"`
	IsSystem       bool      `json:"isSystem,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

type StageInput struct {
	Name  string `json:"name" validate:"required,min=2"`
	Color string `json:"color"`
}

type TelegramIntegration struct {
	OrganizationID string    `json:"-"`
	Enabled        bool      `json:"enabled"`
	WebhookURL     string    `json:"webhookUrl"`
	ChatID         string    `json:"chatId,omitempty"`
	HasToken       bool      `json:"hasToken"`
	EncryptedToken string    `json:"-"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type TelegramInput struct {
	Enabled  bool   `json:"enabled"`
	BotToken string `json:"botToken"`
	ChatID   string `json:"chatId"`
}

type Notification struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Time      string    `json:"time"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type Activity struct {
	ID          string    `json:"id"`
	User        string    `json:"user"`
	Action      string    `json:"action"`
	Target      string    `json:"target"`
	Time        string    `json:"time"`
	IsHighlight bool      `json:"isHighlight"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}
