package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/ligaolin/gin_lin/cache"
)

type Client struct {
	Cache cache.Cache
}

func NewClient(cfg cache.CacheConfig) *Client {
	return &Client{
		Cache: *cache.NewCache(cfg),
	}
}

func (c *Client) Set(k string, v any, expir time.Duration) (string, error) {
	uuid := uuid.New().String()
	c.Cache.Set("client-"+k+uuid, v, expir)
	return uuid, nil
}

func (c *Client) Get(uuid string, k string, t any, clear bool) error {
	err := c.Cache.Get("client-"+k+uuid, t)
	if err != nil && clear {
		c.Cache.Delete("client-" + k + uuid)
	}
	return err
}
