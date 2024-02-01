package config

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Host    string
	DB      int
	Expires time.Duration
	Client  redis.Client
}

func (cache *RedisCache) NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.Host,
		Password: "",
		DB:       cache.DB,
	})
}
