package ports

import (
	"context"

	"github.com/h2hsecure/netbox/internal/core/domain"
)

type MessageQueue interface {
	Sent(context.Context, domain.UserIpTime) error
}
