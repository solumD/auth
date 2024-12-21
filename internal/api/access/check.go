package access

import (
	"context"

	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/access_v1"
)

// Check отправляет запрос в сервисный слой на проверку доступа пользователя к эндпоинту
func (a *API) Check(ctx context.Context, req *desc.CheckRequest) (*desc.CheckResponse, error) {
	username, err := a.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.PermissionDenied)
	}

	return &desc.CheckResponse{
		Username: username,
	}, nil
}
