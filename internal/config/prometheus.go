package config

import (
	"errors"
	"net"
	"os"
)

const (
	promHostEnvName = "PROMETHEUS_HOST"
	promPortEnvName = "PROMETHEUS_PORT"
)

type prometheusConfig struct {
	host string
	port string
}

// NewPrometheusConfig returns new grpc config
func NewPrometheusConfig() (PrometheusConfig, error) {
	host := os.Getenv(promHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(promPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusConfig{
		host: host,
		port: port,
	}, nil
}

// Address returns a full address of a server
func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
