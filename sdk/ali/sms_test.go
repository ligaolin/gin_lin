package ali

import (
	"testing"
)

func TestSms(t *testing.T) {
	client, err := NewAliSms(&AliSmsConfig{
		AccessKeyId:                  "",
		AccessKeySecret:              "",
		TemplateCodeVerificationCode: "",
		SignName:                     "",
	})
	if err != nil {
		t.Error(err)
	}
	code, err := client.SendVerificationCode("")
	if err != nil {
		t.Error(err)
	}
	t.Log(code)
}
