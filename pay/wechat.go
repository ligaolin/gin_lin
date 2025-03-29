package pay

import (
	"context"
	"os"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/go-pay/xlog"
	"github.com/ligaolin/gin_lin/global"
	"github.com/ligaolin/gin_lin/utils"
)

type WechatMerchant struct {
	Client *wechat.ClientV3
}

func WechatMerchantClient() (wm WechatMerchant, err error) {
	// NewClientV3 初始化微信客户端 v3
	// mchid：商户ID 或者服务商模式的 sp_mchid
	// serialNo：商户证书的证书序列号
	// apiV3Key：apiV3Key，商户平台获取
	// privateKey：私钥 apiclient_key.pem 读取后的内容
	privateKey, err := os.ReadFile(global.Config.Pay.Wechat.PrivateKey)
	if err != nil {
		xlog.Error(err)
		return
	}
	wm.Client, err = wechat.NewClientV3(global.Config.Pay.Wechat.MchID, global.Config.Pay.Wechat.SerialNo, global.Config.Pay.Wechat.ApiV3Key, string(privateKey))
	if err != nil {
		xlog.Error(err)
		return
	}

	// 注意：以下两种自动验签方式二选一
	// 微信支付公钥自动同步验签（新微信支付用户推荐）
	// err = wm.Client.AutoVerifySignByPublicKey([]byte("微信支付公钥内容"), "微信支付公钥ID")
	// if err != nil {
	// 	xlog.Error(err)
	// 	return
	// }
	//// 微信平台证书自动获取证书+同步验签（并自动定时更新微信平台API证书）
	err = wm.Client.AutoVerifySign()
	if err != nil {
		xlog.Error(err)
		return
	}

	// 自定义配置http请求接收返回结果body大小，默认 10MB
	// wm.Client.SetBodySize() // 没有特殊需求，可忽略此配置

	// 设置自定义RequestId生成方法，非必须
	// wm.Client.SetRequestIdFunc()

	// 打开Debug开关，输出日志，默认是关闭的
	wm.Client.DebugSwitch = gopay.DebugOn
	return wm, nil
}

func (wm *WechatMerchant) NativePay(c context.Context, tradeNo string, description string, price float32, ip string) (*wechat.NativeRsp, error) {
	// 设置支付参数
	totalFee := int64(price * 100)
	bm := make(gopay.BodyMap)
	bm.Set("appid", global.Config.Pay.Wechat.AppID).
		Set("description", description).                       // 商品描述
		Set("out_trade_no", tradeNo).                          // 商户订单号
		Set("notify_url", global.Config.Pay.Wechat.NotifyUrl). // 支付结果通知地址
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", totalFee).
				Set("currency", "CNY")
		})
	return wm.Client.V3TransactionNative(c, bm)
}

// 生成商户订单号
func GenerateOutTradeNo() string {
	return utils.GenerateRandomAlphanumeric(22)
}
