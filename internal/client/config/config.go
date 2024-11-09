// Package config defines client config.
package config

import (
	"time"

	"github.com/caarlos0/env"
)

// Config describes client config.
type Config struct {
	ServerAddr   string        `env:"SERVER_ADDR" envDefault:"localhost:9090"`
	TargetsFile  string        `env:"TARGETS_FILE,required"`
	SnapInterval time.Duration `env:"SNAP_INTERVAL" envDefault:"10m"`
	LogLevel     string        `env:"LOG_LEVEL" envDefault:"INFO"`
	LogFile      string        `env:"LOG_FILE"`
}

// NewConfig returns client config.
// The function parses environment variables to form config.
func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
