package cmd

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/h2hsecure/netbox/internal/core/domain"
)

func CurrentConfig() (domain.ConfigParams, error) {

	var cfg domain.ConfigParams

	if err := env.Parse(&cfg); err != nil {
		return domain.ConfigParams{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
