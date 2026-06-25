package grpc

import (
	"context"
	"errors"
	"strings"

	userv1 "github.com/prayogopangestu/crm-system/backend/api/protobuf/gen"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/internal/usecase"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	service *usecase.Service
}

func NewUserServer(service *usecase.Service) *UserServer {
	return &UserServer{service: service}
}

func AuthInterceptor(tokens *auth.Manager) googlegrpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *googlegrpc.UnaryServerInfo, handler googlegrpc.UnaryHandler) (any, error) {
		values := metadata.ValueFromIncomingContext(ctx, "authorization")
		if len(values) == 0 || !strings.HasPrefix(values[0], "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "Bearer token diperlukan")
		}
		principal, err := tokens.Parse(strings.TrimSpace(strings.TrimPrefix(values[0], "Bearer ")))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token tidak valid")
		}
		return handler(domain.WithPrincipal(ctx, principal), req)
	}
}

func (s *UserServer) GetProfile(ctx context.Context, _ *userv1.GetProfileRequest) (*userv1.User, error) {
	user, err := s.service.Profile(ctx, grpcPrincipal(ctx))
	if err != nil {
		return nil, grpcError(err)
	}
	return userMessage(user), nil
}

func (s *UserServer) UpdateProfile(ctx context.Context, request *userv1.UpdateProfileRequest) (*userv1.User, error) {
	user, err := s.service.UpdateProfile(ctx, grpcPrincipal(ctx), domain.UpdateProfileInput{
		FirstName: request.GetFirstName(), LastName: request.GetLastName(), Email: request.GetEmail(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return userMessage(user), nil
}

func (s *UserServer) ListTeam(ctx context.Context, _ *userv1.ListTeamRequest) (*userv1.ListTeamResponse, error) {
	users, err := s.service.ListTeam(ctx, grpcPrincipal(ctx))
	if err != nil {
		return nil, grpcError(err)
	}
	response := &userv1.ListTeamResponse{Users: make([]*userv1.User, 0, len(users))}
	for _, user := range users {
		response.Users = append(response.Users, userMessage(user))
	}
	return response, nil
}

func (s *UserServer) InviteMember(ctx context.Context, request *userv1.InviteMemberRequest) (*userv1.InviteMemberResponse, error) {
	result, err := s.service.InviteMember(ctx, grpcPrincipal(ctx), domain.InviteInput{
		Name: request.GetName(), Email: request.GetEmail(), Role: request.GetRole(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return &userv1.InviteMemberResponse{User: userMessage(result.User), InviteUrl: result.InviteURL}, nil
}

func (s *UserServer) RevokeMember(ctx context.Context, request *userv1.RevokeMemberRequest) (*userv1.RevokeMemberResponse, error) {
	if err := s.service.RevokeMember(ctx, grpcPrincipal(ctx), request.GetUserId()); err != nil {
		return nil, grpcError(err)
	}
	return &userv1.RevokeMemberResponse{Success: true}, nil
}

func grpcPrincipal(ctx context.Context) domain.Principal {
	value, _ := domain.PrincipalFromContext(ctx)
	return value
}

func userMessage(user domain.User) *userv1.User {
	return &userv1.User{
		Id: user.ID, FirstName: user.FirstName, LastName: user.LastName, Name: user.Name,
		Email: user.Email, Role: user.Role, Status: user.Status, AvatarUrl: user.AvatarURL,
		Initials: user.Initials,
	}
}

func grpcError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "request tidak valid")
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, "unauthorized")
	case errors.Is(err, domain.ErrForbidden):
		return status.Error(codes.PermissionDenied, "forbidden")
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, "data tidak ditemukan")
	case errors.Is(err, domain.ErrConflict):
		return status.Error(codes.AlreadyExists, "data sudah digunakan")
	default:
		return status.Error(codes.Internal, "terjadi kesalahan internal")
	}
}
