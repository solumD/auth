package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/solumD/auth/pkg/auth_v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serv: %s", err)
	}
}

type server struct {
	desc.UnimplementedAuthV1Server
}

// CreateUser creates new user
func (s *server) CreateUser(_ context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Printf("[Create] request data |\nname: %v, email: %v, password: %v, password_confirm: %v, role: %v",
		req.Info.Info.Name,
		req.Info.Info.Email,
		req.Info.Password,
		req.Info.PasswordConfirm,
		req.Info.Info.Role,
	)

	return &desc.CreateUserResponse{
		Id: gofakeit.Int64(),
	}, nil
}

// GetUser returns user by id
func (s *server) GetUser(_ context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	log.Printf("[Get] request data |\nid: %v", req.Id)

	return &desc.GetUserResponse{
		User: &desc.User{
			Id: gofakeit.Int64(),
			Info: &desc.UserInfo{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Role:  desc.Role(gofakeit.Number(0, 2)),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

// UpdateUser updates user's data by id
func (s *server) UpdateUser(_ context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	log.Printf("[Update] request data |\nid: %v, name: %v, email: %v, role: %v", req.Id, req.Info.Name, req.Info.Email, req.Info.Role)
	return nil, nil
}

// DeleteUser deletes user by id
func (s *server) DeleteUser(_ context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Printf("[Delete] request data |\nid: %v", req.Id)
	return nil, nil
}
