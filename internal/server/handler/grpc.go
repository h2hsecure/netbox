package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/h2hsecure/netbox/internal/core/domain"
	client "github.com/h2hsecure/netbox/internal/repository/grpc"
	"github.com/hashicorp/raft"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	opsProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "total_event_request",
		Help: "The total number of processed events",
	}, []string{"user", "ip"})
)

type ServerHandler struct {
	client.UnimplementedNetworkEventServer
	raft *raft.Raft
}

func NewGrpcHandler(
	raft *raft.Raft) *ServerHandler {
	return &ServerHandler{
		raft: raft,
	}
}

func (s *ServerHandler) Send(ctx context.Context, userIpTime *client.UserIpTime) (*client.Empty, error) {
	log.Info().
		Interface("userIpTime", userIpTime).
		Msg("message get")
	opsProcessed.WithLabelValues("user", userIpTime.User).Inc()
	opsProcessed.WithLabelValues("ip", userIpTime.Ip).Inc()

	buf, err := json.Marshal(domain.UserIpTime{
		Ip:        userIpTime.Ip,
		User:      userIpTime.User,
		Path:      userIpTime.Path,
		Timestamp: int64(userIpTime.Timestamp),
	})

	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	future := s.raft.Apply(buf, 3*time.Second)

	if err := future.Error(); err != nil {
		return nil, fmt.Errorf("server send command: %w", err)
	}

	return &client.Empty{}, nil
}

func (s *ServerHandler) AmILeader() <-chan bool {
	return s.raft.LeaderCh()
}
