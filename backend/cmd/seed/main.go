package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/config"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	postgresrepo "github.com/prayogopangestu/crm-system/backend/internal/repository/postgres"
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
	repo := postgresrepo.New(pool, location)
	defer repo.Close()

	email := env("SEED_ADMIN_EMAIL", "admin@crm.local")
	password := env("SEED_ADMIN_PASSWORD", "Admin123!")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cfg.Auth.BcryptCost)
	if err != nil {
		log.Fatal(err)
	}
	user, err := repo.Register(context.Background(), env("SEED_ORGANIZATION", "CRM Demo"), domain.User{
		FirstName: "Demo", LastName: "Admin", Email: strings.ToLower(email),
		PasswordHash: string(hash), Role: domain.RoleAdmin,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("demo tenant created: user=%s email=%s", user.ID, email)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
