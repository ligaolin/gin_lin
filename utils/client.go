package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin/cache"
)

func ClientSet(c *gin.Context, k string, v interface{}, expir time.Duration) error {
	key, err := setKey(c, k)
	if err != nil {
		return err
	}
	return cache.CacheSet(key, v, expir)
}

func ClientGet[T cache.Base](c *gin.Context, k string) (t T, err error) {
	key, err := setKey(c, k)
	if err != nil {
		return t, err
	}
	return cache.CacheGet[T](key)
}

func setKey(c *gin.Context, k string) (string, error) {
	userAgent := c.Request.Header.Get("User-Agent")
	ip := c.ClientIP()
	return k + "-" + userAgent + "-" + ip, nil
}

func ClientClear(c *gin.Context, k string) error {
	key, err := setKey(c, k)
	if err != nil {
		return err
	}
	return cache.CacheDelete(key)
}
