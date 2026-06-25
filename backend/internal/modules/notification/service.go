package notification

import (
	"context"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type Repository interface {
	List(ctx context.Context, principal shared.Principal) ([]Notification, error)
	Read(ctx context.Context, principal shared.Principal, id string) error
	ReadAll(ctx context.Context, principal shared.Principal) error
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) List(ctx context.Context, principal shared.Principal) ([]Notification, error) {
	return s.repository.List(ctx, principal)
}

func (s *Service) Read(ctx context.Context, principal shared.Principal, id string) error {
	return s.repository.Read(ctx, principal, id)
}

func (s *Service) ReadAll(ctx context.Context, principal shared.Principal) error {
	return s.repository.ReadAll(ctx, principal)
}
