package server

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	client "git.h2hsecure.com/ddos/waf/internal/repository/grpc"
	"git.h2hsecure.com/ddos/waf/internal/server/handler"

	"google.golang.org/grpc"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func CreateGrpcServer(port string, serverHandler *handler.ServerHandler) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Err(err).Msg("failed to listen")
		return err
	}

	s := grpc.NewServer()
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s, healthcheck)

	client.RegisterNetworkEventServer(s, serverHandler)

	go leaderHealthCheck(serverHandler, healthcheck)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func leaderHealthCheck(serverHandler *handler.ServerHandler, healthCheck *health.Server) {
	for b := range serverHandler.AmILeader() {
		next := healthgrpc.HealthCheckResponse_NOT_SERVING

		if b {
			next = healthgrpc.HealthCheckResponse_SERVING
		}

		healthCheck.SetServingStatus("", next)
	}
}
