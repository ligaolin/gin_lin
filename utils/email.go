package utils

import (
	"fmt"

	"github.com/ligaolin/gin_lin/global"
	"gopkg.in/gomail.v2"
)

type Email struct {
	Dialer *gomail.Dialer // SMTP服务器
	Email  string
}

func EmailNew(email string, password string, smtp string, port int) *Email {
	return &Email{
		Dialer: gomail.NewDialer(smtp, port, email, password),
		Email:  email,
	}
}

func (e *Email) Send(to []string, subject string, body string) error {
	// 创建邮件对象
	m := gomail.NewMessage()
	m.SetHeader("From", e.Email)    // 发件人
	m.SetHeader("To", to...)        // 收件人
	m.SetHeader("Subject", subject) // 邮件主题
	m.SetBody("text/plain", body)   // 邮件正文

	return e.Dialer.DialAndSend(m)
}

func (e *Email) SendCode(to string, code int32) error {
	body := fmt.Sprintf(`尊敬的用户：

您好！
您正在进行邮箱验证操作，验证码为：%d。
此验证码有效期为 %d分钟，请尽快完成验证。

如非本人操作，请忽略此邮件。

感谢您的支持！

【系统邮件，请勿直接回复】`, code, global.Config.Captcha.Expir)
	return e.Send([]string{to}, "悟品网站-邮箱验证", body)
}
