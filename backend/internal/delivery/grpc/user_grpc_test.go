package grpc

import (
	"context"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	userv1 "github.com/prayogopangestu/crm-system/backend/api/protobuf/gen"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/internal/usecase"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

type profileRepository struct {
	domain.Repository
	user domain.User
}

func (r *profileRepository) UserByID(context.Context, string, string) (domain.User, error) {
	return r.user, nil
}

func TestGetProfileOverGRPC(t *testing.T) {
	tokens := auth.New("01234567890123456789012345678901", time.Hour)
	repo := &profileRepository{user: domain.User{
		ID: "u1", OrganizationID: "o1", FirstName: "Sarah", LastName: "Jenkins",
		Name: "Sarah Jenkins", Email: "sarah@example.com", Role: domain.RoleAdmin,
	}}
	service := usecase.New(
		repo, nil, tokens, nil, nil, slog.New(slog.NewTextHandler(io.Discard, nil)),
		time.UTC, "http://localhost", 4,
	)
	listener := bufconn.Listen(1024 * 1024)
	server := googlegrpc.NewServer(googlegrpc.UnaryInterceptor(AuthInterceptor(tokens)))
	userv1.RegisterUserServiceServer(server, NewUserServer(service))
	go func() { _ = server.Serve(listener) }()
	defer server.Stop()

	conn, err := googlegrpc.NewClient(
		"passthrough:///bufnet",
		googlegrpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	token, err := tokens.Create(domain.User{ID: "u1", OrganizationID: "o1", Role: domain.RoleAdmin})
	if err != nil {
		t.Fatal(err)
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
	user, err := userv1.NewUserServiceClient(conn).GetProfile(ctx, &userv1.GetProfileRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if user.GetEmail() != "sarah@example.com" {
		t.Fatalf("unexpected profile: %+v", user)
	}
}
