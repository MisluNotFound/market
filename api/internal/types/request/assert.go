package request

type GetAssertReq struct {
	Type  int    `uri:"type"`
	Owner string `uri:"owner" binding:"required"`
	Key   string `uri:"key" binding:"required"`
}
