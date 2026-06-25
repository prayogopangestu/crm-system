package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

func (r *Repository) ListTasks(ctx context.Context, organizationID, date, status string, location *time.Location) ([]domain.Task, error) {
	now := time.Now().In(location)
	args := []any{organizationID}
	where := []string{"t.organization_id=$1", "t.deleted_at IS NULL"}
	if date != "" {
		args = append(args, date)
		where = append(where, fmt.Sprintf("t.due_date=$%d::date", len(args)))
	} else {
		switch status {
		case "overdue":
			args = append(args, now.Format("2006-01-02"))
			where = append(where, fmt.Sprintf("t.due_date < $%d::date AND t.completed=false", len(args)))
		case "today":
			args = append(args, now.Format("2006-01-02"))
			where = append(where, fmt.Sprintf("t.due_date = $%d::date", len(args)))
		case "upcoming":
			args = append(args, now.Format("2006-01-02"))
			where = append(where, fmt.Sprintf("t.due_date > $%d::date", len(args)))
		}
	}
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT t.id,t.organization_id,t.title,t.company,to_char(t.due_time,'HH24:MI'),
		       to_char(t.due_date,'YYYY-MM-DD'),t.type,t.completed,t.notes,t.priority,
		       COALESCE(trim(u.first_name || ' ' || u.last_name),''),
		       COALESCE(u.id::text,''),t.created_at,t.updated_at
		FROM tasks t LEFT JOIN users u ON u.id=t.assignee_id AND u.organization_id=t.organization_id
		WHERE %s ORDER BY t.due_date,t.due_time`, strings.Join(where, " AND ")),
		args...,
	)
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
	var id string
	err = r.pool.QueryRow(ctx, `
		INSERT INTO tasks
		    (organization_id,assignee_id,title,company,due_date,due_time,type,priority,notes,completed,completed_at)
		VALUES ($1,$2,$3,$4,$5::date,$6::time,$7,$8,$9,$10,
		        CASE WHEN $10 THEN now() ELSE NULL END)
		RETURNING id`,
		principal.OrganizationID, assigneeID, input.Title, input.Company, input.Date,
		input.Time, input.Type, input.Priority, input.Notes, input.Completed,
	).Scan(&id)
	if err != nil {
		return domain.Task{}, mapError(err)
	}
	return r.taskByID(ctx, principal.OrganizationID, id)
}

func (r *Repository) UpdateTask(ctx context.Context, principal domain.Principal, id string, input domain.TaskInput) (domain.Task, error) {
	var assigneeID any
	if input.AssigneeID != "" || input.Assignee != "" {
		resolved, err := r.resolveAssignee(ctx, principal, input.AssigneeID, input.Assignee)
		if err != nil {
			return domain.Task{}, err
		}
		assigneeID = resolved
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE tasks SET
		    title=COALESCE(NULLIF($1,''),title),
		    company=COALESCE(NULLIF($2,''),company),
		    due_date=COALESCE(NULLIF($3,'')::date,due_date),
		    due_time=COALESCE(NULLIF($4,'')::time,due_time),
		    type=COALESCE(NULLIF($5,''),type),
		    priority=COALESCE(NULLIF($6,''),priority),
		    notes=CASE WHEN $7='' THEN notes ELSE $7 END,
		    assignee_id=COALESCE($8::uuid,assignee_id),
		    updated_at=now()
		WHERE id=$9 AND organization_id=$10 AND deleted_at IS NULL`,
		input.Title, input.Company, input.Date, input.Time, input.Type, input.Priority,
		input.Notes, assigneeID, id, principal.OrganizationID,
	)
	if err != nil {
		return domain.Task{}, mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.Task{}, domain.ErrNotFound
	}
	return r.taskByID(ctx, principal.OrganizationID, id)
}

func (r *Repository) ToggleTask(ctx context.Context, principal domain.Principal, id string, completed bool) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE tasks SET completed=$1,completed_at=CASE WHEN $1 THEN now() ELSE NULL END,updated_at=now()
		WHERE id=$2 AND organization_id=$3 AND deleted_at IS NULL`,
		completed, id, principal.OrganizationID,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, principal domain.Principal, id string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE tasks SET deleted_at=now(),updated_at=now()
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

func (r *Repository) taskByID(ctx context.Context, organizationID, id string) (domain.Task, error) {
	task, err := scanTask(r.pool.QueryRow(ctx, `
		SELECT t.id,t.organization_id,t.title,t.company,to_char(t.due_time,'HH24:MI'),
		       to_char(t.due_date,'YYYY-MM-DD'),t.type,t.completed,t.notes,t.priority,
		       COALESCE(trim(u.first_name || ' ' || u.last_name),''),
		       COALESCE(u.id::text,''),t.created_at,t.updated_at
		FROM tasks t LEFT JOIN users u ON u.id=t.assignee_id AND u.organization_id=t.organization_id
		WHERE t.id=$1 AND t.organization_id=$2 AND t.deleted_at IS NULL`,
		id, organizationID,
	), time.Now().In(r.location))
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
	var resolved string
	var err error
	if id != "" {
		err = r.pool.QueryRow(ctx, `
			SELECT id FROM users
			WHERE id=$1 AND organization_id=$2 AND revoked_at IS NULL`,
			id, principal.OrganizationID,
		).Scan(&resolved)
	} else {
		err = r.pool.QueryRow(ctx, `
			SELECT id FROM users
			WHERE organization_id=$1 AND lower(trim(first_name || ' ' || last_name))=lower($2)
			  AND revoked_at IS NULL LIMIT 1`,
			principal.OrganizationID, name,
		).Scan(&resolved)
	}
	if err != nil {
		return "", mapError(err)
	}
	return resolved, nil
}

func dayStart(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}
