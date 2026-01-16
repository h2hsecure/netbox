package grpc

import (
	"context"

	"github.com/h2hsecure/netbox/internal/core/domain"
	"github.com/h2hsecure/netbox/internal/core/ports"
)

type mockMq struct {
	store map[string]domain.UserIpTime
}

// Sent implements ports.MessageQueue.
func (m *mockMq) Sent(ctx context.Context, userIpTime domain.UserIpTime) error {
	m.store[userIpTime.User] = userIpTime
	return nil
}

func (m *mockMq) Get(ctx context.Context, id string) (domain.UserIpTime, error) {
	return m.store[id], nil
}

func NewMockMq() ports.MessageQueue {
	return &mockMq{
		store: make(map[string]domain.UserIpTime),
	}
}
