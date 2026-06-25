package deal

import (
	"context"
	"testing"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type fakeRepository struct{ stageCalls int }

func (f *fakeRepository) List(context.Context, string) ([]Deal, error) { return nil, nil }
func (f *fakeRepository) Create(context.Context, shared.Principal, Input) (Deal, error) {
	return Deal{}, nil
}
func (f *fakeRepository) Update(context.Context, shared.Principal, string, Input) (Deal, error) {
	return Deal{}, nil
}
func (f *fakeRepository) UpdateStage(context.Context, shared.Principal, string, StageInput) error {
	f.stageCalls++
	return nil
}
func (f *fakeRepository) Delete(context.Context, shared.Principal, string) error { return nil }

func TestUpdateStageRequiresLostReason(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository, shared.CacheHelper{})
	err := service.UpdateStage(context.Background(), shared.Principal{
		UserID: "u1", OrganizationID: "o1",
	}, "d1", StageInput{Stage: "lost"})
	if err != shared.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
	if repository.stageCalls != 0 {
		t.Fatal("repository must not be called")
	}
}

func TestUpdateStageAcceptsLostReason(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository, shared.CacheHelper{})
	err := service.UpdateStage(context.Background(), shared.Principal{
		UserID: "u1", OrganizationID: "o1",
	}, "d1", StageInput{Stage: "lost", LostReason: "Harga"})
	if err != nil {
		t.Fatal(err)
	}
	if repository.stageCalls != 1 {
		t.Fatalf("expected one call, got %d", repository.stageCalls)
	}
}
