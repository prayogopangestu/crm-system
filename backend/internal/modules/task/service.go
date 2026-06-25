package task

import (
	"context"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type Repository interface {
	List(ctx context.Context, organizationID, date, status string, location *time.Location) ([]Task, error)
	Create(ctx context.Context, principal shared.Principal, input Input) (Task, error)
	Update(ctx context.Context, principal shared.Principal, id string, input Input) (Task, error)
	Toggle(ctx context.Context, principal shared.Principal, id string, completed bool) error
	Delete(ctx context.Context, principal shared.Principal, id string) error
}

type Service struct {
	repository Repository
	cache      shared.CacheHelper
	location   *time.Location
}

func NewService(repository Repository, cache shared.CacheHelper, location *time.Location) *Service {
	return &Service{repository: repository, cache: cache, location: location}
}

func (s *Service) List(ctx context.Context, principal shared.Principal, date, status string) ([]Task, error) {
	if date != "" {
		if _, err := time.ParseInLocation("2006-01-02", date, s.location); err != nil {
			return nil, shared.ErrInvalidInput
		}
	}
	if status != "" && status != "overdue" && status != "today" && status != "upcoming" {
		return nil, shared.ErrInvalidInput
	}
	return s.repository.List(ctx, principal.OrganizationID, date, status, s.location)
}

func (s *Service) Create(ctx context.Context, principal shared.Principal, input Input) (Task, error) {
	if err := validate(input, false); err != nil {
		return Task{}, err
	}
	result, err := s.repository.Create(ctx, principal, input)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return result, err
}

func (s *Service) Update(ctx context.Context, principal shared.Principal, id string, input Input) (Task, error) {
	if err := validate(input, true); err != nil {
		return Task{}, err
	}
	result, err := s.repository.Update(ctx, principal, id, input)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return result, err
}

func (s *Service) Toggle(ctx context.Context, principal shared.Principal, id string, completed bool) error {
	err := s.repository.Toggle(ctx, principal, id, completed)
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

func validate(input Input, partial bool) error {
	if !partial && (len(strings.TrimSpace(input.Title)) < 3 || len(strings.TrimSpace(input.Company)) < 2 ||
		input.Date == "" || input.Time == "" || input.Type == "" || input.Priority == "") {
		return shared.ErrInvalidInput
	}
	if input.Date != "" {
		if _, err := time.Parse("2006-01-02", input.Date); err != nil {
			return shared.ErrInvalidInput
		}
	}
	if input.Time != "" {
		if _, err := time.Parse("15:04", input.Time); err != nil {
			return shared.ErrInvalidInput
		}
	}
	if input.Type != "" && input.Type != "Meeting" && input.Type != "Call" && input.Type != "Proposal" && input.Type != "Other" {
		return shared.ErrInvalidInput
	}
	if input.Priority != "" && input.Priority != "Tinggi" && input.Priority != "Sedang" && input.Priority != "Rendah" {
		return shared.ErrInvalidInput
	}
	return nil
}
