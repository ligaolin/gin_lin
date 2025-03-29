package cache

import (
	"errors"
	"os"

	"github.com/gregjones/httpcache/diskcache"
	"github.com/ligaolin/gin_lin/global"
)

func DiskInit() *diskcache.Cache {
	return diskcache.New(global.Config.Cache.File.Path)
}

func DiskGet(key string) (string, error) {
	c := DiskInit()
	s, ok := c.Get(key)
	if !ok {
		return "", errors.New("从文件缓存获取数据失败")
	}
	return string(s), nil
}

func DiskSet(key string, value []byte) {
	c := DiskInit()
	c.Set(key, value)
}

func DiskDelete(key string) {
	c := DiskInit()
	c.Delete(key)
}

// 删除全部
func DiskDeleteAll() error {
	return os.RemoveAll(global.Config.Cache.File.Path)
}
