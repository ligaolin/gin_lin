package gin_lin

import (
	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Smtp     string `json:"smtp" toml:"smtp" yaml:"smtp"`
	Port     int    `json:"port" toml:"port" yaml:"port"`
	Email    string `json:"email" toml:"email" yaml:"email"`
	Password string `json:"password" toml:"password" yaml:"password"`
}

type Email struct {
	Dialer *gomail.Dialer // SMTP服务器
	Config EmailConfig
}

func NewEmail(cfg EmailConfig) *Email {
	return &Email{
		Dialer: gomail.NewDialer(cfg.Smtp, cfg.Port, cfg.Email, cfg.Password),
		Config: cfg,
	}
}

func (e *Email) Send(to []string, subject string, body string) error {
	// 创建邮件对象
	m := gomail.NewMessage()
	m.SetHeader("From", e.Config.Email) // 发件人
	m.SetHeader("To", to...)            // 收件人
	m.SetHeader("Subject", subject)     // 邮件主题
	m.SetBody("text/plain", body)       // 邮件正文

	return e.Dialer.DialAndSend(m)
}
