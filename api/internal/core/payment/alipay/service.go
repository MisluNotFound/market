package alipay

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mislu/market-api/internal/core/payment/types"
	"github.com/mislu/market-api/internal/utils/app"
	"github.com/smartwalle/alipay/v3"
)

type AlipayService struct {
	client *alipay.Client
}

func NewAlipayClient() (*AlipayService, error) {
	privateKey, err := os.ReadFile("D:\\repository\\m-market\\api\\keys\\private.txt")
	if err != nil {
		return nil, err
	}
	publicKey, err := os.ReadFile("D:\\repository\\m-market\\api\\keys\\public.txt")
	if err != nil {
		return nil, err
	}

	aliPubKey, err := os.ReadFile("D:\\repository\\m-market\\api\\keys\\ali_public.txt")
	if err != nil {
		return nil, err
	}

	cli, err := alipay.New(app.GetConfig().Alipay.APPID, string(privateKey), false)
	if err != nil {
		return nil, err
	}

	err = cli.LoadAliPayPublicKey(string(publicKey))
	if err != nil {
		return nil, err
	}
	err = cli.LoadAliPayPublicKey(string(aliPubKey))
	if err != nil {
		return nil, err
	}

	return &AlipayService{
		client: cli,
	}, nil
}

func (s *AlipayService) Pay(ctx context.Context, req types.PaymentRequest) (*types.PaymentResponse, error) {
	resp := &types.PaymentResponse{}
	if s.client == nil {
		return resp, fmt.Errorf("alipay client not initialized")
	}

	pay := alipay.TradePagePay{
		Trade: alipay.Trade{
			Subject:     req.Subject,
			OutTradeNo:  req.OrderID,
			TotalAmount: req.Amount,
			ProductCode: "FAST_INSTANT_TRADE_PAY",
			// NotifyURL:   app.GetConfig().Server.BaseIP + "/api/order/alipay/notify",
		},
	}

	form, err := s.client.TradePagePay(pay)
	if err != nil {
		return resp, fmt.Errorf("generate pay form failed: %v", err)
	}

	resp.PaymentURL = form.String()

	return resp, nil
}

// Refund 发起退款
func (s *AlipayService) Refund(ctx context.Context, req types.RefundRequest) error {
	if s.client == nil {
		return fmt.Errorf("alipay client not initialized")
	}

	outRequestNo := "REFUND_" + randomString(16)
	refund := alipay.TradeRefund{
		OutTradeNo:   req.OrderID,
		RefundAmount: strconv.FormatFloat(req.Amount, 'f', -1, 64),
		OutRequestNo: outRequestNo,
		RefundReason: req.Reason,
	}

	result, err := s.client.TradeRefund(context.Background(), refund)
	if err != nil {
		return fmt.Errorf("refund failed: %v", err)
	}

	if result.Code != alipay.CodeSuccess {
		return fmt.Errorf("refund failed: %s", result.SubMsg)
	}

	return nil
}

// QueryTrade 查询订单状态
func (s *AlipayService) QueryTrade(outTradeNo string) (*types.QueryTradeResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("alipay client not initialized")
	}

	query := alipay.TradeQuery{
		OutTradeNo: outTradeNo,
	}

	// "crypto/rsa: verification error"
	result, err := s.client.TradeQuery(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("query trade failed: %v", err)
	}

	if result.Code != alipay.CodeSuccess {
		return nil, fmt.Errorf("query trade failed: %s", result.SubMsg)
	}

	return &types.QueryTradeResponse{
		TradeNo:      result.TradeNo,
		OutTradeNo:   result.OutTradeNo,
		BuyerLogonID: result.BuyerLogonId,
		BuyerOpenID:  result.BuyerOpenId,
		TradeStatus:  types.TradeStatus(result.TradeStatus),
	}, nil
}

// CloseTrade 关闭订单
// func CloseTrade(outTradeNo string) error {
// 	if s.client == nil {
// 		return fmt.Errorf("alipay client not initialized")
// 	}

// 	close := alipay.TradeClose{
// 		OutTradeNo: outTradeNo,
// 	}

// 	result, err := s.client.TradeClose(context.Background(), close)
// 	if err != nil {
// 		return fmt.Errorf("close trade failed: %v", err)
// 	}

// 	if result.Code != alipay.CodeSuccess {
// 		return fmt.Errorf("close trade failed: %s", result.SubMsg)
// 	}

// 	return nil
// }

// VerifyNotify 验证异步通知
func (s *AlipayService) VerifyNotify(values url.Values) (*types.NotifyResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("alipay client not initialized")
	}

	notification, err := s.client.DecodeNotification(values)
	if err != nil {
		return nil, fmt.Errorf("verify notify failed: %v", err)
	}

	return &types.NotifyResponse{
		TradeStatus: string(notification.TradeStatus),
		TradeNo:     notification.TradeNo,
		OutTradeNo:  notification.OutTradeNo,
	}, nil
}

// // VerifyReturn 验证同步回调
// func VerifyReturn(r *http.Request) (map[string]string, error) {
// 	if client == nil {
// 		return nil, fmt.Errorf("alipay client not initialized")
// 	}

// 	err := client.VerifyReturn(r)
// 	if err != nil {
// 		return nil, fmt.Errorf("verify return failed: %v", err)
// 	}

// 	params := make(map[string]string)
// 	for key, values := range r.URL.Query() {
// 		if len(values) > 0 {
// 			params[key] = values[0]
// 		}
// 	}

// 	return params, nil
// }

// randomString 生成随机字符串
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
