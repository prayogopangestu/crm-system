package usecase

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, input domain.RegisterInput) (domain.User, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.CompanyName = strings.TrimSpace(input.CompanyName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if len(input.Name) < 2 || len(input.CompanyName) < 2 || len(input.Password) < 6 {
		return domain.User{}, domain.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return domain.User{}, domain.ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		return domain.User{}, err
	}
	first, last := splitName(input.Name)
	return s.users.Register(ctx, input.CompanyName, domain.User{
		FirstName: first, LastName: last, Email: input.Email,
		PasswordHash: string(hash), Role: domain.RoleAdmin,
	})
}

func (s *Service) Login(ctx context.Context, input domain.LoginInput) (domain.LoginResult, error) {
	user, err := s.users.UserByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.LoginResult{}, domain.ErrUnauthorized
		}
		return domain.LoginResult{}, err
	}
	if user.Status != "Aktif" || user.PasswordHash == "" {
		return domain.LoginResult{}, domain.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return domain.LoginResult{}, domain.ErrUnauthorized
	}
	token, err := s.tokens.Create(user)
	if err != nil {
		return domain.LoginResult{}, err
	}
	return domain.LoginResult{Token: token, User: domain.User{
		ID: user.ID, Name: user.Name, Role: user.Role,
	}}, nil
}

func (s *Service) AcceptInvite(ctx context.Context, token, password string) (domain.User, error) {
	if token == "" || len(password) < 6 {
		return domain.User{}, domain.ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return domain.User{}, err
	}
	return s.users.AcceptInvite(ctx, tokenHash(token), string(hash))
}

func (s *Service) Profile(ctx context.Context, principal domain.Principal) (domain.User, error) {
	var user domain.User
	key := "crm:" + principal.OrganizationID + ":profile:" + principal.UserID
	err := s.cached(ctx, key, 5*time.Minute, &user, func() error {
		var err error
		user, err = s.users.UserByID(ctx, principal.OrganizationID, principal.UserID)
		return err
	})
	return user, err
}

func (s *Service) UpdateProfile(ctx context.Context, principal domain.Principal, input domain.UpdateProfileInput) (domain.User, error) {
	if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
		return domain.User{}, domain.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return domain.User{}, domain.ErrInvalidInput
	}
	user, err := s.users.UpdateProfile(ctx, principal, input)
	if err == nil {
		s.invalidateProfile(ctx, principal.OrganizationID, principal.UserID)
	}
	return user, err
}

func (s *Service) ListTeam(ctx context.Context, principal domain.Principal) ([]domain.User, error) {
	if err := requireAdmin(principal); err != nil {
		return nil, err
	}
	return s.users.ListTeam(ctx, principal.OrganizationID)
}

func (s *Service) InviteMember(ctx context.Context, principal domain.Principal, input domain.InviteInput) (domain.InviteResult, error) {
	if err := requireAdmin(principal); err != nil {
		return domain.InviteResult{}, err
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if len(input.Name) < 2 || (input.Role != domain.RoleAdmin && input.Role != domain.RoleSales) {
		return domain.InviteResult{}, domain.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return domain.InviteResult{}, domain.ErrInvalidInput
	}
	plain, hash, err := randomToken()
	if err != nil {
		return domain.InviteResult{}, err
	}
	first, last := splitName(input.Name)
	user, err := s.users.InviteMember(ctx, principal, domain.User{
		FirstName: first, LastName: last, Email: input.Email, Role: input.Role,
	}, domain.Invitation{TokenHash: hash, ExpiresAt: time.Now().Add(72 * time.Hour)})
	if err != nil {
		return domain.InviteResult{}, err
	}
	return domain.InviteResult{
		User: user, InviteURL: s.baseURL + "/accept-invite?token=" + plain,
	}, nil
}

func (s *Service) RevokeMember(ctx context.Context, principal domain.Principal, userID string) error {
	if err := requireAdmin(principal); err != nil {
		return err
	}
	if userID == principal.UserID {
		return domain.ErrInvalidInput
	}
	err := s.users.RevokeMember(ctx, principal.OrganizationID, userID)
	if err == nil {
		s.invalidateProfile(ctx, principal.OrganizationID, userID)
	}
	return err
}
