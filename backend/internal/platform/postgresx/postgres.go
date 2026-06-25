package postgresx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"gorm.io/gorm"
)

type Store struct {
	DB       *gorm.DB
	Location *time.Location
}

type Scanner interface {
	Scan(dest ...any) error
}

func New(db *gorm.DB, location *time.Location) *Store {
	return &Store{DB: db, Location: location}
}

func (s *Store) Query(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

func (s *Store) Ping(ctx context.Context) error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (s *Store) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func MapError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return shared.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return shared.ErrConflict
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return shared.ErrInvalidInput
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return shared.ErrConflict
		case "23503", "23514", "22P02":
			return shared.ErrInvalidInput
		}
	}
	return err
}

func Initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return ""
	}
	value := []rune(strings.ToUpper(parts[0]))
	if len(parts) > 1 {
		value = append(value[:1], []rune(strings.ToUpper(parts[len(parts)-1]))[0])
	}
	if len(value) > 2 {
		value = value[:2]
	}
	return string(value)
}

func HumanTime(value, now time.Time) string {
	duration := now.Sub(value)
	switch {
	case duration < time.Minute:
		return "Baru saja"
	case duration < time.Hour:
		return fmt.Sprintf("%d menit yang lalu", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d jam yang lalu", int(duration.Hours()))
	case duration < 48*time.Hour:
		return "Kemarin"
	default:
		return value.In(now.Location()).Format("02 Jan 2006")
	}
}

func FormatRupiah(value int64) string {
	if value >= 1_000_000_000 {
		return fmt.Sprintf("Rp %.1fM", float64(value)/1_000_000_000)
	}
	if value >= 1_000_000 {
		return fmt.Sprintf("Rp %.1fJt", float64(value)/1_000_000)
	}
	return fmt.Sprintf("Rp %d", value)
}
