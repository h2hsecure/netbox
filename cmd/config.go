package cmd

import (
	"fmt"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/caarlos0/env/v11"
)

func CurrentConfig() (domain.ConfigParams, error) {

	var cfg domain.ConfigParams

	if err := env.Parse(&cfg); err != nil {
		return domain.ConfigParams{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
