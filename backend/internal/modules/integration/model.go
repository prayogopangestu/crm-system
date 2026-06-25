package integration

import "time"

type Telegram struct {
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

type OutboxEvent struct {
	ID             string
	OrganizationID string
	EventType      string
	Payload        []byte
	Attempts       int
}
