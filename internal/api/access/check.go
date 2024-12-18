package access

import (
	"context"

	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/access_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Check отправляет запрос в сервисный слой на проверку доступа пользователя к эндпоинту
func (a *API) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	res, err := a.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.PermissionDenied)
	}

	return res, nil
}
