package payment

import (
	"errors"

	"github.com/mislu/market-api/internal/core/payment/types"
)

func NewPaymentService(paymentType types.PaymentType) (types.PaymentService, error) {
	switch paymentType {
	case types.Alipay:
		// return &alipay.AlipayService{}, nil
	default:
		return nil, errors.New("unsupported payment type")
	}

	return nil, errors.New("unsupported payment type")
}
