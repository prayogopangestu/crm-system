package user

import (
	"context"
	"testing"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type fakeRepository struct{ user User }

func (f *fakeRepository) Register(context.Context, string, User) (User, error) { return User{}, nil }
func (f *fakeRepository) ByEmail(context.Context, string) (User, error)        { return f.user, nil }
func (f *fakeRepository) ByID(context.Context, string, string) (User, error)   { return f.user, nil }
func (f *fakeRepository) UpdateProfile(context.Context, shared.Principal, UpdateProfileInput) (User, error) {
	return f.user, nil
}
func (f *fakeRepository) AcceptInvite(context.Context, string, string) (User, error) {
	return User{}, nil
}
func (f *fakeRepository) ListTeam(context.Context, string) ([]User, error) { return nil, nil }
func (f *fakeRepository) InviteMember(context.Context, shared.Principal, User, Invitation) (User, error) {
	return User{}, nil
}
func (f *fakeRepository) RevokeMember(context.Context, string, string) error { return nil }

func TestInviteMemberRequiresAdmin(t *testing.T) {
	service := NewService(&fakeRepository{}, shared.CacheHelper{}, nil, "http://localhost", 4)
	_, err := service.InviteMember(context.Background(), shared.Principal{
		Role: shared.RoleSales,
	}, InviteInput{Name: "Budi", Email: "budi@example.com", Role: shared.RoleSales})
	if err != shared.ErrForbidden {
		t.Fatalf("expected forbidden, got %v", err)
	}
}
