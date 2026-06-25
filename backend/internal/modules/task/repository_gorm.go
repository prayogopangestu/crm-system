package task

import (
	"context"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
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
	DueDate        time.Time `gorm:"type:date"`
	DueTime        string    `gorm:"type:time"`
	Type           string
	Priority       string
	Notes          string
	Completed      bool
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

func (model) TableName() string { return "tasks" }

const selectTask = `
	t.id,t.organization_id,t.title,t.company,to_char(t.due_time,'HH24:MI'),
	to_char(t.due_date,'YYYY-MM-DD'),t.type,t.completed,t.notes,t.priority,
	COALESCE(trim(u.first_name || ' ' || u.last_name),''),
	COALESCE(u.id::text,''),t.created_at,t.updated_at`

func NewRepository(store *postgresx.Store) Repository {
	return &gormRepository{store: store}
}

func (r *gormRepository) List(ctx context.Context, organizationID, date, status string, location *time.Location) ([]Task, error) {
	now := time.Now().In(location)
	query := r.store.Query(ctx).Table("tasks AS t").Select(selectTask).
		Joins("LEFT JOIN users u ON u.id = t.assignee_id AND u.organization_id = t.organization_id").
		Where("t.organization_id = ? AND t.deleted_at IS NULL", organizationID)
	if date != "" {
		query = query.Where("t.due_date = ?::date", date)
	} else {
		today := now.Format("2006-01-02")
		switch status {
		case "overdue":
			query = query.Where("t.due_date < ?::date AND t.completed = false", today)
		case "today":
			query = query.Where("t.due_date = ?::date", today)
		case "upcoming":
			query = query.Where("t.due_date > ?::date", today)
		}
	}
	rows, err := query.Order("t.due_date,t.due_time").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Task, 0)
	for rows.Next() {
		item, err := scan(rows, now)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *gormRepository) Create(ctx context.Context, principal shared.Principal, input Input) (Task, error) {
	assigneeID, err := r.resolveAssignee(ctx, principal, input.AssigneeID, input.Assignee)
	if err != nil {
		return Task{}, err
	}
	dueDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return Task{}, shared.ErrInvalidInput
	}
	record := model{
		OrganizationID: principal.OrganizationID, AssigneeID: &assigneeID,
		Title: input.Title, Company: input.Company, DueDate: dueDate, DueTime: input.Time,
		Type: input.Type, Priority: input.Priority, Notes: input.Notes, Completed: input.Completed,
	}
	if input.Completed {
		now := time.Now()
		record.CompletedAt = &now
	}
	if err := r.store.Query(ctx).Create(&record).Error; err != nil {
		return Task{}, postgresx.MapError(err)
	}
	return r.byID(ctx, principal.OrganizationID, record.ID)
}

func (r *gormRepository) Update(ctx context.Context, principal shared.Principal, id string, input Input) (Task, error) {
	updates := map[string]any{"updated_at": time.Now()}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Company != "" {
		updates["company"] = input.Company
	}
	if input.Date != "" {
		dueDate, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			return Task{}, shared.ErrInvalidInput
		}
		updates["due_date"] = dueDate
	}
	if input.Time != "" {
		updates["due_time"] = input.Time
	}
	if input.Type != "" {
		updates["type"] = input.Type
	}
	if input.Priority != "" {
		updates["priority"] = input.Priority
	}
	if input.Notes != "" {
		updates["notes"] = input.Notes
	}
	if input.AssigneeID != "" || input.Assignee != "" {
		assigneeID, err := r.resolveAssignee(ctx, principal, input.AssigneeID, input.Assignee)
		if err != nil {
			return Task{}, err
		}
		updates["assignee_id"] = assigneeID
	}
	result := r.store.Query(ctx).Model(&model{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(updates)
	if result.Error != nil {
		return Task{}, postgresx.MapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return Task{}, shared.ErrNotFound
	}
	return r.byID(ctx, principal.OrganizationID, id)
}

func (r *gormRepository) Toggle(ctx context.Context, principal shared.Principal, id string, completed bool) error {
	var completedAt any
	if completed {
		completedAt = time.Now()
	}
	result := r.store.Query(ctx).Model(&model{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(map[string]any{"completed": completed, "completed_at": completedAt, "updated_at": time.Now()})
	if result.Error != nil {
		return postgresx.MapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}
	return nil
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

func (r *gormRepository) byID(ctx context.Context, organizationID, id string) (Task, error) {
	item, err := scan(r.store.Query(ctx).Table("tasks AS t").Select(selectTask).
		Joins("LEFT JOIN users u ON u.id = t.assignee_id AND u.organization_id = t.organization_id").
		Where("t.id = ? AND t.organization_id = ? AND t.deleted_at IS NULL", id, organizationID).Row(),
		time.Now().In(r.store.Location))
	return item, postgresx.MapError(err)
}

func scan(row postgresx.Scanner, now time.Time) (Task, error) {
	var item Task
	err := row.Scan(
		&item.ID, &item.OrganizationID, &item.Title, &item.Company, &item.Time,
		&item.Date, &item.Type, &item.Completed, &item.Notes, &item.Priority,
		&item.Assignee, &item.AssigneeID, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return Task{}, err
	}
	due, _ := time.ParseInLocation("2006-01-02", item.Date, now.Location())
	switch {
	case due.Before(dayStart(now)) && !item.Completed:
		item.Status = "overdue"
	case due.After(dayStart(now)):
		item.Status = "upcoming"
	default:
		item.Status = "today"
	}
	return item, nil
}

func (r *gormRepository) resolveAssignee(ctx context.Context, principal shared.Principal, id, name string) (string, error) {
	if id == "" && name == "" {
		return principal.UserID, nil
	}
	var resolved string
	query := r.store.Query(ctx).Table("users").Select("id").
		Where("organization_id = ? AND revoked_at IS NULL", principal.OrganizationID)
	if id != "" {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("lower(trim(first_name || ' ' || last_name)) = lower(?)", name)
	}
	if err := query.Row().Scan(&resolved); err != nil {
		return "", postgresx.MapError(err)
	}
	return resolved, nil
}

func dayStart(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}
