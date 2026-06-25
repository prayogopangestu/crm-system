package postgres

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func (r *Repository) ListStages(ctx context.Context, organizationID string) ([]domain.PipelineStage, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id,organization_id,key,name,color,position,is_system,created_at
		FROM pipeline_stages WHERE organization_id=$1 ORDER BY position`,
		organizationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stages := make([]domain.PipelineStage, 0)
	for rows.Next() {
		var stage domain.PipelineStage
		if err := rows.Scan(
			&stage.ID, &stage.OrganizationID, &stage.Key, &stage.Name, &stage.Color,
			&stage.Position, &stage.IsSystem, &stage.CreatedAt,
		); err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}
	return stages, rows.Err()
}

func (r *Repository) CreateStage(ctx context.Context, organizationID string, stage domain.PipelineStage) (domain.PipelineStage, error) {
	stage.Key = strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(stage.Name), "-"), "-")
	if stage.Key == "" {
		return domain.PipelineStage{}, domain.ErrInvalidInput
	}
	if stage.Color == "" {
		stage.Color = "bg-surface-variant"
	}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO pipeline_stages (organization_id,key,name,color,position)
		SELECT $1,$2,$3,$4,COALESCE(max(position),0)+1
		FROM pipeline_stages WHERE organization_id=$1
		RETURNING id,organization_id,key,name,color,position,is_system,created_at`,
		organizationID, stage.Key, stage.Name, stage.Color,
	).Scan(
		&stage.ID, &stage.OrganizationID, &stage.Key, &stage.Name, &stage.Color,
		&stage.Position, &stage.IsSystem, &stage.CreatedAt,
	)
	if err != nil {
		return domain.PipelineStage{}, mapError(err)
	}
	return stage, nil
}

func (r *Repository) ReorderStages(ctx context.Context, organizationID string, ids []string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var count int
	if err := tx.QueryRow(ctx,
		`SELECT count(*) FROM pipeline_stages WHERE organization_id=$1`,
		organizationID,
	).Scan(&count); err != nil {
		return err
	}
	if count != len(ids) {
		return domain.ErrInvalidInput
	}
	// Move positions out of the unique range first.
	if _, err := tx.Exec(ctx,
		`UPDATE pipeline_stages SET position=position+1000 WHERE organization_id=$1`,
		organizationID,
	); err != nil {
		return err
	}
	for index, id := range ids {
		tag, err := tx.Exec(ctx, `
			UPDATE pipeline_stages SET position=$1,updated_at=now()
			WHERE id=$2 AND organization_id=$3`,
			index+1, id, organizationID,
		)
		if err != nil {
			return mapError(err)
		}
		if tag.RowsAffected() != 1 {
			return domain.ErrInvalidInput
		}
	}
	return tx.Commit(ctx)
}

func (r *Repository) DeleteStage(ctx context.Context, organizationID, id string) error {
	var key string
	var system bool
	if err := r.pool.QueryRow(ctx, `
		SELECT key,is_system FROM pipeline_stages WHERE id=$1 AND organization_id=$2`,
		id, organizationID,
	).Scan(&key, &system); err != nil {
		return mapError(err)
	}
	if system {
		return domain.ErrForbidden
	}
	var used bool
	if err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM deals
			WHERE organization_id=$1 AND stage_key=$2 AND deleted_at IS NULL)`,
		organizationID, key,
	).Scan(&used); err != nil {
		return err
	}
	if used {
		return domain.ErrStageInUse
	}
	_, err := r.pool.Exec(ctx,
		`DELETE FROM pipeline_stages WHERE id=$1 AND organization_id=$2`,
		id, organizationID,
	)
	return mapError(err)
}

func (r *Repository) GetTelegram(ctx context.Context, organizationID string) (domain.TelegramIntegration, error) {
	var value domain.TelegramIntegration
	err := r.pool.QueryRow(ctx, `
		SELECT organization_id,enabled,chat_id,bot_token_encrypted,updated_at
		FROM telegram_integrations WHERE organization_id=$1`,
		organizationID,
	).Scan(&value.OrganizationID, &value.Enabled, &value.ChatID, &value.EncryptedToken, &value.UpdatedAt)
	if err != nil {
		if mapError(err) == domain.ErrNotFound {
			return domain.TelegramIntegration{OrganizationID: organizationID}, nil
		}
		return domain.TelegramIntegration{}, err
	}
	value.HasToken = value.EncryptedToken != ""
	if value.HasToken {
		value.WebhookURL = "https://api.telegram.org/bot***/sendMessage"
	}
	return value, nil
}

func (r *Repository) UpsertTelegram(ctx context.Context, organizationID string, input domain.TelegramIntegration) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO telegram_integrations
		    (organization_id,enabled,chat_id,bot_token_encrypted)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (organization_id) DO UPDATE SET
		    enabled=EXCLUDED.enabled,
		    chat_id=CASE WHEN EXCLUDED.chat_id='' THEN telegram_integrations.chat_id ELSE EXCLUDED.chat_id END,
		    bot_token_encrypted=CASE WHEN EXCLUDED.bot_token_encrypted='' THEN telegram_integrations.bot_token_encrypted ELSE EXCLUDED.bot_token_encrypted END,
		    updated_at=now()`,
		organizationID, input.Enabled, input.ChatID, input.EncryptedToken,
	)
	return mapError(err)
}

func (r *Repository) Search(ctx context.Context, organizationID, query string) (domain.SearchResult, error) {
	pattern := "%" + query + "%"
	result := domain.SearchResult{Contacts: []domain.Contact{}, Tasks: []domain.Task{}, Deals: []domain.Deal{}}
	rows, err := r.pool.Query(ctx, `
		SELECT id,organization_id,owner_id,name,email,company,role,status,avatar_url,
		       last_contacted_at,created_at,updated_at
		FROM contacts
		WHERE organization_id=$1 AND deleted_at IS NULL
		  AND (name ILIKE $2 OR email ILIKE $2 OR company ILIKE $2)
		ORDER BY updated_at DESC LIMIT 10`,
		organizationID, pattern,
	)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		item, scanErr := r.scanContact(rows)
		if scanErr != nil {
			rows.Close()
			return result, scanErr
		}
		result.Contacts = append(result.Contacts, item)
	}
	rows.Close()

	taskRows, err := r.pool.Query(ctx, `
		SELECT t.id,t.organization_id,t.title,t.company,to_char(t.due_time,'HH24:MI'),
		       to_char(t.due_date,'YYYY-MM-DD'),t.type,t.completed,t.notes,t.priority,
		       COALESCE(trim(u.first_name || ' ' || u.last_name),''),
		       COALESCE(u.id::text,''),t.created_at,t.updated_at
		FROM tasks t LEFT JOIN users u ON u.id=t.assignee_id
		WHERE t.organization_id=$1 AND t.deleted_at IS NULL
		  AND (t.title ILIKE $2 OR t.company ILIKE $2 OR t.notes ILIKE $2)
		ORDER BY t.updated_at DESC LIMIT 10`,
		organizationID, pattern,
	)
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

	dealRows, err := r.pool.Query(ctx, `
		SELECT d.id,d.organization_id,d.title,d.company,d.value,d.priority,d.stage_key,d.lost_reason,
		       COALESCE(u.id::text,''),COALESCE(trim(u.first_name || ' ' || u.last_name),''),
		       COALESCE(u.avatar_url,''),d.created_at,d.updated_at
		FROM deals d LEFT JOIN users u ON u.id=d.assignee_id
		WHERE d.organization_id=$1 AND d.deleted_at IS NULL
		  AND (d.title ILIKE $2 OR d.company ILIKE $2)
		ORDER BY d.updated_at DESC LIMIT 10`,
		organizationID, pattern,
	)
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
	rows, err := r.pool.Query(ctx, `
		SELECT id,title,message,created_at,read_at IS NOT NULL
		FROM notifications
		WHERE organization_id=$1 AND (user_id=$2 OR user_id IS NULL)
		ORDER BY created_at DESC LIMIT 100`,
		principal.OrganizationID, principal.UserID,
	)
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
	tag, err := r.pool.Exec(ctx, `
		UPDATE notifications SET read_at=COALESCE(read_at,now())
		WHERE id=$1 AND organization_id=$2 AND (user_id=$3 OR user_id IS NULL)`,
		id, principal.OrganizationID, principal.UserID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) ReadAllNotifications(ctx context.Context, principal domain.Principal) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE notifications SET read_at=COALESCE(read_at,now())
		WHERE organization_id=$1 AND (user_id=$2 OR user_id IS NULL)`,
		principal.OrganizationID, principal.UserID,
	)
	return err
}

func (r *Repository) ClaimOutbox(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	rows, err := tx.Query(ctx, `
		SELECT id,organization_id,event_type,payload,attempts
		FROM outbox_events
		WHERE processed_at IS NULL AND next_attempt_at <= now()
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	events := make([]domain.OutboxEvent, 0)
	for rows.Next() {
		var event domain.OutboxEvent
		if err := rows.Scan(&event.ID, &event.OrganizationID, &event.EventType, &event.Payload, &event.Attempts); err != nil {
			rows.Close()
			return nil, err
		}
		events = append(events, event)
	}
	rows.Close()
	for _, event := range events {
		if _, err := tx.Exec(ctx, `
			UPDATE outbox_events SET next_attempt_at=now()+interval '1 minute'
			WHERE id=$1`,
			event.ID,
		); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *Repository) CompleteOutbox(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE outbox_events SET processed_at=now(),last_error='' WHERE id=$1`,
		id,
	)
	return err
}

func (r *Repository) RetryOutbox(ctx context.Context, id, reason string, next time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE outbox_events
		SET attempts=attempts+1,last_error=$1,next_attempt_at=$2
		WHERE id=$3`,
		reason, next, id,
	)
	return err
}
