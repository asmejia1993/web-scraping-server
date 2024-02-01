package service

import (
	"context"
	"fmt"
	"time"

	"github.com/asmejia1993/web-scraping-server/pkg/config"
	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/repository"
)

type redisRepository struct {
	redis *config.RedisInfo
}

func NewRedisRepository(redis *config.RedisInfo) repository.IRedisRepository {
	return &redisRepository{redis: redis}
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redis.Client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("error getting key redis: %v", err)
	}
	return val, nil
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}) {
	r.redis.Client.Set(ctx, key, value, time.Hour*24)
}

func (r *redisRepository) Exist(ctx context.Context, key string) (int64, error) {
	return r.redis.Client.Exists(ctx, key).Result()
}
