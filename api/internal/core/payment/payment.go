package payment

import (
	"errors"

	"github.com/mislu/market-api/internal/core/payment/alipay"
	"github.com/mislu/market-api/internal/core/payment/types"
)

func NewPaymentService(paymentType types.PaymentType) (types.PaymentService, error) {
	switch paymentType {
	case types.Alipay:
		return alipayService, nil
	default:
		return nil, errors.New("unsupported payment type")
	}
}

var alipayService *alipay.AlipayService

// var wechatService *wechat.WechatService

func InitPaymentService() error {
	service, err := alipay.NewAlipayClient()
	if err != nil {
		panic(err)
	}
	alipayService = service
	return nil
}
