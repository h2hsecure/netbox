package server

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/raft"
)

func NewRaft(myID, myAddress string, fsm raft.FSM) (*raft.Raft, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)

	ldb, sdb := raft.NewInmemStore(), raft.NewInmemStore()

	fss := raft.NewInmemSnapshotStore()

	transport, err := raft.NewTCPTransport(myAddress, nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, transport)
	if err != nil {
		return nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	cfg := raft.Configuration{
		Servers: []raft.Server{
			{
				Suffrage: raft.Voter,
				ID:       raft.ServerID(myID),
				Address:  raft.ServerAddress(myAddress),
			},
		},
	}
	f := r.BootstrapCluster(cfg)
	if err := f.Error(); err != nil {
		return nil, fmt.Errorf("raft.Raft.BootstrapCluster: %w", err)
	}

	return r, nil
}
