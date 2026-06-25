package deal

import (
	"context"
	"strings"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type Repository interface {
	List(ctx context.Context, organizationID string) ([]Deal, error)
	Create(ctx context.Context, principal shared.Principal, input Input) (Deal, error)
	Update(ctx context.Context, principal shared.Principal, id string, input Input) (Deal, error)
	UpdateStage(ctx context.Context, principal shared.Principal, id string, input StageInput) error
	Delete(ctx context.Context, principal shared.Principal, id string) error
}

type Service struct {
	repository Repository
	cache      shared.CacheHelper
}

func NewService(repository Repository, cache shared.CacheHelper) *Service {
	return &Service{repository: repository, cache: cache}
}

func (s *Service) List(ctx context.Context, principal shared.Principal) ([]Deal, error) {
	return s.repository.List(ctx, principal.OrganizationID)
}

func (s *Service) Create(ctx context.Context, principal shared.Principal, input Input) (Deal, error) {
	if err := validate(input); err != nil {
		return Deal{}, err
	}
	result, err := s.repository.Create(ctx, principal, input)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return result, err
}

func (s *Service) Update(ctx context.Context, principal shared.Principal, id string, input Input) (Deal, error) {
	if err := validate(input); err != nil {
		return Deal{}, err
	}
	result, err := s.repository.Update(ctx, principal, id, input)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return result, err
}

func (s *Service) UpdateStage(ctx context.Context, principal shared.Principal, id string, input StageInput) error {
	if input.Stage == "" || (input.Stage == "lost" && strings.TrimSpace(input.LostReason) == "") {
		return shared.ErrInvalidInput
	}
	if input.Stage != "lost" {
		input.LostReason = ""
	}
	err := s.repository.UpdateStage(ctx, principal, id, input)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func (s *Service) Delete(ctx context.Context, principal shared.Principal, id string) error {
	err := s.repository.Delete(ctx, principal, id)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func validate(input Input) error {
	if len(strings.TrimSpace(input.Title)) < 2 || len(strings.TrimSpace(input.Company)) < 2 ||
		input.Value < 0 || input.Stage == "" {
		return shared.ErrInvalidInput
	}
	if input.Priority != "High" && input.Priority != "Medium" && input.Priority != "Low" {
		return shared.ErrInvalidInput
	}
	if input.Stage == "lost" && strings.TrimSpace(input.LostReason) == "" {
		return shared.ErrInvalidInput
	}
	return nil
}
