package service

import (
	"time"

	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/es"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
)

func SearchProduct(req *request.SearchProductReq) (response.SearchProductResp, exceptions.APIError) {
	var resp response.SearchProductResp

	if len(req.UserID) > 0 {
		searchHistory := models.SearchHistory{
			UserID:     req.UserID,
			Keyword:    req.Keyword,
			SearchTime: time.Now().Unix(),
		}
		db.Create(&searchHistory)
	}

	query := buildSearchReq(req)

	// 从es获取商品搜索结果
	esResp, err := es.Search("m-market", query)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	// TODO 目前需要回数据库查询，后期不依赖于数据库，但是需要保证一致性
	productIDs := make([]string, 0, len(esResp))
	for _, product := range esResp {
		productIDs = append(productIDs, product["id"].(string))
	}

	// TODO 处理排序并应用到es 比如将按距离或者价格转换为es条件
	products, err := db.GetAll[models.Product](
		db.InArray("id", productIDs),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	userIDs := make([]string, 0, len(products))
	for _, product := range products {
		userIDs = append(userIDs, product.UserID)
	}

	users, err := db.GetAll[models.User](
		db.Fields("id", "username", "avatar"),
		db.InArray("id", userIDs),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	userProduct := make([]response.UserProduct, 0, len(products))
	for _, product := range products {
		for _, user := range users {
			if user.ID == product.UserID {
				userProduct = append(userProduct, response.UserProduct{
					User:    user,
					Product: product,
				})
			}
		}
	}

	history := &models.SearchHistory{
		UserID:     req.UserID,
		Keyword:    req.Keyword,
		SearchTime: time.Now().Unix(),
	}

	db.Create(history)

	resp.Products = userProduct
	resp.Page = req.Page
	resp.Size = req.Size
	return resp, nil
}

func buildSearchReq(req *request.SearchProductReq) map[string]interface{} {
	mustQueries := []map[string]interface{}{}

	if req.Keyword != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"match": map[string]interface{}{
				"describe": req.Keyword,
			},
		})
	}

	query := map[string]interface{}{
		"from": (req.Page - 1) * req.Size,
		"size": req.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustQueries,
			},
		},
	}

	// if req.Sort.Field != "" {
	// 	order := "asc"
	// 	if req.Sort.Decs {
	// 		order = "desc"
	// 	}
	// 	query["sort"] = []map[string]interface{}{
	// 		{
	// 			req.Sort.Field: map[string]interface{}{
	// 				"order": order,
	// 			},
	// 		},
	// 	}
	// }

	return query
}

func GetSearchHistory(req *request.GetSearchHistoryReq) (response.SearchHistoryResp, exceptions.APIError) {
	var resp response.SearchHistoryResp

	var pageQuery db.GenericQuery
	if !req.ShowAll {
		pageQuery = db.Page(1, 10)
	}

	history, err := db.GetAll[models.SearchHistory](
		db.Equal("user_id", req.UserID),
		db.OrderBy("search_time", true),
		db.GreaterThan("search_time", time.Now().AddDate(0, -1, 0).Unix()),
		pageQuery,
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.History = history
	return resp, nil
}
