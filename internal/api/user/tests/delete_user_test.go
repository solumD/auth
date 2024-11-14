package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/solumD/auth/internal/api/user"
	"github.com/solumD/auth/internal/service"
	serviceMocks "github.com/solumD/auth/internal/service/mocks"
	desc "github.com/solumD/auth/pkg/user_v1"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mn *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.DeleteUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		req = &desc.DeleteUserRequest{
			Id: id,
		}

		res = &emptypb.Empty{}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
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
				mock.DeleteUserMock.Expect(ctx, id).Return(&emptypb.Empty{}, nil)
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
				mock.DeleteUserMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			UserServiceMock := tt.UserServiceMock(mc)
			api := user.NewAPI(UserServiceMock)

			res, err := api.DeleteUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
