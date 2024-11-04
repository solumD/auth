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
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/repository"
	repoMocks "github.com/solumD/auth/internal/repository/mocks"
	"github.com/solumD/auth/internal/service/user"
	"github.com/solumD/auth/internal/validation"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
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
		password        = gofakeit.Password(true, true, true, false, false, 8)
		passwordConfirm = password
		role            = gofakeit.RandomInt([]int{0, 1, 2})

		repoErr            = fmt.Errorf("repo error")
		differentPassesErr = fmt.Errorf("password and passwordConfirm do not match")

		validReq = &model.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            int64(role),
		}

		nameWithSpacesReq = &model.User{
			Name:            gofakeit.Username() + " " + gofakeit.Username(),
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            int64(role),
		}

		invalidEmailReq = &model.User{
			Name:            gofakeit.Username(),
			Email:           "emailgmail",
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            int64(role),
		}

		shortPassReq = &model.User{
			Name:            name,
			Email:           email,
			Password:        "12345",
			PasswordConfirm: "12345",
			Role:            int64(role),
		}

		differentPassesReq = &model.User{
			Name:            name,
			Email:           email,
			Password:        "12345678",
			PasswordConfirm: "87654321",
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
				req: validReq,
			},
			want: id,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, validReq).Return(id, nil)
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
				req: validReq,
			},
			want: 0,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, validReq).Return(0, repoErr)
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
			name: "error name contains spaces",
			args: args{
				ctx: ctx,
				req: nameWithSpacesReq,
			},
			want: 0,
			err:  validation.ErrNameContainsSpaces,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				return mock
			},
		},
		{
			name: "error invalid email",
			args: args{
				ctx: ctx,
				req: invalidEmailReq,
			},
			want: 0,
			err:  validation.ErrInvalidEmail,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				return mock
			},
		},
		{
			name: "error short password",
			args: args{
				ctx: ctx,
				req: shortPassReq,
			},
			want: 0,
			err:  validation.ErrPassTooShort,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				return mock
			},
			authCacheMock: func(mc *minimock.Controller) cache.AuthCache {
				mock := cacheMocks.NewAuthCacheMock(mc)
				return mock
			},
		},
		{
			name: "error passwords do not match",
			args: args{
				ctx: ctx,
				req: differentPassesReq,
			},
			want: 0,
			err:  differentPassesErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
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