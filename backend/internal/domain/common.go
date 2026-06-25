package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource conflict")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidInput   = errors.New("invalid input")
	ErrStageInUse     = errors.New("pipeline stage is in use")
	ErrInviteExpired  = errors.New("invitation expired")
	ErrInviteUsed     = errors.New("invitation already used")
	ErrRedisRateLimit = errors.New("rate limit exceeded")
)

const (
	RoleAdmin = "Admin"
	RoleSales = "Staf Sales"
)

type Principal struct {
	UserID         string
	OrganizationID string
	Role           string
	Name           string
}

type principalKey struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok
}

type Page struct {
	Page  int
	Limit int
}

type DashboardStats struct {
	TotalLeads       int64  `json:"totalLeads"`
	LeadsTrend       string `json:"leadsTrend"`
	DealWonCount     int64  `json:"dealWonCount"`
	WonTrend         string `json:"wonTrend"`
	TotalRevenue     string `json:"totalRevenue"`
	RevenueTrend     string `json:"revenueTrend"`
	UrgentTasksCount int64  `json:"urgentTasksCount"`
}

type ConversionPoint struct {
	Name       string  `json:"name"`
	Conversion float64 `json:"Konversi"`
}

type SearchResult struct {
	Contacts []Contact `json:"contacts"`
	Tasks    []Task    `json:"tasks"`
	Deals    []Deal    `json:"deals"`
}

type OutboxEvent struct {
	ID             string
	OrganizationID string
	EventType      string
	Payload        []byte
	Attempts       int
}

type Cache interface {
	GetJSON(ctx context.Context, key string, dst any) (bool, error)
	SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error
	DeletePattern(ctx context.Context, pattern string) error
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
	Close() error
}
