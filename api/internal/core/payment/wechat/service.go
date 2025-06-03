package wechat

import (
	"context"
	"net/url"

	"github.com/mislu/market-api/internal/core/payment/types"
)

type wechatService struct {
	// 微信支付配置参数
	appID    string
	mchID    string
	apiKey   string
	certPath string
	keyPath  string
}

func NewWechatService(cfg types.WechatConfig) (*wechatService, error) {
	return &wechatService{
		appID:    cfg.AppID,
		mchID:    cfg.MchID,
		apiKey:   cfg.ApiKey,
		certPath: cfg.CertPath,
		keyPath:  cfg.KeyPath,
	}, nil
}

func (s *wechatService) Pay(ctx context.Context, req types.PaymentRequest) (*types.PaymentResponse, error) {
	// if req.Amount <= 0 {
	// 	return nil, types.ErrInvalidAmount
	// }

	// 调用微信支付API创建交易
	// 这里简化实现，实际应该调用微信支付SDK
	paymentURL := "https://pay.weixin.qq.com/wxpay/pay?" +
		"appid=" + s.appID +
		"&mch_id=" + s.mchID +
		"&out_trade_no=" + req.OrderID +
		"&body=" + url.QueryEscape(req.Subject)

	return &types.PaymentResponse{
		PaymentURL: paymentURL,
		QRCode:     "https://wechat.com/qr/" + req.OrderID,
	}, nil
}

func (s *wechatService) Refund(ctx context.Context, req types.RefundRequest) error {
	if req.Amount <= 0 {
		return types.ErrInvalidAmount
	}

	// 调用微信支付API发起退款
	// 这里简化实现，实际应该调用微信支付SDK
	// 需要证书文件(s.certPath, s.keyPath)
	// 检查订单状态、退款金额等
	// 返回退款结果

	return nil
}

func (s *wechatService) Transfer(ctx context.Context, req types.TransferRequest) error {
	if req.Amount <= 0 {
		return types.ErrInvalidAmount
	}

	// 调用微信支付API发起企业付款
	// 这里简化实现，实际应该调用微信支付SDK
	// 需要证书文件(s.certPath, s.keyPath)
	// 检查收款方账户、转账金额等
	// 返回转账结果

	return nil
}
