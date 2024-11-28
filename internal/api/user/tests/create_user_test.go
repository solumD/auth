package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/solumD/auth/internal/api/user"
	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/service"
	serviceMocks "github.com/solumD/auth/internal/service/mocks"
	desc "github.com/solumD/auth/pkg/user_v1"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mn *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = gofakeit.Int64()
		name            = gofakeit.Username()
		email           = gofakeit.Email()
		password        = gofakeit.Animal()
		passwordConfirm = password
		role            = 1

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateUserRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
		}

		info = &model.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            int64(role),
		}

		res = &desc.CreateUserResponse{
			Id: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
		{
			name: "error req is nil",
			args: args{
				ctx: ctx,
				req: nil,
			},
			want: nil,
			err:  errors.ErrDescUserIsNil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	logger.MockInit()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewAPI(userServiceMock)

			res, err := api.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
