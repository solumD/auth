package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

const (
	tableName = "users"

	nameColumn     = "name"
	passwordColumn = "password"
	roleColumn     = "role"
)

type repo struct {
	db db.Client
}

// NewRepository возвращает новый объект репо слоя
func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{
		db: db,
	}
}

// GetUser получает из БД информацию пользователя
func (r *repo) GetUser(ctx context.Context, name string) (*model.UserInfo, error) {
	exist, err := r.isExistByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("user with name %s doesn't exist", name)
	}

	userInfo, err := r.getByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// isExistById проверяет, существует ли в БД пользователь с указанным name
func (r *repo) isExistByName(ctx context.Context, name string) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "auth_repository.isExistByName",
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

// getByName получает юзера из БД по его name
func (r *repo) getByName(ctx context.Context, name string) (*model.UserInfo, error) {

	query, args, err := sq.Select(nameColumn, passwordColumn, roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.getByName",
		QueryRaw: query,
	}

	user := new(model.UserInfo)
	err = r.db.DB().ScanOneContext(ctx, user, q, args...)
	if err != nil {
		return user, err
	}

	return user, nil
}
