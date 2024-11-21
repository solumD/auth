package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	accessTokenSecretKeyEnvName  = "ACCESS_TOKEN_SECRET_KEY"  // #nosec G101
	accessTokenExpEnvName        = "ACCESS_TOKEN_EXP"         // #nosec G101
	refreshTokenSecretKeyEnvName = "REFRESH_TOKEN_SECRET_KEY" // #nosec G101
	refreshTokenExpEnvName       = "REFRESH_TOKEN_EXP"        // #nosec G101
)

type authConfig struct {
	refreshTokenSecretKey []byte
	refreshTokenExp       time.Duration

	accessTokenSecretKey []byte
	accessTokenExp       time.Duration
}

// NewAuthConfig returns new auth service config
func NewAuthConfig() (AuthConfig, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, fmt.Errorf("refresh token secret not found")
	}

	refreshTokenExp, err := strconv.Atoi(os.Getenv(refreshTokenExpEnvName))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token exp. time")
	}

	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, fmt.Errorf("access token secret not found")
	}

	accessTokenExp, err := strconv.Atoi(os.Getenv(accessTokenExpEnvName))
	if err != nil {
		return nil, fmt.Errorf("invalid access token exp. time")
	}

	return &authConfig{
		refreshTokenSecretKey: []byte(refreshTokenSecretKey),
		refreshTokenExp:       time.Minute * time.Duration(refreshTokenExp),
		accessTokenSecretKey:  []byte(accessTokenSecretKey),
		accessTokenExp:        time.Minute * time.Duration(accessTokenExp),
	}, nil
}

func (cfg *authConfig) RefreshTokenSecretKey() []byte {
	return cfg.refreshTokenSecretKey
}

func (cfg *authConfig) RefreshTokenExp() time.Duration {
	return cfg.refreshTokenExp
}

func (cfg *authConfig) AccessTokenSecretKey() []byte {
	return cfg.accessTokenSecretKey
}

func (cfg *authConfig) AccessTokenExp() time.Duration {
	return cfg.accessTokenExp
}
