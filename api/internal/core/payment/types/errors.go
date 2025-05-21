package types

import "errors"

var (
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrPaymentFailed     = errors.New("payment failed")
	ErrRefundFailed      = errors.New("refund failed")
	ErrTransferFailed    = errors.New("transfer failed")
	ErrUnsupportedMethod = errors.New("unsupported payment method")
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidSign       = errors.New("invalid signature")
	ErrNetworkError      = errors.New("network error")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
