package interceptor

import (
	"context"
	"time"

	"github.com/solumD/auth/internal/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// LogInterceptor логирует содержание запроса
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("req", req))
	}

	logger.Info("request", zap.String("method", info.FullMethod), zap.Any("req", req), zap.Any("res", res), zap.Duration("duration", time.Since(now)))

	return res, err
}
