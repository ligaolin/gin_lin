package captcha

import (
	"errors"
	"fmt"
	"time"

	"github.com/ligaolin/gin_lin"
	"github.com/ligaolin/gin_lin/email"
)

func (c *Captcha) EmailSend(email string, cfg *email.EmailConfig, subject string) (string, error) {
	code := gin_lin.Random(6)
	uuid, err := c.Client.Set("captcha-email", code, time.Minute*time.Duration(c.Config.Expir))
	if err != nil {
		return "", err
	}

	err = c.SendEmailCode(cfg, email, code, subject)
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

func (e *Captcha) SendEmailCode(cfg *email.EmailConfig, to string, code int32, subject string) error {
	return email.NewEmail(cfg).Send([]string{to}, subject, fmt.Sprintf(`尊敬的用户：

	您好！
	您正在进行邮箱验证操作，验证码为：%d。
	此验证码有效期为 %d分钟，请尽快完成验证。
	
	如非本人操作，请忽略此邮件。
	
	感谢您的支持！
	
	【系统邮件，请勿直接回复】`, code, e.Config.Expir))
}
