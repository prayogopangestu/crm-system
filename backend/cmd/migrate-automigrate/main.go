package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prayogopangestu/crm-system/backend/pkg/database"
	"gorm.io/gorm"
)

type Organization struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (Organization) TableName() string { return "organizations" }

type User struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string     `gorm:"type:uuid;not null"`
	FirstName      string     `gorm:"column:first_name;type:text;not null"`
	LastName       string     `gorm:"column:last_name;type:text;not null;default:''"`
	Email          string     `gorm:"type:text;not null"`
	PasswordHash   *string    `gorm:"column:password_hash;type:text"`
	Role           string     `gorm:"type:text;not null"`
	Status         string     `gorm:"type:text;not null;default:'Aktif'"`
	AvatarURL      string     `gorm:"column:avatar_url;type:text;not null;default:''"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	RevokedAt      *time.Time `gorm:"column:revoked_at;type:timestamptz"`
}

func (User) TableName() string { return "users" }

type UserInvitation struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string     `gorm:"type:uuid;not null"`
	UserID         string     `gorm:"type:uuid;not null"`
	Email          string     `gorm:"type:text;not null"`
	Role           string     `gorm:"type:text;not null"`
	TokenHash      string     `gorm:"column:token_hash;type:text;not null"`
	ExpiresAt      time.Time  `gorm:"column:expires_at;type:timestamptz;not null"`
	AcceptedAt     *time.Time `gorm:"column:accepted_at;type:timestamptz"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

func (UserInvitation) TableName() string { return "user_invitations" }

type Contact struct {
	ID              string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID  string     `gorm:"type:uuid;not null"`
	OwnerID         *string    `gorm:"column:owner_id;type:uuid"`
	Name            string     `gorm:"type:text;not null"`
	Email           string     `gorm:"type:text;not null"`
	Company         string     `gorm:"type:text;not null"`
	Role            string     `gorm:"type:text;not null;default:''"`
	Status          string     `gorm:"type:text;not null"`
	AvatarURL       string     `gorm:"column:avatar_url;type:text;not null;default:''"`
	LastContactedAt time.Time  `gorm:"column:last_contacted_at;type:timestamptz;not null;default:now()"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamptz"`
}

func (Contact) TableName() string { return "contacts" }

type PipelineStage struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string    `gorm:"type:uuid;not null"`
	Key            string    `gorm:"type:text;not null"`
	Name           string    `gorm:"type:text;not null"`
	Color          string    `gorm:"type:text;not null;default:'bg-surface-variant'"`
	Position       int       `gorm:"type:integer;not null"`
	IsSystem       bool      `gorm:"column:is_system;type:boolean;not null;default:false"`
	CreatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (PipelineStage) TableName() string { return "pipeline_stages" }

type Deal struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string     `gorm:"type:uuid;not null"`
	AssigneeID     *string    `gorm:"column:assignee_id;type:uuid"`
	Title          string     `gorm:"type:text;not null"`
	Company        string     `gorm:"type:text;not null"`
	Value          int64      `gorm:"type:bigint;not null"`
	Priority       string     `gorm:"type:text;not null"`
	StageKey       string     `gorm:"column:stage_key;type:text;not null"`
	LostReason     string     `gorm:"column:lost_reason;type:text;not null;default:''"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;type:timestamptz"`
}

func (Deal) TableName() string { return "deals" }

type Task struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string     `gorm:"type:uuid;not null"`
	AssigneeID     *string    `gorm:"column:assignee_id;type:uuid"`
	Title          string     `gorm:"type:text;not null"`
	Company        string     `gorm:"type:text;not null"`
	DueDate        time.Time  `gorm:"column:due_date;type:date;not null"`
	DueTime        string     `gorm:"column:due_time;type:time;not null"`
	Type           string     `gorm:"type:text;not null"`
	Priority       string     `gorm:"type:text;not null"`
	Notes          string     `gorm:"type:text;not null;default:''"`
	Completed      bool       `gorm:"type:boolean;not null;default:false"`
	CompletedAt    *time.Time `gorm:"column:completed_at;type:timestamptz"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;type:timestamptz"`
}

func (Task) TableName() string { return "tasks" }

type Activity struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string    `gorm:"type:uuid;not null"`
	ActorID        string    `gorm:"column:actor_id;type:uuid"`
	ActorName      string    `gorm:"column:actor_name;type:text;not null"`
	Action         string    `gorm:"type:text;not null"`
	Target         string    `gorm:"type:text;not null"`
	IsHighlight    bool      `gorm:"column:is_highlight;type:boolean;not null;default:false"`
	CreatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (Activity) TableName() string { return "activities" }

type PerformanceGoal struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string    `gorm:"type:uuid;not null"`
	Month          time.Time `gorm:"type:date;not null"`
	Goal           int64     `gorm:"type:bigint;not null"`
	CreatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (PerformanceGoal) TableName() string { return "performance_goals" }

type Notification struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string     `gorm:"type:uuid;not null"`
	UserID         *string    `gorm:"column:user_id;type:uuid"`
	Title          string     `gorm:"type:text;not null"`
	Message        string     `gorm:"type:text;not null"`
	ReadAt         *time.Time `gorm:"column:read_at;type:timestamptz"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

func (Notification) TableName() string { return "notifications" }

type TelegramIntegration struct {
	OrganizationID    string    `gorm:"type:uuid;primaryKey"`
	BotTokenEncrypted string    `gorm:"column:bot_token_encrypted;type:text;not null;default:''"`
	ChatID            string    `gorm:"column:chat_id;type:text;not null;default:''"`
	Enabled           bool      `gorm:"type:boolean;not null;default:false"`
	CreatedAt         time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt         time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (TelegramIntegration) TableName() string { return "telegram_integrations" }

type OutboxEvent struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID string     `gorm:"type:uuid;not null"`
	EventType      string     `gorm:"column:event_type;type:text;not null"`
	Payload        []byte     `gorm:"type:jsonb;not null"`
	Attempts       int        `gorm:"type:integer;not null;default:0"`
	NextAttemptAt  time.Time  `gorm:"column:next_attempt_at;type:timestamptz;not null;default:now()"`
	ProcessedAt    *time.Time `gorm:"column:processed_at;type:timestamptz"`
	LastError      string     `gorm:"column:last_error;type:text;not null;default:''"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

func (OutboxEvent) TableName() string { return "outbox_events" }

func main() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://crm:crm@localhost:5432/crm?sslmode=disable"
	}
	log.Printf("menghubungkan ke database: %s", maskURL(url))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := database.OpenPostgres(ctx, url, 1, 5)
	if err != nil {
		log.Fatalf("gagal membuka database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("gagal mendapatkan *sql.DB: %v", err)
	}
	defer sqlDB.Close()

	if _, err := sqlDB.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS pgcrypto`); err != nil {
		log.Fatalf("gagal membuat ekstensi pgcrypto: %v", err)
	}
	log.Println("ekstensi pgcrypto siap")

	models := []any{
		Organization{},
		User{},
		UserInvitation{},
		Contact{},
		PipelineStage{},
		Deal{},
		Task{},
		Activity{},
		PerformanceGoal{},
		Notification{},
		TelegramIntegration{},
		OutboxEvent{},
	}
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("AutoMigrate gagal: %v", err)
	}
	for _, m := range models {
		log.Printf("tabel disinkronkan: %s", tableName(m))
	}

	if err := applySupplemental(ctx, db); err != nil {
		log.Fatalf("gagal menerapkan index & foreign key tambahan: %v", err)
	}
	log.Println("AutoMigrate selesai. Semua tabel siap digunakan.")
}

func tableName(model any) string {
	type namer interface{ TableName() string }
	if n, ok := model.(namer); ok {
		return n.TableName()
	}
	return fmt.Sprintf("%T", model)
}

func applySupplemental(ctx context.Context, db *gorm.DB) error {
	indexes := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users (lower(email))`,
		`CREATE INDEX IF NOT EXISTS users_org_status_idx ON users (organization_id, status)`,
		`CREATE INDEX IF NOT EXISTS invitations_org_email_idx ON user_invitations (organization_id, lower(email))`,
		`CREATE INDEX IF NOT EXISTS contacts_org_status_idx ON contacts (organization_id, status) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS contacts_org_search_idx ON contacts (organization_id, lower(name), lower(email), lower(company)) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS pipeline_stages_org_key_uniq ON pipeline_stages (organization_id, key)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS pipeline_stages_org_position_uniq ON pipeline_stages (organization_id, position)`,
		`CREATE INDEX IF NOT EXISTS deals_org_stage_idx ON deals (organization_id, stage_key) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS deals_org_assignee_idx ON deals (organization_id, assignee_id) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS tasks_org_date_idx ON tasks (organization_id, due_date) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS tasks_org_assignee_idx ON tasks (organization_id, assignee_id) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS activities_org_created_idx ON activities (organization_id, created_at DESC)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS performance_goals_org_month_uniq ON performance_goals (organization_id, month)`,
		`CREATE INDEX IF NOT EXISTS notifications_user_read_idx ON notifications (organization_id, user_id, read_at, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS outbox_pending_idx ON outbox_events (next_attempt_at, created_at) WHERE processed_at IS NULL`,
	}

	foreignKeys := []struct {
		table, column, refTable, refColumn, onDelete string
	}{
		{"users", "organization_id", "organizations", "id", "CASCADE"},
		{"user_invitations", "organization_id", "organizations", "id", "CASCADE"},
		{"user_invitations", "user_id", "users", "id", "CASCADE"},
		{"contacts", "organization_id", "organizations", "id", "CASCADE"},
		{"contacts", "owner_id", "users", "id", "SET NULL"},
		{"pipeline_stages", "organization_id", "organizations", "id", "CASCADE"},
		{"deals", "organization_id", "organizations", "id", "CASCADE"},
		{"deals", "assignee_id", "users", "id", "SET NULL"},
		{"tasks", "organization_id", "organizations", "id", "CASCADE"},
		{"tasks", "assignee_id", "users", "id", "SET NULL"},
		{"activities", "organization_id", "organizations", "id", "CASCADE"},
		{"activities", "actor_id", "users", "id", "SET NULL"},
		{"performance_goals", "organization_id", "organizations", "id", "CASCADE"},
		{"notifications", "organization_id", "organizations", "id", "CASCADE"},
		{"notifications", "user_id", "users", "id", "CASCADE"},
		{"telegram_integrations", "organization_id", "organizations", "id", "CASCADE"},
		{"outbox_events", "organization_id", "organizations", "id", "CASCADE"},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, stmt := range indexes {
			if err := tx.Exec(stmt).Error; err != nil {
				return fmt.Errorf("index: %w", err)
			}
		}
		for _, fk := range foreignKeys {
			name := fmt.Sprintf("%s_%s_fkey", fk.table, fk.column)
			stmt := fmt.Sprintf(
				`DO $$ BEGIN
					IF NOT EXISTS (
						SELECT 1 FROM pg_constraint
						WHERE conname = '%s' AND conrelid = '%s'::regclass
					) THEN
						ALTER TABLE %s ADD CONSTRAINT %s
							FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE %s;
					END IF;
				END $$;`,
				name, fk.table, fk.table, name, fk.column, fk.refTable, fk.refColumn, fk.onDelete,
			)
			if err := tx.Exec(stmt).Error; err != nil {
				return fmt.Errorf("foreign key %s.%s: %w", fk.table, fk.column, err)
			}
		}
		return nil
	})
}

func maskURL(rawURL string) string {
	if host, _, ok := splitURL(rawURL); ok {
		return host
	}
	return "[unparseable]"
}

func splitURL(rawURL string) (host, _ string, _ bool) {
	for i := 0; i < len(rawURL); i++ {
		if rawURL[i] == '@' {
			rest := rawURL[i+1:]
			for j := 0; j < len(rest); j++ {
				if rest[j] == '/' || rest[j] == '?' {
					return rest[:j], "", true
				}
			}
			return rest, "", true
		}
	}
	return rawURL, "", true
}
