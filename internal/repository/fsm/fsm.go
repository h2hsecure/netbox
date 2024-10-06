package fsm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
)

type StateMachine struct {
	cache ports.Cache
}

func NewStateMachine(cache ports.Cache) raft.FSM {
	return &StateMachine{
		cache: cache,
	}
}

// Apply implements raft.FSM.
func (s *StateMachine) Apply(l *raft.Log) interface{} {
	log.Info().Interface("raft log", l).Msg("log came")

	var userIpTime domain.UserIpTime
	err := json.Unmarshal(l.Data, &userIpTime)

	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	ctx := context.Background()

	count, err := s.cache.Inc(ctx, "c:"+userIpTime.User, 1)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	if count > 100 {
		s.cache.Set(ctx, userIpTime.User, "1", 10*time.Second)
		s.cache.Set(ctx, "c:"+userIpTime.User, "0", 0)
	}
	s.cache.Inc(ctx, "c:"+userIpTime.Ip, 1)

	//s.cache.Set(ctx, userIpTime.User, "1")

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
	return nil, nil
}
