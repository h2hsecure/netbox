package fsm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/h2hsecure/netbox/internal/core/domain"
	"github.com/h2hsecure/netbox/internal/core/ports"
	"github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
)

type StateMachine struct {
	mx                      *sync.Mutex
	cache                   ports.Cache
	maxUser, maxIp, maxPath uint64
}

func NewStateMachine(cfg domain.CacheParams, cache ports.Cache) raft.FSM {
	return &StateMachine{
		cache:   cache,
		mx:      &sync.Mutex{},
		maxUser: cfg.MaxUser,
		maxIp:   cfg.MaxIp,
		maxPath: cfg.MaxPath,
	}
}

// Apply implements raft.FSM.
func (s *StateMachine) Apply(l *raft.Log) interface{} {
	log.Info().
		Interface("raft log", l).
		Msg("log came")

	var userIpTime domain.UserIpTime

	err := json.Unmarshal(l.Data, &userIpTime)

	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	ctx := context.Background()

	count, err := s.cache.Inc(ctx, "c:"+userIpTime.User, 1)

	if err != nil {
		return fmt.Errorf("cache user: %w", err)
	}

	if s.maxUser != 0 && count > s.maxUser {
		s.cache.Set(ctx, userIpTime.User, "1", 10*time.Second)
		s.cache.Set(ctx, "c:"+userIpTime.User, "0", 0)
	}

	count, err = s.cache.Inc(ctx, "c:"+userIpTime.Ip, 1)

	if err != nil {
		return fmt.Errorf("cache ip: %w", err)
	}

	if s.maxIp != 0 && count > s.maxIp {
		s.cache.Set(ctx, userIpTime.Ip, "1", 10*time.Second)
		s.cache.Set(ctx, "c:"+userIpTime.Ip, "0", 0)
	}

	count, err = s.cache.Inc(ctx, "c:"+userIpTime.Path, 1)

	if err != nil {
		return fmt.Errorf("cache path: %w", err)
	}

	if s.maxPath != 0 && count > s.maxPath {
		s.cache.Set(ctx, userIpTime.Path, "1", 10*time.Second)
		s.cache.Set(ctx, "c:"+userIpTime.Path, "0", 0)
	}

	return nil
}

// Restore implements raft.FSM.
func (s *StateMachine) Restore(snapshot io.ReadCloser) error {
	log.Info().Msg("restore")
	return nil
}

// Snapshot implements raft.FSM.
func (s *StateMachine) Snapshot() (raft.FSMSnapshot, error) {
	log.Info().Msg("snapshot")
	return s, nil
}

func (s *StateMachine) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (s *StateMachine) Release() {

}
