//go:build integration

package contact_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/modules/contact"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/user"
	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/pkg/database"
)

func TestTenantIsolationAgainstRealPostgres(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := database.OpenPostgres(context.Background(), url, 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	store := postgresx.New(db, time.UTC)
	defer store.Close()
	users := user.NewRepository(store)
	contacts := contact.NewRepository(store)
	first, err := users.Register(context.Background(), "Tenant A", user.User{
		FirstName: "Admin", Email: "tenant-a@example.test", PasswordHash: "hash", Role: shared.RoleAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := users.Register(context.Background(), "Tenant B", user.User{
		FirstName: "Admin", Email: "tenant-b@example.test", PasswordHash: "hash", Role: shared.RoleAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	created, err := contacts.Create(context.Background(), shared.Principal{
		UserID: first.ID, OrganizationID: first.OrganizationID, Name: first.Name,
	}, contact.Input{Name: "Private Contact", Email: "private@example.test", Company: "A", Status: "Prospek Awal"})
	if err != nil {
		t.Fatal(err)
	}
	if created.OrganizationID != first.OrganizationID {
		t.Fatal("contact was created in wrong tenant")
	}
	list, err := contacts.List(context.Background(), second.OrganizationID, "Private Contact", "", contact.Page{Page: 1, Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 0 {
		t.Fatal("tenant B can see tenant A data")
	}
}
