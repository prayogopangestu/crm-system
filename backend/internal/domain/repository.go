package domain

import (
	"context"
	"time"
)

// HealthRepository contains lifecycle operations for the persistence layer.
// It is intentionally separate from feature repositories so health checks do
// not need to depend on every business capability.
type HealthRepository interface {
	Ping(ctx context.Context) error
	Close() error
}

type UserRepository interface {
	Register(ctx context.Context, orgName string, user User) (User, error)
	UserByEmail(ctx context.Context, email string) (User, error)
	UserByID(ctx context.Context, organizationID, userID string) (User, error)
	UpdateProfile(ctx context.Context, principal Principal, input UpdateProfileInput) (User, error)
	AcceptInvite(ctx context.Context, tokenHash, passwordHash string) (User, error)
	ListTeam(ctx context.Context, organizationID string) ([]User, error)
	InviteMember(ctx context.Context, principal Principal, user User, invitation Invitation) (User, error)
	RevokeMember(ctx context.Context, organizationID, userID string) error
}

type ContactRepository interface {
	ListContacts(ctx context.Context, organizationID, search, status string, page Page) (ContactList, error)
	CreateContact(ctx context.Context, principal Principal, input ContactInput) (Contact, error)
	UpdateContact(ctx context.Context, principal Principal, id string, input ContactInput) (Contact, error)
	DeleteContact(ctx context.Context, principal Principal, id string) error
}

type DealRepository interface {
	ListDeals(ctx context.Context, organizationID string) ([]Deal, error)
	CreateDeal(ctx context.Context, principal Principal, input DealInput) (Deal, error)
	UpdateDeal(ctx context.Context, principal Principal, id string, input DealInput) (Deal, error)
	UpdateDealStage(ctx context.Context, principal Principal, id string, input StageUpdateInput) error
	DeleteDeal(ctx context.Context, principal Principal, id string) error
}

type TaskRepository interface {
	ListTasks(ctx context.Context, organizationID, date, status string, location *time.Location) ([]Task, error)
	CreateTask(ctx context.Context, principal Principal, input TaskInput) (Task, error)
	UpdateTask(ctx context.Context, principal Principal, id string, input TaskInput) (Task, error)
	ToggleTask(ctx context.Context, principal Principal, id string, completed bool) error
	DeleteTask(ctx context.Context, principal Principal, id string) error
}

type AnalyticsRepository interface {
	DashboardStats(ctx context.Context, organizationID string, now time.Time) (DashboardStats, error)
	ConversionChart(ctx context.Context, organizationID string, now time.Time) ([]ConversionPoint, error)
	Activities(ctx context.Context, organizationID string, limit int) ([]Activity, error)
	Leaderboard(ctx context.Context, organizationID string, month time.Time) ([]LeaderboardEntry, error)
	LostReasons(ctx context.Context, organizationID string) ([]LostReason, error)
	Goals(ctx context.Context, organizationID string) ([]PerformanceGoal, error)
}

type PipelineRepository interface {
	ListStages(ctx context.Context, organizationID string) ([]PipelineStage, error)
	CreateStage(ctx context.Context, organizationID string, stage PipelineStage) (PipelineStage, error)
	ReorderStages(ctx context.Context, organizationID string, ids []string) error
	DeleteStage(ctx context.Context, organizationID, id string) error
}

type IntegrationRepository interface {
	GetTelegram(ctx context.Context, organizationID string) (TelegramIntegration, error)
	UpsertTelegram(ctx context.Context, organizationID string, input TelegramIntegration) error
}

type SearchRepository interface {
	Search(ctx context.Context, organizationID, query string) (SearchResult, error)
}

type NotificationRepository interface {
	ListNotifications(ctx context.Context, principal Principal) ([]Notification, error)
	ReadNotification(ctx context.Context, principal Principal, id string) error
	ReadAllNotifications(ctx context.Context, principal Principal) error
}

type OutboxRepository interface {
	ClaimOutbox(ctx context.Context, limit int) ([]OutboxEvent, error)
	CompleteOutbox(ctx context.Context, id string) error
	RetryOutbox(ctx context.Context, id, reason string, next time.Time) error
}
