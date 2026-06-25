package usecase

import (
	"context"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, input domain.RegisterInput) (domain.User, error)
	Login(ctx context.Context, input domain.LoginInput) (domain.LoginResult, error)
	AcceptInvite(ctx context.Context, token, password string) (domain.User, error)
	Profile(ctx context.Context, principal domain.Principal) (domain.User, error)
	UpdateProfile(ctx context.Context, principal domain.Principal, input domain.UpdateProfileInput) (domain.User, error)
	ListTeam(ctx context.Context, principal domain.Principal) ([]domain.User, error)
	InviteMember(ctx context.Context, principal domain.Principal, input domain.InviteInput) (domain.InviteResult, error)
	RevokeMember(ctx context.Context, principal domain.Principal, userID string) error
}

type ContactUseCase interface {
	ListContacts(ctx context.Context, principal domain.Principal, search, status string, page, limit int) (domain.ContactList, error)
	CreateContact(ctx context.Context, principal domain.Principal, input domain.ContactInput) (domain.Contact, error)
	UpdateContact(ctx context.Context, principal domain.Principal, id string, input domain.ContactInput) (domain.Contact, error)
	DeleteContact(ctx context.Context, principal domain.Principal, id string) error
}

type DealUseCase interface {
	ListDeals(ctx context.Context, principal domain.Principal) ([]domain.Deal, error)
	CreateDeal(ctx context.Context, principal domain.Principal, input domain.DealInput) (domain.Deal, error)
	UpdateDeal(ctx context.Context, principal domain.Principal, id string, input domain.DealInput) (domain.Deal, error)
	UpdateDealStage(ctx context.Context, principal domain.Principal, id string, input domain.StageUpdateInput) error
	DeleteDeal(ctx context.Context, principal domain.Principal, id string) error
}

type TaskUseCase interface {
	ListTasks(ctx context.Context, principal domain.Principal, date, status string) ([]domain.Task, error)
	CreateTask(ctx context.Context, principal domain.Principal, input domain.TaskInput) (domain.Task, error)
	UpdateTask(ctx context.Context, principal domain.Principal, id string, input domain.TaskInput) (domain.Task, error)
	ToggleTask(ctx context.Context, principal domain.Principal, id string, completed bool) error
	DeleteTask(ctx context.Context, principal domain.Principal, id string) error
}

type AnalyticsUseCase interface {
	DashboardStats(ctx context.Context, principal domain.Principal) (domain.DashboardStats, error)
	ConversionChart(ctx context.Context, principal domain.Principal) ([]domain.ConversionPoint, error)
	Activities(ctx context.Context, principal domain.Principal, limit int) ([]domain.Activity, error)
	Leaderboard(ctx context.Context, principal domain.Principal, period string) ([]domain.LeaderboardEntry, error)
	LostReasons(ctx context.Context, principal domain.Principal) ([]domain.LostReason, error)
	Goals(ctx context.Context, principal domain.Principal) ([]domain.PerformanceGoal, error)
	ExportCSV(ctx context.Context, principal domain.Principal) ([]byte, error)
	ExportPDF(ctx context.Context, principal domain.Principal) ([]byte, error)
}

type SettingsUseCase interface {
	ListStages(ctx context.Context, principal domain.Principal) ([]domain.PipelineStage, error)
	CreateStage(ctx context.Context, principal domain.Principal, input domain.StageInput) (domain.PipelineStage, error)
	ReorderStages(ctx context.Context, principal domain.Principal, ids []string) error
	DeleteStage(ctx context.Context, principal domain.Principal, id string) error
	GetTelegram(ctx context.Context, principal domain.Principal) (domain.TelegramIntegration, error)
	UpdateTelegram(ctx context.Context, principal domain.Principal, input domain.TelegramInput) (domain.TelegramIntegration, error)
	TestTelegram(ctx context.Context, principal domain.Principal) error
}

type SearchUseCase interface {
	Search(ctx context.Context, principal domain.Principal, query string) (domain.SearchResult, error)
}

type NotificationUseCase interface {
	Notifications(ctx context.Context, principal domain.Principal) ([]domain.Notification, error)
	ReadNotification(ctx context.Context, principal domain.Principal, id string) error
	ReadAllNotifications(ctx context.Context, principal domain.Principal) error
}
