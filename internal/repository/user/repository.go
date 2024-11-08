package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/repository/user/converter"
	modelRepo "github.com/solumD/auth/internal/repository/user/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
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
	exist, err := r.IsNameExist(ctx, user.Name)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, fmt.Errorf("user with name %s already exists", user.Name)
	}

	exist, err = r.IsEmailExist(ctx, user.Email)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, fmt.Errorf("user with email %s already exists", user.Email)
	}

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(user.Name, user.Email, user.Password, user.Role).
		Suffix("RETURNING id").ToSql()

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
	exist, err := r.IsExistByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user with id %d doesn't exist", userID)
	}

	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	u := converter.ToUserFromRepo(&user)
	if u == nil {
		return nil, fmt.Errorf("convertion failed, user model is nil")
	}

	return u, nil
}

// UpdateUser обновляет данные пользователя по id
func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	exist, err := r.IsExistByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user with id %d doesn't exist", user.ID)
	}

	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(roleColumn, user.Role)

	if user.Name != nil {
		builderUpdate = builderUpdate.Set(nameColumn, *user.Name)
	}

	if user.Email != nil {
		builderUpdate = builderUpdate.Set(emailColumn, *user.Email)
	}

	builderUpdate = builderUpdate.Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
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
	exist, err := r.IsExistByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user with id %d doesn't exist", userID)
	}

	query, args, err := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetById получает юзера из БД по его id
func (r *repo) GetByID(ctx context.Context, userID int64) (modelRepo.User, error) {
	var user modelRepo.User

	query, args, err := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		Limit(1).ToSql()

	if err != nil {
		return user, err
	}

	q := db.Query{
		Name:     "user_repository.GetById",
		QueryRaw: query,
	}

	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return user, err
	}

	return user, nil
}

// IsExistById проверяет, существует ли в БД пользователь с указанным ID
func (r *repo) IsExistByID(ctx context.Context, userID int64) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsExistById",
		QueryRaw: query,
	}

	var one int

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// IsEmailExist проверяет, существует ли в БД указанный email
func (r *repo) IsEmailExist(ctx context.Context, email string) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{emailColumn: email}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsEmailExist",
		QueryRaw: query,
	}

	var one int

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// IsNameExist проверяет, существует ли в БД указанный name
func (r *repo) IsNameExist(ctx context.Context, name string) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsNameExist",
		QueryRaw: query,
	}

	var one int

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
