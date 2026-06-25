package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

type claims struct {
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	Name           string `json:"name"`
	jwt.RegisteredClaims
}

func New(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

func (m *Manager) Create(userID, organizationID, role, name string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		OrganizationID: organizationID,
		Role:           role,
		Name:           name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	})
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(raw string) (shared.Principal, error) {
	parsed, err := jwt.ParseWithClaims(raw, &claims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil || !parsed.Valid {
		return shared.Principal{}, shared.ErrUnauthorized
	}
	value, ok := parsed.Claims.(*claims)
	if !ok || value.Subject == "" || value.OrganizationID == "" {
		return shared.Principal{}, shared.ErrUnauthorized
	}
	return shared.Principal{
		UserID:         value.Subject,
		OrganizationID: value.OrganizationID,
		Role:           value.Role,
		Name:           value.Name,
	}, nil
}
