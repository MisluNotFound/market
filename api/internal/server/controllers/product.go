package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

// POST /api/product/{userID}
func CreateProduct() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateProductReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.CreateProduct(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// GET /api/product/{productID}
// TODO verify if userID or productID is empty
func GetProduct() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetProductReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetProduct(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// PUT /api/product/{userID}/{productID}
func UpdateProduct() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateProductReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.UpdateProduct(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// PUT /api/product/{userID}/{productID}/off-shelves
func OffShelves() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateProductStatusReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.UpdateProductSellingStatus(req.UserID, req.ProductID, false)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// PUT /api/product/{userID}/{productID}/on-shelves
func OnShelves() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateProductStatusReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.UpdateProductSellingStatus(req.UserID, req.ProductID, true)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// PUT /api/product/{userID}/{productID}/selling
func NotSold() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateProductStatusReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.UpdateProductSoldStatus(req.UserID, req.ProductID, false)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// PUT /api/product/{userID}/{productID}/sold
func SoldOut() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateProductStatusReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.UpdateProductSoldStatus(req.UserID, req.ProductID, true)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// GET /api/product/{userID}/products
func GetUserProducts() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetUserProductsReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		req.PageReq.Fill()

		resp, err := service.GetUserProducts(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// /api/product/products
func GetProductList() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetProductListReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		req.PageReq.Fill()

		resp, err := service.GetProductList(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// // get /api/product/search
// func SearchProduct() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		req := &request.SearchProductReq{}
// 		if err := BindRequest(c, req); err != nil {
// 			AbortWithError(c, err)
// 			return
// 		}

// 		req.PageReq.Fill()

// 		resp, err := service.SearchProduct(req)
// 		if err != nil {
// 			AbortWithError(c, err)
// 			return
// 		}

// 		Success(c, ResponseTypeJSON, resp)
// 	}
// }

// get /api/product/category
func GetAllCategory() func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := service.GetAllCategory()
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// put /api/product/{userID}/{productID}/price
func UpdateProductPrice() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateProductPriceReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.UpdateProductPrice(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}
