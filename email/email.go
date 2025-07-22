package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Smtp     string `json:"smtp" toml:"smtp" yaml:"smtp"`
	Port     int    `json:"port" toml:"port" yaml:"port"`
	Email    string `json:"email" toml:"email" yaml:"email"`
	Password string `json:"password" toml:"password" yaml:"password"`
	FromName string `json:"from_name" toml:"from_name" yaml:"from_name"` // 发件人名称
}

type Email struct {
	dialer *gomail.Dialer
	config *EmailConfig
}

func New(cfg *EmailConfig) (*Email, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return &Email{
		dialer: gomail.NewDialer(cfg.Smtp, cfg.Port, cfg.Email, cfg.Password),
		config: cfg,
	}, nil
}

func validateConfig(cfg *EmailConfig) error {
	if cfg.Smtp == "" {
		return fmt.Errorf("smtp服务器地址必须")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("无效的端口号")
	}
	if cfg.Email == "" {
		return fmt.Errorf("邮箱地址必须")
	}
	return nil
}

func (e *Email) Send(to []string, subject, body string, opts ...Option) error {
	// 默认选项
	options := &emailOptions{
		contentType: "text/plain",
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(options)
	}

	m := gomail.NewMessage()

	// 设置发件人
	from := e.config.Email
	if e.config.FromName != "" {
		m.SetHeader("From", m.FormatAddress(from, e.config.FromName))
	} else {
		m.SetHeader("From", from)
	}

	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody(options.contentType, body)

	// 设置抄送和密送
	if len(options.cc) > 0 {
		m.SetHeader("Cc", options.cc...)
	}
	if len(options.bcc) > 0 {
		m.SetHeader("Bcc", options.bcc...)
	}

	// 添加附件
	for _, attachment := range options.attachments {
		m.Attach(attachment)
	}

	// 发送邮件
	if err := e.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}

// 选项模式配置
type emailOptions struct {
	contentType string   // 内容类型
	cc          []string // 抄送
	bcc         []string // 密送
	attachments []string // 附件路径
}

type Option func(*emailOptions)

func WithHTML() Option {
	return func(o *emailOptions) {
		o.contentType = "text/html"
	}
}

func WithCC(cc []string) Option {
	return func(o *emailOptions) {
		o.cc = cc
	}
}

func WithBCC(bcc []string) Option {
	return func(o *emailOptions) {
		o.bcc = bcc
	}
}

func WithAttachments(files []string) Option {
	return func(o *emailOptions) {
		o.attachments = files
	}
}
