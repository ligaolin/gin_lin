package captcha

import (
	"errors"
	"fmt"
	"time"

	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ligaolin/gin_lin"
	"github.com/ligaolin/gin_lin/sdk/ali"
)

func (c *Captcha) SmsSend(cfg *ali.AliSmsConfig, phone string) (string, error) {
	value := Value{
		Code:    gin_lin.Random(6),
		Carrier: phone,
	}
	uuid, err := c.Client.Set("captcha-sms", value, time.Minute*time.Duration(c.Config.Expir))
	if err != nil {
		return "", err
	}

	err = c.SendSmsCode(cfg, phone, value.Code)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (c *Captcha) SmsCodeVerify(uuid string, code int32, clear bool) (string, error) {
	var val Value
	if err := c.Client.Get(uuid, "captcha-sms", &val, clear); err != nil {
		return "", errors.New("验证码不存在或过期")
	}
	if val.Code == code {
		return val.Carrier, nil
	} else {
		return "", errors.New("验证码错误")
	}
}

func (as Captcha) SendSmsCode(cfg *ali.AliSmsConfig, mobile string, code int32) error {
	alisms, err := ali.NewAliSms(cfg)
	if err != nil {
		return err
	}
	if _, err := alisms.Client.SendSmsWithOptions(&dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(mobile),
		SignName:      tea.String(cfg.SignName),
		TemplateCode:  tea.String(cfg.TemplateCodeVerificationCode),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":"%d"}`, code)),
	}, &service.RuntimeOptions{}); err != nil {
		return err
	}
	return nil
}
