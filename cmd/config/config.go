package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	DSN           string `env:"DATABASE_DSN" envDefault:""`
}

func GetConfig() (*Config, error) {

	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddress, "f", cfg.ServerAddress, "SERVER_ADDRESS")
	flag.StringVar(&cfg.DSN, "d", cfg.DSN, "DATABASE_DSN")

	flag.Parse()
	err := env.Parse(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
