package alipay

import (
	"fmt"
	"os"
	"testing"

	"github.com/mislu/market-api/internal/core/payment/types"
)

func init() {
	os.Setenv("m_market_config", "D:\\repository\\m-market\\api\\config.yaml")
}

func TestNewService(t *testing.T) {
	s, err := NewAlipayClient()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := s.Pay(nil, types.PaymentRequest{
		Subject:     "支付给商家",
		OrderID:     "99eae5d2-395c-11f0-970d-0242ac150002",
		Amount:      "199",
		Description: "Xbox 冰雪白游戏手柄",
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(resp.PaymentURL)
}
