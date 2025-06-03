package request

type AddressIDReq struct {
	AddressID string `uri:"addressID"`
}

type CreateAddressReq struct {
	Address      string `form:"address" binding:"required"`
	City         string `form:"city" binding:"required"`
	District     string `form:"district" binding:"required"`
	Province     string `form:"province" binding:"required"`
	Street       string `form:"street" binding:"required"`
	StreetNumber string `form:"streetNumber" binding:"required"`
	IsDefault    *bool  `form:"isDefault" binding:"required"`
	Detail       string `form:"detail" binding:"required"`
	UserIDReq
	Phone     string  `form:"phone" binding:"required"`
	Receiver  string  `form:"name" binding:"required"`
	Latitude  float64 `form:"latitude" binding:"required"`
	Longitude float64 `form:"longitude" binding:"required"`
}

type UpdateAddressReq struct {
	AddressIDReq
	Address      string `form:"address" binding:"required"`
	City         string `form:"city" binding:"required"`
	District     string `form:"district" binding:"required"`
	Province     string `form:"province" binding:"required"`
	Street       string `form:"street" binding:"required"`
	StreetNumber string `form:"streetNumber" binding:"required"`
	Detail       string `form:"detail" binding:"required"`
	UserIDReq
	Phone     string  `form:"phone" binding:"required"`
	Receiver  string  `form:"name" binding:"required"`
	Latitude  float64 `form:"latitude" binding:"required"`
	Longitude float64 `form:"longitude" binding:"required"`
}

type DeleteAddressReq struct {
	AddressIDReq
}

type GetAddressReq struct {
	UserIDReq
	PageReq
}

type SetDefaultAddressReq struct {
	UserIDReq
	AddressIDReq
	IsDefault *bool `form:"isDefault" binding:"required"`
}
