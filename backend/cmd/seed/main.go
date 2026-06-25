package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/config"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/user"
	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "configs/config.yaml"
	}
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	location, _ := time.LoadLocation(cfg.App.Timezone)
	pool, err := database.OpenPostgres(context.Background(), cfg.Database.URL, 1, 2)
	if err != nil {
		log.Fatal(err)
	}
	store := postgresx.New(pool, location)
	defer store.Close()
	repository := user.NewRepository(store)

	email := env("SEED_ADMIN_EMAIL", "admin@crm.local")
	password := env("SEED_ADMIN_PASSWORD", "Admin123!")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cfg.Auth.BcryptCost)
	if err != nil {
		log.Fatal(err)
	}
	created, err := repository.Register(context.Background(), env("SEED_ORGANIZATION", "CRM Demo"), user.User{
		FirstName: "Demo", LastName: "Admin", Email: strings.ToLower(email),
		PasswordHash: string(hash), Role: shared.RoleAdmin,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("demo tenant created: user=%s email=%s", created.ID, email)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
