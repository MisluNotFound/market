package service

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"slices"
	"strings"
	"time"

	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/es"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
	"github.com/mislu/market-api/internal/utils/lib"
	"gorm.io/gorm"
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

	attributes := make(map[uint]string)
	err = json.Unmarshal([]byte(req.AttributesJson), &attributes)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
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
		IsPublished: true,
		PublishAt:   time.Now(),
	}

	switch req.Condition {
	case "new":
		product.Condition = "全新"
	case "good":
		product.Condition = "八成新"
	case "excellent":
		product.Condition = "九成新"
	default:
		product.Condition = "使用过"
		product.UsedTime = req.UsedTime
	}

	pics, resp, err := uploadProductPics(req.Pics, user.ID)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	deletePics := func(pics []string) {
		for _, pic := range pics {
			key := lib.SplitResourceURL(pic)

			// 忽略删除错误，由定时清理任务清除多余的文件
			resourcemanager.DeleteFile(resourcemanager.ProductBucket, resourcemanager.GetObjectPath(resourcemanager.ProductBucket, user.ID, key))
		}
	}

	product.Pics = strings.Join(pics, ",")

	// 数据库创建商品，分类
	err = db.WithTransaction(func(tx *gorm.DB) error {
		if err := db.Create(product, tx); err != nil {
			return err
		}

		productAttributes := make([]models.ProductAttribute, 0, len(attributes))
		for id, value := range attributes {
			productAttributes = append(productAttributes, models.ProductAttribute{
				ProductID:   product.ID,
				AttributeID: id,
				Value:       value,
			})
		}

		productCategories := make([]models.ProductCategory, 0, len(req.Categories))
		for _, categoryID := range req.Categories {
			productCategories = append(productCategories, models.ProductCategory{
				ProductID:  product.ID,
				CategoryID: categoryID,
			})
		}

		if err := db.Create(productCategories, tx); err != nil {
			return err
		}

		return db.Create(productAttributes, tx)
	})

	if err != nil {
		deletePics(pics)

		return resp, exceptions.InternalServerError(err)
	}

	productDocument := &request.ProductDocument{
		ID:        product.ID,
		Describe:  product.Describe,
		CreatedAt: product.CreatedAt,
		Price:     product.Price,
	}

	productCategories, err := db.GetAll[models.Category](
		db.InArray("id", req.Categories),
	)

	categoriesString := make([]string, 0, len(productCategories))
	for _, category := range productCategories {
		categoriesString = append(categoriesString, category.TypeName)
	}

	if err != nil {
		deletePics(pics)

		return resp, exceptions.InternalServerError(err)
	}

	productDocument.Category = categoriesString

	for id, value := range attributes {
		attribute, err := db.GetOne[models.AttributeTemplate](
			db.Equal("id", id),
		)

		if err != nil {
			deletePics(pics)

			return resp, exceptions.InternalServerError(err)
		}

		productDocument.Attributes = append(productDocument.Attributes, request.AttributeES{
			Key:   attribute.Name,
			Value: value,
		})
	}
	// 写入es
	err = es.IndexDocument("m-market", productDocument.ID, productDocument)
	if err != nil {
		deletePics(pics)

		return resp, exceptions.InternalServerError(err)
	}

	// TODO 写入gorse
	return resp, nil
}

func GetProduct(req *request.GetProductReq) (response.GetProductResp, exceptions.APIError) {
	var resp response.GetProductResp

	product, err := db.GetOne[models.Product](
		db.Equal("id", req.ProductID),
		db.Equal("is_published", true),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return resp, exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	user, err := db.GetOne[models.User](
		db.Equal("id", product.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	resp.User = user
	resp.Product = product

	// TODO get comment

	productCategories, err := db.GetAll[models.ProductCategory](
		db.Equal("product_id", req.ProductID),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	productAttributes, err := db.GetAll[models.ProductAttribute](
		db.Equal("product_id", req.ProductID),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	for _, productCategory := range productCategories {
		resp.Categories = append(resp.Categories, productCategory.CategoryID)
	}

	attributes := make(map[uint]string)
	for _, productAttribute := range productAttributes {
		attributes[productAttribute.AttributeID] = productAttribute.Value
	}
	resp.Attributes = attributes

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

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return resp, exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	if product.IsSold {
		return resp, exceptions.BadRequestError(errProductSold, exceptions.ProductSoldError)
	}

	if !product.IsOwner(user.ID) {
		return resp, exceptions.BadRequestError(errNotOwner, exceptions.UserNotProductOwnerError)
	}

	attributes := make(map[uint]string)
	err = json.Unmarshal([]byte(req.AttributesJson), &attributes)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
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

	err = db.WithTransaction(func(tx *gorm.DB) error {
		if err := db.Update(product, tx); err != nil {
			return err
		}

		productAttributes := make([]models.ProductAttribute, 0, len(attributes))
		for id, value := range attributes {
			productAttributes = append(productAttributes, models.ProductAttribute{
				ProductID:   product.ID,
				AttributeID: id,
				Value:       value,
			})
		}

		productCategories := make([]models.ProductCategory, 0, len(req.Categories))
		for _, categoryID := range req.Categories {
			productCategories = append(productCategories, models.ProductCategory{
				ProductID:  product.ID,
				CategoryID: categoryID,
			})
		}

		for _, productCategory := range productCategories {
			if err := db.FirstOrCreate(&productCategory, tx); err != nil {
				return err
			}
		}

		return db.FirstOrCreate(productAttributes, tx)
	})
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

	user, err := db.GetOne[models.User](
		db.Fields("id", "username", "avatar"),
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	credit, err := db.GetOne[models.Credit](
		db.Equal("user_id", req.UserID),
	)

	products, err := db.GetAll[models.Product](
		db.OrderBy("publish_at", true),
		db.Equal("user_id", req.UserID),
		db.Equal("is_published", true),
		db.Page(req.Page, req.Size),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	userProducts := make([]response.UserProduct, 0, len(products))
	for _, product := range products {
		userProducts = append(userProducts, response.UserProduct{
			User:    user,
			Product: product,
			Credit:  credit,
		})
	}

	resp.Products = userProducts
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
		user, err := db.GetOne[models.User](
			db.Equal("id", product.UserID),
		)

		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		if !user.Exists() {
			continue
		}

		credit, err := db.GetOne[models.Credit](
			db.Equal("user_id", product.UserID),
		)

		resp.Products = append(resp.Products, response.UserProduct{
			User:    user,
			Product: product,
			Credit:  credit,
		})
	}

	resp.Page = req.Page
	resp.Size = req.Size

	return resp, nil
}

func GetAllCategory() (response.GetAllCategoryResp, exceptions.APIError) {
	flatCategories, err := db.GetAll[models.Category](
		db.OrderBy("level", false),
	)

	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	categoryMap := make(map[uint]*response.WrappedCategory)
	var rootNodes []*response.WrappedCategory

	for _, cat := range flatCategories {
		node := &response.WrappedCategory{
			Category: cat,
			Children: []*response.WrappedCategory{},
		}
		categoryMap[cat.ID] = node
		if cat.ParentID == 0 {
			rootNodes = append(rootNodes, node)
		}
	}

	for _, cat := range flatCategories {
		if cat.ParentID != 0 {
			if parent, ok := categoryMap[cat.ParentID]; ok {
				parent.Children = append(parent.Children, categoryMap[cat.ID])
			}
		}
	}

	for _, cat := range flatCategories {
		if cat.IsLeaf {
			categoryAttributes, err := db.GetAll[models.CategoryAttribute](
				db.Equal("category_id", cat.ID),
			)
			if err != nil {
				return nil, exceptions.InternalServerError(err)
			}

			attributeIDs := make([]uint, 0, len(categoryAttributes))
			for _, ca := range categoryAttributes {
				attributeIDs = append(attributeIDs, ca.AttributeID)
			}

			attributes, err := db.GetAll[models.AttributeTemplate](
				db.InArray("id", attributeIDs),
			)

			if err != nil {
				return nil, exceptions.InternalServerError(err)
			}

			categoryMap[cat.ID].Attributes = attributes
		}
	}

	return response.GetAllCategoryResp(rootNodes), nil
}

func UpdateProductPrice(req *request.UpdateProductPriceReq) exceptions.APIError {
	product, err := db.GetOne[models.Product](
		db.Equal("id", req.ProductID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	product.Price = req.Price
	err = db.Update(&product)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func LikeProduct(req *request.LikeProductReq) exceptions.APIError {
	product, err := db.GetOne[models.Product](
		db.Equal("id", req.ProductID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !product.Exists() {
		return exceptions.BadRequestError(errProductNotFound, exceptions.ProductNotFoundError)
	}

	user, err := db.GetOne[models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return exceptions.BadRequestError(errors.New("user does not exist"), exceptions.UserNotExistsError)
	}

	like := models.Like{
		UserID:    user.ID,
		ProductID: product.ID,
	}

	err = db.Create(&like)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func DislikeProduct(req *request.DislikeProductReq) exceptions.APIError {
	like, err := db.GetOne[models.Like](
		db.Equal("user_id", req.UserID),
		db.Equal("product_id", req.ProductID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	err = db.Delete(&like)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func GetUserLikes(req *request.GetUserLikesReq) (response.GetUserLikesResp, exceptions.APIError) {
	var resp response.GetUserLikesResp

	user, err := db.GetOne[models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	likes, err := db.GetAll[models.Like](
		db.Equal("user_id", req.UserID),
		db.Page(req.Page, req.Size),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	for _, like := range likes {
		product, err := db.GetOne[models.Product](
			db.Equal("id", like.ProductID),
		)

		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		if !product.Exists() {
			continue
		}

		resp.Products = append(resp.Products, response.UserProduct{
			User:    user,
			Product: product,
			IsLiked: true,
		})
	}

	resp.Page = req.Page
	resp.Size = req.Size

	return resp, nil
}

func GetInterestTags() (response.GetInterestTagsResp, exceptions.APIError) {
	var resp response.GetInterestTagsResp

	tags, err := db.GetAll[models.InterestTag](
		db.OrderBy("tag_name", true),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp = tags

	return resp, nil
}
