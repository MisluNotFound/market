package response

import "github.com/mislu/market-api/internal/types/models"

type GetAddressResp struct {
	Addresses []WrappedAddress `json:"addresses"`
	PageResp
}

type WrappedAddress struct {
	models.UserAddress
	models.Address
}
