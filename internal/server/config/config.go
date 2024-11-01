package config

import "github.com/caarlos0/env"

type Config struct {
	HTTPAddr    string `env:"HTTP_ADDR" envDefault:"localhost:8080"`
	GRPCAddr    string `env:"GRPC_ADDR" envDefault:"localhost:9090"`
	DatabaseDSN string `env:"DATABASE_DSN,required"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	LogFile     string `env:"LOG_FILE"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
