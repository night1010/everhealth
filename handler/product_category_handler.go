package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ProductCategoryHandler struct {
	productCategoryUsecase usecase.ProductCategoryUsecase
}

func NewProductCategoryHandler(pcu usecase.ProductCategoryUsecase) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		productCategoryUsecase: pcu,
	}
}

func (h *ProductCategoryHandler) GetProductCategories(c *gin.Context) {
	var queryParam dto.ProductCategoryParams
	err := c.ShouldBindQuery(&queryParam)
	if err != nil {
		_ = c.Error(err)
		return
	}
	query, err := queryParam.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pagedResult, err := h.productCategoryUsecase.GetProductCategories(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	productCategoriesRes := []dto.ProductCategoryRes{}
	for _, productCategory := range pagedResult.Data.([]*entity.ProductCategory) {
		productCategoryRes := dto.NewProductCategoryRes(productCategory)
		productCategoriesRes = append(productCategoriesRes, productCategoryRes)
	}
	c.JSON(http.StatusOK, dto.Response{Data: productCategoriesRes, TotalPage: &pagedResult.TotalPage, TotalItem: &pagedResult.TotalItem, CurrentPage: &pagedResult.CurrentPage, CurrentItem: &pagedResult.CurrentItems})
}

func (h *ProductCategoryHandler) GetProductCategoriesDetail(c *gin.Context) {
	var requestUri dto.RequestUri
	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	productCategory, err := h.productCategoryUsecase.GetProductCategoriesDetail(c.Request.Context(), &entity.ProductCategory{Id: uint(requestUri.Id)})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewProductCategoryRes(productCategory)})

}

func (h *ProductCategoryHandler) PostProductCategory(c *gin.Context) {
	var productCategoryReq dto.ProductCategoryReq
	err := c.ShouldBindWith(&productCategoryReq, binding.Form)
	if err != nil {
		_ = c.Error(err)
		return
	}
	productCategory := productCategoryReq.ToModel()
	newproductCategory, err := h.productCategoryUsecase.CreateProductCategory(c.Request.Context(), productCategory)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewProductCategoryRes(newproductCategory)})
}

func (h *ProductCategoryHandler) PutProductCategory(c *gin.Context) {
	var productCategoryUri dto.ProductCategoryUri
	var productCategoryReq dto.ProductCategoryReq
	err := c.ShouldBindUri(&productCategoryUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = c.ShouldBindWith(&productCategoryReq, binding.Form)
	if err != nil {
		_ = c.Error(err)
		return
	}
	productCategory := productCategoryReq.ToModel()
	productCategory.Id = uint(productCategoryUri.Id)
	updatedproductCategory, err := h.productCategoryUsecase.UpdateProductCategory(c.Request.Context(), productCategory)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewProductCategoryRes(updatedproductCategory)})
}

func (h *ProductCategoryHandler) DeleteProductCategory(c *gin.Context) {
	var productCategoryUri dto.ProductCategoryUri
	var productCategory entity.ProductCategory
	err := c.ShouldBindUri(&productCategoryUri)
	if err != nil {
		_ = c.Error(err)
		return
	}

	productCategory.Id = uint(productCategoryUri.Id)
	err = h.productCategoryUsecase.DeleteProductCategories(c.Request.Context(), productCategory)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "delete success"})
}
