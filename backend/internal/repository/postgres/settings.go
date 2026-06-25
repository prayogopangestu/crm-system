package postgres

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func (r *Repository) ListStages(ctx context.Context, organizationID string) ([]domain.PipelineStage, error) {
	var models []pipelineStageModel
	if err := r.query(ctx).
		Where("organization_id = ?", organizationID).
		Order("position").
		Find(&models).Error; err != nil {
		return nil, err
	}
	stages := make([]domain.PipelineStage, 0, len(models))
	for _, model := range models {
		stages = append(stages, stageFromModel(model))
	}
	return stages, nil
}

func (r *Repository) CreateStage(ctx context.Context, organizationID string, stage domain.PipelineStage) (domain.PipelineStage, error) {
	stage.Key = strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(stage.Name), "-"), "-")
	if stage.Key == "" {
		return domain.PipelineStage{}, domain.ErrInvalidInput
	}
	if stage.Color == "" {
		stage.Color = "bg-surface-variant"
	}
	var model pipelineStageModel
	err := r.query(ctx).Transaction(func(tx *gorm.DB) error {
		var maxPosition int
		if err := tx.Model(&pipelineStageModel{}).
			Select("COALESCE(MAX(position), 0)").
			Where("organization_id = ?", organizationID).
			Scan(&maxPosition).Error; err != nil {
			return err
		}
		model = pipelineStageModel{
			OrganizationID: organizationID,
			Key:            stage.Key,
			Name:           stage.Name,
			Color:          stage.Color,
			Position:       maxPosition + 1,
		}
		return mapError(tx.Create(&model).Error)
	})
	if err != nil {
		return domain.PipelineStage{}, err
	}
	return stageFromModel(model), nil
}

func (r *Repository) ReorderStages(ctx context.Context, organizationID string, ids []string) error {
	return r.query(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&pipelineStageModel{}).
			Where("organization_id = ?", organizationID).
			Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(ids)) {
			return domain.ErrInvalidInput
		}
		if err := tx.Model(&pipelineStageModel{}).
			Where("organization_id = ?", organizationID).
			UpdateColumn("position", gorm.Expr("position + 1000")).Error; err != nil {
			return mapError(err)
		}
		for index, id := range ids {
			result := tx.Model(&pipelineStageModel{}).
				Where("id = ? AND organization_id = ?", id, organizationID).
				Updates(map[string]any{"position": index + 1, "updated_at": time.Now()})
			if result.Error != nil {
				return mapError(result.Error)
			}
			if result.RowsAffected != 1 {
				return domain.ErrInvalidInput
			}
		}
		return nil
	})
}

func (r *Repository) DeleteStage(ctx context.Context, organizationID, id string) error {
	var model pipelineStageModel
	if err := r.query(ctx).
		Where("id = ? AND organization_id = ?", id, organizationID).
		First(&model).Error; err != nil {
		return mapError(err)
	}
	if model.IsSystem {
		return domain.ErrForbidden
	}
	var used int64
	if err := r.query(ctx).Model(&dealModel{}).
		Where("organization_id = ? AND stage_key = ? AND deleted_at IS NULL", organizationID, model.Key).
		Count(&used).Error; err != nil {
		return err
	}
	if used > 0 {
		return domain.ErrStageInUse
	}
	result := r.query(ctx).Where("id = ? AND organization_id = ?", id, organizationID).Delete(&pipelineStageModel{})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) GetTelegram(ctx context.Context, organizationID string) (domain.TelegramIntegration, error) {
	var model telegramIntegrationModel
	err := r.query(ctx).Where("organization_id = ?", organizationID).First(&model).Error
	if err != nil {
		if mapError(err) == domain.ErrNotFound {
			return domain.TelegramIntegration{OrganizationID: organizationID}, nil
		}
		return domain.TelegramIntegration{}, err
	}
	value := domain.TelegramIntegration{
		OrganizationID: model.OrganizationID,
		Enabled:        model.Enabled,
		ChatID:         model.ChatID,
		EncryptedToken: model.BotTokenEncrypted,
		UpdatedAt:      model.UpdatedAt,
	}
	value.HasToken = value.EncryptedToken != ""
	if value.HasToken {
		value.WebhookURL = "https://api.telegram.org/bot***/sendMessage"
	}
	return value, nil
}

func (r *Repository) UpsertTelegram(ctx context.Context, organizationID string, input domain.TelegramIntegration) error {
	model := telegramIntegrationModel{
		OrganizationID:    organizationID,
		Enabled:           input.Enabled,
		ChatID:            input.ChatID,
		BotTokenEncrypted: input.EncryptedToken,
	}
	return mapError(r.query(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "organization_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"enabled":             gorm.Expr("EXCLUDED.enabled"),
			"chat_id":             gorm.Expr("CASE WHEN EXCLUDED.chat_id = '' THEN telegram_integrations.chat_id ELSE EXCLUDED.chat_id END"),
			"bot_token_encrypted": gorm.Expr("CASE WHEN EXCLUDED.bot_token_encrypted = '' THEN telegram_integrations.bot_token_encrypted ELSE EXCLUDED.bot_token_encrypted END"),
			"updated_at":          time.Now(),
		}),
	}).Create(&model).Error)
}

func (r *Repository) Search(ctx context.Context, organizationID, query string) (domain.SearchResult, error) {
	pattern := "%" + query + "%"
	result := domain.SearchResult{Contacts: []domain.Contact{}, Tasks: []domain.Task{}, Deals: []domain.Deal{}}

	contactRows, err := r.query(ctx).Model(&contactModel{}).
		Where("organization_id = ? AND deleted_at IS NULL", organizationID).
		Where("name ILIKE ? OR email ILIKE ? OR company ILIKE ?", pattern, pattern, pattern).
		Order("updated_at DESC").Limit(10).Rows()
	if err != nil {
		return result, err
	}
	for contactRows.Next() {
		var model contactModel
		if err := r.query(ctx).ScanRows(contactRows, &model); err != nil {
			contactRows.Close()
			return result, err
		}
		item := contactFromModel(model)
		r.decorateContact(&item)
		result.Contacts = append(result.Contacts, item)
	}
	contactRows.Close()

	taskRows, err := r.query(ctx).
		Table("tasks AS t").
		Select(taskSelect).
		Joins("LEFT JOIN users u ON u.id = t.assignee_id AND u.organization_id = t.organization_id").
		Where("t.organization_id = ? AND t.deleted_at IS NULL", organizationID).
		Where("t.title ILIKE ? OR t.company ILIKE ? OR t.notes ILIKE ?", pattern, pattern, pattern).
		Order("t.updated_at DESC").Limit(10).Rows()
	if err != nil {
		return result, err
	}
	for taskRows.Next() {
		item, scanErr := scanTask(taskRows, time.Now().In(r.location))
		if scanErr != nil {
			taskRows.Close()
			return result, scanErr
		}
		result.Tasks = append(result.Tasks, item)
	}
	taskRows.Close()

	dealRows, err := r.query(ctx).Raw(
		dealSelect+` WHERE d.organization_id = ? AND d.deleted_at IS NULL
			AND (d.title ILIKE ? OR d.company ILIKE ?)
			ORDER BY d.updated_at DESC LIMIT 10`,
		organizationID, pattern, pattern,
	).Rows()
	if err != nil {
		return result, err
	}
	defer dealRows.Close()
	for dealRows.Next() {
		item, scanErr := scanDeal(dealRows)
		if scanErr != nil {
			return result, scanErr
		}
		result.Deals = append(result.Deals, item)
	}
	return result, dealRows.Err()
}

func (r *Repository) ListNotifications(ctx context.Context, principal domain.Principal) ([]domain.Notification, error) {
	rows, err := r.query(ctx).Raw(`
		SELECT id,title,message,created_at,read_at IS NOT NULL
		FROM notifications
		WHERE organization_id = ? AND (user_id = ? OR user_id IS NULL)
		ORDER BY created_at DESC LIMIT 100`,
		principal.OrganizationID, principal.UserID,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	now := time.Now().In(r.location)
	items := make([]domain.Notification, 0)
	for rows.Next() {
		var item domain.Notification
		if err := rows.Scan(&item.ID, &item.Title, &item.Message, &item.CreatedAt, &item.Read); err != nil {
			return nil, err
		}
		item.Time = humanTime(item.CreatedAt, now)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) ReadNotification(ctx context.Context, principal domain.Principal, id string) error {
	result := r.query(ctx).Model(&notificationModel{}).
		Where("id = ? AND organization_id = ? AND (user_id = ? OR user_id IS NULL)", id, principal.OrganizationID, principal.UserID).
		UpdateColumn("read_at", gorm.Expr("COALESCE(read_at, now())"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) ReadAllNotifications(ctx context.Context, principal domain.Principal) error {
	return r.query(ctx).Model(&notificationModel{}).
		Where("organization_id = ? AND (user_id = ? OR user_id IS NULL)", principal.OrganizationID, principal.UserID).
		UpdateColumn("read_at", gorm.Expr("COALESCE(read_at, now())")).Error
}

func (r *Repository) ClaimOutbox(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	events := make([]domain.OutboxEvent, 0)
	err := r.query(ctx).Transaction(func(tx *gorm.DB) error {
		var models []outboxEventModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("processed_at IS NULL AND next_attempt_at <= now()").
			Order("created_at").
			Limit(limit).
			Find(&models).Error; err != nil {
			return err
		}
		if len(models) == 0 {
			return nil
		}
		ids := make([]string, 0, len(models))
		for _, model := range models {
			ids = append(ids, model.ID)
			events = append(events, domain.OutboxEvent{
				ID: model.ID, OrganizationID: model.OrganizationID,
				EventType: model.EventType, Payload: model.Payload, Attempts: model.Attempts,
			})
		}
		return tx.Model(&outboxEventModel{}).
			Where("id IN ?", ids).
			UpdateColumn("next_attempt_at", gorm.Expr("now() + interval '1 minute'")).Error
	})
	return events, err
}

func (r *Repository) CompleteOutbox(ctx context.Context, id string) error {
	return r.query(ctx).Model(&outboxEventModel{}).
		Where("id = ?", id).
		Updates(map[string]any{"processed_at": time.Now(), "last_error": ""}).Error
}

func (r *Repository) RetryOutbox(ctx context.Context, id, reason string, next time.Time) error {
	return r.query(ctx).Model(&outboxEventModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"attempts":        gorm.Expr("attempts + 1"),
			"last_error":      reason,
			"next_attempt_at": next,
		}).Error
}
