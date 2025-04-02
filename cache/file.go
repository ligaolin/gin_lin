package cache

import (
	"errors"
	"os"

	"github.com/gregjones/httpcache/diskcache"
)

type FileCache struct {
	Client *diskcache.Cache
	Path   string
}

func NewFileCache(path string) *FileCache {
	return &FileCache{
		Client: diskcache.New(path),
		Path:   path,
	}
}

func (f *FileCache) Get(key string) (string, error) {
	s, ok := f.Client.Get(key)
	if !ok {
		return "", errors.New("从文件缓存获取数据失败")
	}
	return string(s), nil
}

func (f *FileCache) Set(key string, value []byte) {
	f.Client.Set(key, value)
}

func (f *FileCache) Delete(key string) {
	f.Client.Delete(key)
}

// 删除全部
func (f *FileCache) DiskDeleteAll() error {
	return os.RemoveAll(f.Path)
}
