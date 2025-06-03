package request

type CreateOrderCommentReq struct {
	OrderIDReq
	Comment string `json:"comment" binding:"required,max=255"`
	Pics    string `json:"pics" binding:"omitempty,max=500"`
	IsGood  bool   `json:"isGood"`
}

type ReplyOrderCommentReq struct {
	CommentID int    `json:"commentID" binding:"required"`
	Comment   string `json:"comment" binding:"required"`
}

type GetOrderCommentsReq struct {
	OrderIDReq
	PageReq
}

type GetSellerCommentsReq struct {
	UserIDReq
	PageReq
	IsGood *bool `json:"isGood"`
}
