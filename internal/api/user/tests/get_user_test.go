package tests

import (
	"context"
	"testing"
	"time"

	"github.com/solumD/auth/internal/api/user"
	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/service"
	serviceMocks "github.com/solumD/auth/internal/service/mocks"
	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/user_v1"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mn *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Username()
		email     = gofakeit.Email()
		role      = desc.Role(gofakeit.RandomInt([]int{0, 1, 2}))
		createdAt = time.Now()

		serviceErr = sys.NewCommonError("service error", codes.Aborted)

		req = &desc.GetUserRequest{
			Id: id,
		}

		info = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      int64(role),
			CreatedAt: createdAt,
		}

		res = &desc.GetUserResponse{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: timestamppb.New(createdAt),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetUserResponse
		err             error
		UserServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			UserServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(info, nil)
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
			UserServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "error user is nil",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  sys.NewCommonError(errors.ErrUserModelIsNil.Error(), codes.Internal),
			UserServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, nil)
				return mock
			},
		},
	}

	logger.MockInit()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			UserServiceMock := tt.UserServiceMock(mc)
			api := user.NewAPI(UserServiceMock)

			res, err := api.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
