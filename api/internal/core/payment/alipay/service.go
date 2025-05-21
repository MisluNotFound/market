package alipay

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mislu/market-api/internal/core/payment/types"
	"github.com/mislu/market-api/internal/utils/app"
	"github.com/smartwalle/alipay/v3"
)

var client *alipay.Client

type AlipayService struct {
}

func NewAlipayClient() {
	privateKey, err := os.ReadFile("D:\\repository\\m-market\\api\\keys\\private.txt")
	if err != nil {
		log.Fatal("read private key file failed", err)
	}

	cli, err := alipay.New(app.GetConfig().Alipay.APPID, string(privateKey), false)
	if err != nil {
		log.Fatal("new client failed", err)
	}

	client = cli
}

func (s *AlipayService) Pay(req types.PaymentRequest) (string, error) {
	if client == nil {
		return "", fmt.Errorf("alipay client not initialized")
	}

	pay := alipay.TradePagePay{
		Trade: alipay.Trade{
			Subject:     req.Subject,
			OutTradeNo:  req.OrderID,
			TotalAmount: req.Amount,
			ProductCode: "FAST_INSTANT_TRADE_PAY",
		},
	}

	form, err := client.TradePagePay(pay)
	if err != nil {
		return "", fmt.Errorf("generate pay form failed: %v", err)
	}

	return form.String(), nil
}

// Refund 发起退款
func Refund(outTradeNo, refundAmount, refundReason string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("alipay client not initialized")
	}

	outRequestNo := "REFUND_" + randomString(16)
	refund := alipay.TradeRefund{
		OutTradeNo:   outTradeNo,
		RefundAmount: refundAmount,
		OutRequestNo: outRequestNo,
		RefundReason: refundReason,
	}

	result, err := client.TradeRefund(context.Background(), refund)
	if err != nil {
		return "", fmt.Errorf("refund failed: %v", err)
	}

	if result.Code != alipay.CodeSuccess {
		return "", fmt.Errorf("refund failed: %s", result.SubMsg)
	}

	return outRequestNo, nil
}

// QueryTrade 查询订单状态
func QueryTrade(outTradeNo string) (*alipay.TradeQueryRsp, error) {
	if client == nil {
		return nil, fmt.Errorf("alipay client not initialized")
	}

	query := alipay.TradeQuery{
		OutTradeNo: outTradeNo,
	}

	result, err := client.TradeQuery(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("query trade failed: %v", err)
	}

	if result.Code != alipay.CodeSuccess {
		return nil, fmt.Errorf("query trade failed: %s", result.SubMsg)
	}

	return result, nil
}

// CloseTrade 关闭订单
func CloseTrade(outTradeNo string) error {
	if client == nil {
		return fmt.Errorf("alipay client not initialized")
	}

	close := alipay.TradeClose{
		OutTradeNo: outTradeNo,
	}

	result, err := client.TradeClose(context.Background(), close)
	if err != nil {
		return fmt.Errorf("close trade failed: %v", err)
	}

	if result.Code != alipay.CodeSuccess {
		return fmt.Errorf("close trade failed: %s", result.SubMsg)
	}

	return nil
}

// VerifyNotify 验证异步通知
// func VerifyNotify(r *http.Request) (map[string]string, error) {
// 	if client == nil {
// 		return nil, fmt.Errorf("alipay client not initialized")
// 	}

// 	err := client.NotifyVerify(r, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("verify notify failed: %v", err)
// 	}

// 	params := make(map[string]string)
// 	r.ParseForm()
// 	for key, values := range r.Form {
// 		if len(values) > 0 {
// 			params[key] = values[0]
// 		}
// 	}

// 	return params, nil
// }

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
