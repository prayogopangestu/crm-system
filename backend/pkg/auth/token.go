package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
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

func (m *Manager) Create(user domain.User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
		Name:           user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	})
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(raw string) (domain.Principal, error) {
	parsed, err := jwt.ParseWithClaims(raw, &claims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil || !parsed.Valid {
		return domain.Principal{}, domain.ErrUnauthorized
	}
	value, ok := parsed.Claims.(*claims)
	if !ok || value.Subject == "" || value.OrganizationID == "" {
		return domain.Principal{}, domain.ErrUnauthorized
	}
	return domain.Principal{
		UserID:         value.Subject,
		OrganizationID: value.OrganizationID,
		Role:           value.Role,
		Name:           value.Name,
	}, nil
}
