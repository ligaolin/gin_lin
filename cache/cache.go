package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
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
	Config     *CacheConfig
}

// 创建缓存
func NewCache(cfg *CacheConfig) *Cache {
	c := &Cache{
		RedisCache: NewRedis(cfg.Redis.Addr, cfg.Redis.Port, cfg.Redis.Password),
		FileCache:  NewFileCache(cfg.File.Path),
		Config:     cfg,
	}
	go c.CleanDiskCache()
	return c
}

// 获取缓存
func (c *Cache) Get(key string, t any) (err error) {
	if c.Config.Use == "redis" {
		s, err := c.RedisCache.Get(key)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(s), t)
		return err
	} else {
		fc, err := c.getFileCacheValue(key)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(fc.Value.(string)), t)
		return err
	}
}

// 设置缓存，expir小于等于-1为永久缓存
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

// 删除缓存
func (c *Cache) Delete(key string) error {
	if c.Config.Use == "redis" {
		return c.RedisCache.Delete(key)
	} else {
		c.FileCache.Delete(key)
		return nil
	}
}

// 获取文件缓存的数据
func (c *Cache) getFileCacheValue(key string) (fc fileCacheValue, err error) {
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

// 清理过期文件缓存
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
