package cache

import (
	"context"
	"time"

	"github.com/gregjones/httpcache/diskcache"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	// 从缓存中获取数据
	Get(key string, value any) error

	// 将数据存入缓存
	Set(key string, value any, expir time.Duration) error

	// 从缓存中删除数据
	Delete(key string) error
}

type Config struct {
	Use   string       `json:"use" toml:"use" yaml:"use"`
	File  *FileConfig  `json:"file" toml:"file" yaml:"file"`
	Redis *RedisConfig `json:"redis" toml:"redis" yaml:"redis"`
}

type FileConfig struct {
	Path string `json:"path" toml:"path" yaml:"path"`
}

type RedisConfig struct {
	Addr     string `json:"addr" toml:"addr" yaml:"addr"`
	Password string `json:"password" toml:"password" yaml:"password"`
}

type Factory interface {
	New(config *Config) (Cache, error)
}

type FileFactory struct{}

func (f *FileFactory) New(config *Config) (Cache, error) {
	return &File{
		Client: diskcache.New(config.File.Path),
		Path:   config.File.Path,
	}, nil
}

type RedisFactory struct{}

func (f *RedisFactory) New(config *Config) (Cache, error) {
	return &Redis{
		Context: context.Background(),
		Client:  redis.NewClient(&redis.Options{Addr: config.Redis.Addr, Password: config.Redis.Password}),
	}, nil
}
