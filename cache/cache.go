package cache

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

type CacheConfig struct {
	Use   string     `json:"use" toml:"use" yaml:"use"`
	Redis CacheRedis `json:"redis" toml:"redis" yaml:"redis"`
	File  CacheFile  `json:"file" toml:"file" yaml:"file"`
}
type CacheRedis struct {
	Addr     string `json:"addr" toml:"addr" yaml:"addr"`
	Port     int    `json:"port" toml:"port" yaml:"port"`
	Password string `json:"password" toml:"password" yaml:"password"`
}
type CacheFile struct {
	Path string `json:"path" toml:"path" yaml:"path"`
}

type fileCacheValue struct {
	Expir   time.Time
	IsExpir bool // 是否会过期，true会过期，false不会过期
	Value   any
}

type Cache struct {
	RedisCache *RedisCache
	FileCache  *FileCache
	Config     CacheConfig
}

func NewCache(cfg CacheConfig) *Cache {
	return &Cache{
		RedisCache: NewRedis(cfg.Redis.Addr, cfg.Redis.Port, cfg.Redis.Password),
		FileCache:  NewFileCache(cfg.File.Path),
		Config:     cfg,
	}
}

func (c *Cache) Get(key string, t any) (err error) {
	if c.Config.Use == "redis" {
		s, err := c.RedisCache.Get(key)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(s), t)
		return err
	} else {
		fc, err := c.Expir(key)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(fc.Value.(string)), t)
		return err
	}
}

func (c *Cache) Set(key string, value any, expir time.Duration) error {
	s, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if c.Config.Use == "redis" {
		return c.RedisCache.Set(key, string(s), expir)
	} else {
		is_expir := true
		if expir <= -1 {
			is_expir = false
		}
		v, err := json.Marshal(fileCacheValue{Value: string(s), Expir: time.Now().Add(expir), IsExpir: is_expir})
		if err != nil {
			return err
		}
		c.FileCache.Set(key, v)
		return nil
	}
}

func (c *Cache) Delete(key string) error {
	if c.Config.Use == "redis" {
		return c.RedisCache.Delete(key)
	} else {
		c.FileCache.Delete(key)
		return nil
	}
}

func (c *Cache) Expir(key string) (fc fileCacheValue, err error) {
	s, err := c.FileCache.Get(key)
	if err != nil {
		return fc, err
	}

	if err = json.Unmarshal([]byte(s), &fc); err != nil {
		return fc, err
	}

	// 过期删除
	if fc.IsExpir {
		if time.Now().After(fc.Expir) {
			c.FileCache.Delete(key)
			return fc, errors.New("缓存数据已过期")
		}
	}

	return fc, nil
}

func (c *Cache) CleanDiskCache() {
	path := c.Config.File.Path
	files, err := os.ReadDir(path)
	var fc fileCacheValue
	if err != nil {
		return
	}
	for _, file := range files {
		file_path := filepath.Join(path, file.Name())
		b, err := os.ReadFile(file_path)
		if err != nil {
			continue
		}

		if err = json.Unmarshal(b, &fc); err != nil {
			continue
		}

		if fc.IsExpir {
			if time.Now().After(fc.Expir) {
				os.Remove(file_path)
			}
		}
	}
}

func (c *Cache) CleanDiskCacheCron() {
	if c.Config.Use == "file" {
		cron := cron.New()

		// @daily 每天凌晨 0 点
		// @every 1m 每分钟
		_, err := cron.AddFunc("@daily", func() {
			// 清除文件缓存
			c.CleanDiskCache()
		})
		if err != nil {
			log.Printf("Failed to add cron job: %v", err)
			return
		}

		// 启动 Cron
		cron.Start()
	}
}
