package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
	"github.com/solumD/auth/internal/config"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	ctx := context.Background()
	pgPool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pgPool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pgPool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serv: %s", err)
	}
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

// CreateUser creates new user
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	fn := "CreateUser"
	log.Printf("[%s] request data | name: %v, email: %v, password: %v, password_confirm: %v, role: %v",
		fn,
		req.Name,
		req.Email,
		req.Password,
		req.PasswordConfirm,
		req.GetRole(),
	)

	builderInsertUser := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("username", "email", "password", "role", "created_at").
		Values(req.Name, req.Email,
			req.Password, req.GetRole(), time.Now()).
		Suffix("RETURNING id")

	query, args, err := builderInsertUser.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", fn, err)
		return nil, err
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("%s: failed to insert user: %v", fn, err)
		return nil, err
	}

	log.Printf("%s: inserted user with id: %d", fn, userID)
	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUser returns user by id
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	fn := "GetUser"
	log.Printf("[%s] request data | id: %v", fn, req.Id)

	builderGetUser := sq.Select("id", "username", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", fn, err)
		return nil, err
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("%s: failed to select user: %v", fn, err)
		return nil, err
	}

	var userID int64
	var username, email string
	var role int32
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&userID, &username, &email, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("%s: failed to scan user: %v", fn, err)
			return nil, err
		}
	}

	log.Printf("%s: selected user %d", fn, req.Id)
	return &desc.GetUserResponse{
		Id:        userID,
		Name:      username,
		Email:     email,
		Role:      desc.Role(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt.Time),
	}, nil
}

// UpdateUser updates user's data by id
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	fn := "UpdateUser"
	log.Printf("[%s] request data | id: %v, name: %v, email: %v, role: %v", fn, req.Id, req.Name, req.Email, req.Role)

	builderUpdateUser := sq.Update("users").
		PlaceholderFormat(sq.Dollar)

	if len(req.Name.Value) > 0 {
		builderUpdateUser = builderUpdateUser.Set("username", req.GetName().Value)
	}

	if len(req.Email.Value) > 0 {
		builderUpdateUser = builderUpdateUser.Set("email", req.GetEmail().Value)
	}

	if req.GetRole() >= 0 && req.GetRole() <= 2 {
		builderUpdateUser = builderUpdateUser.Set("role", req.GetRole())
	}

	builderUpdateUser = builderUpdateUser.Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderUpdateUser.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", fn, err)
		return nil, err
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("%s: failed to update user: %v", fn, err)
		return nil, err
	}

	log.Printf("%s: updated %d rows", fn, res.RowsAffected())
	return &emptypb.Empty{}, nil
}

// DeleteUser deletes user by id
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	fn := "DeleteUser"
	log.Printf("[%s] request data | id: %v", fn, req.Id)

	builderDeletUser := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDeletUser.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", fn, err)
		return nil, err
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("%s: failed to delete user: %v", fn, err)
		return nil, err
	}

	log.Printf("%s: deleted %d row", fn, res.RowsAffected())

	return &emptypb.Empty{}, nil
}
