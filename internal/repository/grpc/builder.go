package grpc

import (
	"context"
	"fmt"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

var serviceConfig = `{
	"loadBalancingPolicy": "round_robin",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

type clientImpl struct {
	pb NetworkEventClient
}

// Sent implements ports.MessageQueue.
func (c *clientImpl) Sent(ctx context.Context, event domain.UserIpTime) error {
	log.Info().Msg("send event client")
	_, err := c.pb.Send(ctx, &UserIpTime{
		Ip:        event.Ip,
		User:      event.User,
		Timestamp: event.Timestamp,
	})

	if err != nil {
		return fmt.Errorf("sent client data: %w", err)
	}

	return nil
}

func NewEnforceClient(address []domain.ConnectionItem) (ports.MessageQueue, error) {
	r := manual.NewBuilderWithScheme("whatever")
	r.InitialState(resolver.State{
		Addresses: lo.Map(address, func(item domain.ConnectionItem, _ int) resolver.Address {
			return resolver.Address{
				Addr: item.GrpcAddress(),
			}
		}),
	})

	nullAddress := fmt.Sprintf("%s:///unused", r.Scheme())

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(r),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	conn, err := grpc.NewClient(nullAddress, options...)
	if err != nil {
		log.Err(err).Msgf("grpc.NewClient(%q)", nullAddress)
		return nil, err
	}

	log.Info().Msg("starting grpc client")

	pb := NewNetworkEventClient(conn)

	return &clientImpl{
		pb: pb,
	}, nil
}
