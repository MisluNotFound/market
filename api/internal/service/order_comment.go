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

// CreateOrderComment 创建订单评论
func CreateOrderComment(req *request.CreateOrderCommentReq, userID string) exceptions.APIError {
	// 检查订单是否存在且属于当前用户
	order, err := db.GetOne[models.Order](db.Equal("id", req.OrderID))
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	if !order.Exists() {
		return exceptions.BadRequestError(errOrderNotFound, exceptions.OrderNotFoundError)
	}
	if !order.IsOwner(userID) {
		return exceptions.BadRequestError(errors.New("not the owner of the order"), exceptions.UserNotOrderOwnerError)
	}
	if order.IsEvaluated {
		return exceptions.BadRequestError(errors.New("order already evaluated"), exceptions.OrderAlreadyEvaluatedError)
	}

	if order.FinishTime.Before(time.Now().AddDate(0, -1, 0)) {
		return exceptions.BadRequestError(errors.New("order older than 30 days"), exceptions.OrderOlderThan30DaysError)
	}

	// 创建评论
	comment := models.OrderComment{
		OrderID:   req.OrderID,
		UserID:    userID,
		ProductID: order.ProductID,
		Comment:   req.Comment,
		IsGood:    req.IsGood,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&comment); err != nil {
		return exceptions.InternalServerError(err)
	}

	order.IsEvaluated = true
	err = db.Update(&order)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

// ReplyOrderComment 回复订单评论
func ReplyOrderComment(req *request.ReplyOrderCommentReq, userID string) exceptions.APIError {
	// 获取原评论
	parentComment, err := db.GetOne[models.OrderComment](db.Equal("id", req.CommentID))
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	if parentComment.ID == 0 {
		return exceptions.BadRequestError(errors.New("comment not found"), exceptions.CommentNotFoundError)
	}

	order, err := db.GetOne[models.Order](db.Equal("id", parentComment.OrderID))
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !order.Exists() {
		return exceptions.BadRequestError(errOrderNotFound, exceptions.OrderNotFoundError)
	}

	// 创建回复
	reply := models.OrderComment{
		OrderID:   parentComment.OrderID,
		UserID:    userID,
		ProductID: parentComment.ProductID,
		Comment:   req.Comment,
		ParentID:  req.CommentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&reply); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

// GetOrderComments 获取订单评论
func GetOrderComments(req *request.GetOrderCommentsReq) (response.GetOrderCommentsResp, exceptions.APIError) {
	var resp response.GetOrderCommentsResp
	resp.Page = req.Page
	resp.Size = req.Size

	// 获取主评论
	comments, err := db.GetAll[models.OrderComment](
		db.Equal("order_id", req.OrderID),
		db.Equal("is_top", true),
		db.OrderBy("created_at", true),
		db.Page(req.Page, req.Size),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	total, err := db.GetCount[models.OrderComment](
		db.Equal("order_id", req.OrderID),
		db.Equal("is_top", true),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}
	resp.Total = total

	resp.Comments = make([]response.OrderCommentDetail, 0, len(comments))

	// 获取评论用户信息和回复
	for _, comment := range comments {
		detail := response.OrderCommentDetail{
			OrderComment: comment,
		}

		user, err := db.GetOne[models.User](
			db.Equal("id", comment.UserID),
			db.Fields("username", "avatar"),
		)
		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}
		detail.Username = user.Username
		detail.Avatar = user.Avatar

		// 获取回复
		replies, err := db.GetAll[models.OrderComment](
			db.Equal("parent_id", comment.ID),
			db.OrderBy("created_at", true),
		)
		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		// 获取回复用户信息
		detail.Replies = make([]response.OrderCommentDetail, 0, len(replies))
		for _, reply := range replies {
			replyDetail := response.OrderCommentDetail{
				OrderComment: reply,
			}

			replyUser, err := db.GetOne[models.User](
				db.Equal("id", reply.UserID),
				db.Fields("username", "avatar"),
			)
			if err != nil {
				return resp, exceptions.InternalServerError(err)
			}

			replyToUser, err := db.GetOne[models.User](
				db.Equal("id", reply.ReplyTo),
				db.Fields("username"),
			)
			replyDetail.ReplyTo = replyToUser.Username
			replyDetail.Username = replyUser.Username
			replyDetail.Avatar = replyUser.Avatar

			detail.Replies = append(detail.Replies, replyDetail)
		}

		resp.Comments = append(resp.Comments, detail)
	}

	return resp, nil
}

// GetSellerComments 获取卖家评论
func GetSellerComments(req *request.GetSellerCommentsReq) (response.GetSellerCommentsResp, exceptions.APIError) {
	var resp response.GetSellerCommentsResp
	resp.Page = req.Page
	resp.Size = req.Size

	orders, err := db.GetAll[models.Order](
		db.Fields("id"),
		db.Equal("seller_id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	orderIDs := make([]string, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	var goodQuery db.GenericQuery
	if req.IsGood != nil {
		goodQuery = db.Equal("is_good", *req.IsGood)
	}

	comments, err := db.GetAll[models.OrderComment](
		db.InArray("order_id", orderIDs),
		db.OrderBy("created_at", true),
		goodQuery,
		db.Page(req.Page, req.Size),
	)

	total, err := db.GetCount[models.OrderComment](
		db.InArray("order_id", orderIDs),
		goodQuery,
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.Comments = make([]response.OrderCommentDetail, 0, len(comments))

	for _, comment := range comments {
		detail := response.OrderCommentDetail{
			OrderComment: comment,
		}

		user, err := db.GetOne[models.User](
			db.Equal("id", comment.UserID),
			db.Fields("username", "avatar"),
		)
		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}
		detail.Username = user.Username
		detail.Avatar = user.Avatar
		resp.Comments = append(resp.Comments, detail)
	}

	resp.Total = total

	return resp, nil
}
