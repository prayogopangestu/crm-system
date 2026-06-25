package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const dealSelect = `
	SELECT d.id,d.organization_id,d.title,d.company,d.value,d.priority,d.stage_key,d.lost_reason,
	       COALESCE(u.id::text,''),COALESCE(trim(u.first_name || ' ' || u.last_name),''),
	       COALESCE(u.avatar_url,''),d.created_at,d.updated_at
	FROM deals d LEFT JOIN users u ON u.id=d.assignee_id AND u.organization_id=d.organization_id`

func (r *Repository) ListDeals(ctx context.Context, organizationID string) ([]domain.Deal, error) {
	rows, err := r.query(ctx).Raw(
		dealSelect+` WHERE d.organization_id = ? AND d.deleted_at IS NULL ORDER BY d.created_at`,
		organizationID,
	).Rows()
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
	if err := r.ensureUser(ctx, principal.OrganizationID, assigneeID); err != nil {
		return domain.Deal{}, err
	}
	model := dealModel{
		OrganizationID: principal.OrganizationID,
		AssigneeID:     &assigneeID,
		Title:          input.Title,
		Company:        input.Company,
		Value:          input.Value,
		Priority:       input.Priority,
		StageKey:       input.Stage,
		LostReason:     input.LostReason,
	}
	if err := r.query(ctx).Create(&model).Error; err != nil {
		return domain.Deal{}, mapError(err)
	}
	return r.dealByID(ctx, principal.OrganizationID, model.ID)
}

func (r *Repository) UpdateDeal(ctx context.Context, principal domain.Principal, id string, input domain.DealInput) (domain.Deal, error) {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return domain.Deal{}, err
	}
	assigneeID := input.AssigneeID
	if assigneeID == "" {
		assigneeID = principal.UserID
	}
	if err := r.ensureUser(ctx, principal.OrganizationID, assigneeID); err != nil {
		return domain.Deal{}, err
	}
	err := r.query(ctx).Transaction(func(tx *gorm.DB) error {
		var current dealModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
			First(&current).Error; err != nil {
			return mapError(err)
		}
		result := tx.Model(&current).Updates(map[string]any{
			"title":       input.Title,
			"company":     input.Company,
			"value":       input.Value,
			"priority":    input.Priority,
			"stage_key":   input.Stage,
			"lost_reason": input.LostReason,
			"assignee_id": assigneeID,
			"updated_at":  time.Now(),
		})
		if result.Error != nil {
			return mapError(result.Error)
		}
		if current.StageKey != input.Stage {
			return r.recordDealStageChange(tx, principal, input.Title, input.Stage)
		}
		return nil
	})
	if err != nil {
		return domain.Deal{}, err
	}
	return r.dealByID(ctx, principal.OrganizationID, id)
}

func (r *Repository) UpdateDealStage(ctx context.Context, principal domain.Principal, id string, input domain.StageUpdateInput) error {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return err
	}
	return r.query(ctx).Transaction(func(tx *gorm.DB) error {
		var current dealModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
			First(&current).Error; err != nil {
			return mapError(err)
		}
		if err := tx.Model(&current).Updates(map[string]any{
			"stage_key":   input.Stage,
			"lost_reason": input.LostReason,
			"updated_at":  time.Now(),
		}).Error; err != nil {
			return mapError(err)
		}
		if current.StageKey == input.Stage {
			return nil
		}
		return r.recordDealStageChange(tx, principal, current.Title, input.Stage)
	})
}

func (r *Repository) recordDealStageChange(tx *gorm.DB, principal domain.Principal, title, stage string) error {
	action := "memindahkan deal ke " + stage
	highlight := stage == "won"
	if err := tx.Create(&activityModel{
		OrganizationID: principal.OrganizationID,
		ActorID:        principal.UserID,
		ActorName:      principal.Name,
		Action:         action,
		Target:         title,
		IsHighlight:    highlight,
	}).Error; err != nil {
		return mapError(err)
	}
	if stage != "won" {
		return nil
	}
	message := "Deal " + title + " berhasil dimenangkan oleh " + principal.Name
	if err := tx.Exec(`
		INSERT INTO notifications (organization_id,user_id,title,message)
		SELECT ?,id,'Deal Won!',? FROM users
		WHERE organization_id = ? AND revoked_at IS NULL AND status = 'Aktif'`,
		principal.OrganizationID, message, principal.OrganizationID,
	).Error; err != nil {
		return mapError(err)
	}
	payload, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		return err
	}
	return mapError(tx.Create(&outboxEventModel{
		OrganizationID: principal.OrganizationID,
		EventType:      "telegram.deal_won",
		Payload:        payload,
	}).Error)
}

func (r *Repository) DeleteDeal(ctx context.Context, principal domain.Principal, id string) error {
	now := time.Now()
	result := r.query(ctx).Model(&dealModel{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(map[string]any{"deleted_at": now, "updated_at": now})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) dealByID(ctx context.Context, organizationID, id string) (domain.Deal, error) {
	row := r.query(ctx).Raw(
		dealSelect+` WHERE d.id = ? AND d.organization_id = ? AND d.deleted_at IS NULL`,
		id, organizationID,
	).Row()
	deal, err := scanDeal(row)
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
	var count int64
	if err := r.query(ctx).Model(&pipelineStageModel{}).
		Where("organization_id = ? AND key = ?", organizationID, key).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrInvalidInput
	}
	return nil
}

func (r *Repository) ensureUser(ctx context.Context, organizationID, userID string) error {
	var count int64
	if err := r.query(ctx).Model(&userModel{}).
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", userID, organizationID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrInvalidInput
	}
	return nil
}
