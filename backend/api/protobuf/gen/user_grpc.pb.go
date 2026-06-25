// Code generated from api/protobuf/user.proto. DO NOT EDIT.
package userv1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	UserService_GetProfile_FullMethodName    = "/crm.user.v1.UserService/GetProfile"
	UserService_UpdateProfile_FullMethodName = "/crm.user.v1.UserService/UpdateProfile"
	UserService_ListTeam_FullMethodName      = "/crm.user.v1.UserService/ListTeam"
	UserService_InviteMember_FullMethodName  = "/crm.user.v1.UserService/InviteMember"
	UserService_RevokeMember_FullMethodName  = "/crm.user.v1.UserService/RevokeMember"
)

type UserServiceClient interface {
	GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*User, error)
	UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*User, error)
	ListTeam(ctx context.Context, in *ListTeamRequest, opts ...grpc.CallOption) (*ListTeamResponse, error)
	InviteMember(ctx context.Context, in *InviteMemberRequest, opts ...grpc.CallOption) (*InviteMemberResponse, error)
	RevokeMember(ctx context.Context, in *RevokeMemberRequest, opts ...grpc.CallOption) (*RevokeMemberResponse, error)
}

type userServiceClient struct{ cc grpc.ClientConnInterface }

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserService_GetProfile_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserService_UpdateProfile_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) ListTeam(ctx context.Context, in *ListTeamRequest, opts ...grpc.CallOption) (*ListTeamResponse, error) {
	out := new(ListTeamResponse)
	err := c.cc.Invoke(ctx, UserService_ListTeam_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) InviteMember(ctx context.Context, in *InviteMemberRequest, opts ...grpc.CallOption) (*InviteMemberResponse, error) {
	out := new(InviteMemberResponse)
	err := c.cc.Invoke(ctx, UserService_InviteMember_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) RevokeMember(ctx context.Context, in *RevokeMemberRequest, opts ...grpc.CallOption) (*RevokeMemberResponse, error) {
	out := new(RevokeMemberResponse)
	err := c.cc.Invoke(ctx, UserService_RevokeMember_FullMethodName, in, out, opts...)
	return out, err
}

type UserServiceServer interface {
	GetProfile(context.Context, *GetProfileRequest) (*User, error)
	UpdateProfile(context.Context, *UpdateProfileRequest) (*User, error)
	ListTeam(context.Context, *ListTeamRequest) (*ListTeamResponse, error)
	InviteMember(context.Context, *InviteMemberRequest) (*InviteMemberResponse, error)
	RevokeMember(context.Context, *RevokeMemberRequest) (*RevokeMemberResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) GetProfile(context.Context, *GetProfileRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}
func (UnimplementedUserServiceServer) UpdateProfile(context.Context, *UpdateProfileRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
func (UnimplementedUserServiceServer) ListTeam(context.Context, *ListTeamRequest) (*ListTeamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTeam not implemented")
}
func (UnimplementedUserServiceServer) InviteMember(context.Context, *InviteMemberRequest) (*InviteMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteMember not implemented")
}
func (UnimplementedUserServiceServer) RevokeMember(context.Context, *RevokeMemberRequest) (*RevokeMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeMember not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_GetProfile_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(GetProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_GetProfile_FullMethodName}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UserServiceServer).GetProfile(ctx, req.(*GetProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _UserService_UpdateProfile_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(UpdateProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_UpdateProfile_FullMethodName}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UserServiceServer).UpdateProfile(ctx, req.(*UpdateProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _UserService_ListTeam_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ListTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ListTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_ListTeam_FullMethodName}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UserServiceServer).ListTeam(ctx, req.(*ListTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _UserService_InviteMember_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(InviteMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).InviteMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_InviteMember_FullMethodName}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UserServiceServer).InviteMember(ctx, req.(*InviteMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _UserService_RevokeMember_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(RevokeMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RevokeMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_RevokeMember_FullMethodName}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UserServiceServer).RevokeMember(ctx, req.(*RevokeMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "crm.user.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetProfile", Handler: _UserService_GetProfile_Handler},
		{MethodName: "UpdateProfile", Handler: _UserService_UpdateProfile_Handler},
		{MethodName: "ListTeam", Handler: _UserService_ListTeam_Handler},
		{MethodName: "InviteMember", Handler: _UserService_InviteMember_Handler},
		{MethodName: "RevokeMember", Handler: _UserService_RevokeMember_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/protobuf/user.proto",
}
