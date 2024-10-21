package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/repository/user/converter"
	modelRepo "github.com/solumD/auth/internal/repository/user/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "username"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.AuthRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	builderInsertUser := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn).
		Values(user.Name, user.Email, user.Password, user.Role, user.CreatedAt).
		Suffix("RETURNING id")

	query, args, err := builderInsertUser.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *repo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	builderGetUser := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID})

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var user modelRepo.User

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	builderUpdateUser := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar)

	if len(user.Name) > 0 {
		builderUpdateUser = builderUpdateUser.Set(nameColumn, user.Name)
	}

	if len(user.Email) > 0 {
		builderUpdateUser = builderUpdateUser.Set(emailColumn, user.Email)
	}

	if user.Role >= 0 && user.Role <= 2 {
		builderUpdateUser = builderUpdateUser.Set(roleColumn, user.Role)
	}

	builderUpdateUser = builderUpdateUser.Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderUpdateUser.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (r *repo) DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error) {
	builderDeletUser := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID})

	query, args, err := builderDeletUser.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
