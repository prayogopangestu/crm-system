package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/cryptoutil"
)

func (s *Service) ListStages(ctx context.Context, principal domain.Principal) ([]domain.PipelineStage, error) {
	return s.repo.ListStages(ctx, principal.OrganizationID)
}

func (s *Service) CreateStage(ctx context.Context, principal domain.Principal, input domain.StageInput) (domain.PipelineStage, error) {
	if err := requireAdmin(principal); err != nil {
		return domain.PipelineStage{}, err
	}
	if len(strings.TrimSpace(input.Name)) < 2 {
		return domain.PipelineStage{}, domain.ErrInvalidInput
	}
	stage, err := s.repo.CreateStage(ctx, principal.OrganizationID, domain.PipelineStage{Name: input.Name, Color: input.Color})
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return stage, err
}

func (s *Service) ReorderStages(ctx context.Context, principal domain.Principal, ids []string) error {
	if err := requireAdmin(principal); err != nil {
		return err
	}
	if len(ids) == 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.ReorderStages(ctx, principal.OrganizationID, ids)
}

func (s *Service) DeleteStage(ctx context.Context, principal domain.Principal, id string) error {
	if err := requireAdmin(principal); err != nil {
		return err
	}
	return s.repo.DeleteStage(ctx, principal.OrganizationID, id)
}

func (s *Service) GetTelegram(ctx context.Context, principal domain.Principal) (domain.TelegramIntegration, error) {
	if err := requireAdmin(principal); err != nil {
		return domain.TelegramIntegration{}, err
	}
	return s.repo.GetTelegram(ctx, principal.OrganizationID)
}

func (s *Service) UpdateTelegram(ctx context.Context, principal domain.Principal, input domain.TelegramInput) (domain.TelegramIntegration, error) {
	if err := requireAdmin(principal); err != nil {
		return domain.TelegramIntegration{}, err
	}
	var encrypted string
	var err error
	if input.BotToken != "" {
		encrypted, err = s.cipher.Encrypt(input.BotToken)
		if err != nil {
			return domain.TelegramIntegration{}, err
		}
	}
	current, err := s.repo.GetTelegram(ctx, principal.OrganizationID)
	if err != nil {
		return domain.TelegramIntegration{}, err
	}
	if input.Enabled && input.BotToken == "" && !current.HasToken {
		return domain.TelegramIntegration{}, domain.ErrInvalidInput
	}
	if input.Enabled && input.ChatID == "" && current.ChatID == "" {
		return domain.TelegramIntegration{}, domain.ErrInvalidInput
	}
	if err := s.repo.UpsertTelegram(ctx, principal.OrganizationID, domain.TelegramIntegration{
		Enabled: input.Enabled, ChatID: input.ChatID, EncryptedToken: encrypted,
	}); err != nil {
		return domain.TelegramIntegration{}, err
	}
	return s.repo.GetTelegram(ctx, principal.OrganizationID)
}

func (s *Service) TestTelegram(ctx context.Context, principal domain.Principal) error {
	if err := requireAdmin(principal); err != nil {
		return err
	}
	integration, err := s.repo.GetTelegram(ctx, principal.OrganizationID)
	if err != nil {
		return err
	}
	if !integration.Enabled || !integration.HasToken || integration.ChatID == "" {
		return domain.ErrInvalidInput
	}
	token, err := s.cipher.Decrypt(integration.EncryptedToken)
	if err != nil {
		return err
	}
	return s.telegram.Send(ctx, token, integration.ChatID, "CRM Enterprise: koneksi Telegram berhasil.")
}

func (s *Service) Notifications(ctx context.Context, principal domain.Principal) ([]domain.Notification, error) {
	return s.repo.ListNotifications(ctx, principal)
}

func (s *Service) ReadNotification(ctx context.Context, principal domain.Principal, id string) error {
	return s.repo.ReadNotification(ctx, principal, id)
}

func (s *Service) ReadAllNotifications(ctx context.Context, principal domain.Principal) error {
	return s.repo.ReadAllNotifications(ctx, principal)
}

type OutboxWorker struct {
	repo     domain.Repository
	cipher   *cryptoutil.Cipher
	sender   TelegramSender
	logger   *slog.Logger
	interval time.Duration
	batch    int
}

func NewOutboxWorker(repo domain.Repository, cipher *cryptoutil.Cipher, sender TelegramSender, logger *slog.Logger, interval time.Duration, batch int) *OutboxWorker {
	return &OutboxWorker{repo: repo, cipher: cipher, sender: sender, logger: logger, interval: interval, batch: batch}
}

func (w *OutboxWorker) Run(ctx context.Context) {
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

func (w *OutboxWorker) process(ctx context.Context) error {
	events, err := w.repo.ClaimOutbox(ctx, w.batch)
	if err != nil {
		return err
	}
	for _, event := range events {
		if event.EventType != "telegram.deal_won" {
			_ = w.repo.CompleteOutbox(ctx, event.ID)
			continue
		}
		integration, err := w.repo.GetTelegram(ctx, event.OrganizationID)
		if err == nil && (!integration.Enabled || !integration.HasToken || integration.ChatID == "") {
			_ = w.repo.CompleteOutbox(ctx, event.ID)
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
			token, err = w.cipher.Decrypt(integration.EncryptedToken)
		}
		if err == nil {
			err = w.sender.Send(ctx, token, integration.ChatID, payload.Message)
		}
		if err == nil {
			_ = w.repo.CompleteOutbox(ctx, event.ID)
			continue
		}
		delay := time.Duration(1<<min(event.Attempts, 8)) * time.Minute
		_ = w.repo.RetryOutbox(ctx, event.ID, err.Error(), time.Now().Add(delay))
	}
	return nil
}
