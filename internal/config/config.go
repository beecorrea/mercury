package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBPath string `envconfig:"MERCURY_DB_PATH" default:"mercury.db"`
	Domain string `envconfig:"MERCURY_DOMAIN" default:"communist.mom"`
	Port   string `envconfig:"MERCURY_PORT" default:"45800"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("loading configuration: %w", err)
	}
	return &cfg, nil
}
