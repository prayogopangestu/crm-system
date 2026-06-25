package integration

import (
	"context"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct{ store *postgresx.Store }

type telegramModel struct {
	OrganizationID    string `gorm:"type:uuid;primaryKey"`
	BotTokenEncrypted string
	ChatID            string
	Enabled           bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (telegramModel) TableName() string { return "telegram_integrations" }

type outboxModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	EventType      string
	Payload        []byte `gorm:"type:jsonb"`
	Attempts       int
	NextAttemptAt  time.Time
	ProcessedAt    *time.Time
	LastError      string
	CreatedAt      time.Time
}

func (outboxModel) TableName() string { return "outbox_events" }

func NewRepository(store *postgresx.Store) Repository { return &gormRepository{store: store} }

func (r *gormRepository) GetTelegram(ctx context.Context, organizationID string) (Telegram, error) {
	var record telegramModel
	err := r.store.Query(ctx).Where("organization_id = ?", organizationID).First(&record).Error
	if err != nil {
		if postgresx.MapError(err) == shared.ErrNotFound {
			return Telegram{OrganizationID: organizationID}, nil
		}
		return Telegram{}, err
	}
	value := Telegram{
		OrganizationID: record.OrganizationID, Enabled: record.Enabled, ChatID: record.ChatID,
		EncryptedToken: record.BotTokenEncrypted, UpdatedAt: record.UpdatedAt,
	}
	value.HasToken = value.EncryptedToken != ""
	if value.HasToken {
		value.WebhookURL = "https://api.telegram.org/bot***/sendMessage"
	}
	return value, nil
}

func (r *gormRepository) UpsertTelegram(ctx context.Context, organizationID string, input Telegram) error {
	record := telegramModel{
		OrganizationID: organizationID, Enabled: input.Enabled,
		ChatID: input.ChatID, BotTokenEncrypted: input.EncryptedToken,
	}
	return postgresx.MapError(r.store.Query(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "organization_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"enabled":             gorm.Expr("EXCLUDED.enabled"),
			"chat_id":             gorm.Expr("CASE WHEN EXCLUDED.chat_id = '' THEN telegram_integrations.chat_id ELSE EXCLUDED.chat_id END"),
			"bot_token_encrypted": gorm.Expr("CASE WHEN EXCLUDED.bot_token_encrypted = '' THEN telegram_integrations.bot_token_encrypted ELSE EXCLUDED.bot_token_encrypted END"),
			"updated_at":          time.Now(),
		}),
	}).Create(&record).Error)
}

func (r *gormRepository) ClaimOutbox(ctx context.Context, limit int) ([]OutboxEvent, error) {
	events := make([]OutboxEvent, 0)
	err := r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		var records []outboxModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("processed_at IS NULL AND next_attempt_at <= now()").
			Order("created_at").Limit(limit).Find(&records).Error; err != nil {
			return err
		}
		if len(records) == 0 {
			return nil
		}
		ids := make([]string, 0, len(records))
		for _, record := range records {
			ids = append(ids, record.ID)
			events = append(events, OutboxEvent{
				ID: record.ID, OrganizationID: record.OrganizationID,
				EventType: record.EventType, Payload: record.Payload, Attempts: record.Attempts,
			})
		}
		return tx.Model(&outboxModel{}).Where("id IN ?", ids).
			UpdateColumn("next_attempt_at", gorm.Expr("now() + interval '1 minute'")).Error
	})
	return events, err
}

func (r *gormRepository) CompleteOutbox(ctx context.Context, id string) error {
	return r.store.Query(ctx).Model(&outboxModel{}).Where("id = ?", id).
		Updates(map[string]any{"processed_at": time.Now(), "last_error": ""}).Error
}

func (r *gormRepository) RetryOutbox(ctx context.Context, id, reason string, next time.Time) error {
	return r.store.Query(ctx).Model(&outboxModel{}).Where("id = ?", id).
		Updates(map[string]any{
			"attempts": gorm.Expr("attempts + 1"), "last_error": reason, "next_attempt_at": next,
		}).Error
}
