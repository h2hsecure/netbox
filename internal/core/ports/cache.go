package ports

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(context.Context, string, string, time.Duration) error

	Inc(ctx context.Context, key string, delta int) (uint64, error)
	Dec(ctx context.Context, key string, delta int) (uint64, error)
}
