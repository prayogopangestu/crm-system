package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Register(ctx context.Context, orgName string, user domain.User) (domain.User, error) {
	err := r.query(ctx).Transaction(func(tx *gorm.DB) error {
		organization := organizationModel{Name: orgName}
		if err := tx.Create(&organization).Error; err != nil {
			return mapError(err)
		}
		model := userModel{
			OrganizationID: organization.ID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Email:          strings.ToLower(user.Email),
			PasswordHash:   &user.PasswordHash,
			Role:           user.Role,
			Status:         "Aktif",
		}
		if err := tx.Create(&model).Error; err != nil {
			return mapError(err)
		}
		user = userFromModel(model)

		stages := []pipelineStageModel{
			{OrganizationID: organization.ID, Key: "lead", Name: "Lead Masuk", Color: "bg-primary-container", Position: 1, IsSystem: true},
			{OrganizationID: organization.ID, Key: "contacted", Name: "Dihubungi", Color: "bg-secondary-container", Position: 2, IsSystem: true},
			{OrganizationID: organization.ID, Key: "meeting", Name: "Meeting", Color: "bg-tertiary-container", Position: 3, IsSystem: true},
			{OrganizationID: organization.ID, Key: "negotiation", Name: "Negosiasi", Color: "bg-primary-fixed", Position: 4, IsSystem: true},
			{OrganizationID: organization.ID, Key: "won", Name: "Deal Won", Color: "bg-surface-tint", Position: 5, IsSystem: true},
			{OrganizationID: organization.ID, Key: "lost", Name: "Deal Lost", Color: "bg-error-container", Position: 6, IsSystem: true},
		}
		if err := tx.Create(&stages).Error; err != nil {
			return mapError(err)
		}

		now := time.Now().In(r.location)
		goals := make([]performanceGoalModel, 0, 3)
		for offset, goal := range []int64{1_000_000_000, 900_000_000, 900_000_000} {
			goals = append(goals, performanceGoalModel{
				OrganizationID: organization.ID,
				Month:          time.Date(now.Year(), now.Month()-time.Month(offset), 1, 0, 0, 0, 0, r.location),
				Goal:           goal,
			})
		}
		return mapError(tx.Create(&goals).Error)
	})
	return user, err
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	var model userModel
	err := r.query(ctx).
		Where("lower(email) = lower(?) AND revoked_at IS NULL", email).
		First(&model).Error
	if err != nil {
		return domain.User{}, mapError(err)
	}
	return userFromModel(model), nil
}

func (r *Repository) UserByID(ctx context.Context, organizationID, userID string) (domain.User, error) {
	var model userModel
	err := r.query(ctx).
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", userID, organizationID).
		First(&model).Error
	if err != nil {
		return domain.User{}, mapError(err)
	}
	return userFromModel(model), nil
}

func (r *Repository) UpdateProfile(ctx context.Context, principal domain.Principal, input domain.UpdateProfileInput) (domain.User, error) {
	result := r.query(ctx).Model(&userModel{}).
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", principal.UserID, principal.OrganizationID).
		Updates(map[string]any{
			"first_name": input.FirstName,
			"last_name":  input.LastName,
			"email":      strings.ToLower(input.Email),
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return domain.User{}, mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.User{}, domain.ErrNotFound
	}
	return r.UserByID(ctx, principal.OrganizationID, principal.UserID)
}

func (r *Repository) AcceptInvite(ctx context.Context, tokenHash, passwordHash string) (domain.User, error) {
	var invitation invitationModel
	err := r.query(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", tokenHash).
			First(&invitation).Error; err != nil {
			return mapError(err)
		}
		if invitation.AcceptedAt != nil {
			return domain.ErrInviteUsed
		}
		if time.Now().After(invitation.ExpiresAt) {
			return domain.ErrInviteExpired
		}
		result := tx.Model(&userModel{}).
			Where("id = ? AND organization_id = ? AND revoked_at IS NULL", invitation.UserID, invitation.OrganizationID).
			Updates(map[string]any{"password_hash": passwordHash, "status": "Aktif", "updated_at": time.Now()})
		if result.Error != nil {
			return mapError(result.Error)
		}
		if result.RowsAffected == 0 {
			return domain.ErrNotFound
		}
		return mapError(tx.Model(&invitationModel{}).
			Where("id = ?", invitation.ID).
			Update("accepted_at", time.Now()).Error)
	})
	if err != nil {
		return domain.User{}, err
	}
	return r.UserByID(ctx, invitation.OrganizationID, invitation.UserID)
}

func (r *Repository) ListTeam(ctx context.Context, organizationID string) ([]domain.User, error) {
	var models []userModel
	err := r.query(ctx).
		Where("organization_id = ? AND revoked_at IS NULL", organizationID).
		Order("created_at").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, 0, len(models))
	for _, model := range models {
		users = append(users, userFromModel(model))
	}
	return users, nil
}

func (r *Repository) InviteMember(ctx context.Context, principal domain.Principal, user domain.User, invitation domain.Invitation) (domain.User, error) {
	err := r.query(ctx).Transaction(func(tx *gorm.DB) error {
		model := userModel{
			OrganizationID: principal.OrganizationID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Email:          strings.ToLower(user.Email),
			Role:           user.Role,
			Status:         "Menunggu",
		}
		if err := tx.Create(&model).Error; err != nil {
			return mapError(err)
		}
		user = userFromModel(model)
		if err := tx.Create(&invitationModel{
			OrganizationID: principal.OrganizationID,
			UserID:         model.ID,
			Email:          strings.ToLower(user.Email),
			Role:           user.Role,
			TokenHash:      invitation.TokenHash,
			ExpiresAt:      invitation.ExpiresAt,
		}).Error; err != nil {
			return mapError(err)
		}
		return mapError(tx.Create(&activityModel{
			OrganizationID: principal.OrganizationID,
			ActorID:        principal.UserID,
			ActorName:      principal.Name,
			Action:         "mengundang anggota tim",
			Target:         user.Email,
		}).Error)
	})
	return user, err
}

func (r *Repository) RevokeMember(ctx context.Context, organizationID, userID string) error {
	now := time.Now()
	result := r.query(ctx).Model(&userModel{}).
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", userID, organizationID).
		Updates(map[string]any{"status": "Dicabut", "revoked_at": now, "updated_at": now})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
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
