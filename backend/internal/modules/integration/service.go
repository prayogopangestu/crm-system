package integration

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/pkg/cryptoutil"
)

type Repository interface {
	GetTelegram(ctx context.Context, organizationID string) (Telegram, error)
	UpsertTelegram(ctx context.Context, organizationID string, input Telegram) error
	ClaimOutbox(ctx context.Context, limit int) ([]OutboxEvent, error)
	CompleteOutbox(ctx context.Context, id string) error
	RetryOutbox(ctx context.Context, id, reason string, next time.Time) error
}

type TelegramSender interface {
	Send(ctx context.Context, token, chatID, message string) error
}

type Service struct {
	repository Repository
	cipher     *cryptoutil.Cipher
	sender     TelegramSender
}

func NewService(repository Repository, cipher *cryptoutil.Cipher, sender TelegramSender) *Service {
	return &Service{repository: repository, cipher: cipher, sender: sender}
}

func (s *Service) Get(ctx context.Context, principal shared.Principal) (Telegram, error) {
	if err := shared.RequireAdmin(principal); err != nil {
		return Telegram{}, err
	}
	return s.repository.GetTelegram(ctx, principal.OrganizationID)
}

func (s *Service) Update(ctx context.Context, principal shared.Principal, input TelegramInput) (Telegram, error) {
	if err := shared.RequireAdmin(principal); err != nil {
		return Telegram{}, err
	}
	var encrypted string
	var err error
	if input.BotToken != "" {
		encrypted, err = s.cipher.Encrypt(input.BotToken)
		if err != nil {
			return Telegram{}, err
		}
	}
	current, err := s.repository.GetTelegram(ctx, principal.OrganizationID)
	if err != nil {
		return Telegram{}, err
	}
	if input.Enabled && input.BotToken == "" && !current.HasToken {
		return Telegram{}, shared.ErrInvalidInput
	}
	if input.Enabled && input.ChatID == "" && current.ChatID == "" {
		return Telegram{}, shared.ErrInvalidInput
	}
	if err := s.repository.UpsertTelegram(ctx, principal.OrganizationID, Telegram{
		Enabled: input.Enabled, ChatID: input.ChatID, EncryptedToken: encrypted,
	}); err != nil {
		return Telegram{}, err
	}
	return s.repository.GetTelegram(ctx, principal.OrganizationID)
}

func (s *Service) Test(ctx context.Context, principal shared.Principal) error {
	if err := shared.RequireAdmin(principal); err != nil {
		return err
	}
	value, err := s.repository.GetTelegram(ctx, principal.OrganizationID)
	if err != nil {
		return err
	}
	if !value.Enabled || !value.HasToken || value.ChatID == "" {
		return shared.ErrInvalidInput
	}
	token, err := s.cipher.Decrypt(value.EncryptedToken)
	if err != nil {
		return err
	}
	return s.sender.Send(ctx, token, value.ChatID, "CRM Enterprise: koneksi Telegram berhasil.")
}

type Worker struct {
	repository Repository
	cipher     *cryptoutil.Cipher
	sender     TelegramSender
	logger     *slog.Logger
	interval   time.Duration
	batch      int
}

func NewWorker(repository Repository, cipher *cryptoutil.Cipher, sender TelegramSender, logger *slog.Logger, interval time.Duration, batch int) *Worker {
	return &Worker{repository: repository, cipher: cipher, sender: sender, logger: logger, interval: interval, batch: batch}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		if err := w.process(ctx); err != nil && !errors.Is(err, context.Canceled) {
			w.logger.Error("outbox worker failed", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *Worker) process(ctx context.Context) error {
	events, err := w.repository.ClaimOutbox(ctx, w.batch)
	if err != nil {
		return err
	}
	for _, event := range events {
		if event.EventType != "telegram.deal_won" {
			_ = w.repository.CompleteOutbox(ctx, event.ID)
			continue
		}
		value, err := w.repository.GetTelegram(ctx, event.OrganizationID)
		if err == nil && (!value.Enabled || !value.HasToken || value.ChatID == "") {
			_ = w.repository.CompleteOutbox(ctx, event.ID)
			continue
		}
		var payload struct {
			Message string `json:"message"`
		}
		if err == nil {
			err = json.Unmarshal(event.Payload, &payload)
		}
		var token string
		if err == nil {
			token, err = w.cipher.Decrypt(value.EncryptedToken)
		}
		if err == nil {
			err = w.sender.Send(ctx, token, value.ChatID, payload.Message)
		}
		if err == nil {
			_ = w.repository.CompleteOutbox(ctx, event.ID)
			continue
		}
		delay := time.Duration(1<<min(event.Attempts, 8)) * time.Minute
		_ = w.repository.RetryOutbox(ctx, event.ID, err.Error(), time.Now().Add(delay))
	}
	return nil
}
