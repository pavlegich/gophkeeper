// Package config contains server and client configuration
// objects and methods
package config

import (
	"context"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

// ServerConfig contains values of server flags and environments.
type ServerConfig struct {
	Address string `env:"ADDRESS" json:"address"`
	DSN     string `env:"DATABASE_DSN" json:"database_dsn"`
}

// NewServerConfig returns new server config.
func NewServerConfig(ctx context.Context) *ServerConfig {
	return &ServerConfig{}
}

// ParseFlags handles and processes flags and environments values
// when launching the server.
func (cfg *ServerConfig) ParseFlags(ctx context.Context) error {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP-server endpoint address host:port")
	flag.StringVar(&cfg.DSN, "d", "postgresql://localhost:5432/postgres", "URI (DSN) to database")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("ParseFlags: wrong environment values %w", err)
	}

	return nil
}
