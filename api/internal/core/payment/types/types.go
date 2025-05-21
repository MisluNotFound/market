package types

import "context"

// PaymentType 支付类型枚举
type PaymentType string

const (
	Alipay PaymentType = "alipay"
	Wechat PaymentType = "wechat"
)

// PaymentRequest 支付请求参数
type PaymentRequest struct {
	Type        PaymentType // 支付类型
	Amount      string      // 金额(元)
	OrderID     string      // 订单ID
	Subject     string      // 商品标题
	Description string      // 商品描述
	ReturnURL   string      // 回调URL
}

// PaymentResponse 支付响应
type PaymentResponse struct {
	PaymentURL string // 支付跳转链接
	QRCode     string // 支付二维码(可选)
}

// PaymentConfig 支付模块配置
type PaymentConfig struct {
	Alipay AlipayConfig `yaml:"alipay"`
	Wechat WechatConfig `yaml:"wechat"`
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID      string `yaml:"app_id"`
	PrivateKey string `yaml:"private_key"`
	NotifyURL  string `yaml:"notify_url"`
}

// WechatConfig 微信支付配置
type WechatConfig struct {
	AppID     string `yaml:"app_id"`
	MchID     string `yaml:"mch_id"`
	ApiKey    string `yaml:"api_key"`
	CertPath  string `yaml:"cert_path"`
	KeyPath   string `yaml:"key_path"`
	NotifyURL string `yaml:"notify_url"`
}

// RefundRequest 退款请求参数
type RefundRequest struct {
	Type     PaymentType // 支付类型
	Amount   float64     // 退款金额(元)
	OrderID  string      // 原订单ID
	RefundID string      // 退款单ID
	Reason   string      // 退款原因
}

// TransferRequest 转账请求参数
type TransferRequest struct {
	Type      PaymentType // 支付类型
	Amount    float64     // 转账金额(元)
	OutPayNo  string      // 商户转账单号
	Payee     string      // 收款方账号
	PayeeName string      // 收款方姓名
	Remark    string      // 转账备注
}

// PaymentService 支付服务接口
type PaymentService interface {
	Pay(ctx context.Context, req PaymentRequest) (*PaymentResponse, error)
	Refund(ctx context.Context, req RefundRequest) error
	Transfer(ctx context.Context, req TransferRequest) error
}
