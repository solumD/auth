package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/repository/user/converter"
	modelRepo "github.com/solumD/auth/internal/repository/user/model"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

// Структура репо с клиентом базы данных (интерфейсом)
type repo struct {
	db db.Client
}

// NewRepository возвращает новый объект репо слоя
func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{
		db: db,
	}
}

// CreateUser создает пользователя
func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	builderInsertUser := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(user.Name, user.Email, user.Password, user.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsertUser.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().ScanOneContext(ctx, &userID, q, args...)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUser возвращает пользователя по id
func (r *repo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	builderGetUser := sq.Select(nameColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		Limit(1)

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}

	var name string

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d doesn't exist", userID)
		}
		return nil, err
	}

	builderGetUser = sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		Limit(1)

	query, args, err = builderGetUser.ToSql()
	if err != nil {
		return nil, err
	}

	q = db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	u, err := converter.ToUserFromRepo(&user)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdateUser обновляет данные пользователя по id
func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	builderGetUser := sq.Select(nameColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: user.ID}).
		Limit(1)

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	var name string

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d doesn't exist", user.ID)
		}
		return nil, err
	}

	builderUpdateUser := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, user.Name).
		Set(emailColumn, user.Email).
		Set(roleColumn, user.Role).
		Where(sq.Eq{idColumn: user.ID})

	query, args, err = builderUpdateUser.ToSql()
	if err != nil {
		return nil, err
	}

	q = db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser удаляет пользователя по id
func (r *repo) DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error) {
	builderGetUser := sq.Select(nameColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		Limit(1)

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}

	var name string

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d doesn't exist", userID)
		}
		return nil, err
	}
	builderDeleteUser := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID})

	query, args, err = builderDeleteUser.ToSql()
	if err != nil {
		return nil, err
	}

	q = db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
