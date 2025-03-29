package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/ligaolin/gin_lin/global"
	"github.com/redis/go-redis/v9"
)

func RedisInit() (context.Context, *redis.Client) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.Cache.Redis.Addr, global.Config.Cache.Redis.Port),
		Password: global.Config.Cache.Redis.Password, // no password set
		DB:       0,                                  // use default DB
	})
	return ctx, rdb
}

func RedisSet(key string, value interface{}, expir time.Duration) error {
	ctx, rdb := RedisInit()
	return rdb.Set(ctx, key, value, expir).Err()
}

func RedisGet(key string) (string, error) {
	ctx, rdb := RedisInit()
	return rdb.Get(ctx, key).Result()
}

func RedisDelete(key string) error {
	ctx, rdb := RedisInit()
	return rdb.Del(ctx, key).Err()
}
