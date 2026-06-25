package contact

import (
	"context"
	"net/mail"
	"strings"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

var statuses = map[string]bool{
	"Negosiasi": true, "Menang": true, "Prospek Awal": true,
	"Proposal": true, "Kalah": true, "Kualifikasi": true,
}

type Repository interface {
	List(ctx context.Context, organizationID, search, status string, page Page) (List, error)
	Create(ctx context.Context, principal shared.Principal, input Input) (Contact, error)
	Update(ctx context.Context, principal shared.Principal, id string, input Input) (Contact, error)
	Delete(ctx context.Context, principal shared.Principal, id string) error
}

type Service struct {
	repository Repository
	cache      shared.CacheHelper
}

func NewService(repository Repository, cache shared.CacheHelper) *Service {
	return &Service{repository: repository, cache: cache}
}

func (s *Service) List(ctx context.Context, principal shared.Principal, search, status string, page, limit int) (List, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if status != "" && !statuses[status] {
		return List{}, shared.ErrInvalidInput
	}
	return s.repository.List(ctx, principal.OrganizationID, search, status, Page{Page: page, Limit: limit})
}

func (s *Service) Create(ctx context.Context, principal shared.Principal, input Input) (Contact, error) {
	if err := validate(input); err != nil {
		return Contact{}, err
	}
	result, err := s.repository.Create(ctx, principal, input)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return result, err
}

func (s *Service) Update(ctx context.Context, principal shared.Principal, id string, input Input) (Contact, error) {
	if err := validate(input); err != nil {
		return Contact{}, err
	}
	result, err := s.repository.Update(ctx, principal, id, input)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return result, err
}

func (s *Service) Delete(ctx context.Context, principal shared.Principal, id string) error {
	err := s.repository.Delete(ctx, principal, id)
	if err == nil {
		s.cache.InvalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func validate(input Input) error {
	if len(strings.TrimSpace(input.Name)) < 2 || len(strings.TrimSpace(input.Company)) < 2 || !statuses[input.Status] {
		return shared.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return shared.ErrInvalidInput
	}
	return nil
}
