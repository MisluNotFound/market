package response

type SearchProductResp struct {
	Products []UserProduct `json:"products"`
	PageResp
}