package alipay

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("m_market_config", "D:\\repository\\m-market\\api\\config.yaml")
}

func TestNewService(t *testing.T) {
	NewAlipayClient()
}
