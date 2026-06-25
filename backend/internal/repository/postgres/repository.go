package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

type Repository struct {
	pool     *pgxpool.Pool
	location *time.Location
}

func New(pool *pgxpool.Pool, location *time.Location) *Repository {
	return &Repository{pool: pool, location: location}
}

func (r *Repository) Ping(ctx context.Context) error { return r.pool.Ping(ctx) }
func (r *Repository) Close()                         { r.pool.Close() }

func mapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return domain.ErrConflict
		case "23503", "23514", "22P02":
			return domain.ErrInvalidInput
		}
	}
	return err
}

func initials(name string) string {
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

func humanTime(t time.Time, now time.Time) string {
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return "Baru saja"
	case d < time.Hour:
		return fmt.Sprintf("%d menit yang lalu", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d jam yang lalu", int(d.Hours()))
	case d < 48*time.Hour:
		return "Kemarin"
	default:
		return t.In(now.Location()).Format("02 Jan 2006")
	}
}

func formatRupiah(value int64) string {
	if value >= 1_000_000_000 {
		return fmt.Sprintf("Rp %.1fM", float64(value)/1_000_000_000)
	}
	if value >= 1_000_000 {
		return fmt.Sprintf("Rp %.1fJt", float64(value)/1_000_000)
	}
	return fmt.Sprintf("Rp %d", value)
}
