package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

type scanner interface {
	Scan(dest ...any) error
}

func (r *Repository) ListContacts(ctx context.Context, organizationID, search, status string, page domain.Page) (domain.ContactList, error) {
	args := []any{organizationID}
	where := []string{"organization_id=$1", "deleted_at IS NULL"}
	if search != "" {
		args = append(args, "%"+search+"%")
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d OR company ILIKE $%d)", len(args), len(args), len(args)))
	}
	if status != "" {
		args = append(args, status)
		where = append(where, fmt.Sprintf("status=$%d", len(args)))
	}
	var total int64
	if err := r.pool.QueryRow(ctx,
		"SELECT count(*) FROM contacts WHERE "+strings.Join(where, " AND "),
		args...,
	).Scan(&total); err != nil {
		return domain.ContactList{}, err
	}
	args = append(args, page.Limit, (page.Page-1)*page.Limit)
	query := fmt.Sprintf(`
		SELECT id, organization_id, owner_id, name, email, company, role, status,
		       avatar_url, last_contacted_at, created_at, updated_at
		FROM contacts WHERE %s
		ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "), len(args)-1, len(args),
	)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return domain.ContactList{}, err
	}
	defer rows.Close()
	contacts := make([]domain.Contact, 0)
	for rows.Next() {
		contact, err := r.scanContact(rows)
		if err != nil {
			return domain.ContactList{}, err
		}
		contacts = append(contacts, contact)
	}
	return domain.ContactList{Data: contacts, Total: total, Page: page.Page}, rows.Err()
}

func (r *Repository) CreateContact(ctx context.Context, principal domain.Principal, input domain.ContactInput) (domain.Contact, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Contact{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var contact domain.Contact
	err = tx.QueryRow(ctx, `
		INSERT INTO contacts (organization_id,owner_id,name,email,company,role,status,avatar_url)
		VALUES ($1,$2,$3,lower($4),$5,$6,$7,$8)
		RETURNING id, organization_id, owner_id, name, email, company, role, status,
		          avatar_url, last_contacted_at, created_at, updated_at`,
		principal.OrganizationID, principal.UserID, input.Name, input.Email, input.Company,
		input.Role, input.Status, input.AvatarURL,
	).Scan(
		&contact.ID, &contact.OrganizationID, &contact.OwnerID, &contact.Name, &contact.Email,
		&contact.Company, &contact.Role, &contact.Status, &contact.AvatarURL,
		&contact.LastContactedAt, &contact.CreatedAt, &contact.UpdatedAt,
	)
	if err != nil {
		return domain.Contact{}, mapError(err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO activities (organization_id,actor_id,actor_name,action,target)
		VALUES ($1,$2,$3,'menambahkan kontak baru',$4)`,
		principal.OrganizationID, principal.UserID, principal.Name, contact.Name,
	); err != nil {
		return domain.Contact{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Contact{}, err
	}
	r.decorateContact(&contact)
	return contact, nil
}

func (r *Repository) UpdateContact(ctx context.Context, principal domain.Principal, id string, input domain.ContactInput) (domain.Contact, error) {
	var contact domain.Contact
	err := r.pool.QueryRow(ctx, `
		UPDATE contacts
		SET name=$1,email=lower($2),company=$3,role=$4,status=$5,avatar_url=$6,
		    last_contacted_at=now(),updated_at=now()
		WHERE id=$7 AND organization_id=$8 AND deleted_at IS NULL
		RETURNING id, organization_id, owner_id, name, email, company, role, status,
		          avatar_url, last_contacted_at, created_at, updated_at`,
		input.Name, input.Email, input.Company, input.Role, input.Status, input.AvatarURL,
		id, principal.OrganizationID,
	).Scan(
		&contact.ID, &contact.OrganizationID, &contact.OwnerID, &contact.Name, &contact.Email,
		&contact.Company, &contact.Role, &contact.Status, &contact.AvatarURL,
		&contact.LastContactedAt, &contact.CreatedAt, &contact.UpdatedAt,
	)
	if err != nil {
		return domain.Contact{}, mapError(err)
	}
	r.decorateContact(&contact)
	return contact, nil
}

func (r *Repository) DeleteContact(ctx context.Context, principal domain.Principal, id string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE contacts SET deleted_at=now(),updated_at=now()
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

func (r *Repository) scanContact(row scanner) (domain.Contact, error) {
	var contact domain.Contact
	err := row.Scan(
		&contact.ID, &contact.OrganizationID, &contact.OwnerID, &contact.Name, &contact.Email,
		&contact.Company, &contact.Role, &contact.Status, &contact.AvatarURL,
		&contact.LastContactedAt, &contact.CreatedAt, &contact.UpdatedAt,
	)
	if err != nil {
		return domain.Contact{}, mapError(err)
	}
	r.decorateContact(&contact)
	return contact, nil
}

func (r *Repository) decorateContact(contact *domain.Contact) {
	now := time.Now().In(r.location)
	contact.LastContacted = humanTime(contact.LastContactedAt, now)
	contact.Initials = initials(contact.Name)
}
