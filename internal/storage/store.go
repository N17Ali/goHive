package storage

import (
	"context"
)

// Store interface abstracts data storage
type Store interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
}
