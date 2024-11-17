package config

import (
	"errors"
	"os"
	"strings"
)

const (
	userEndpointsEnvName = "USER_ENDPOINTS"
)

type accessConfig struct {
	userAccesses map[string]struct{}
}

// NewAccessConfig returns new access config
func NewAccessConfig() (AccessConfig, error) {
	accessStr := os.Getenv(userEndpointsEnvName)
	if len(accessStr) == 0 {
		return nil, errors.New("user accesses endpoints not found")
	}

	accesses := strings.Split(accessStr, ",")
	uMap := make(map[string]struct{})
	for _, e := range accesses {
		uMap[strings.TrimSpace(e)] = struct{}{}
	}

	return &accessConfig{
		userAccesses: uMap,
	}, nil
}

// UserAccessesMap returns map of endpoints allowed to users
func (cfg *accessConfig) UserAccessesMap() (map[string]struct{}, error) {
	if cfg.userAccesses == nil {
		return nil, errors.New("user access map is nil")
	}

	return cfg.userAccesses, nil
}
