package postgres

import (
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

type organizationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (organizationModel) TableName() string { return "organizations" }

type userModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid;index"`
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

type contactModel struct {
	ID              string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID  string  `gorm:"type:uuid;index"`
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

func (contactModel) TableName() string { return "contacts" }

type pipelineStageModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid;index"`
	Key            string
	Name           string
	Color          string
	Position       int
	IsSystem       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (pipelineStageModel) TableName() string { return "pipeline_stages" }

type dealModel struct {
	ID             string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string  `gorm:"type:uuid;index"`
	AssigneeID     *string `gorm:"type:uuid"`
	Title          string
	Company        string
	Value          int64
	Priority       string
	StageKey       string
	LostReason     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

func (dealModel) TableName() string { return "deals" }

type taskModel struct {
	ID             string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string  `gorm:"type:uuid;index"`
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

func (taskModel) TableName() string { return "tasks" }

type activityModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid;index"`
	ActorID        string `gorm:"type:uuid"`
	ActorName      string
	Action         string
	Target         string
	IsHighlight    bool
	CreatedAt      time.Time
}

func (activityModel) TableName() string { return "activities" }

type performanceGoalModel struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string    `gorm:"type:uuid;index"`
	Month          time.Time `gorm:"type:date"`
	Goal           int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (performanceGoalModel) TableName() string { return "performance_goals" }

type telegramIntegrationModel struct {
	OrganizationID    string `gorm:"type:uuid;primaryKey"`
	BotTokenEncrypted string
	ChatID            string
	Enabled           bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (telegramIntegrationModel) TableName() string { return "telegram_integrations" }

type notificationModel struct {
	ID             string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string  `gorm:"type:uuid;index"`
	UserID         *string `gorm:"type:uuid"`
	Title          string
	Message        string
	ReadAt         *time.Time
	CreatedAt      time.Time
}

func (notificationModel) TableName() string { return "notifications" }

type outboxEventModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string `gorm:"type:uuid;index"`
	EventType      string
	Payload        []byte `gorm:"type:jsonb"`
	Attempts       int
	NextAttemptAt  time.Time `gorm:"default:now()"`
	ProcessedAt    *time.Time
	LastError      string
	CreatedAt      time.Time
}

func (outboxEventModel) TableName() string { return "outbox_events" }

func userFromModel(value userModel) domain.User {
	user := domain.User{
		ID: value.ID, OrganizationID: value.OrganizationID,
		FirstName: value.FirstName, LastName: value.LastName,
		Email: value.Email, PasswordHash: stringValue(value.PasswordHash),
		Role: value.Role, Status: value.Status, AvatarURL: value.AvatarURL,
		CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}
	user.Name = joinName(user.FirstName, user.LastName)
	user.Initials = initials(user.Name)
	return user
}

func contactFromModel(value contactModel) domain.Contact {
	return domain.Contact{
		ID: value.ID, OrganizationID: value.OrganizationID, OwnerID: stringValue(value.OwnerID),
		Name: value.Name, Email: value.Email, Company: value.Company, Role: value.Role,
		Status: value.Status, AvatarURL: value.AvatarURL, LastContactedAt: value.LastContactedAt,
		CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func stageFromModel(value pipelineStageModel) domain.PipelineStage {
	return domain.PipelineStage{
		ID: value.ID, OrganizationID: value.OrganizationID, Key: value.Key,
		Name: value.Name, Color: value.Color, Position: value.Position,
		IsSystem: value.IsSystem, CreatedAt: value.CreatedAt,
	}
}
