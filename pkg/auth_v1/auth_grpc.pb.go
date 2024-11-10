// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: auth.proto

package auth_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthV1Client is the client API for AuthV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthV1Client interface {
	// Создает нового пользователя
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	// Получает пользователя по id
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	// Обновляет данные пользователя по id
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Удаляет пользователя по id
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type authV1Client struct {
	cc grpc.ClientConnInterface
}

func NewAuthV1Client(cc grpc.ClientConnInterface) AuthV1Client {
	return &authV1Client{cc}
}

func (c *authV1Client) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/auth_v1.AuthV1/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authV1Client) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/auth_v1.AuthV1/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authV1Client) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/auth_v1.AuthV1/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authV1Client) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/auth_v1.AuthV1/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthV1Server is the server API for AuthV1 service.
// All implementations must embed UnimplementedAuthV1Server
// for forward compatibility
type AuthV1Server interface {
	// Создает нового пользователя
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	// Получает пользователя по id
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	// Обновляет данные пользователя по id
	UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error)
	// Удаляет пользователя по id
	DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedAuthV1Server()
}

// UnimplementedAuthV1Server must be embedded to have forward compatible implementations.
type UnimplementedAuthV1Server struct {
}

func (UnimplementedAuthV1Server) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedAuthV1Server) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedAuthV1Server) UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedAuthV1Server) DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedAuthV1Server) mustEmbedUnimplementedAuthV1Server() {}

// UnsafeAuthV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthV1Server will
// result in compilation errors.
type UnsafeAuthV1Server interface {
	mustEmbedUnimplementedAuthV1Server()
}

func RegisterAuthV1Server(s grpc.ServiceRegistrar, srv AuthV1Server) {
	s.RegisterService(&AuthV1_ServiceDesc, srv)
}

func _AuthV1_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthV1Server).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_v1.AuthV1/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthV1Server).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthV1_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthV1Server).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_v1.AuthV1/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthV1Server).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthV1_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthV1Server).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_v1.AuthV1/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthV1Server).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthV1_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthV1Server).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_v1.AuthV1/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthV1Server).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthV1_ServiceDesc is the grpc.ServiceDesc for AuthV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth_v1.AuthV1",
	HandlerType: (*AuthV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _AuthV1_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _AuthV1_GetUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _AuthV1_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _AuthV1_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
