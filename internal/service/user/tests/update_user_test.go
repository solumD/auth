package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/solumD/auth/internal/cache"
	cacheMocks "github.com/solumD/auth/internal/cache/mocks"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/client/db/mocks"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/repository"
	repoMocks "github.com/solumD/auth/internal/repository/mocks"
	"github.com/solumD/auth/internal/service/user"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type authCacheMockFunc func(mc *minimock.Controller) cache.AuthCache
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UserUpdate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Username()
		email = gofakeit.Email()
		role  = gofakeit.RandomInt([]int{0, 1, 2})

		repoErr  = fmt.Errorf("repo error")
		cacheErr = fmt.Errorf("cache error")

		req = &model.UserUpdate{
			ID:    id,
			Name:  name,
			Email: email,
			Role:  int64(role),
		}

		res = &emptypb.Empty{}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *emptypb.Empty
		err                error
		authRepositoryMock authRepositoryMockFunc
		txManagerMock      txManagerMockFunc
		authCacheMock      authCacheMockFunc
	}{
		{
			name: "success from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(res, nil)
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
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
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
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(nil, repoErr)
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
			want: res,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(res, nil)
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
				mock.DeleteUserMock.Expect(ctx, id).Return(cacheErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authCacheMock := tt.authCacheMock(mc)
			authRepoMock := tt.authRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := user.NewMockService(authRepoMock, txManagerMock, authCacheMock)

			res, err := service.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
