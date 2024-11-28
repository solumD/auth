package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/solumD/auth/internal/cache"
	cacheMocks "github.com/solumD/auth/internal/cache/mocks"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/client/db/mocks"
	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/repository"
	repoMocks "github.com/solumD/auth/internal/repository/mocks"
	"github.com/solumD/auth/internal/service/user"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type UserRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type authCacheMockFunc func(mc *minimock.Controller) cache.AuthCache
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Username()
		email     = gofakeit.Email()
		role      = gofakeit.RandomInt([]int{0, 1, 2})
		createdAt = time.Now()

		repoErr  = fmt.Errorf("repo error")
		cacheErr = fmt.Errorf("cache error")

		req = id

		res = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      int64(role),
			CreatedAt: createdAt,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		UserRepositoryMock UserRepositoryMockFunc
		txManagerMock      txManagerMockFunc
		authCacheMock      authCacheMockFunc
	}{
		{
			name: "success from cache",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(res, nil)
				return mock
			},
		},
		{
			name: "error from cache, success from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(res, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(nil, cacheErr)
				mock.CreateUserMock.Expect(ctx, res).Return(nil)
				return mock
			},
		},
		{
			name: "error from cache, error from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repoErr,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(nil, repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(nil, cacheErr)
				return mock
			},
		},
	}

	logger.MockInit()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authCacheMock := tt.authCacheMock(mc)
			authRepoMock := tt.UserRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := user.NewMockService(authRepoMock, txManagerMock, authCacheMock)

			res, err := service.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
