package postgres

import (
	"context"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

const taskSelect = `
	t.id,t.organization_id,t.title,t.company,to_char(t.due_time,'HH24:MI') AS due_time,
	to_char(t.due_date,'YYYY-MM-DD') AS due_date,t.type,t.completed,t.notes,t.priority,
	COALESCE(trim(u.first_name || ' ' || u.last_name),'') AS assignee,
	COALESCE(u.id::text,'') AS assignee_id,t.created_at,t.updated_at`

func (r *Repository) ListTasks(ctx context.Context, organizationID, date, status string, location *time.Location) ([]domain.Task, error) {
	now := time.Now().In(location)
	query := r.query(ctx).
		Table("tasks AS t").
		Select(taskSelect).
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
	tasks := make([]domain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows, now)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *Repository) CreateTask(ctx context.Context, principal domain.Principal, input domain.TaskInput) (domain.Task, error) {
	assigneeID, err := r.resolveAssignee(ctx, principal, input.AssigneeID, input.Assignee)
	if err != nil {
		return domain.Task{}, err
	}
	dueDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return domain.Task{}, domain.ErrInvalidInput
	}
	model := taskModel{
		OrganizationID: principal.OrganizationID,
		AssigneeID:     &assigneeID,
		Title:          input.Title,
		Company:        input.Company,
		DueDate:        dueDate,
		DueTime:        input.Time,
		Type:           input.Type,
		Priority:       input.Priority,
		Notes:          input.Notes,
		Completed:      input.Completed,
	}
	if input.Completed {
		now := time.Now()
		model.CompletedAt = &now
	}
	if err := r.query(ctx).Create(&model).Error; err != nil {
		return domain.Task{}, mapError(err)
	}
	return r.taskByID(ctx, principal.OrganizationID, model.ID)
}

func (r *Repository) UpdateTask(ctx context.Context, principal domain.Principal, id string, input domain.TaskInput) (domain.Task, error) {
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
			return domain.Task{}, domain.ErrInvalidInput
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
			return domain.Task{}, err
		}
		updates["assignee_id"] = assigneeID
	}
	result := r.query(ctx).Model(&taskModel{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(updates)
	if result.Error != nil {
		return domain.Task{}, mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.Task{}, domain.ErrNotFound
	}
	return r.taskByID(ctx, principal.OrganizationID, id)
}

func (r *Repository) ToggleTask(ctx context.Context, principal domain.Principal, id string, completed bool) error {
	var completedAt any
	if completed {
		completedAt = time.Now()
	}
	result := r.query(ctx).Model(&taskModel{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(map[string]any{
			"completed":    completed,
			"completed_at": completedAt,
			"updated_at":   time.Now(),
		})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, principal domain.Principal, id string) error {
	now := time.Now()
	result := r.query(ctx).Model(&taskModel{}).
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

func (r *Repository) taskByID(ctx context.Context, organizationID, id string) (domain.Task, error) {
	row := r.query(ctx).
		Table("tasks AS t").
		Select(taskSelect).
		Joins("LEFT JOIN users u ON u.id = t.assignee_id AND u.organization_id = t.organization_id").
		Where("t.id = ? AND t.organization_id = ? AND t.deleted_at IS NULL", id, organizationID).
		Row()
	task, err := scanTask(row, time.Now().In(r.location))
	return task, mapError(err)
}

func scanTask(row scanner, now time.Time) (domain.Task, error) {
	var task domain.Task
	err := row.Scan(
		&task.ID, &task.OrganizationID, &task.Title, &task.Company, &task.Time,
		&task.Date, &task.Type, &task.Completed, &task.Notes, &task.Priority,
		&task.Assignee, &task.AssigneeID, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return domain.Task{}, err
	}
	due, _ := time.ParseInLocation("2006-01-02", task.Date, now.Location())
	switch {
	case due.Before(dayStart(now)) && !task.Completed:
		task.Status = "overdue"
	case due.After(dayStart(now)):
		task.Status = "upcoming"
	default:
		task.Status = "today"
	}
	return task, nil
}

func (r *Repository) resolveAssignee(ctx context.Context, principal domain.Principal, id, name string) (string, error) {
	if id == "" && name == "" {
		return principal.UserID, nil
	}
	var user userModel
	query := r.query(ctx).Select("id").
		Where("organization_id = ? AND revoked_at IS NULL", principal.OrganizationID)
	if id != "" {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("lower(trim(first_name || ' ' || last_name)) = lower(?)", name)
	}
	if err := query.First(&user).Error; err != nil {
		return "", mapError(err)
	}
	return user.ID, nil
}

func dayStart(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}
