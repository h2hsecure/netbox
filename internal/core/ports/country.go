package ports

import (
	"context"
	"net/netip"
)

type CountryAdater interface {
	FindCountryByIp(context.Context, netip.Addr) (string, error)
	Close() error
}
