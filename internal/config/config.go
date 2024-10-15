package config

import (
	"github.com/joho/godotenv"
)

// Load reads ,env file from path and loads
// variables into a project
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}
