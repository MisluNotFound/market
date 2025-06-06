package service

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/mislu/market-api/internal/core/payment"
	"github.com/mislu/market-api/internal/core/payment/types"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
)

const (
	orderStatusPending = iota + 1
	orderStatusPaid
	orderStatusShipped
	orderStatusDone
	orderStatusRefunded
	orderStatusReShipped
	orderStatusClosed
	orderStatusCancelled
)

var (
	errOrderNotFound = errors.New("order not found")
)

func PurchaseProduct(req *request.PurchaseProductReq) (response.PurchaseProductResp, exceptions.APIError) {
	var resp response.PurchaseProductResp
	// create an order

	// TODO check if the product is available in cache

	product, err := db.GetOne[*models.Product](
		db.Equal("id", req.ProductID),
		db.Equal("is_published", true),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return resp, exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	if product.IsSold {
		return resp, exceptions.BadRequestError(errProductSold, exceptions.ProductSoldError)
	}

	if !product.IsSelling {
		return resp, exceptions.BadRequestError(errors.New("product is off shelves"), exceptions.ProductNotAvailableError)
	}

	order, err := db.GetOne[*models.Order](
		db.Equal("product_id", req.ProductID),
		db.NotEqual("status", orderStatusCancelled),
		db.NotEqual("status", orderStatusClosed),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if order.UserID == req.UserID {
		return resp, nil
	}

	if order.Exists() {
		return resp, exceptions.BadRequestError(errors.New("order already exists"), exceptions.AvatarSizeExceededError)
	}

	order = &models.Order{
		ProductID:   req.ProductID,
		UserID:      req.UserID,
		SellerID:    product.UserID,
		Status:      orderStatusPending,
		TotalAmount: req.TotalAmount,
		ShipAmount:  req.ShipAmount,
	}

	// TODO use rabbit mq dead queue to delete the order if not paid in 15 minutes
	if err := db.Create(order); err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	product.IsSold = true
	if err := db.Update(product); err != nil {
		return resp, exceptions.InternalServerError(err)
	}
	resp.OrderID = order.ID
	return resp, nil
}

func GetAllOrderStatus(req *request.GetAllOrderStatusReq) (response.GetAllOrderStatusResp, exceptions.APIError) {
	var resp response.GetAllOrderStatusResp

	bought, err := db.GetCount[models.Order](
		db.Equal("user_id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	sold, err := db.GetCount[models.Order](
		db.Equal("seller_id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	beEvaluated, err := db.GetCount[models.Order](
		db.Equal("seller_id", req.UserID),
		db.EqualOr("user_id", req.UserID),
		db.Equal("status", orderStatusDone),
		db.Equal("is_evaluated", false),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.Bought = int(bought)
	resp.Sold = int(sold)
	resp.BeEvaluated = int(beEvaluated)
	return resp, nil
}

func GetOrderList(req *request.GetOrderListReq) (response.GetOrderListResp, exceptions.APIError) {
	var resp response.GetOrderListResp

	userField := "user_id"
	if !req.IsBought {
		userField = "seller_id"
	}

	queries := []db.GenericQuery{
		db.Equal(userField, req.UserID),
		db.Page(req.Page, req.Size),
		db.OrderBy("created_at", true),
	}

	if req.Status > 0 {
		queries = append(queries, db.Equal("status", req.Status))
	}

	orders, err := db.GetAll[models.Order](queries...)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	sellerOrders := make([]response.UserOrder, 0, len(orders))
	for _, order := range orders {
		userID := order.SellerID
		if !req.IsBought {
			userID = order.UserID
		}

		sellerOrder := response.UserOrder{}
		user, err := db.GetOne[models.User](
			db.Fields("avatar", "username", "id"),
			db.Equal("id", userID),
		)
		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		product, err := db.GetOne[models.Product](
			db.Equal("id", order.ProductID),
		)

		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		sellerOrder.Product = product
		sellerOrder.Order = order
		sellerOrder.User = user
		sellerOrders = append(sellerOrders, sellerOrder)
	}

	resp.Orders = sellerOrders
	resp.Page = req.Page
	resp.Size = req.Size
	return resp, nil
}

func ConfirmOrderSigned(req *request.ConfirmOrderReq) exceptions.APIError {
	order, err := db.GetOne[*models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return exceptions.BadRequestError(errOrderNotFound, exceptions.OrderNotFoundError)
	}

	if !req.Refound && !order.IsOwner(req.UserID) {
		return exceptions.BadRequestError(errors.New("not the owner of the order"), exceptions.UserNotOrderOwnerError)
	}

	if req.Refound && !order.IsSeller(req.UserID) {
		return exceptions.BadRequestError(errors.New("not the seller of the order"), exceptions.UserNotOrderSellerError)
	}

	if !req.Refound && order.Status != orderStatusShipped {
		return exceptions.BadRequestError(errors.New("order has not been shipped"), exceptions.OrderHasNotShipped)
	}

	if req.Refound && order.Status != orderStatusReShipped {
		return exceptions.BadRequestError(errors.New("order has not been shipped"), exceptions.OrderHasNotShipped)
	}

	if req.Refound {
		order.Status = orderStatusClosed
	} else {
		order.Status = orderStatusDone
	}

	order.FinishTime = time.Now()
	if order.Status == orderStatusClosed {
		product, err := db.GetOne[models.Product](
			db.Equal("id", order.ProductID),
		)

		product.IsSold = false
		if err != nil {
			return exceptions.InternalServerError(err)
		}

		err = db.Update(product)
		if err != nil {
			return exceptions.InternalServerError(err)
		}
	}

	if err := db.Update(order); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func ConfirmOrderShipped(req *request.ConfirmOrderReq) exceptions.APIError {
	order, err := db.GetOne[models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return exceptions.BadRequestError(errOrderNotFound, exceptions.OrderNotFoundError)
	}

	if req.Refound && !order.IsOwner(req.UserID) {
		return exceptions.BadRequestError(errors.New("not the owner of the order"), exceptions.UserNotOrderOwnerError)
	}

	if !req.Refound && !order.IsSeller(req.UserID) {
		return exceptions.BadRequestError(errors.New("not the seller of the order"), exceptions.UserNotOrderSellerError)
	}

	if req.Refound && order.Status != orderStatusShipped {
		return exceptions.BadRequestError(errors.New("order has not been shipped"), exceptions.OrderHasNotShipped)
	}

	if !req.Refound && order.Status != orderStatusPaid {
		return exceptions.BadRequestError(errors.New("order has not been paid"), exceptions.OrderNotPaidError)
	}

	if req.Refound {
		order.Status = orderStatusReShipped
	} else {
		order.Status = orderStatusShipped
	}

	order.ShipTime = time.Now()

	if err := db.Update(order); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func PayOrder(req *request.PayOrderReq) (response.PayOrderResp, exceptions.APIError) {
	var resp response.PayOrderResp
	order, err := db.GetOne[models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return resp, exceptions.BadRequestError(errOrderNotFound, exceptions.OrderNotFoundError)
	}

	if !order.IsOwner(req.UserID) {
		return resp, exceptions.BadRequestError(errors.New("not the owner of the order"), exceptions.UserNotOrderOwnerError)
	}

	if order.Status != orderStatusPending {
		return resp, exceptions.BadRequestError(errors.New("order is not to be paid"), exceptions.OrderNotToBePaidError)
	}

	if err := db.Update(&order); err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	paymentService := payment.GetGlobalPaymentService()

	product, err := db.GetOne[models.Product](
		db.Equal("id", order.ProductID),
		db.Fields("describe"),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	payResp, err := paymentService.Pay(context.Background(), types.PaymentRequest{
		OrderID: req.OrderID,
		Amount:  strconv.FormatFloat(order.TotalAmount, 'f', -1, 64),
		Subject: product.Describe,
	})
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	order.PayMethod = req.Method
	resp.PayURL = payResp.PaymentURL
	err = db.Update(order)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}
	return resp, nil
}

func GetOrder(req *request.GetOrderReq) (response.GetOrderResp, exceptions.APIError) {
	var resp response.GetOrderResp

	order, err := db.GetOne[models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	var toQueryID string
	if !order.IsOwner(req.UserID) && !order.IsSeller(req.UserID) {
		return resp, exceptions.BadRequestError(errors.New("order not related"), exceptions.OrderNotRelatedError)
	} else if order.IsOwner(req.UserID) {
		toQueryID = order.SellerID
	} else {
		toQueryID = order.UserID
	}

	if time.Since(order.CreatedAt) > time.Minute*15 && order.Status == orderStatusPending {
		order.Status = orderStatusCancelled
	}

	user, err := db.GetOne[models.User](
		db.Fields("avatar", "username", "id"),
		db.Equal("id", toQueryID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if err := db.Update(order); err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	product, err := db.GetOne[models.Product](
		db.Equal("id", order.ProductID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.Product = product
	resp.User = user
	resp.Order = order
	return resp, nil
}

func RefoundOrder(req *request.ConfirmOrderReq) exceptions.APIError {
	order, err := db.GetOne[models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return exceptions.BadRequestError(errOrderNotFound, exceptions.OrderNotFoundError)
	}

	if !order.IsOwner(req.UserID) {
		return exceptions.BadRequestError(errors.New("not the owner of the order"), exceptions.UserNotOrderOwnerError)
	}

	if order.Status != orderStatusShipped && order.Status != orderStatusDone {
		return exceptions.BadRequestError(errors.New("can not refound"), exceptions.OrderHasNotShipped)
	}

	order.Status = orderStatusReShipped
	if err := db.Update(order); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func CancelOrder(req *request.CancelOrderReq) exceptions.APIError {
	order, err := db.GetOne[models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return exceptions.BadRequestError(errOrderNotFound, exceptions.OrderNotFoundError)
	}

	if !order.IsOwner(req.UserID) {
		return exceptions.BadRequestError(errors.New("not the owner of the order"), exceptions.UserNotOrderOwnerError)
	}

	if order.Status != orderStatusPending {
		return exceptions.BadRequestError(errors.New("order can not be canceled"), exceptions.OrderCanNotCanceled)
	}

	order.Status = orderStatusCancelled
	order.FinishTime = time.Now()

	// TODO 事务封装
	if err := db.Update(order); err != nil {
		return exceptions.InternalServerError(err)
	}

	product, err := db.GetOne[models.Product](
		db.Equal("id", order.ProductID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	product.IsSold = false
	if err := db.Update(product); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func AlipayNotify(values url.Values) exceptions.APIError {
	paymentService, err := payment.NewPaymentService(types.Alipay)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	notify, err := paymentService.VerifyNotify(values)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	order, err := db.GetOne[models.Order](
		db.Equal("id", notify.OutTradeNo),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return exceptions.BadRequestError(errors.New("order not found"), exceptions.OrderNotFoundError)
	}

	if notify.TradeStatus == "TRADE_SUCCESS" {
		order.Status = orderStatusPaid
		order.PayTime = time.Now()
		if err := db.Update(order); err != nil {
			return exceptions.InternalServerError(err)
		}
	}

	return nil
}

func GetUnCommentOrder(req *request.GetUncommentOrder) (response.GetUnCommentOrderResp, exceptions.APIError) {
	var resp response.GetUnCommentOrderResp

	queries := []db.GenericQuery{
		db.Equal("user_id", req.UserID),
		db.OrderBy("created_at", true),
		db.Equal("status", orderStatusDone),
		db.Equal("is_evaluated", false),
	}

	orders, err := db.GetAll[models.Order](queries...)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	sellerOrders := make([]response.UserOrder, 0, len(orders))
	for _, order := range orders {
		userID := order.SellerID

		sellerOrder := response.UserOrder{}
		user, err := db.GetOne[models.User](
			db.Fields("avatar", "username", "id"),
			db.Equal("id", userID),
		)
		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		product, err := db.GetOne[models.Product](
			db.Equal("id", order.ProductID),
		)

		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		sellerOrder.Product = product
		sellerOrder.Order = order
		sellerOrder.User = user
		sellerOrders = append(sellerOrders, sellerOrder)
	}

	resp.Orders = sellerOrders

	return resp, nil
}

func GetOrderStatus(req *request.GetOrderStatusReq) (response.GetOrderStatusResp, exceptions.APIError) {
	order, err := db.GetOne[models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return response.GetOrderStatusResp{}, exceptions.InternalServerError(err)
	}

	if order.Status == orderStatusPaid {
		return response.GetOrderStatusResp{Status: order.Status}, nil
	}

	paymentService := payment.GetGlobalPaymentService()

	gap := time.Second * 2
	maxRetries := 3
	retryCount := 0
	for {
		if retryCount >= maxRetries {
			break
		}

		tradeResp, err := paymentService.QueryTrade(order.ID)
		if err != nil {
			time.Sleep(gap)
			gap = gap * 2
			if gap > time.Second*60 {
				gap = time.Second * 60
			}
			retryCount++
			continue
		}

		switch tradeResp.TradeStatus {
		case types.TradeStatusWaitBuyerPay:
			time.Sleep(gap)
			gap = gap * 2
			if gap > time.Second*60 {
				gap = time.Second * 60
			}
			retryCount++
			continue
		case types.TradeStatusSuccess:
			order.Status = orderStatusPaid
			order.PayTime = time.Now()
			if err := db.Update(&order); err != nil {
				return response.GetOrderStatusResp{Status: orderStatusPaid}, exceptions.InternalServerError(err)
			}
			return response.GetOrderStatusResp{Status: orderStatusPaid}, exceptions.InternalServerError(err)
		default:
			log.Printf("Order %s: Unexpected trade status: %s", order.ID, tradeResp.TradeStatus)
		}
	}

	resp := response.GetOrderStatusResp{}
	resp.Status = order.Status
	return resp, nil
}
