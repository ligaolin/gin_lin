package cache

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ligaolin/gin_lin/global"
	"github.com/robfig/cron/v3"
)

type Base interface {
}

type fileCache struct {
	Expir   time.Time
	IsExpir bool // 是否会过期，true会过期，false不会过期
	Value   interface{}
}

func CacheGet[T Base](key string) (t T, err error) {
	if global.Config.Cache.Use == "redis" {
		s, err := RedisGet(key)
		if err != nil {
			return t, err
		}
		err = json.Unmarshal([]byte(s), &t)
		return t, err
	} else {
		fc, err := diskExpir(key)
		if err != nil {
			return t, err
		}
		err = json.Unmarshal([]byte(fc.Value.(string)), &t)
		return t, err
	}
}

func CacheSet(key string, value interface{}, expir time.Duration) error {
	s, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if global.Config.Cache.Use == "redis" {
		return RedisSet(key, string(s), expir)
	} else {
		is_expir := true
		if expir <= -1 {
			is_expir = false
		}
		v, err := json.Marshal(fileCache{Value: string(s), Expir: time.Now().Add(expir), IsExpir: is_expir})
		if err != nil {
			return err
		}
		DiskSet(key, v)
		return nil
	}
}

func CacheDelete(key string) error {
	if global.Config.Cache.Use == "redis" {
		return RedisDelete(key)
	} else {
		DiskDelete(key)
		return nil
	}
}

func diskExpir(key string) (fc fileCache, err error) {
	s, err := DiskGet(key)
	if err != nil {
		return fc, err
	}

	if err = json.Unmarshal([]byte(s), &fc); err != nil {
		return fc, err
	}

	// 过期删除
	if fc.IsExpir {
		if time.Now().After(fc.Expir) {
			DiskDelete(key)
			return fc, errors.New("缓存数据已过期")
		}
	}

	return fc, nil
}

func CleanDiskCache() {
	path := global.Config.Cache.File.Path
	files, err := os.ReadDir(path)
	var fc fileCache
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

func CleanDiskCacheCron() {
	if global.Config.Cache.Use == "file" {
		c := cron.New()

		// @daily 每天凌晨 0 点
		// @every 1m 每分钟
		_, err := c.AddFunc("@daily", func() {
			// 清除文件缓存
			CleanDiskCache()
		})
		if err != nil {
			log.Printf("Failed to add cron job: %v", err)
			return
		}

		// 启动 Cron
		c.Start()
	}
}
