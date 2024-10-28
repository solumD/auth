package redis

import (
	"context"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/solumD/auth/internal/cache"
	"github.com/solumD/auth/internal/cache/redis/converter"
	modelCache "github.com/solumD/auth/internal/cache/redis/model"
	cacheCl "github.com/solumD/auth/internal/client/cache"
	"github.com/solumD/auth/internal/model"
)

// время жизни записи
const redisTTL = 900 * time.Second

type redisCache struct {
	cl cacheCl.Client
}

// NewRedisCache возвращает объект кэша redis с клиентом
func NewRedisCache(cl cacheCl.Client) cache.AuthCache {
	return &redisCache{
		cl: cl,
	}
}

// CreateUser добавляет данные пользователя в кэш
func (c *redisCache) CreateUser(ctx context.Context, user *model.User) error {
	userIDstr := strconv.FormatInt(user.ID, 10)
	var updatedAt int64
	if user.UpdatedAt.Valid {
		updatedAt = user.UpdatedAt.Time.UnixNano()
	}

	u := modelCache.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAtNs: user.CreatedAt.UnixNano(),
		UpdatedAtNs: &updatedAt,
	}

	err := c.cl.HashSet(ctx, userIDstr, u)
	if err != nil {
		return err
	}

	// устанавливаем время жизни для записи
	err = c.cl.Expire(ctx, userIDstr, redisTTL)
	if err != nil {
		return err
	}

	return nil
}

// GetUser получает данные пользователя из кэша
func (c *redisCache) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	userIDstr := strconv.FormatInt(userID, 10)

	values, err := c.cl.HGetAll(ctx, userIDstr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, errors.Errorf("no user with id %d in cache", userID)
	}

	var user modelCache.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromCache(&user), nil
}

// DeleteUser удаляет данные пользователя из кэша
func (c *redisCache) DeleteUser(ctx context.Context, userID int64) error {
	userIDstr := strconv.FormatInt(userID, 10)

	err := c.cl.HDel(ctx, userIDstr)
	if err != nil {
		return err
	}

	return nil
}
