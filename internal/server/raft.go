package server

import (
	"fmt"
	"os"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/hashicorp/raft"
	"github.com/samber/lo"
)

func NewRaft(myAddress domain.ConnectionItem, clusterAddress []domain.ConnectionItem, fsm raft.FSM) (*raft.Raft, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myAddress.GetId())

	ldb, sdb := raft.NewInmemStore(), raft.NewInmemStore()

	fss := raft.NewInmemSnapshotStore()

	transport, err := raft.NewTCPTransport(myAddress.RaftAddress(), nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, transport)
	if err != nil {
		return nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	cfg := raft.Configuration{
		Servers: lo.Map(clusterAddress, func(item domain.ConnectionItem, _ int) raft.Server {
			return raft.Server{
				Suffrage: raft.Voter,
				ID:       raft.ServerID(item.GetId()),
				Address:  raft.ServerAddress(item.RaftAddress()),
			}
		}),
	}
	f := r.BootstrapCluster(cfg)
	if err := f.Error(); err != nil {
		return nil, fmt.Errorf("raft.Raft.BootstrapCluster: %w", err)
	}

	return r, nil
}
