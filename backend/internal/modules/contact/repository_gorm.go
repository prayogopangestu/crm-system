package contact

import (
	"context"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"gorm.io/gorm"
)

type gormRepository struct {
	store *postgresx.Store
}

type model struct {
	ID              string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID  string  `gorm:"type:uuid"`
	OwnerID         *string `gorm:"type:uuid"`
	Name            string
	Email           string
	Company         string
	Role            string
	Status          string
	AvatarURL       string
	LastContactedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

func (model) TableName() string { return "contacts" }

type activityModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	ActorID        string `gorm:"type:uuid"`
	ActorName      string
	Action         string
	Target         string
	CreatedAt      time.Time
}

func (activityModel) TableName() string { return "activities" }

func NewRepository(store *postgresx.Store) Repository {
	return &gormRepository{store: store}
}

func (r *gormRepository) List(ctx context.Context, organizationID, search, status string, page Page) (List, error) {
	query := r.store.Query(ctx).Model(&model{}).Where("organization_id = ? AND deleted_at IS NULL", organizationID)
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR company ILIKE ?", pattern, pattern, pattern)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return List{}, err
	}
	var records []model
	if err := query.Order("created_at DESC").Limit(page.Limit).Offset((page.Page - 1) * page.Limit).Find(&records).Error; err != nil {
		return List{}, err
	}
	items := make([]Contact, 0, len(records))
	for _, record := range records {
		items = append(items, r.toEntity(record))
	}
	return List{Data: items, Total: total, Page: page.Page}, nil
}

func (r *gormRepository) Create(ctx context.Context, principal shared.Principal, input Input) (Contact, error) {
	var item Contact
	err := r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		ownerID := principal.UserID
		record := model{
			OrganizationID: principal.OrganizationID, OwnerID: &ownerID,
			Name: input.Name, Email: strings.ToLower(input.Email), Company: input.Company,
			Role: input.Role, Status: input.Status, AvatarURL: input.AvatarURL, LastContactedAt: time.Now(),
		}
		if err := tx.Create(&record).Error; err != nil {
			return postgresx.MapError(err)
		}
		if err := tx.Create(&activityModel{
			OrganizationID: principal.OrganizationID, ActorID: principal.UserID,
			ActorName: principal.Name, Action: "menambahkan kontak baru", Target: record.Name,
		}).Error; err != nil {
			return postgresx.MapError(err)
		}
		item = r.toEntity(record)
		return nil
	})
	return item, err
}

func (r *gormRepository) Update(ctx context.Context, principal shared.Principal, id string, input Input) (Contact, error) {
	now := time.Now()
	result := r.store.Query(ctx).Model(&model{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(map[string]any{
			"name": input.Name, "email": strings.ToLower(input.Email), "company": input.Company,
			"role": input.Role, "status": input.Status, "avatar_url": input.AvatarURL,
			"last_contacted_at": now, "updated_at": now,
		})
	if result.Error != nil {
		return Contact{}, postgresx.MapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return Contact{}, shared.ErrNotFound
	}
	var record model
	if err := r.store.Query(ctx).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		First(&record).Error; err != nil {
		return Contact{}, postgresx.MapError(err)
	}
	return r.toEntity(record), nil
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

func (r *gormRepository) toEntity(record model) Contact {
	ownerID := ""
	if record.OwnerID != nil {
		ownerID = *record.OwnerID
	}
	item := Contact{
		ID: record.ID, OrganizationID: record.OrganizationID, OwnerID: ownerID,
		Name: record.Name, Email: record.Email, Company: record.Company, Role: record.Role,
		Status: record.Status, AvatarURL: record.AvatarURL, LastContactedAt: record.LastContactedAt,
		CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt,
	}
	item.Initials = postgresx.Initials(item.Name)
	item.LastContacted = postgresx.HumanTime(item.LastContactedAt, time.Now().In(r.store.Location))
	return item
}
