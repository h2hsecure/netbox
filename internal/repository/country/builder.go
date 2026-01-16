package country

import (
	"context"
	"fmt"

	"net/netip"

	"github.com/h2hsecure/netbox/internal/core/domain"
	"github.com/h2hsecure/netbox/internal/core/ports"
	"github.com/oschwald/maxminddb-golang/v2"
	"github.com/rs/zerolog/log"
)

type countryDatabaseImpl struct {
	db *maxminddb.Reader
}

var record struct {
	CountryISOCode string `maxminddb:"country_name"`
}

func CreateCountryDatabae(cfg domain.CountryDbParams) (ports.CountryAdater, error) {
	if cfg.FilePath == "" {
		return nil, nil
	}

	log.Info().
		Str("path", cfg.FilePath).
		Msg("country database loading")

	db, err := maxminddb.Open(cfg.FilePath)
	if err != nil {
		return nil, fmt.Errorf("opening mmdb file: %w", err)
	}

	return &countryDatabaseImpl{
		db: db,
	}, nil

}

func (c *countryDatabaseImpl) FindCountryByIp(_ context.Context, ip netip.Addr) (string, error) {
	result := c.db.Lookup(ip)

	if !result.Found() {
		return "", domain.ErrNotFound
	}
	err := result.Decode(&record)
	if err != nil {
		return "", fmt.Errorf("decoding result: %w", err)
	}

	return record.CountryISOCode, nil

}

func (c *countryDatabaseImpl) Close() error {
	return c.db.Close()
}
