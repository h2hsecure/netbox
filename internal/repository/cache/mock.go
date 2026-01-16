package cache

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/h2hsecure/netbox/internal/core/domain"
	"github.com/h2hsecure/netbox/internal/core/ports"
)

var _ ports.Cache = &MockCache{}

type MockCache struct {
	mx *sync.Mutex
	m  map[string]string
}

// Dec implements ports.Cache.
func (m *MockCache) Dec(ctx context.Context, key string, delta int) (uint64, error) {
	value, err := m.Get(ctx, key)

	if err != nil {
		return 0, err
	}

	intValue, err := strconv.ParseInt(value, 0, 10)

	if err != nil {
		return 0, domain.ErrNotFound
	}
	intValue -= int64(delta)

	return uint64(intValue), m.Set(ctx, key, strconv.FormatInt(intValue, 10), 0*time.Second)
}

// Get implements ports.Cache.
func (m *MockCache) Get(_ context.Context, key string) (string, error) {
	m.mx.Lock()
	defer m.mx.Unlock()
	value, has := m.m[key]

	if !has {
		return "", domain.ErrNotFound
	}

	return value, nil
}

// Inc implements ports.Cache.
func (m *MockCache) Inc(ctx context.Context, key string, delta int) (uint64, error) {
	value, err := m.Get(ctx, key)

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return 0, err
	}

	if errors.Is(err, domain.ErrNotFound) {
		value = "0"
		if err := m.Set(ctx, key, "0", 0*time.Second); err != nil {
			return 0, err
		}
	}

	intValue, err := strconv.ParseInt(value, 0, 10)

	if err != nil {
		return 0, domain.ErrNotFound
	}
	intValue += int64(delta)

	return uint64(intValue), m.Set(ctx, key, strconv.FormatInt(intValue, 10), 0*time.Second)
}

// Set implements ports.Cache.
func (m *MockCache) Set(_ context.Context, key string, value string, _ time.Duration) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.m[key] = value

	return nil
}

// Set implements ports.Cache.
func (m *MockCache) Clear() error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.m = make(map[string]string)

	return nil
}

func CreateMockCache() *MockCache {
	return &MockCache{
		mx: &sync.Mutex{},
		m:  make(map[string]string),
	}
}
