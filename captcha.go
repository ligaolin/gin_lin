package gin_lin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ligaolin/gin_lin/cache"
	"github.com/mojocn/base64Captcha"
)

type CaptchaConfig struct {
	Expir      int64 `json:"expir" toml:"expir" yaml:"expir"` // 过期时间
	Width      int   `json:"width" toml:"width" yaml:"width"`
	Height     int   `json:"height" toml:"height" yaml:"height"`
	Length     int   `json:"length" toml:"length" yaml:"length"`
	NoiseCount int   `json:"noise_count" toml:"noise_count" yaml:"noise_count"` // 噪点数量
}

type Captcha struct {
	Client *cache.Client
	Config CaptchaConfig
}

func NewCaptcha(cfg CaptchaConfig, cacheCfg cache.CacheConfig) *Captcha {
	return &Captcha{
		Client: cache.NewClient(cacheCfg),
		Config: cfg,
	}
}

// 创建图片验证码
func (c *Captcha) GenerateImageCode() (string, string, error) {
	driver := base64Captcha.NewDriverString(
		c.Config.Height,                    // 高度
		c.Config.Width,                     // 宽度
		c.Config.NoiseCount,                // 噪点数量
		base64Captcha.OptionShowHollowLine, // 显示线条选项
		c.Config.Length,                    // 验证码长度
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

func (c *Captcha) EmailSend(email string, cfg EmailConfig, subject string) (string, error) {
	code := Random(6)
	uuid, err := c.Client.Set("captcha-email", code, time.Minute*time.Duration(c.Config.Expir))
	if err != nil {
		return "", err
	}

	err = c.SendCode(cfg, email, code, subject)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (c *Captcha) EmailCodeVerify(uuid string, code int32, clear bool) error {
	var val int32
	if err := c.Client.Get(uuid, "captcha-email", &val, clear); err != nil {
		return errors.New("验证码不存在或过期")
	}
	if val == code {
		return nil
	} else {
		return errors.New("验证码错误")
	}
}

func (e *Captcha) SendCode(cfg EmailConfig, to string, code int32, subject string) error {
	return NewEmail(cfg).Send([]string{to}, subject, fmt.Sprintf(`尊敬的用户：

	您好！
	您正在进行邮箱验证操作，验证码为：%d。
	此验证码有效期为 %d分钟，请尽快完成验证。
	
	如非本人操作，请忽略此邮件。
	
	感谢您的支持！
	
	【系统邮件，请勿直接回复】`, code, e.Config.Expir))
}
