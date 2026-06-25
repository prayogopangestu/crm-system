package usecase

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
)

type fakeRepository struct {
	domain.UserRepository
	domain.DealRepository
	user       domain.User
	stageCalls int
}

func (f *fakeRepository) UserByEmail(context.Context, string) (domain.User, error) {
	return f.user, nil
}

func (f *fakeRepository) UpdateDealStage(context.Context, domain.Principal, string, domain.StageUpdateInput) error {
	f.stageCalls++
	return nil
}

func testService(repo *fakeRepository) *Service {
	return New(
		Repositories{Users: repo, Deals: repo},
		nil, auth.New("01234567890123456789012345678901", time.Hour),
		nil, nil, slog.New(slog.NewTextHandler(io.Discard, nil)),
		time.FixedZone("WIB", 7*60*60), "http://localhost:3000", 4,
	)
}

func TestUpdateDealStageRequiresLostReason(t *testing.T) {
	repo := &fakeRepository{}
	service := testService(repo)
	err := service.UpdateDealStage(context.Background(), domain.Principal{
		UserID: "u1", OrganizationID: "o1",
	}, "d1", domain.StageUpdateInput{Stage: "lost"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
	if repo.stageCalls != 0 {
		t.Fatal("repository must not be called for invalid lost transition")
	}
}

func TestUpdateDealStageAcceptsLostReason(t *testing.T) {
	repo := &fakeRepository{}
	service := testService(repo)
	err := service.UpdateDealStage(context.Background(), domain.Principal{
		UserID: "u1", OrganizationID: "o1",
	}, "d1", domain.StageUpdateInput{Stage: "lost", LostReason: "Harga"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.stageCalls != 1 {
		t.Fatalf("expected one repository call, got %d", repo.stageCalls)
	}
}

func TestInviteMemberRequiresAdmin(t *testing.T) {
	service := testService(&fakeRepository{})
	_, err := service.InviteMember(context.Background(), domain.Principal{
		Role: domain.RoleSales,
	}, domain.InviteInput{Name: "Budi", Email: "budi@example.com", Role: domain.RoleSales})
	if err != domain.ErrForbidden {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestLoginRejectsPendingUser(t *testing.T) {
	service := testService(&fakeRepository{user: domain.User{
		Status: "Menunggu", PasswordHash: "$2a$04$invalid",
	}})
	_, err := service.Login(context.Background(), domain.LoginInput{Email: "pending@example.com", Password: "secret"})
	if err != domain.ErrUnauthorized {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}
