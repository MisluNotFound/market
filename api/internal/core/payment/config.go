package payment

import "github.com/mislu/market-api/internal/core/payment/types"

// NewPaymentConfig 创建支付配置
func NewPaymentConfig() *types.PaymentConfig {
	return &types.PaymentConfig{
		Alipay: types.AlipayConfig{
			AppID:      "your-alipay-appid",
			PrivateKey: "your-alipay-private-key",
			NotifyURL:  "/payment/notify/alipay",
		},
		Wechat: types.WechatConfig{
			AppID:     "your-wechat-appid",
			MchID:     "your-wechat-mchid",
			ApiKey:    "your-wechat-apikey",
			CertPath:  "path/to/wechat/cert",
			KeyPath:   "path/to/wechat/key",
			NotifyURL: "/payment/notify/wechat",
		},
	}
}