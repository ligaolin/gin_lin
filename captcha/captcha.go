package captcha

import (
	"github.com/ligaolin/gin_lin/cache"
)

type CaptchaConfig struct {
	Expir int64 `json:"expir" toml:"expir" yaml:"expir"` // 过期时间
}

type Captcha struct {
	Client *cache.Client
	Config *CaptchaConfig
}

func NewCaptcha(cfg *CaptchaConfig, cacheCfg *cache.CacheConfig) *Captcha {
	return &Captcha{
		Client: cache.NewClient(cacheCfg),
		Config: cfg,
	}
}
