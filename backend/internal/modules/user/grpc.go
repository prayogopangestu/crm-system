package user

import (
	"context"
	"errors"

	userv1 "github.com/prayogopangestu/crm-system/backend/api/protobuf/gen"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	userv1.UnimplementedUserServiceServer
	service *Service
}

func NewGRPCServer(service *Service) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) GetProfile(ctx context.Context, _ *userv1.GetProfileRequest) (*userv1.User, error) {
	value, err := s.service.Profile(ctx, grpcPrincipal(ctx))
	if err != nil {
		return nil, grpcError(err)
	}
	return userMessage(value), nil
}

func (s *GRPCServer) UpdateProfile(ctx context.Context, request *userv1.UpdateProfileRequest) (*userv1.User, error) {
	value, err := s.service.UpdateProfile(ctx, grpcPrincipal(ctx), UpdateProfileInput{
		FirstName: request.GetFirstName(), LastName: request.GetLastName(), Email: request.GetEmail(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return userMessage(value), nil
}

func (s *GRPCServer) ListTeam(ctx context.Context, _ *userv1.ListTeamRequest) (*userv1.ListTeamResponse, error) {
	users, err := s.service.ListTeam(ctx, grpcPrincipal(ctx))
	if err != nil {
		return nil, grpcError(err)
	}
	result := &userv1.ListTeamResponse{Users: make([]*userv1.User, 0, len(users))}
	for _, value := range users {
		result.Users = append(result.Users, userMessage(value))
	}
	return result, nil
}

func (s *GRPCServer) InviteMember(ctx context.Context, request *userv1.InviteMemberRequest) (*userv1.InviteMemberResponse, error) {
	result, err := s.service.InviteMember(ctx, grpcPrincipal(ctx), InviteInput{
		Name: request.GetName(), Email: request.GetEmail(), Role: request.GetRole(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return &userv1.InviteMemberResponse{User: userMessage(result.User), InviteUrl: result.InviteURL}, nil
}

func (s *GRPCServer) RevokeMember(ctx context.Context, request *userv1.RevokeMemberRequest) (*userv1.RevokeMemberResponse, error) {
	if err := s.service.RevokeMember(ctx, grpcPrincipal(ctx), request.GetUserId()); err != nil {
		return nil, grpcError(err)
	}
	return &userv1.RevokeMemberResponse{Success: true}, nil
}

func grpcPrincipal(ctx context.Context) shared.Principal {
	value, _ := shared.PrincipalFromContext(ctx)
	return value
}

func userMessage(value User) *userv1.User {
	return &userv1.User{
		Id: value.ID, FirstName: value.FirstName, LastName: value.LastName, Name: value.Name,
		Email: value.Email, Role: value.Role, Status: value.Status, AvatarUrl: value.AvatarURL,
		Initials: value.Initials,
	}
}

func grpcError(err error) error {
	switch {
	case errors.Is(err, shared.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "request tidak valid")
	case errors.Is(err, shared.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, "unauthorized")
	case errors.Is(err, shared.ErrForbidden):
		return status.Error(codes.PermissionDenied, "forbidden")
	case errors.Is(err, shared.ErrNotFound):
		return status.Error(codes.NotFound, "data tidak ditemukan")
	case errors.Is(err, shared.ErrConflict):
		return status.Error(codes.AlreadyExists, "data sudah digunakan")
	default:
		return status.Error(codes.Internal, "terjadi kesalahan internal")
	}
}
