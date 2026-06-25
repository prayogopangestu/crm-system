package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/mail"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Register(ctx context.Context, orgName string, user User) (User, error)
	ByEmail(ctx context.Context, email string) (User, error)
	ByID(ctx context.Context, organizationID, userID string) (User, error)
	UpdateProfile(ctx context.Context, principal shared.Principal, input UpdateProfileInput) (User, error)
	AcceptInvite(ctx context.Context, tokenHash, passwordHash string) (User, error)
	ListTeam(ctx context.Context, organizationID string) ([]User, error)
	InviteMember(ctx context.Context, principal shared.Principal, user User, invitation Invitation) (User, error)
	RevokeMember(ctx context.Context, organizationID, userID string) error
}

type Service struct {
	repository Repository
	cache      shared.CacheHelper
	tokens     *auth.Manager
	baseURL    string
	bcryptCost int
}

func NewService(repository Repository, cache shared.CacheHelper, tokens *auth.Manager, baseURL string, bcryptCost int) *Service {
	return &Service{
		repository: repository, cache: cache, tokens: tokens,
		baseURL: strings.TrimRight(baseURL, "/"), bcryptCost: bcryptCost,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (User, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.CompanyName = strings.TrimSpace(input.CompanyName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if len(input.Name) < 2 || len(input.CompanyName) < 2 || len(input.Password) < 6 {
		return User{}, shared.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return User{}, shared.ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		return User{}, err
	}
	first, last := splitName(input.Name)
	return s.repository.Register(ctx, input.CompanyName, User{
		FirstName: first, LastName: last, Email: input.Email,
		PasswordHash: string(hash), Role: shared.RoleAdmin,
	})
}

func (s *Service) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	value, err := s.repository.ByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if err == shared.ErrNotFound {
			return LoginResult{}, shared.ErrUnauthorized
		}
		return LoginResult{}, err
	}
	if value.Status != "Aktif" || value.PasswordHash == "" {
		return LoginResult{}, shared.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(value.PasswordHash), []byte(input.Password)); err != nil {
		return LoginResult{}, shared.ErrUnauthorized
	}
	token, err := s.tokens.Create(value.ID, value.OrganizationID, value.Role, value.Name)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{Token: token, User: User{ID: value.ID, Name: value.Name, Role: value.Role}}, nil
}

func (s *Service) AcceptInvite(ctx context.Context, token, password string) (User, error) {
	if token == "" || len(password) < 6 {
		return User{}, shared.ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return User{}, err
	}
	return s.repository.AcceptInvite(ctx, tokenHash(token), string(hash))
}

func (s *Service) Profile(ctx context.Context, principal shared.Principal) (User, error) {
	var value User
	key := "crm:" + principal.OrganizationID + ":profile:" + principal.UserID
	err := s.cache.Load(ctx, key, 5*time.Minute, &value, func() error {
		var err error
		value, err = s.repository.ByID(ctx, principal.OrganizationID, principal.UserID)
		return err
	})
	return value, err
}

func (s *Service) UpdateProfile(ctx context.Context, principal shared.Principal, input UpdateProfileInput) (User, error) {
	if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
		return User{}, shared.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return User{}, shared.ErrInvalidInput
	}
	value, err := s.repository.UpdateProfile(ctx, principal, input)
	if err == nil {
		s.cache.InvalidateProfile(ctx, principal.OrganizationID, principal.UserID)
	}
	return value, err
}

func (s *Service) ListTeam(ctx context.Context, principal shared.Principal) ([]User, error) {
	if err := shared.RequireAdmin(principal); err != nil {
		return nil, err
	}
	return s.repository.ListTeam(ctx, principal.OrganizationID)
}

func (s *Service) InviteMember(ctx context.Context, principal shared.Principal, input InviteInput) (InviteResult, error) {
	if err := shared.RequireAdmin(principal); err != nil {
		return InviteResult{}, err
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if len(input.Name) < 2 || (input.Role != shared.RoleAdmin && input.Role != shared.RoleSales) {
		return InviteResult{}, shared.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return InviteResult{}, shared.ErrInvalidInput
	}
	plain, hash, err := randomToken()
	if err != nil {
		return InviteResult{}, err
	}
	first, last := splitName(input.Name)
	value, err := s.repository.InviteMember(ctx, principal, User{
		FirstName: first, LastName: last, Email: input.Email, Role: input.Role,
	}, Invitation{TokenHash: hash, ExpiresAt: time.Now().Add(72 * time.Hour)})
	if err != nil {
		return InviteResult{}, err
	}
	return InviteResult{User: value, InviteURL: s.baseURL + "/accept-invite?token=" + plain}, nil
}

func (s *Service) RevokeMember(ctx context.Context, principal shared.Principal, userID string) error {
	if err := shared.RequireAdmin(principal); err != nil {
		return err
	}
	if userID == principal.UserID {
		return shared.ErrInvalidInput
	}
	err := s.repository.RevokeMember(ctx, principal.OrganizationID, userID)
	if err == nil {
		s.cache.InvalidateProfile(ctx, principal.OrganizationID, userID)
	}
	return err
}

func splitName(name string) (string, string) {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func randomToken() (plain, hash string, err error) {
	value := make([]byte, 32)
	if _, err = rand.Read(value); err != nil {
		return "", "", err
	}
	plain = base64.RawURLEncoding.EncodeToString(value)
	sum := sha256.Sum256([]byte(plain))
	return plain, hex.EncodeToString(sum[:]), nil
}

func tokenHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
