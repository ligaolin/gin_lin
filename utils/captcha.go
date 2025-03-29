package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin/global"
	"github.com/mojocn/base64Captcha"
)

type CaptchaData struct {
	Id     string
	Answer string
}
type Captcha struct {
	Driver base64Captcha.Driver
	Store  base64Captcha.Store
}

// 创建图片验证码
func CaptchaGenerate(c *gin.Context, w int, h int) (string, string, error) {
	if w == 0 {
		w = global.Config.Captcha.Width
	}
	if h == 0 {
		h = global.Config.Captcha.Height
	}
	captcha := &Captcha{Driver: base64Captcha.NewDriverString(
		h,                                  // 高度
		w,                                  // 宽度
		global.Config.Captcha.NoiseCount,   // 噪点数量
		base64Captcha.OptionShowHollowLine, // 显示线条选项
		global.Config.Captcha.Length,       // 验证码长度
		base64Captcha.TxtNumbers+base64Captcha.TxtAlphabet, // 数据源
		nil,        // &color.RGBA{R: 255, G: 255, B: 0, A: 255}, &color.RGBA{R: 195, G: 245, B: 237, A: 255}// 背景颜
		nil,        // 字体存储（可以根据需要设置）
		[]string{}, // 字体列表
	)}

	id, content, answer := captcha.Driver.GenerateIdQuestionAnswer()
	item, err := captcha.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	b64s := item.EncodeB64string()
	ClientSet(c, "captcha", CaptchaData{Id: id, Answer: answer}, time.Minute*time.Duration(global.Config.Captcha.Expir))
	return id, b64s, err
}

// 验证
func CaptchaVerify(c *gin.Context, id string, answer string, clear bool) error {
	v, err := ClientGet[CaptchaData](c, "captcha")
	if err != nil {
		return errors.New("验证码不存在或过期")
	}
	if clear {
		ClientClear(c, "captcha")
	}
	if strings.EqualFold(v.Answer, answer) && v.Id == id {
		return nil
	} else {
		return errors.New("验证码错误")
	}
}

func CaptchaEmailSend(c *gin.Context, e *Email, email string) error {
	code := Random(6)
	ClientSet(c, "email_captcha", code, time.Minute*time.Duration(global.Config.Captcha.Expir))
	return e.SendCode(email, code)
}

func CaptchaEmailVerify(c *gin.Context, code int, clear bool) error {
	v, err := ClientGet[int](c, "email_captcha")
	if err != nil {
		return errors.New("验证码不存在或过期")
	}
	if clear {
		ClientClear(c, "captcha")
	}
	if v == code {
		return nil
	} else {
		return errors.New("验证码错误")
	}
}
