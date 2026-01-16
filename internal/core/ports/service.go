package ports

import (
	"context"
	"net/netip"

	"github.com/h2hsecure/netbox/internal/core/domain"
)

type Service interface {
	OpenSession(context.Context, domain.UserIpTime) (string, error)
	AccessAtempt(context.Context, string, domain.AttemptRequest) domain.AttemptOperation

	ManageIps(context.Context, []netip.Addr, domain.IpOperation) error
	QueryIp(context.Context, netip.Addr) (domain.IpOperation, error)

	CurrentConfig(context.Context) domain.ConfigParams
	Health(context.Context) error
	Stop() error
}
