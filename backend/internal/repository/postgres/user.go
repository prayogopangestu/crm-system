package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

func (r *Repository) Register(ctx context.Context, orgName string, user domain.User) (domain.User, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := tx.QueryRow(ctx,
		`INSERT INTO organizations (name) VALUES ($1) RETURNING id`,
		orgName,
	).Scan(&user.OrganizationID); err != nil {
		return domain.User{}, mapError(err)
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO users (organization_id, first_name, last_name, email, password_hash, role, status)
		VALUES ($1,$2,$3,lower($4),$5,$6,'Aktif')
		RETURNING id, created_at, updated_at`,
		user.OrganizationID, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, mapError(err)
	}

	stages := []struct {
		Key, Name, Color string
	}{
		{"lead", "Lead Masuk", "bg-primary-container"},
		{"contacted", "Dihubungi", "bg-secondary-container"},
		{"meeting", "Meeting", "bg-tertiary-container"},
		{"negotiation", "Negosiasi", "bg-primary-fixed"},
		{"won", "Deal Won", "bg-surface-tint"},
		{"lost", "Deal Lost", "bg-error-container"},
	}
	for position, stage := range stages {
		if _, err := tx.Exec(ctx, `
			INSERT INTO pipeline_stages (organization_id,key,name,color,position,is_system)
			VALUES ($1,$2,$3,$4,$5,true)`,
			user.OrganizationID, stage.Key, stage.Name, stage.Color, position+1,
		); err != nil {
			return domain.User{}, mapError(err)
		}
	}

	now := time.Now().In(r.location)
	for offset, goal := range []int64{1_000_000_000, 900_000_000, 900_000_000} {
		month := time.Date(now.Year(), now.Month()-time.Month(offset), 1, 0, 0, 0, 0, r.location)
		if _, err := tx.Exec(ctx,
			`INSERT INTO performance_goals (organization_id,month,goal) VALUES ($1,$2,$3)`,
			user.OrganizationID, month, goal,
		); err != nil {
			return domain.User{}, mapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	user.Name = joinName(user.FirstName, user.LastName)
	user.Status = "Aktif"
	user.Initials = initials(user.Name)
	return user, nil
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, organization_id, first_name, last_name, email, password_hash,
		       role, status, avatar_url, created_at, updated_at
		FROM users
		WHERE lower(email)=lower($1) AND revoked_at IS NULL`,
		email,
	).Scan(
		&user.ID, &user.OrganizationID, &user.FirstName, &user.LastName, &user.Email,
		&user.PasswordHash, &user.Role, &user.Status, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, mapError(err)
	}
	user.Name = joinName(user.FirstName, user.LastName)
	user.Initials = initials(user.Name)
	return user, nil
}

func (r *Repository) UserByID(ctx context.Context, organizationID, userID string) (domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, organization_id, first_name, last_name, email, role, status,
		       avatar_url, created_at, updated_at
		FROM users
		WHERE id=$1 AND organization_id=$2 AND revoked_at IS NULL`,
		userID, organizationID,
	).Scan(
		&user.ID, &user.OrganizationID, &user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.Status, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, mapError(err)
	}
	user.Name = joinName(user.FirstName, user.LastName)
	user.Initials = initials(user.Name)
	return user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, principal domain.Principal, input domain.UpdateProfileInput) (domain.User, error) {
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET first_name=$1,last_name=$2,email=lower($3),updated_at=now()
		WHERE id=$4 AND organization_id=$5 AND revoked_at IS NULL`,
		input.FirstName, input.LastName, input.Email, principal.UserID, principal.OrganizationID,
	)
	if err != nil {
		return domain.User{}, mapError(err)
	}
	return r.UserByID(ctx, principal.OrganizationID, principal.UserID)
}

func (r *Repository) AcceptInvite(ctx context.Context, tokenHash, passwordHash string) (domain.User, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var invitation domain.Invitation
	err = tx.QueryRow(ctx, `
		SELECT id, organization_id, user_id, email, role, token_hash, expires_at, accepted_at
		FROM user_invitations WHERE token_hash=$1 FOR UPDATE`,
		tokenHash,
	).Scan(
		&invitation.ID, &invitation.OrganizationID, &invitation.UserID, &invitation.Email,
		&invitation.Role, &invitation.TokenHash, &invitation.ExpiresAt, &invitation.AcceptedAt,
	)
	if err != nil {
		return domain.User{}, mapError(err)
	}
	if invitation.AcceptedAt != nil {
		return domain.User{}, domain.ErrInviteUsed
	}
	if time.Now().After(invitation.ExpiresAt) {
		return domain.User{}, domain.ErrInviteExpired
	}
	if _, err := tx.Exec(ctx, `
		UPDATE users SET password_hash=$1,status='Aktif',updated_at=now()
		WHERE id=$2 AND organization_id=$3 AND revoked_at IS NULL`,
		passwordHash, invitation.UserID, invitation.OrganizationID,
	); err != nil {
		return domain.User{}, mapError(err)
	}
	if _, err := tx.Exec(ctx, `UPDATE user_invitations SET accepted_at=now() WHERE id=$1`, invitation.ID); err != nil {
		return domain.User{}, mapError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	return r.UserByID(ctx, invitation.OrganizationID, invitation.UserID)
}

func (r *Repository) ListTeam(ctx context.Context, organizationID string) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, organization_id, first_name, last_name, email, role, status,
		       avatar_url, created_at, updated_at
		FROM users
		WHERE organization_id=$1 AND revoked_at IS NULL
		ORDER BY created_at`,
		organizationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID, &user.OrganizationID, &user.FirstName, &user.LastName, &user.Email,
			&user.Role, &user.Status, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		user.Name = joinName(user.FirstName, user.LastName)
		user.Initials = initials(user.Name)
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *Repository) InviteMember(ctx context.Context, principal domain.Principal, user domain.User, invitation domain.Invitation) (domain.User, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = tx.QueryRow(ctx, `
		INSERT INTO users (organization_id, first_name, last_name, email, role, status)
		VALUES ($1,$2,$3,lower($4),$5,'Menunggu')
		RETURNING id, created_at, updated_at`,
		principal.OrganizationID, user.FirstName, user.LastName, user.Email, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, mapError(err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO user_invitations (organization_id,user_id,email,role,token_hash,expires_at)
		VALUES ($1,$2,lower($3),$4,$5,$6)`,
		principal.OrganizationID, user.ID, user.Email, user.Role, invitation.TokenHash, invitation.ExpiresAt,
	); err != nil {
		return domain.User{}, mapError(err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO activities (organization_id,actor_id,actor_name,action,target)
		VALUES ($1,$2,$3,'mengundang anggota tim', $4)`,
		principal.OrganizationID, principal.UserID, principal.Name, user.Email,
	); err != nil {
		return domain.User{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	user.OrganizationID = principal.OrganizationID
	user.Status = "Menunggu"
	user.Name = joinName(user.FirstName, user.LastName)
	user.Initials = initials(user.Name)
	return user, nil
}

func (r *Repository) RevokeMember(ctx context.Context, organizationID, userID string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE users SET status='Dicabut',revoked_at=now(),updated_at=now()
		WHERE id=$1 AND organization_id=$2 AND revoked_at IS NULL`,
		userID, organizationID,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func joinName(first, last string) string {
	if last == "" {
		return first
	}
	return first + " " + last
}
