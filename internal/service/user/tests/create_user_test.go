package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

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
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type authCacheMockFunc func(mc *minimock.Controller) cache.AuthCache
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = gofakeit.Int64()
		name            = gofakeit.Username()
		email           = gofakeit.Email()
		password        = gofakeit.Animal()
		passwordConfirm = password
		role            = gofakeit.RandomInt([]int{0, 1, 2})

		repoErr = fmt.Errorf("repo error")

		req = &model.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            int64(role),
		}

		cacheUser = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      int64(role),
			CreatedAt: time.Now(),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
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
			want: id,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(id, nil)
				mock.GetUserMock.Expect(ctx, id).Return(cacheUser, nil)
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
				mock.CreateUserMock.Expect(ctx, cacheUser).Return(nil)
				return mock
			},
		},
		{
			name: "error from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(0, repoErr)
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authCacheMock := tt.authCacheMock(mc)
			authRepoMock := tt.authRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := user.NewMockService(authRepoMock, txManagerMock, authCacheMock)

			newID, err := service.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
