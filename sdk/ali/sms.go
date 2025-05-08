package ali

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ligaolin/gin_lin"
)

type AliSms struct {
	Client *dysmsapi20170525.Client
	Config *AliSmsConfig
}

type AliSmsConfig struct {
	AccessKeyId                  string
	AccessKeySecret              string
	TemplateCodeVerificationCode string
	SignName                     string
}

func NewAliSms(cfg *AliSmsConfig) (*AliSms, error) {
	client, err := dysmsapi20170525.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyId),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	})
	return &AliSms{Client: client, Config: cfg}, err
}

func (as AliSms) SendVerificationCode(phone string) (int32, error) {
	code := gin_lin.Random(4)
	if _, err := as.Client.SendSmsWithOptions(&dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(as.Config.SignName),
		TemplateCode:  tea.String(as.Config.TemplateCodeVerificationCode),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":"%d"}`, code)),
	}, &service.RuntimeOptions{}); err != nil {
		return 0, err
	}
	return code, nil
}
