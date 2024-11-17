package access

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/solumD/auth/internal/utils/jwt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	authHeader = "authorization"
	authPrefix = "Bearer "
)

// Check проверяет, имеет ли пользователь доступ к эндпоинту
func (s *srv) Check(ctx context.Context, endpointAddress string) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

	authHeader, ok := md[authHeader]
	if !ok || len(authHeader) == 0 {
		return nil, fmt.Errorf("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := jwt.VerifyToken(accessToken, s.authConfig.AccessTokenSecretKey())
	if err != nil {
		return nil, fmt.Errorf("access token is invalid")
	}

	log.Println(claims.Role)
	// админ имеет доступ ко всем эндпоинтам
	if claims.Role == 2 {
		return &emptypb.Empty{}, nil
	}

	log.Println(endpointAddress)
	log.Println(s.userAccesses)
	// смотрим, есть ли доступ у пользователя
	_, ok = s.userAccesses[endpointAddress]
	if !ok {
		return nil, fmt.Errorf("access denied")
	}

	return &emptypb.Empty{}, nil
}
