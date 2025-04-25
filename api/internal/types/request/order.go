package request

type OrderIDReq struct {
	OrderID string `uri:"orderID" binding:"required"`
}

type PurchaseProductReq struct {
	ProductIDReq
	UserIDReq
	TotalAmount float64 `form:"totalAmount" binding:"required"`
	ShipAmount  float64 `form:"shipAmount"`
}

type GetAllOrderStatusReq struct {
	UserIDReq
}

type GetOrderListReq struct {
	UserIDReq
	PageReq
	IsBought bool `form:"isBought"`
	Status   int  `form:"status"`
}

type ConfirmOrderReq struct {
	UserIDReq
	OrderIDReq
	Refound bool `form:"refund"`
}

type PayOrderReq struct {
	UserIDReq
	OrderIDReq
}

type GetOrderReq struct {
	UserIDReq
	OrderIDReq
}

type CancelOrderReq struct {
	UserIDReq
	OrderIDReq
}
