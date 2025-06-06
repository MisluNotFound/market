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
		history, _ := db.GetOne[models.SearchHistory](
			db.Equal("user_id", req.UserID),
			db.Equal("keyword", req.Keyword),
		)

		if len(history.Keyword) > 0 {
			// 之前搜索过
			history.SearchTime = time.Now().Unix()
			db.Update(&history)
		} else {
			searchHistory := models.SearchHistory{
				UserID:     req.UserID,
				Keyword:    req.Keyword,
				SearchTime: time.Now().Unix(),
			}
			db.Create(&searchHistory)
		}
	}

	query := buildSearchReq(req)

	esResp, err := es.Search("m-market", query)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	productIDs := make([]string, 0, len(esResp))
	for _, product := range esResp {
		productIDs = append(productIDs, product["id"].(string))
	}

	products, err := db.GetAll[models.Product](
		db.InArray("id", productIDs),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	productMap := make(map[string]models.Product, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}
	orderedProducts := make([]models.Product, 0, len(productIDs))
	for _, id := range productIDs {
		if product, exists := productMap[id]; exists {
			orderedProducts = append(orderedProducts, product)
		}
	}

	userIDs := make([]string, 0, len(orderedProducts))
	for _, product := range orderedProducts {
		userIDs = append(userIDs, product.UserID)
	}

	users, err := db.GetAll[models.User](
		db.Fields("id", "username", "avatar"),
		db.InArray("id", userIDs),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	userProduct := make([]response.UserProduct, 0, len(orderedProducts))
	for _, product := range orderedProducts {
		for _, user := range users {
			if user.ID == product.UserID {
				address, err := getProductAddress(product.Location)
				if err != nil {
					continue
				}
				userProduct = append(userProduct, response.UserProduct{
					User:    user,
					Product: product,
					Address: address.Address,
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

	if req.Sort.Field != "" {
		order := "asc"
		if req.Sort.Desc {
			order = "desc"
		}
		query["sort"] = []map[string]interface{}{
			{
				req.Sort.Field: map[string]interface{}{
					"order": order,
				},
			},
		}
	}

	return query
}

func GetSearchHistory(req *request.GetSearchHistoryReq) (response.SearchHistoryResp, exceptions.APIError) {
	var resp response.SearchHistoryResp

	history, err := db.GetAll[models.SearchHistory](
		db.Equal("user_id", req.UserID),
		db.OrderBy("search_time", true),
		db.GreaterThan("search_time", time.Now().AddDate(0, -1, 0).Unix()),
		db.Page(1, 20),
		db.Distinct("keyword"),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.History = history
	return resp, nil
}
