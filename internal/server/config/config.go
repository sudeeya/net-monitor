package config

import "github.com/caarlos0/env"

type Config struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseDSN string `env:"DATABASE_DSN,required"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	LogFile     string `env:"LOG_FILE,required"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
