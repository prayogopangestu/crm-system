//go:build integration

package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/database"
)

func TestTenantIsolationAgainstRealPostgres(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	pool, err := database.OpenPostgres(context.Background(), url, 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	repo := New(pool, time.UTC)
	defer repo.Close()
	first, err := repo.Register(context.Background(), "Tenant A", domain.User{
		FirstName: "Admin", Email: "tenant-a@example.test", PasswordHash: "hash", Role: domain.RoleAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := repo.Register(context.Background(), "Tenant B", domain.User{
		FirstName: "Admin", Email: "tenant-b@example.test", PasswordHash: "hash", Role: domain.RoleAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	contact, err := repo.CreateContact(context.Background(), domain.Principal{
		UserID: first.ID, OrganizationID: first.OrganizationID, Name: first.Name,
	}, domain.ContactInput{
		Name: "Private Contact", Email: "private@example.test", Company: "A", Status: "Prospek Awal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if contact.OrganizationID != first.OrganizationID {
		t.Fatal("contact was created in wrong tenant")
	}
	list, err := repo.ListContacts(context.Background(), second.OrganizationID, "Private Contact", "", domain.Page{Page: 1, Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 0 {
		t.Fatal("tenant B can see tenant A data")
	}
}
