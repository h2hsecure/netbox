package ports

import (
	"context"
	"time"
)

type Cache interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string, time.Duration) error

	Inc(context.Context, string, int) (uint64, error)
	Dec(context.Context, string, int) (uint64, error)
}
