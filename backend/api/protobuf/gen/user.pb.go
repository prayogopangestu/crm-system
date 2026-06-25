// Code generated from api/protobuf/user.proto. DO NOT EDIT.
package userv1

import "fmt"

type GetProfileRequest struct{}

func (x *GetProfileRequest) Reset()         { *x = GetProfileRequest{} }
func (x *GetProfileRequest) String() string { return "{}" }
func (*GetProfileRequest) ProtoMessage()    {}

type ListTeamRequest struct{}

func (x *ListTeamRequest) Reset()         { *x = ListTeamRequest{} }
func (x *ListTeamRequest) String() string { return "{}" }
func (*ListTeamRequest) ProtoMessage()    {}

type User struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName string `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName  string `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Name      string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Email     string `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	Role      string `protobuf:"bytes,6,opt,name=role,proto3" json:"role,omitempty"`
	Status    string `protobuf:"bytes,7,opt,name=status,proto3" json:"status,omitempty"`
	AvatarUrl string `protobuf:"bytes,8,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"`
	Initials  string `protobuf:"bytes,9,opt,name=initials,proto3" json:"initials,omitempty"`
}

func (x *User) Reset()         { *x = User{} }
func (x *User) String() string { return fmt.Sprintf("%+v", *x) }
func (*User) ProtoMessage()    {}
func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *User) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}
func (x *User) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}
func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}
func (x *User) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}
func (x *User) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}
func (x *User) GetAvatarUrl() string {
	if x != nil {
		return x.AvatarUrl
	}
	return ""
}
func (x *User) GetInitials() string {
	if x != nil {
		return x.Initials
	}
	return ""
}

type UpdateProfileRequest struct {
	FirstName string `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName  string `protobuf:"bytes,2,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Email     string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *UpdateProfileRequest) Reset()         { *x = UpdateProfileRequest{} }
func (x *UpdateProfileRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*UpdateProfileRequest) ProtoMessage()    {}
func (x *UpdateProfileRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}
func (x *UpdateProfileRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}
func (x *UpdateProfileRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type ListTeamResponse struct {
	Users []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *ListTeamResponse) Reset()         { *x = ListTeamResponse{} }
func (x *ListTeamResponse) String() string { return fmt.Sprintf("%+v", *x) }
func (*ListTeamResponse) ProtoMessage()    {}
func (x *ListTeamResponse) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

type InviteMemberRequest struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Role  string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *InviteMemberRequest) Reset()         { *x = InviteMemberRequest{} }
func (x *InviteMemberRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*InviteMemberRequest) ProtoMessage()    {}
func (x *InviteMemberRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *InviteMemberRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}
func (x *InviteMemberRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type InviteMemberResponse struct {
	User      *User  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	InviteUrl string `protobuf:"bytes,2,opt,name=invite_url,json=inviteUrl,proto3" json:"invite_url,omitempty"`
}

func (x *InviteMemberResponse) Reset()         { *x = InviteMemberResponse{} }
func (x *InviteMemberResponse) String() string { return fmt.Sprintf("%+v", *x) }
func (*InviteMemberResponse) ProtoMessage()    {}
func (x *InviteMemberResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}
func (x *InviteMemberResponse) GetInviteUrl() string {
	if x != nil {
		return x.InviteUrl
	}
	return ""
}

type RevokeMemberRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *RevokeMemberRequest) Reset()         { *x = RevokeMemberRequest{} }
func (x *RevokeMemberRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*RevokeMemberRequest) ProtoMessage()    {}
func (x *RevokeMemberRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type RevokeMemberResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *RevokeMemberResponse) Reset()           { *x = RevokeMemberResponse{} }
func (x *RevokeMemberResponse) String() string   { return fmt.Sprintf("%+v", *x) }
func (*RevokeMemberResponse) ProtoMessage()      {}
func (x *RevokeMemberResponse) GetSuccess() bool { return x != nil && x.Success }
