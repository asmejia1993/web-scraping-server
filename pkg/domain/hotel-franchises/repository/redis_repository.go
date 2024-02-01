package repository

import (
	"context"
)

type IRedisRepository interface {
	Set(ctx context.Context, key string, value interface{})
	Get(ctx context.Context, key string) (string, error)
	Exist(ctx context.Context, key string) (int64, error)
}
