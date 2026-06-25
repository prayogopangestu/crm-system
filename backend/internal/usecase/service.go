package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	"github.com/prayogopangestu/crm-system/backend/pkg/cryptoutil"
)

type TelegramSender interface {
	Send(ctx context.Context, token, chatID, message string) error
}

type Service struct {
	repo       domain.Repository
	cache      domain.Cache
	tokens     *auth.Manager
	cipher     *cryptoutil.Cipher
	telegram   TelegramSender
	logger     *slog.Logger
	location   *time.Location
	baseURL    string
	bcryptCost int
}

func New(
	repo domain.Repository,
	cache domain.Cache,
	tokens *auth.Manager,
	cipher *cryptoutil.Cipher,
	telegram TelegramSender,
	logger *slog.Logger,
	location *time.Location,
	baseURL string,
	bcryptCost int,
) *Service {
	return &Service{
		repo: repo, cache: cache, tokens: tokens, cipher: cipher, telegram: telegram,
		logger: logger, location: location, baseURL: strings.TrimRight(baseURL, "/"),
		bcryptCost: bcryptCost,
	}
}

func (s *Service) invalidateCRM(ctx context.Context, organizationID string) {
	if s.cache == nil {
		return
	}
	for _, pattern := range []string{
		fmt.Sprintf("crm:%s:dashboard:*", organizationID),
		fmt.Sprintf("crm:%s:reports:*", organizationID),
		fmt.Sprintf("crm:%s:search:*", organizationID),
	} {
		if err := s.cache.DeletePattern(ctx, pattern); err != nil {
			s.logger.Warn("cache invalidation failed", "pattern", pattern, "error", err)
		}
	}
}

func (s *Service) invalidateProfile(ctx context.Context, organizationID, userID string) {
	if s.cache == nil {
		return
	}
	if err := s.cache.DeletePattern(ctx, fmt.Sprintf("crm:%s:profile:%s", organizationID, userID)); err != nil {
		s.logger.Warn("profile cache invalidation failed", "error", err)
	}
}

func (s *Service) cached(ctx context.Context, key string, ttl time.Duration, dst any, loader func() error) error {
	if s.cache != nil {
		hit, err := s.cache.GetJSON(ctx, key, dst)
		if err != nil {
			s.logger.Warn("cache read failed", "key", key, "error", err)
		} else if hit {
			return nil
		}
	}
	if err := loader(); err != nil {
		return err
	}
	if s.cache != nil {
		if err := s.cache.SetJSON(ctx, key, dst, ttl); err != nil {
			s.logger.Warn("cache write failed", "key", key, "error", err)
		}
	}
	return nil
}

func requireAdmin(principal domain.Principal) error {
	if principal.Role != domain.RoleAdmin {
		return domain.ErrForbidden
	}
	return nil
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
