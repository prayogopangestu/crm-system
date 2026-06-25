package user

import (
	"context"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	userv1 "github.com/prayogopangestu/crm-system/backend/api/protobuf/gen"
	"github.com/prayogopangestu/crm-system/backend/internal/server/grpcserver"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

func TestGetProfileOverGRPC(t *testing.T) {
	tokens := auth.New("01234567890123456789012345678901", time.Hour)
	repository := &fakeRepository{user: User{
		ID: "u1", OrganizationID: "o1", FirstName: "Sarah", LastName: "Jenkins",
		Name: "Sarah Jenkins", Email: "sarah@example.com", Role: shared.RoleAdmin,
	}}
	service := NewService(
		repository, shared.CacheHelper{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))},
		tokens, "http://localhost", 4,
	)
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer(grpc.UnaryInterceptor(grpcserver.AuthInterceptor(tokens)))
	userv1.RegisterUserServiceServer(server, NewGRPCServer(service))
	go func() { _ = server.Serve(listener) }()
	defer server.Stop()

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	token, err := tokens.Create("u1", "o1", shared.RoleAdmin, "Sarah Jenkins")
	if err != nil {
		t.Fatal(err)
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
	value, err := userv1.NewUserServiceClient(conn).GetProfile(ctx, &userv1.GetProfileRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if value.GetEmail() != "sarah@example.com" {
		t.Fatalf("unexpected profile: %+v", value)
	}
}
