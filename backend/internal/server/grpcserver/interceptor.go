package grpcserver

import (
	"context"
	"strings"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(tokens *auth.Manager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		values := metadata.ValueFromIncomingContext(ctx, "authorization")
		if len(values) == 0 || !strings.HasPrefix(values[0], "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "Bearer token diperlukan")
		}
		principal, err := tokens.Parse(strings.TrimSpace(strings.TrimPrefix(values[0], "Bearer ")))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token tidak valid")
		}
		return handler(shared.WithPrincipal(ctx, principal), req)
	}
}
