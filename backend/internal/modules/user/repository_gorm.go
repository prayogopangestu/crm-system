package user

import (
	"context"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct {
	store *postgresx.Store
}

type organizationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (organizationModel) TableName() string { return "organizations" }

type userModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	FirstName      string
	LastName       string
	Email          string
	PasswordHash   *string
	Role           string
	Status         string
	AvatarURL      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RevokedAt      *time.Time
}

func (userModel) TableName() string { return "users" }

type invitationModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	UserID         string `gorm:"type:uuid"`
	Email          string
	Role           string
	TokenHash      string
	ExpiresAt      time.Time
	AcceptedAt     *time.Time
	CreatedAt      time.Time
}

func (invitationModel) TableName() string { return "user_invitations" }

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

type stageModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid"`
	Key            string
	Name           string
	Color          string
	Position       int
	IsSystem       bool
}

func (stageModel) TableName() string { return "pipeline_stages" }

type goalModel struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string    `gorm:"type:uuid"`
	Month          time.Time `gorm:"type:date"`
	Goal           int64
}

func (goalModel) TableName() string { return "performance_goals" }

func NewRepository(store *postgresx.Store) Repository {
	return &gormRepository{store: store}
}

func (r *gormRepository) Register(ctx context.Context, orgName string, value User) (User, error) {
	err := r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		organization := organizationModel{Name: orgName}
		if err := tx.Create(&organization).Error; err != nil {
			return postgresx.MapError(err)
		}
		passwordHash := value.PasswordHash
		record := userModel{
			OrganizationID: organization.ID, FirstName: value.FirstName, LastName: value.LastName,
			Email: strings.ToLower(value.Email), PasswordHash: &passwordHash, Role: value.Role, Status: "Aktif",
		}
		if err := tx.Create(&record).Error; err != nil {
			return postgresx.MapError(err)
		}
		value = toEntity(record)
		stages := []stageModel{
			{OrganizationID: organization.ID, Key: "lead", Name: "Lead Masuk", Color: "bg-primary-container", Position: 1, IsSystem: true},
			{OrganizationID: organization.ID, Key: "contacted", Name: "Dihubungi", Color: "bg-secondary-container", Position: 2, IsSystem: true},
			{OrganizationID: organization.ID, Key: "meeting", Name: "Meeting", Color: "bg-tertiary-container", Position: 3, IsSystem: true},
			{OrganizationID: organization.ID, Key: "negotiation", Name: "Negosiasi", Color: "bg-primary-fixed", Position: 4, IsSystem: true},
			{OrganizationID: organization.ID, Key: "won", Name: "Deal Won", Color: "bg-surface-tint", Position: 5, IsSystem: true},
			{OrganizationID: organization.ID, Key: "lost", Name: "Deal Lost", Color: "bg-error-container", Position: 6, IsSystem: true},
		}
		if err := tx.Create(&stages).Error; err != nil {
			return postgresx.MapError(err)
		}
		now := time.Now().In(r.store.Location)
		goals := make([]goalModel, 0, 3)
		for offset, target := range []int64{1_000_000_000, 900_000_000, 900_000_000} {
			goals = append(goals, goalModel{
				OrganizationID: organization.ID,
				Month:          time.Date(now.Year(), now.Month()-time.Month(offset), 1, 0, 0, 0, 0, r.store.Location),
				Goal:           target,
			})
		}
		return postgresx.MapError(tx.Create(&goals).Error)
	})
	return value, err
}

func (r *gormRepository) ByEmail(ctx context.Context, email string) (User, error) {
	var record userModel
	err := r.store.Query(ctx).Where("lower(email) = lower(?) AND revoked_at IS NULL", email).First(&record).Error
	if err != nil {
		return User{}, postgresx.MapError(err)
	}
	return toEntity(record), nil
}

func (r *gormRepository) ByID(ctx context.Context, organizationID, userID string) (User, error) {
	var record userModel
	err := r.store.Query(ctx).
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", userID, organizationID).
		First(&record).Error
	if err != nil {
		return User{}, postgresx.MapError(err)
	}
	return toEntity(record), nil
}

func (r *gormRepository) UpdateProfile(ctx context.Context, principal shared.Principal, input UpdateProfileInput) (User, error) {
	result := r.store.Query(ctx).Model(&userModel{}).
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", principal.UserID, principal.OrganizationID).
		Updates(map[string]any{
			"first_name": input.FirstName, "last_name": input.LastName,
			"email": strings.ToLower(input.Email), "updated_at": time.Now(),
		})
	if result.Error != nil {
		return User{}, postgresx.MapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return User{}, shared.ErrNotFound
	}
	return r.ByID(ctx, principal.OrganizationID, principal.UserID)
}

func (r *gormRepository) AcceptInvite(ctx context.Context, tokenHash, passwordHash string) (User, error) {
	var invitation invitationModel
	err := r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", tokenHash).First(&invitation).Error; err != nil {
			return postgresx.MapError(err)
		}
		if invitation.AcceptedAt != nil {
			return shared.ErrInviteUsed
		}
		if time.Now().After(invitation.ExpiresAt) {
			return shared.ErrInviteExpired
		}
		result := tx.Model(&userModel{}).
			Where("id = ? AND organization_id = ? AND revoked_at IS NULL", invitation.UserID, invitation.OrganizationID).
			Updates(map[string]any{"password_hash": passwordHash, "status": "Aktif", "updated_at": time.Now()})
		if result.Error != nil {
			return postgresx.MapError(result.Error)
		}
		if result.RowsAffected == 0 {
			return shared.ErrNotFound
		}
		return postgresx.MapError(tx.Model(&invitationModel{}).
			Where("id = ?", invitation.ID).Update("accepted_at", time.Now()).Error)
	})
	if err != nil {
		return User{}, err
	}
	return r.ByID(ctx, invitation.OrganizationID, invitation.UserID)
}

func (r *gormRepository) ListTeam(ctx context.Context, organizationID string) ([]User, error) {
	var records []userModel
	if err := r.store.Query(ctx).
		Where("organization_id = ? AND revoked_at IS NULL", organizationID).
		Order("created_at").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]User, 0, len(records))
	for _, record := range records {
		items = append(items, toEntity(record))
	}
	return items, nil
}

func (r *gormRepository) InviteMember(ctx context.Context, principal shared.Principal, value User, invitation Invitation) (User, error) {
	err := r.store.Query(ctx).Transaction(func(tx *gorm.DB) error {
		record := userModel{
			OrganizationID: principal.OrganizationID, FirstName: value.FirstName, LastName: value.LastName,
			Email: strings.ToLower(value.Email), Role: value.Role, Status: "Menunggu",
		}
		if err := tx.Create(&record).Error; err != nil {
			return postgresx.MapError(err)
		}
		value = toEntity(record)
		if err := tx.Create(&invitationModel{
			OrganizationID: principal.OrganizationID, UserID: record.ID,
			Email: strings.ToLower(value.Email), Role: value.Role,
			TokenHash: invitation.TokenHash, ExpiresAt: invitation.ExpiresAt,
		}).Error; err != nil {
			return postgresx.MapError(err)
		}
		return postgresx.MapError(tx.Create(&activityModel{
			OrganizationID: principal.OrganizationID, ActorID: principal.UserID,
			ActorName: principal.Name, Action: "mengundang anggota tim", Target: value.Email,
		}).Error)
	})
	return value, err
}

func (r *gormRepository) RevokeMember(ctx context.Context, organizationID, userID string) error {
	now := time.Now()
	result := r.store.Query(ctx).Model(&userModel{}).
		Where("id = ? AND organization_id = ? AND revoked_at IS NULL", userID, organizationID).
		Updates(map[string]any{"status": "Dicabut", "revoked_at": now, "updated_at": now})
	if result.Error != nil {
		return postgresx.MapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func toEntity(record userModel) User {
	passwordHash := ""
	if record.PasswordHash != nil {
		passwordHash = *record.PasswordHash
	}
	value := User{
		ID: record.ID, OrganizationID: record.OrganizationID,
		FirstName: record.FirstName, LastName: record.LastName, Email: record.Email,
		PasswordHash: passwordHash, Role: record.Role, Status: record.Status,
		AvatarURL: record.AvatarURL, CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt,
	}
	value.Name = strings.TrimSpace(value.FirstName + " " + value.LastName)
	value.Initials = postgresx.Initials(value.Name)
	return value
}
