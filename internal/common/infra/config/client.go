// Package config contains server and client configuration
// objects and methods
package config

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"
)

// ClientConfig contains values of client flags and environments.
type ClientConfig struct {
	Address string `env:"ADDRESS" json:"address"`
	Cookie  *http.Cookie
}

// NewClientConfig returns new client config.
func NewClientConfig(ctx context.Context) *ClientConfig {
	return &ClientConfig{}
}

// ParseFlags handles and processes flags and environments values
// when launching the client.
func (cfg *ClientConfig) ParseFlags(ctx context.Context) error {
	flag.StringVar(&cfg.Address, "a", "http://localhost:8080", "HTTP-server endpoint address 'protocol://host:port'")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("ParseFlags: wrong environment values %w", err)
	}

	return nil
}
