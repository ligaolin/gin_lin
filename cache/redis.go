package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Context context.Context
	Client  *redis.Client
}

func NewRedis(addr string, port int, password string) *RedisCache {
	return &RedisCache{
		Context: context.Background(),
		Client:  redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", addr, port), Password: password, DB: 0}),
	}
}

func (r *RedisCache) Set(key string, value any, expir time.Duration) error {
	return r.Client.Set(r.Context, key, value, expir).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.Client.Get(r.Context, key).Result()
}

func (r *RedisCache) Delete(key string) error {
	return r.Client.Del(r.Context, key).Err()
}
