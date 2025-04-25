package service

import (
	"errors"
	"mime/multipart"
	"slices"
	"strings"
	"time"

	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
	"github.com/mislu/market-api/internal/utils/lib"
)

var (
	errProductNotFound = errors.New("product not found")
	errProductSold     = errors.New("product sold")
	errNotOwner        = errors.New("user not owner")
)

func CreateProduct(req *request.CreateProductReq) (response.CreateProductResp, exceptions.APIError) {
	var resp response.CreateProductResp

	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, "User does not exist")
	}

	if req.ShipMethod == "included" {
		req.ShipPrice = 0
	}

	product := &models.Product{
		UserID:         req.UserID,
		Price:          req.Price,
		Describe:       req.Describe,
		ShippingMethod: req.ShipMethod,
		ShippingPrise:  req.ShipPrice,
		CanSelfPickup:  req.CanSelfPickup,
		OriginalPrice:  req.OriginalPrice,

		// TODO location

		// TODO 审核
		IsPublished: true,
		PublishAt:   time.Now(),
	}

	pics, resp, err := uploadProductPics(req.Pics, user.ID)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	product.Pics = strings.Join(pics, ",")

	err = db.Create(product)
	if err != nil {
		for _, pic := range pics {
			key := lib.SplitResourceURL(pic)

			// 忽略删除错误，由定时清理任务清除多余的文件
			resourcemanager.DeleteFile(resourcemanager.ProductBucket, resourcemanager.GetObjectPath(resourcemanager.ProductBucket, user.ID, key))
		}

		return resp, exceptions.InternalServerError(err)
	}

	return resp, nil
}

func GetProduct(req *request.GetProductReq) (response.GetProductResp, exceptions.APIError) {
	var resp response.GetProductResp

	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	product, err := db.GetOne[*models.Product](
		db.Equal("id", req.ProductID),
		db.Equal("user_id", req.UserID),
		db.Equal("is_published", true),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return resp, exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	resp.User = *user
	resp.Product = *product

	// TODO get comment

	// TODO get types and attributes

	return resp, nil
}

func UpdateProduct(req *request.UpdateProductReq) (response.CreateProductResp, exceptions.APIError) {
	var resp response.CreateProductResp

	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	product, err := db.GetOne[*models.Product](
		db.Equal("id", req.ProductID),
		db.Equal("user_id", req.UserID),
		db.Equal("is_published", true),
	)

	if !product.Exists() {
		return resp, exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	if product.IsSold {
		return resp, exceptions.BadRequestError(errProductSold, exceptions.ProductSoldError)
	}

	if !product.IsOwner(user.ID) {
		return resp, exceptions.BadRequestError(errNotOwner, exceptions.UserNotProductOwnerError)
	}

	if req.ShipMethod == "included" {
		product.ShippingPrise = 0
	}

	added, resp, err := uploadProductPics(req.AddedPics, user.ID)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	pics := strings.Split(product.Pics, ",")
	for _, pic := range req.DeletedPics {
		key := lib.SplitResourceURL(pic)

		resourcemanager.DeleteFile(resourcemanager.ProductBucket, resourcemanager.GetObjectPath(resourcemanager.ProductBucket, user.ID, key))

		idx := slices.Index(pics, pic)
		if idx != -1 {
			pics = slices.Delete(pics, idx, idx+1)
		}
	}

	pics = append(pics, added...)

	product.Pics = strings.Join(pics, ",")
	product.Price = req.Price
	product.Describe = req.Describe
	product.ShippingMethod = req.ShipMethod
	product.CanSelfPickup = req.CanSelfPickup
	product.OriginalPrice = req.OriginalPrice

	err = db.Update(product)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	return resp, nil
}

func uploadProductPics(pics []*multipart.FileHeader, userID string) (added []string, resp response.CreateProductResp, err error) {
	for _, file := range pics {
		if len(added) >= 5 {
			resp.Failures = append(resp.Failures, response.UploadFileFailure{
				FileName: file.Filename,
				Error:    "Too many files",
			})
			continue
		}

		if file.Size > picMaxSize {
			resp.Failures = append(resp.Failures, response.UploadFileFailure{
				FileName: file.Filename,
				Error:    "File size exceeds limit",
			})
			continue
		}

		picFile, err := file.Open()
		if err != nil {
			return added, resp, err
		}

		data := make([]byte, file.Size)
		_, err = picFile.Read(data)
		if err != nil {
			return added, resp, err
		}

		key := resourcemanager.GenerateObjectKey(file.Filename)
		path := resourcemanager.GetObjectPath(resourcemanager.ProductBucket, userID, key)
		err = resourcemanager.UploadFile(resourcemanager.ProductBucket, path, data)
		if err != nil {
			resp.Failures = append(resp.Failures, response.UploadFileFailure{
				FileName: file.Filename,
				Error:    "Failed to upload file",
			})
			continue
		}

		added = append(added, lib.GetResourceURL(int(resourcemanager.ProductBucket), userID, key))
	}

	return added, resp, nil
}

// 修该是否上架
func UpdateProductSellingStatus(userID, productID string, status bool) exceptions.APIError {
	product, err := db.GetOne[*models.Product](
		db.Equal("id", productID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	if !product.IsOwner(userID) {
		return exceptions.BadRequestError(errNotOwner, exceptions.UserNotProductOwnerError)
	}

	product.IsSelling = status

	err = db.Update(product)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

// 修改商品是否售出
func UpdateProductSoldStatus(userID, productID string, status bool) exceptions.APIError {
	product, err := db.GetOne[*models.Product](
		db.Equal("id", productID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	if !product.IsOwner(userID) {
		return exceptions.BadRequestError(errNotOwner, exceptions.UserNotProductOwnerError)
	}

	product.IsSold = status
	// 上架已下架的商品
	product.IsSelling = true

	err = db.Update(product)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func GetUserProducts(req *request.GetUserProductsReq) (response.GetUserProductsResp, exceptions.APIError) {
	var resp response.GetUserProductsResp

	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	products, err := db.GetAll[models.Product](
		db.Equal("user_id", req.UserID),
		db.Equal("is_published", true),
		db.Page(req.Page, req.Size),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.Products = products
	resp.Page = req.Page
	resp.Size = req.Size

	return resp, nil
}

func GetProductList(req *request.GetProductListReq) (response.GetProductListResp, exceptions.APIError) {
	var resp response.GetProductListResp

	products, err := db.GetAll[models.Product](
		db.Equal("is_published", true),
		db.Page(req.Page, req.Size),
		db.OrderBy("created_at", true),
		db.Equal("is_selling", true),
		db.Equal("is_sold", false),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	for _, product := range products {
		user, err := db.GetOne[*models.User](
			db.Equal("id", product.UserID),
		)

		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		if !user.Exists() {
			continue
		}

		resp.Products = append(resp.Products, response.UserProduct{
			User:    *user,
			Product: product,
		})
	}

	resp.Page = req.Page
	resp.Size = req.Size

	return resp, nil
}

func SearchProduct(req *request.SearchProductReq) (response.GetProductListResp, exceptions.APIError) {
	var resp response.GetProductListResp

	products, err := db.GetAll[models.Product](
		db.Equal("is_published", true),
		db.Page(req.Page, req.Size),
		db.OrderBy("created_at", true),
		db.Equal("is_selling", true),
		db.Equal("is_sold", false),
		db.Like("`describe`", req.Keyword),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	for _, product := range products {
		user, err := db.GetOne[*models.User](
			db.Equal("id", product.UserID),
		)

		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		if !user.Exists() {
			continue
		}

		resp.Products = append(resp.Products, response.UserProduct{
			User:    *user,
			Product: product,
		})
	}

	resp.Page = req.Page
	resp.Size = req.Size

	return resp, nil
}
