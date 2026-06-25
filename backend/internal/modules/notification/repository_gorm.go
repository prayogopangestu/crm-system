package notification

import (
	"context"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"gorm.io/gorm"
)

type gormRepository struct{ store *postgresx.Store }

type model struct {
	ID             string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string  `gorm:"type:uuid"`
	UserID         *string `gorm:"type:uuid"`
	Title          string
	Message        string
	ReadAt         *time.Time
	CreatedAt      time.Time
}

func (model) TableName() string { return "notifications" }

func NewRepository(store *postgresx.Store) Repository { return &gormRepository{store: store} }

func (r *gormRepository) List(ctx context.Context, principal shared.Principal) ([]Notification, error) {
	rows, err := r.store.Query(ctx).Raw(`
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
	now := time.Now().In(r.store.Location)
	items := make([]Notification, 0)
	for rows.Next() {
		var item Notification
		if err := rows.Scan(&item.ID, &item.Title, &item.Message, &item.CreatedAt, &item.Read); err != nil {
			return nil, err
		}
		item.Time = postgresx.HumanTime(item.CreatedAt, now)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *gormRepository) Read(ctx context.Context, principal shared.Principal, id string) error {
	result := r.store.Query(ctx).Model(&model{}).
		Where("id = ? AND organization_id = ? AND (user_id = ? OR user_id IS NULL)", id, principal.OrganizationID, principal.UserID).
		UpdateColumn("read_at", gorm.Expr("COALESCE(read_at, now())"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (r *gormRepository) ReadAll(ctx context.Context, principal shared.Principal) error {
	return r.store.Query(ctx).Model(&model{}).
		Where("organization_id = ? AND (user_id = ? OR user_id IS NULL)", principal.OrganizationID, principal.UserID).
		UpdateColumn("read_at", gorm.Expr("COALESCE(read_at, now())")).Error
}
