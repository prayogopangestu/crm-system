package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up, down, status")
	dir := flag.String("dir", "migrations", "migration directory")
	flag.Parse()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://crm:crm@localhost:5432/crm?sslmode=disable"
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := ensureTable(context.Background(), db); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	switch *direction {
	case "up":
		err = up(ctx, db, *dir)
	case "down":
		err = down(ctx, db, *dir)
	case "status":
		err = status(ctx, db, *dir)
	default:
		log.Fatalf("unknown direction %q", *direction)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func ensureTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	return err
}

func migrationFiles(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func up(ctx context.Context, db *sql.DB, dir string) error {
	files, err := migrationFiles(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		version := filepath.Base(file)
		var exists bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, version,
		).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		upSQL, _, err := readMigration(file)
		if err != nil {
			return err
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, upSQL); err == nil {
			_, err = tx.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, version)
		}
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply %s: %w", version, err)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("applied %s", version)
	}
	return nil
}

func down(ctx context.Context, db *sql.DB, dir string) error {
	var version string
	err := db.QueryRowContext(ctx,
		`SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`,
	).Scan(&version)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	_, downSQL, err := readMigration(filepath.Join(dir, version))
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, downSQL); err == nil {
		_, err = tx.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version=$1`, version)
	}
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("rollback %s: %w", version, err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("rolled back %s", version)
	return nil
}

func status(ctx context.Context, db *sql.DB, dir string) error {
	files, err := migrationFiles(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		version := filepath.Base(file)
		var exists bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, version,
		).Scan(&exists); err != nil {
			return err
		}
		log.Printf("%-8s %s", map[bool]string{true: "applied", false: "pending"}[exists], version)
	}
	return nil
}

func readMigration(path string) (string, string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(raw), "-- +goose Down", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("%s has no down section", path)
	}
	upSQL := strings.Replace(parts[0], "-- +goose Up", "", 1)
	return strings.TrimSpace(upSQL), strings.TrimSpace(parts[1]), nil
}
