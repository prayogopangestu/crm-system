package usecase

import (
	"context"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

var contactStatuses = map[string]bool{
	"Negosiasi": true, "Menang": true, "Prospek Awal": true,
	"Proposal": true, "Kalah": true, "Kualifikasi": true,
}

func (s *Service) ListContacts(ctx context.Context, principal domain.Principal, search, status string, page, limit int) (domain.ContactList, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if status != "" && !contactStatuses[status] {
		return domain.ContactList{}, domain.ErrInvalidInput
	}
	return s.repo.ListContacts(ctx, principal.OrganizationID, search, status, domain.Page{Page: page, Limit: limit})
}

func (s *Service) CreateContact(ctx context.Context, principal domain.Principal, input domain.ContactInput) (domain.Contact, error) {
	if err := validateContact(input); err != nil {
		return domain.Contact{}, err
	}
	contact, err := s.repo.CreateContact(ctx, principal, input)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return contact, err
}

func (s *Service) UpdateContact(ctx context.Context, principal domain.Principal, id string, input domain.ContactInput) (domain.Contact, error) {
	if err := validateContact(input); err != nil {
		return domain.Contact{}, err
	}
	contact, err := s.repo.UpdateContact(ctx, principal, id, input)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return contact, err
}

func (s *Service) DeleteContact(ctx context.Context, principal domain.Principal, id string) error {
	err := s.repo.DeleteContact(ctx, principal, id)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func validateContact(input domain.ContactInput) error {
	if len(strings.TrimSpace(input.Name)) < 2 || len(strings.TrimSpace(input.Company)) < 2 || !contactStatuses[input.Status] {
		return domain.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return domain.ErrInvalidInput
	}
	return nil
}

func (s *Service) ListDeals(ctx context.Context, principal domain.Principal) ([]domain.Deal, error) {
	return s.repo.ListDeals(ctx, principal.OrganizationID)
}

func (s *Service) CreateDeal(ctx context.Context, principal domain.Principal, input domain.DealInput) (domain.Deal, error) {
	if err := validateDeal(input); err != nil {
		return domain.Deal{}, err
	}
	deal, err := s.repo.CreateDeal(ctx, principal, input)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return deal, err
}

func (s *Service) UpdateDeal(ctx context.Context, principal domain.Principal, id string, input domain.DealInput) (domain.Deal, error) {
	if err := validateDeal(input); err != nil {
		return domain.Deal{}, err
	}
	deal, err := s.repo.UpdateDeal(ctx, principal, id, input)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return deal, err
}

func (s *Service) UpdateDealStage(ctx context.Context, principal domain.Principal, id string, input domain.StageUpdateInput) error {
	if input.Stage == "" || (input.Stage == "lost" && strings.TrimSpace(input.LostReason) == "") {
		return domain.ErrInvalidInput
	}
	if input.Stage != "lost" {
		input.LostReason = ""
	}
	err := s.repo.UpdateDealStage(ctx, principal, id, input)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func (s *Service) DeleteDeal(ctx context.Context, principal domain.Principal, id string) error {
	err := s.repo.DeleteDeal(ctx, principal, id)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func validateDeal(input domain.DealInput) error {
	if len(strings.TrimSpace(input.Title)) < 2 || len(strings.TrimSpace(input.Company)) < 2 ||
		input.Value < 0 || input.Stage == "" {
		return domain.ErrInvalidInput
	}
	if input.Priority != "High" && input.Priority != "Medium" && input.Priority != "Low" {
		return domain.ErrInvalidInput
	}
	if input.Stage == "lost" && strings.TrimSpace(input.LostReason) == "" {
		return domain.ErrInvalidInput
	}
	return nil
}

func (s *Service) ListTasks(ctx context.Context, principal domain.Principal, date, status string) ([]domain.Task, error) {
	if date != "" {
		if _, err := time.ParseInLocation("2006-01-02", date, s.location); err != nil {
			return nil, domain.ErrInvalidInput
		}
	}
	if status != "" && status != "overdue" && status != "today" && status != "upcoming" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.ListTasks(ctx, principal.OrganizationID, date, status, s.location)
}

func (s *Service) CreateTask(ctx context.Context, principal domain.Principal, input domain.TaskInput) (domain.Task, error) {
	if err := validateTask(input, false); err != nil {
		return domain.Task{}, err
	}
	task, err := s.repo.CreateTask(ctx, principal, input)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return task, err
}

func (s *Service) UpdateTask(ctx context.Context, principal domain.Principal, id string, input domain.TaskInput) (domain.Task, error) {
	if err := validateTask(input, true); err != nil {
		return domain.Task{}, err
	}
	task, err := s.repo.UpdateTask(ctx, principal, id, input)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return task, err
}

func (s *Service) ToggleTask(ctx context.Context, principal domain.Principal, id string, completed bool) error {
	err := s.repo.ToggleTask(ctx, principal, id, completed)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func (s *Service) DeleteTask(ctx context.Context, principal domain.Principal, id string) error {
	err := s.repo.DeleteTask(ctx, principal, id)
	if err == nil {
		s.invalidateCRM(ctx, principal.OrganizationID)
	}
	return err
}

func validateTask(input domain.TaskInput, partial bool) error {
	if !partial && (len(strings.TrimSpace(input.Title)) < 3 || len(strings.TrimSpace(input.Company)) < 2 ||
		input.Date == "" || input.Time == "" || input.Type == "" || input.Priority == "") {
		return domain.ErrInvalidInput
	}
	if input.Date != "" {
		if _, err := time.Parse("2006-01-02", input.Date); err != nil {
			return domain.ErrInvalidInput
		}
	}
	if input.Time != "" {
		if _, err := time.Parse("15:04", input.Time); err != nil {
			return domain.ErrInvalidInput
		}
	}
	if input.Type != "" && input.Type != "Meeting" && input.Type != "Call" && input.Type != "Proposal" && input.Type != "Other" {
		return domain.ErrInvalidInput
	}
	if input.Priority != "" && input.Priority != "Tinggi" && input.Priority != "Sedang" && input.Priority != "Rendah" {
		return domain.ErrInvalidInput
	}
	return nil
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
