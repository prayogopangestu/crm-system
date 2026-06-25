package pipeline

import (
	"context"
	"strings"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type Repository interface {
	List(ctx context.Context, organizationID string) ([]Stage, error)
	Create(ctx context.Context, organizationID string, stage Stage) (Stage, error)
	Reorder(ctx context.Context, organizationID string, ids []string) error
	Delete(ctx context.Context, organizationID, id string) error
}

type Service struct {
	repository Repository
	cache      shared.CacheHelper
}

func NewService(repository Repository, cache shared.CacheHelper) *Service {
	return &Service{repository: repository, cache: cache}
}

func (s *Service) List(ctx context.Context, principal shared.Principal) ([]Stage, error) {
	return s.repository.List(ctx, principal.OrganizationID)
}

func (s *Service) Create(ctx context.Context, principal shared.Principal, input Input) (Stage, error) {
	if err := shared.RequireAdmin(principal); err != nil {
		return Stage{}, err
	}
	if len(strings.TrimSpace(input.Name)) < 2 {
		return Stage{}, shared.ErrInvalidInput
	}
	result, err := s.repository.Create(ctx, principal.OrganizationID, Stage{Name: input.Name, Color: input.Color})
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return result, err
}

func (s *Service) Reorder(ctx context.Context, principal shared.Principal, ids []string) error {
	if err := shared.RequireAdmin(principal); err != nil {
		return err
	}
	if len(ids) == 0 {
		return shared.ErrInvalidInput
	}
	return s.repository.Reorder(ctx, principal.OrganizationID, ids)
}

func (s *Service) Delete(ctx context.Context, principal shared.Principal, id string) error {
	if err := shared.RequireAdmin(principal); err != nil {
		return err
	}
	return s.repository.Delete(ctx, principal.OrganizationID, id)
}
