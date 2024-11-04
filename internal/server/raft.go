package server

import (
	"fmt"
	"math/rand"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/hashicorp/raft"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var (
	clusterGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "current_leader",
		Help: "Current Leader in Cluster",
	}, []string{"node"})
)

func NewRaft(myAddress domain.ConnectionItem, clusterAddress []domain.ConnectionItem, fsm raft.FSM) (*raft.Raft, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myAddress.GetId())

	ldb, sdb := raft.NewInmemStore(), raft.NewInmemStore()

	fss := raft.NewInmemSnapshotStore()

	transport, err := raft.NewTCPTransportWithLogger(myAddress.RaftAddress(), nil, 3, 10*time.Second, &internalLogger{
		log: &log.Logger,
	})

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

	filteredCluster := lo.Filter(clusterAddress, func(item domain.ConnectionItem, _ int) bool {
		return item.GetId() != string(myAddress.GetId())
	})

	go scheduleLeader(r, raft.ServerID(myAddress.GetId()), filteredCluster)

	log.Info().
		Str("address", myAddress.RaftAddress()).
		Msg("Raft Cluster Instance started")

	return r, nil
}

func scheduleLeader(r *raft.Raft, myId raft.ServerID, cluster []domain.ConnectionItem) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if _, id := r.LeaderWithID(); id == myId {
			id := rand.Int() % len(cluster)

			futuer := r.LeadershipTransferToServer(raft.ServerID(cluster[id].GetId()), raft.ServerAddress(cluster[id].RaftAddress()))

			if err := futuer.Error(); err != nil {
				log.Err(err).
					Str("id", string(myId)).
					Msg("leader shift")
			}

			clusterGauge.WithLabelValues(string(myId)).Set(0)
			lo.ForEach(cluster, func(item domain.ConnectionItem, i int) {
				if i == id {
					clusterGauge.WithLabelValues(item.GetId()).Set(1)
				} else {
					clusterGauge.WithLabelValues(item.GetId()).Set(1)
				}
			})

			log.Info().
				Interface("to node", cluster[id].RaftAddress()).
				Interface("id", cluster[id].GetId()).
				Msg("leadership transfered")
		}
	}
}
