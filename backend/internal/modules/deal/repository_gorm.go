package deal

import (
	"context"
	"encoding/json"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct {
	store *postgresx.Store
}

type model struct {
	ID             string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string  `gorm:"type:uuid"`
	AssigneeID     *string `gorm:"type:uuid"`
	Title          string
	Company        string
	Value          int64
	Priority       string
	StageKey       string
	LostReason     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

func (model) TableName() string { return "deals" }

type activityModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	ActorID        string `gorm:"type:uuid"`
	ActorName      string
	Action         string
	Target         string
	IsHighlight    bool
	CreatedAt      time.Time
}

func (activityModel) TableName() string { return "activities" }

type outboxModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	EventType      string
	Payload        []byte    `gorm:"type:jsonb"`
	NextAttemptAt  time.Time `gorm:"default:now()"`
	CreatedAt      time.Time
}

func (outboxModel) TableName() string { return "outbox_events" }

const selectDeal = `
	SELECT d.id,d.organization_id,d.title,d.company,d.value,d.priority,d.stage_key,d.lost_reason,
	       COALESCE(u.id::text,''),COALESCE(trim(u.first_name || ' ' || u.last_name),''),
	       COALESCE(u.avatar_url,''),d.created_at,d.updated_at
	FROM deals d LEFT JOIN users u ON u.id=d.assignee_id AND u.organization_id=d.organization_id`

func NewRepository(store *postgresx.Store) Repository {
	return &gormRepository{store: store}
}

func (r *gormRepository) List(ctx context.Context, organizationID string) ([]Deal, error) {
	rows, err := r.store.Query(ctx).Raw(
		selectDeal+` WHERE d.organization_id = ? AND d.deleted_at IS NULL ORDER BY d.created_at`,
		organizationID,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Deal, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *gormRepository) Create(ctx context.Context, principal shared.Principal, input Input) (Deal, error) {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return Deal{}, err
	}
	assigneeID := input.AssigneeID
	if assigneeID == "" {
		assigneeID = principal.UserID
	}
	if err := r.ensureUser(ctx, principal.OrganizationID, assigneeID); err != nil {
		return Deal{}, err
	}
	record := model{
		OrganizationID: principal.OrganizationID, AssigneeID: &assigneeID,
		Title: input.Title, Company: input.Company, Value: input.Value,
		Priority: input.Priority, StageKey: input.Stage, LostReason: input.LostReason,
	}
	if err := r.store.Query(ctx).Create(&record).Error; err != nil {
		return Deal{}, postgresx.MapError(err)
	}
	return r.byID(ctx, principal.OrganizationID, record.ID)
}

func (r *gormRepository) Update(ctx context.Context, principal shared.Principal, id string, input Input) (Deal, error) {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return Deal{}, err
	}
	assigneeID := input.AssigneeID
	if assigneeID == "" {
		assigneeID = principal.UserID
	}
	if err := r.ensureUser(ctx, principal.OrganizationID, assigneeID); err != nil {
		return Deal{}, err
	}
	err := r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		var current model
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
			First(&current).Error; err != nil {
			return postgresx.MapError(err)
		}
		if err := tx.Model(&current).Updates(map[string]any{
			"title": input.Title, "company": input.Company, "value": input.Value,
			"priority": input.Priority, "stage_key": input.Stage,
			"lost_reason": input.LostReason, "assignee_id": assigneeID, "updated_at": time.Now(),
		}).Error; err != nil {
			return postgresx.MapError(err)
		}
		if current.StageKey != input.Stage {
			return r.recordStageChange(tx, principal, input.Title, input.Stage)
		}
		return nil
	})
	if err != nil {
		return Deal{}, err
	}
	return r.byID(ctx, principal.OrganizationID, id)
}

func (r *gormRepository) UpdateStage(ctx context.Context, principal shared.Principal, id string, input StageInput) error {
	if err := r.ensureStage(ctx, principal.OrganizationID, input.Stage); err != nil {
		return err
	}
	return r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		var current model
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
			First(&current).Error; err != nil {
			return postgresx.MapError(err)
		}
		if err := tx.Model(&current).Updates(map[string]any{
			"stage_key": input.Stage, "lost_reason": input.LostReason, "updated_at": time.Now(),
		}).Error; err != nil {
			return postgresx.MapError(err)
		}
		if current.StageKey == input.Stage {
			return nil
		}
		return r.recordStageChange(tx, principal, current.Title, input.Stage)
	})
}

func (r *gormRepository) Delete(ctx context.Context, principal shared.Principal, id string) error {
	now := time.Now()
	result := r.store.Query(ctx).Model(&model{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(map[string]any{"deleted_at": now, "updated_at": now})
	if result.Error != nil {
		return postgresx.MapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (r *gormRepository) recordStageChange(tx *gorm.DB, principal shared.Principal, title, stage string) error {
	if err := tx.Create(&activityModel{
		OrganizationID: principal.OrganizationID, ActorID: principal.UserID,
		ActorName: principal.Name, Action: "memindahkan deal ke " + stage,
		Target: title, IsHighlight: stage == "won",
	}).Error; err != nil {
		return postgresx.MapError(err)
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
		return postgresx.MapError(err)
	}
	payload, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		return err
	}
	return postgresx.MapError(tx.Create(&outboxModel{
		OrganizationID: principal.OrganizationID, EventType: "telegram.deal_won", Payload: payload,
	}).Error)
}

func (r *gormRepository) byID(ctx context.Context, organizationID, id string) (Deal, error) {
	item, err := scan(r.store.Query(ctx).Raw(
		selectDeal+` WHERE d.id = ? AND d.organization_id = ? AND d.deleted_at IS NULL`,
		id, organizationID,
	).Row())
	return item, postgresx.MapError(err)
}

func scan(row postgresx.Scanner) (Deal, error) {
	var item Deal
	err := row.Scan(
		&item.ID, &item.OrganizationID, &item.Title, &item.Company, &item.Value,
		&item.Priority, &item.Stage, &item.LostReason, &item.Assignee.ID,
		&item.Assignee.Name, &item.Assignee.AvatarURL, &item.CreatedAt, &item.UpdatedAt,
	)
	return item, err
}

func (r *gormRepository) ensureStage(ctx context.Context, organizationID, key string) error {
	var count int64
	if err := r.store.Query(ctx).Table("pipeline_stages").
		Where("organization_id = ? AND key = ?", organizationID, key).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return shared.ErrInvalidInput
	}
	return nil
}

func (r *gormRepository) ensureUser(ctx context.Context, organizationID, userID string) error {
	var count int64
	if err := r.store.Query(ctx).Table("users").
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", userID, organizationID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return shared.ErrInvalidInput
	}
	return nil
}
