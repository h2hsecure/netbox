package service_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/service"
)

func TestDispatcher(t *testing.T) {

	dispatcher := service.NewDispatcher(10, 10000)
	dispatcher.Run()

	var ops atomic.Uint64

	for range 10000 {
		dispatcher.Push(func(ctx context.Context) error {
			ops.Add(1)
			return nil
		})
	}

	time.Sleep(5 * time.Second)

	dispatcher.Close()

	if ops.Load() != 10000 {
		t.Fatalf("count is different: %d", ops.Load())
	}
}
