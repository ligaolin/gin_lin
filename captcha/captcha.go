package captcha

import (
	"time"

	"github.com/ligaolin/gin_lin/cache"
)

type InputCaptcha interface {
	Generate(carrier string, expir time.Duration) (string, error)
	Verify(carrier string, uuid, code string) error
}

type OutputCaptcha interface {
	Generate(expir time.Duration) (string, string, error)
	Verify(uuid, code string) error
}

type CaptchaConfig struct {
	Expir int64 `json:"expir" toml:"expir" yaml:"expir"` // 过期时间
}

type CaptchaFactory interface {
	New(cfg *CaptchaConfig, c cache.Cache) InputCaptcha
	New(cfg *CaptchaConfig, c cache.Cache) OutputCaptcha
}

// func NewCaptcha(cfg *CaptchaConfig, c cache.Cache) *Captcha {
// 	return &Captcha{
// 		Client: cache.NewClient(c),
// 		Config: cfg,
// 	}
// }

type Value struct {
	Code    string
	Carrier string
}
