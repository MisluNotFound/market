package request

type SearchProductReq struct {
	Keyword    string            `json:"keyword"`
	Categories []string          `json:"categories"`
	Attributes []AttributeFilter `json:"attributes"`
	Sort       SortOption        `json:"sort"`
	UserID     string
	PageReq
}

type AttributeFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SortOption struct {
	Field string `json:"field"`
	Decs  bool   `json:"decs"`
}

type GetSearchHistoryReq struct {
	UserIDReq
	ShowAll bool `json:"showAll"`
}
