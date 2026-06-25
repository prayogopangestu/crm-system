package search

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type Repository interface {
	Search(ctx context.Context, organizationID, query string) (Result, error)
}

type Service struct {
	repository Repository
	cache      shared.CacheHelper
}

func NewService(repository Repository, cache shared.CacheHelper) *Service {
	return &Service{repository: repository, cache: cache}
}

func (s *Service) Search(ctx context.Context, principal shared.Principal, query string) (Result, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return Result{}, shared.ErrInvalidInput
	}
	sum := sha256.Sum256([]byte(strings.ToLower(query)))
	key := "crm:" + principal.OrganizationID + ":search:" + hex.EncodeToString(sum[:8])
	var result Result
	err := s.cache.Load(ctx, key, 30*time.Second, &result, func() error {
		var err error
		result, err = s.repository.Search(ctx, principal.OrganizationID, query)
		return err
	})
	return result, err
}
