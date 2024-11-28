package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/solumD/auth/internal/cache"
	cacheMocks "github.com/solumD/auth/internal/cache/mocks"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/client/db/mocks"
	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/repository"
	repoMocks "github.com/solumD/auth/internal/repository/mocks"
	"github.com/solumD/auth/internal/service/user"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDeleteUser(t *testing.T) {
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

		id = gofakeit.Int64()

		repoErr  = fmt.Errorf("repo error")
		cacheErr = fmt.Errorf("cache error")

		req = id

		res = &emptypb.Empty{}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *emptypb.Empty
		err                error
		UserRepositoryMock UserRepositoryMockFunc
		txManagerMock      txManagerMockFunc
		authCacheMock      authCacheMockFunc
	}{
		{
			name: "success from repo, success from cache",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(res, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			name: "error from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repoErr,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(nil, repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				return mock
			},
		},
		{
			name: "success from repo, error from cache",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(res, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(cacheErr)
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

			res, err := service.DeleteUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
