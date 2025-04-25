package service

import (
	"errors"
	"time"

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
	errOrderOntFound = errors.New("order not found")
)

func PurchaseProduct(req *request.PurchaseProductReq) exceptions.APIError {
	// create an order

	// TODO check if the product is available in cache

	product, err := db.GetOne[*models.Product](
		db.Equal("id", req.ProductID),
		db.Equal("is_published", true),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	if product.IsSold {
		return exceptions.BadRequestError(errProductSold, exceptions.ProductSoldError)
	}

	if !product.IsSelling {
		return exceptions.BadRequestError(errors.New("product is off shelves"), exceptions.ProductNotAvailableError)
	}

	order, err := db.GetOne[*models.Order](
		db.Equal("product_id", req.ProductID),
		db.NotEqual("status", orderStatusCancelled),
		db.NotEqual("status", orderStatusClosed),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if order.UserID == req.UserID {
		return nil
	}

	if order.Exists() {
		return exceptions.BadRequestError(errors.New("order already exists"), exceptions.AvatarSizeExceededError)
	}

	order = &models.Order{
		ProductID:   req.ProductID,
		UserID:      req.UserID,
		SellerID:    product.UserID,
		Status:      orderStatusPending,
		TotalAmount: req.TotalAmount,
		ShipAmount:  req.ShipAmount,
		PayTime:     time.Now(),
	}

	// TODO use rabbit mq dead queue to delete the order if not paid in 15 minutes
	if err := db.Create(order); err != nil {
		return exceptions.InternalServerError(err)
	}

	product.IsSold = true
	if err := db.Update(product); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
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
		return exceptions.BadRequestError(errOrderOntFound, exceptions.OrderNotFoundError)
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
		return exceptions.BadRequestError(errOrderOntFound, exceptions.OrderNotFoundError)
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

func PayOrder(req *request.PayOrderReq) exceptions.APIError {
	order, err := db.GetOne[models.Order](
		db.Equal("id", req.OrderID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return exceptions.BadRequestError(errOrderOntFound, exceptions.OrderNotFoundError)
	}

	if !order.IsOwner(req.UserID) {
		return exceptions.BadRequestError(errors.New("not the owner of the order"), exceptions.UserNotOrderOwnerError)
	}

	// TODO implement
	if order.Status != orderStatusPending {
		return exceptions.BadRequestError(errors.New("order is not to be paid"), exceptions.OrderNotToBePaidError)
	}

	order.Status = orderStatusPaid
	if err := db.Update(order); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
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
		return exceptions.BadRequestError(errOrderOntFound, exceptions.OrderNotFoundError)
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
		return exceptions.BadRequestError(errOrderOntFound, exceptions.OrderNotFoundError)
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
