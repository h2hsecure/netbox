package ports

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error

	Inc(ctx context.Context, key string, delta int) (uint64, error)
	Dec(ctx context.Context, key string, delta int) (uint64, error)
}
