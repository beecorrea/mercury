package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBPath             string        `envconfig:"MERCURY_DB_PATH" default:"mercury.db"`
	Domain             string        `envconfig:"MERCURY_DOMAIN" default:"communist.mom"`
	Port               string        `envconfig:"MERCURY_PORT" default:"45800"`
	// ScrapeRetryInterval controls how often the retrier polls for failed scrapes.
	ScrapeRetryInterval time.Duration `envconfig:"SCRAPE_RETRY_INTERVAL" default:"30s"`
	// ScrapeMaxAttempts is the maximum number of retry attempts before a link is
	// considered permanently un-scrapable and excluded from future retry cycles.
	ScrapeMaxAttempts  int           `envconfig:"SCRAPE_MAX_ATTEMPTS" default:"5"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("loading configuration: %w", err)
	}
	return &cfg, nil
}
