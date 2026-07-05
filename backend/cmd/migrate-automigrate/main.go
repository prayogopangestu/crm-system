package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/migrations"
	"github.com/prayogopangestu/crm-system/backend/pkg/database"
)

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

	if err := migrations.Run(ctx, db); err != nil {
		log.Fatalf("AutoMigrate gagal: %v", err)
	}
	log.Println("AutoMigrate selesai. Semua tabel siap digunakan.")
}

func maskURL(rawURL string) string {
	for i := 0; i < len(rawURL); i++ {
		if rawURL[i] == '@' {
			rest := rawURL[i+1:]
			for j := 0; j < len(rest); j++ {
				if rest[j] == '/' || rest[j] == '?' {
					return rest[:j]
				}
			}
			return rest
		}
	}
	return rawURL
}
