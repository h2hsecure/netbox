package cache

import (
	"context"
	"errors"
	"fmt"

	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"github.com/bradfitz/gomemcache/memcache"
)

type Impl struct {
	m *memcache.Client
}

func NewMemcache(servers ...string) ports.Cache {
	mc := memcache.New(servers...)

	return &Impl{
		m: mc,
	}
}

func (i *Impl) Set(ctx context.Context, key, value string) error {
	err := i.m.Set(&memcache.Item{Key: key, Value: []byte(value)})

	if err != nil {
		return fmt.Errorf("memcache-set: %w", err)
	}

	return nil
}

func (i *Impl) Get(ctx context.Context, key string) (string, error) {
	item, err := i.m.Get(key)

	if err != nil && err == memcache.ErrCacheMiss {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("memcache-get: %w", err)
	}

	return string(item.Value), nil
}

func (i *Impl) Inc(ctx context.Context, key string, delta int) (uint64, error) {
	last, err := i.m.Increment(key, uint64(delta))

	if err != nil && errors.Is(err, memcache.ErrCacheMiss) {
		i.Set(ctx, key, "0")
		last = 0
	} else if err != nil {
		return 0, fmt.Errorf("memcache-increment: %w", err)
	}

	return last, nil
}

func (i *Impl) Dec(ctx context.Context, key string, delta int) (uint64, error) {
	last, err := i.m.Decrement(key, uint64(delta))

	if err != nil && errors.Is(err, memcache.ErrCacheMiss) {
		i.Set(ctx, key, "0")
		last = 0
	} else if err != nil {
		return 0, fmt.Errorf("memcache-decrement: %w", err)
	}

	return last, nil
}
