package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/solumD/auth/internal/api/user"
	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/service"
	serviceMocks "github.com/solumD/auth/internal/service/mocks"
	desc "github.com/solumD/auth/pkg/auth_v1"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mn *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Username()
		email = gofakeit.Email()
		role  = desc.Role(gofakeit.RandomInt([]int{0, 1, 2}))

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateUserRequest{
			Id:    id,
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(email),
			Role:  role,
		}

		info = &model.UserUpdate{
			ID:    id,
			Name:  &name,
			Email: &email,
			Role:  int64(role),
		}

		res = &emptypb.Empty{}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, info).Return(&emptypb.Empty{}, nil)
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
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, info).Return(nil, serviceErr)
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
			err:  errors.ErrDescUserUpdateIsNil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			api := user.NewAuthAPI(authServiceMock)

			res, err := api.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
