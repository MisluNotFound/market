package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

// GET /api/order/{userID}/{productID}
func PurchaseProduct() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.PurchaseProductReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.PurchaseProduct(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// GET /api/order/{userID}/status
func GetAllOrderStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetAllOrderStatusReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		orderStatusList, err := service.GetAllOrderStatus(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, orderStatusList)
	}
}

// GET /api/order/{userID}/list
func GetOrderList() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetOrderListReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		req.PageReq.Fill()

		orderList, err := service.GetOrderList(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, orderList)
	}
}

// PUT /api/order/{userID}/{orderID}/singed
func ConfirmOrderSigned() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.ConfirmOrderReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.ConfirmOrderSigned(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// PUT /api/order/{userID}/{orderID}/shipped
func ConfirmOrderShipped() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.ConfirmOrderReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.ConfirmOrderShipped(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// POST /api/order/{userID}/{orderID}/pay
func PayOrder() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.PayOrderReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		// TODO implement wechat pay
		req.Method = "alipay"
		resp, err := service.PayOrder(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// GET /api/order/{userID}/{orderID}
func GetOrder() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetOrderReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetOrder(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// PUT /api/order/{userID}/{orderID}/refound
func RefoundOrder() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.ConfirmOrderReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.RefoundOrder(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// /api/order/cancel/userID/orderID
func CancelOrder() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CancelOrderReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.CancelOrder(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// POST /api/order/comment/orderID
func CreateOrderComment() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateOrderCommentReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		userID, _ := GetContextUserID(c)
		err := service.CreateOrderComment(req, userID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// POST /api/order/comment/orderID/reply
func ReplyOrderComment() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.ReplyOrderCommentReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		userID, _ := GetContextUserID(c)
		err := service.ReplyOrderComment(req, userID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// GET /api/order/comment/orderID
func GetOrderComments() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetOrderCommentsReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetOrderComments(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// GET /api/order/comment/seller/sellerID
func GetSellerComments() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetSellerCommentsReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetSellerComments(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// POST /api/order/alipay/notify
func AliPayNotify() func(c *gin.Context) {
	return func(c *gin.Context) {
		values := c.Request.PostForm

		service.AlipayNotify(values)
		c.Writer.Write([]byte("success"))
	}
}

func GetUnCommentOrder() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetUncommentOrder{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetUnCommentOrder(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

func GetOrderStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetOrderStatusReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetOrderStatus(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}
