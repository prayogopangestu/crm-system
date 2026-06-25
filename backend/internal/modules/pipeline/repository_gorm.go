package pipeline

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"gorm.io/gorm"
)

type gormRepository struct{ store *postgresx.Store }

type model struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	Key            string
	Name           string
	Color          string
	Position       int
	IsSystem       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (model) TableName() string { return "pipeline_stages" }

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func NewRepository(store *postgresx.Store) Repository { return &gormRepository{store: store} }

func (r *gormRepository) List(ctx context.Context, organizationID string) ([]Stage, error) {
	var records []model
	if err := r.store.Query(ctx).Where("organization_id = ?", organizationID).Order("position").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]Stage, 0, len(records))
	for _, record := range records {
		items = append(items, toEntity(record))
	}
	return items, nil
}

func (r *gormRepository) Create(ctx context.Context, organizationID string, stage Stage) (Stage, error) {
	stage.Key = strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(stage.Name), "-"), "-")
	if stage.Key == "" {
		return Stage{}, shared.ErrInvalidInput
	}
	if stage.Color == "" {
		stage.Color = "bg-surface-variant"
	}
	var record model
	err := r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		var maxPosition int
		if err := tx.Model(&model{}).Select("COALESCE(MAX(position), 0)").
			Where("organization_id = ?", organizationID).Scan(&maxPosition).Error; err != nil {
			return err
		}
		record = model{
			OrganizationID: organizationID, Key: stage.Key, Name: stage.Name,
			Color: stage.Color, Position: maxPosition + 1,
		}
		return postgresx.MapError(tx.Create(&record).Error)
	})
	return toEntity(record), err
}

func (r *gormRepository) Reorder(ctx context.Context, organizationID string, ids []string) error {
	return r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model{}).Where("organization_id = ?", organizationID).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(ids)) {
			return shared.ErrInvalidInput
		}
		if err := tx.Model(&model{}).Where("organization_id = ?", organizationID).
			UpdateColumn("position", gorm.Expr("position + 1000")).Error; err != nil {
			return postgresx.MapError(err)
		}
		for index, id := range ids {
			result := tx.Model(&model{}).Where("id = ? AND organization_id = ?", id, organizationID).
				Updates(map[string]any{"position": index + 1, "updated_at": time.Now()})
			if result.Error != nil {
				return postgresx.MapError(result.Error)
			}
			if result.RowsAffected != 1 {
				return shared.ErrInvalidInput
			}
		}
		return nil
	})
}

func (r *gormRepository) Delete(ctx context.Context, organizationID, id string) error {
	var record model
	if err := r.store.Query(ctx).Where("id = ? AND organization_id = ?", id, organizationID).First(&record).Error; err != nil {
		return postgresx.MapError(err)
	}
	if record.IsSystem {
		return shared.ErrForbidden
	}
	var used int64
	if err := r.store.Query(ctx).Table("deals").
		Where("organization_id = ? AND stage_key = ? AND deleted_at IS NULL", organizationID, record.Key).
		Count(&used).Error; err != nil {
		return err
	}
	if used > 0 {
		return shared.ErrStageInUse
	}
	result := r.store.Query(ctx).Where("id = ? AND organization_id = ?", id, organizationID).Delete(&model{})
	if result.Error != nil {
		return postgresx.MapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func toEntity(record model) Stage {
	return Stage{
		ID: record.ID, OrganizationID: record.OrganizationID, Key: record.Key,
		Name: record.Name, Color: record.Color, Position: record.Position,
		IsSystem: record.IsSystem, CreatedAt: record.CreatedAt,
	}
}
