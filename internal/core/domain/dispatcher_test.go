package domain_test

import (
	"context"
	"testing"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/rs/zerolog/log"
)

type testPayload struct {
	count int
}

func (t *testPayload) Send(ctx context.Context) error {
	log.Info().Msg("test Playload")
	t.count++

	return nil
}

func TestDispatcher(t *testing.T) {

	dispatcher := domain.NewDispatcher(10, 100)
	dispatcher.Run()
	defer dispatcher.Close()

	payload := &testPayload{}

	for i := 0; i < 10000; i++ {
		dispatcher.Push(payload)
	}

	time.Sleep(5 * time.Second)

	if payload.count != 10000 {
		t.Fatalf("count is different: %d", payload.count)
	}
}
