package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) ListProduct(c *gin.Context) {
	var request dto.ListProductQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}

	pagedResult, products, top, floor, err := h.productUsecase.ListAllProduct(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var response []*dto.SimpleProductResponse
	for _, product := range products {
		response = append(response, &dto.SimpleProductResponse{
			Id:          product.Id,
			Name:        product.Name,
			TopPrice:    top[product.Id],
			FloorPrice:  floor[product.Id],
			SellingUnit: product.SellingUnit,
			Image:       product.Image,
		})
	}
	c.JSON(200, dto.Response{
		Data:        response,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *ProductHandler) ListProductAdmin(c *gin.Context) {
	var request dto.ListProductQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}

	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}

	pagedResult, err := h.productUsecase.ListAllProductAdmin(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	products := pagedResult.Data.([]*entity.Product)

	var response []*dto.ProductResponse
	for _, product := range products {
		response = append(response, dto.NewFromProduct(product))
	}
	c.JSON(200, dto.Response{
		Data:        response,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *ProductHandler) ListNearbyProduct(c *gin.Context) {
	var request dto.ListProductQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}

	pagedResult, products, top, floor, err := h.productUsecase.ListNearbyProduct(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var response []*dto.SimpleProductResponse
	for _, product := range products {
		response = append(response, &dto.SimpleProductResponse{
			Id:          product.Id,
			Name:        product.Name,
			TopPrice:    top[product.Id],
			FloorPrice:  floor[product.Id],
			SellingUnit: product.SellingUnit,
			Image:       product.Image,
		})
	}
	c.JSON(200, dto.Response{
		Data:        response,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var request dto.AddProductRequest
	if err := c.ShouldBindWith(&request, binding.Form); err != nil {
		_ = c.Error(err)
		return
	}
	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	product := request.ToProduct()
	drug := request.ToDrug()

	createdProduct, err := h.productUsecase.AddProduct(c.Request.Context(), product, drug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.NewFromProduct(createdProduct),
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var request dto.AddProductRequest
	var requestUri dto.RequestUri
	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindWith(&request, binding.Form); err != nil {
		_ = c.Error(err)
		return
	}
	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	product := request.ToProduct()
	drug := request.ToDrug()
	product.Id = uint(requestUri.Id)
	if drug != nil {
		drug.ProductId = uint(requestUri.Id)
	}

	updatedProduct, err := h.productUsecase.UpdateProduct(c.Request.Context(), product, drug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.NewFromProduct(updatedProduct),
	})
}

func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	var uri dto.RequestUri

	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(err)
		return
	}

	fetchedProduct, topPrice, floorPrice, err := h.productUsecase.GetProductDetail(c.Request.Context(), uint(uri.Id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Data: dto.NewFromPharmacyProduct(fetchedProduct, topPrice, floorPrice),
	})
}

func (h *ProductHandler) GetProductDetailAdmin(c *gin.Context) {
	var uri dto.RequestUri

	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(err)
		return
	}

	fetchedProduct, err := h.productUsecase.GetProductDetailAdmin(c.Request.Context(), uint(uri.Id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Data: dto.NewFromProduct(fetchedProduct),
	})
}
