package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

func (r *Repository) ListDeals(ctx context.Context, organizationID string) ([]domain.Deal, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT d.id,d.organization_id,d.title,d.company,d.value,d.priority,d.stage_key,d.lost_reason,
		       COALESCE(u.id::text,''),COALESCE(trim(u.first_name || ' ' || u.last_name),''),
		       COALESCE(u.avatar_url,''),d.created_at,d.updated_at
		FROM deals d LEFT JOIN users u ON u.id=d.assignee_id AND u.organization_id=d.organization_id
		WHERE d.organization_id=$1 AND d.deleted_at IS NULL
		ORDER BY d.created_at`,
		organizationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	deals := make([]domain.Deal, 0)
	for rows.Next() {
		deal, err := scanDeal(rows)
		if err != nil {
			return nil, err
		}
		deals = append(deals, deal)
	}
	return deals, rows.Err()
}

func (r *Repository) CreateDeal(ctx context.Context, principal domain.Principal, input domain.DealInput) (domain.Deal, error) {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return domain.Deal{}, err
	}
	assigneeID := input.AssigneeID
	if assigneeID == "" {
		assigneeID = principal.UserID
	}
	var id string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO deals (organization_id,assignee_id,title,company,value,priority,stage_key,lost_reason)
		SELECT $1,u.id,$3,$4,$5,$6,$7,$8
		FROM users u WHERE u.id=$2 AND u.organization_id=$1 AND u.revoked_at IS NULL
		RETURNING id`,
		principal.OrganizationID, assigneeID, input.Title, input.Company, input.Value,
		input.Priority, input.Stage, input.LostReason,
	).Scan(&id)
	if err != nil {
		return domain.Deal{}, mapError(err)
	}
	return r.dealByID(ctx, principal.OrganizationID, id)
}

func (r *Repository) UpdateDeal(ctx context.Context, principal domain.Principal, id string, input domain.DealInput) (domain.Deal, error) {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return domain.Deal{}, err
	}
	assigneeID := input.AssigneeID
	if assigneeID == "" {
		assigneeID = principal.UserID
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Deal{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var previousStage, title string
	if err := tx.QueryRow(ctx, `
		SELECT stage_key,title FROM deals
		WHERE id=$1 AND organization_id=$2 AND deleted_at IS NULL
		FOR UPDATE`,
		id, principal.OrganizationID,
	).Scan(&previousStage, &title); err != nil {
		return domain.Deal{}, mapError(err)
	}
	tag, err := tx.Exec(ctx, `
		UPDATE deals d
		SET title=$1,company=$2,value=$3,priority=$4,stage_key=$5,lost_reason=$6,
		    assignee_id=$7,updated_at=now()
		WHERE d.id=$8 AND d.organization_id=$9 AND d.deleted_at IS NULL
		  AND EXISTS (SELECT 1 FROM users u WHERE u.id=$7 AND u.organization_id=$9 AND u.revoked_at IS NULL)`,
		input.Title, input.Company, input.Value, input.Priority, input.Stage, input.LostReason,
		assigneeID, id, principal.OrganizationID,
	)
	if err != nil {
		return domain.Deal{}, mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.Deal{}, domain.ErrNotFound
	}
	if previousStage != input.Stage {
		if err := r.recordDealStageChange(ctx, tx, principal, input.Title, input.Stage); err != nil {
			return domain.Deal{}, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Deal{}, err
	}
	return r.dealByID(ctx, principal.OrganizationID, id)
}

func (r *Repository) UpdateDealStage(ctx context.Context, principal domain.Principal, id string, input domain.StageUpdateInput) error {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return err
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var title, previousStage string
	err = tx.QueryRow(ctx, `
		SELECT title,stage_key FROM deals
		WHERE id=$1 AND organization_id=$2 AND deleted_at IS NULL
		FOR UPDATE`,
		id, principal.OrganizationID,
	).Scan(&title, &previousStage)
	if err != nil {
		return mapError(err)
	}
	_, err = tx.Exec(ctx, `
		UPDATE deals SET stage_key=$1,lost_reason=$2,updated_at=now()
		WHERE id=$3 AND organization_id=$4 AND deleted_at IS NULL
		`,
		input.Stage, input.LostReason, id, principal.OrganizationID,
	)
	if err != nil {
		return mapError(err)
	}
	if previousStage == input.Stage {
		return tx.Commit(ctx)
	}
	if err := r.recordDealStageChange(ctx, tx, principal, title, input.Stage); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) recordDealStageChange(ctx context.Context, tx pgx.Tx, principal domain.Principal, title, stage string) error {
	action := "memindahkan deal ke " + stage
	highlight := stage == "won"
	if _, err := tx.Exec(ctx, `
		INSERT INTO activities (organization_id,actor_id,actor_name,action,target,is_highlight)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		principal.OrganizationID, principal.UserID, principal.Name, action, title, highlight,
	); err != nil {
		return err
	}
	if stage == "won" {
		message := "Deal " + title + " berhasil dimenangkan oleh " + principal.Name
		if _, err := tx.Exec(ctx, `
			INSERT INTO notifications (organization_id,user_id,title,message)
			SELECT $1,id,'Deal Won!',$2 FROM users
			WHERE organization_id=$1 AND revoked_at IS NULL AND status='Aktif'`,
			principal.OrganizationID, message,
		); err != nil {
			return err
		}
		payload, _ := json.Marshal(map[string]string{"message": message})
		if _, err := tx.Exec(ctx, `
			INSERT INTO outbox_events (organization_id,event_type,payload)
			VALUES ($1,'telegram.deal_won',$2)`,
			principal.OrganizationID, payload,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeleteDeal(ctx context.Context, principal domain.Principal, id string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE deals SET deleted_at=now(),updated_at=now()
		WHERE id=$1 AND organization_id=$2 AND deleted_at IS NULL`,
		id, principal.OrganizationID,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) dealByID(ctx context.Context, organizationID, id string) (domain.Deal, error) {
	deal, err := scanDeal(r.pool.QueryRow(ctx, `
		SELECT d.id,d.organization_id,d.title,d.company,d.value,d.priority,d.stage_key,d.lost_reason,
		       COALESCE(u.id::text,''),COALESCE(trim(u.first_name || ' ' || u.last_name),''),
		       COALESCE(u.avatar_url,''),d.created_at,d.updated_at
		FROM deals d LEFT JOIN users u ON u.id=d.assignee_id AND u.organization_id=d.organization_id
		WHERE d.id=$1 AND d.organization_id=$2 AND d.deleted_at IS NULL`,
		id, organizationID,
	))
	return deal, mapError(err)
}

func scanDeal(row scanner) (domain.Deal, error) {
	var deal domain.Deal
	err := row.Scan(
		&deal.ID, &deal.OrganizationID, &deal.Title, &deal.Company, &deal.Value,
		&deal.Priority, &deal.Stage, &deal.LostReason, &deal.Assignee.ID,
		&deal.Assignee.Name, &deal.Assignee.AvatarURL, &deal.CreatedAt, &deal.UpdatedAt,
	)
	return deal, err
}

func (r *Repository) ensureStage(ctx context.Context, organizationID, key string) error {
	var exists bool
	if err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM pipeline_stages WHERE organization_id=$1 AND key=$2)`,
		organizationID, key,
	).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return domain.ErrInvalidInput
	}
	return nil
}
