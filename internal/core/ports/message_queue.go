package ports

import (
	"context"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
)

type MessageQueue interface {
	Sent(ctx context.Context, event domain.UserIpTime) error
}
