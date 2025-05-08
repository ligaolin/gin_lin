package captcha

import (
	"errors"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
)

type ImageCaptchaConfig struct {
	Width      int `json:"width" toml:"width" yaml:"width"`
	Height     int `json:"height" toml:"height" yaml:"height"`
	Length     int `json:"length" toml:"length" yaml:"length"`
	NoiseCount int `json:"noise_count" toml:"noise_count" yaml:"noise_count"` // 噪点数量
}

// 创建图片验证码
func (c *Captcha) GenerateImageCode(cfg *ImageCaptchaConfig) (string, string, error) {
	driver := base64Captcha.NewDriverString(
		cfg.Height,                         // 高度
		cfg.Width,                          // 宽度
		cfg.NoiseCount,                     // 噪点数量
		base64Captcha.OptionShowHollowLine, // 显示线条选项
		cfg.Length,                         // 验证码长度
		base64Captcha.TxtNumbers+base64Captcha.TxtAlphabet, // 数据源
		nil,        // &color.RGBA{R: 255, G: 255, B: 0, A: 255}, &color.RGBA{R: 195, G: 245, B: 237, A: 255}// 背景颜
		nil,        // 字体存储（可以根据需要设置）
		[]string{}, // 字体列表
	)

	_, content, answer := driver.GenerateIdQuestionAnswer()
	item, err := driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	b64s := item.EncodeB64string()
	uuid, err := c.Client.Set("captcha-image", answer, time.Minute*time.Duration(c.Config.Expir))
	return uuid, b64s, err
}

// 验证
func (c *Captcha) VerifyImageCode(uuid string, answer string, clear bool) error {
	var v string
	if err := c.Client.Get(uuid, "captcha-image", &v, clear); err != nil {
		return errors.New("验证码不存在或过期")
	}
	if strings.EqualFold(v, answer) {
		return nil
	} else {
		return errors.New("验证码错误")
	}
}
