package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

const (
	defaultHTTPAddr         = ":8080"
	defaultShutdownTimeout  = 5 * time.Second
	defaultReadHeaderTimout = 5 * time.Second
)

type Config struct {
	HTTP struct {
		Addr            string        `env:"HTTP_ADDR" envDefault:":8080"`
		ReadHeader      time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"5s"`
	}
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = defaultHTTPAddr
	}
	if cfg.HTTP.ReadHeader == 0 {
		cfg.HTTP.ReadHeader = defaultReadHeaderTimout
	}
	if cfg.HTTP.ShutdownTimeout == 0 {
		cfg.HTTP.ShutdownTimeout = defaultShutdownTimeout
	}

	return cfg, nil
}



