package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/night1010/everhealth/valueobject"
	"github.com/gin-gonic/gin"
)

type PharmacyProductHandler struct {
	pharmacyProductUsecase usecase.PharmacyProductUsecase
}

func NewPharmacyProductHandler(u usecase.PharmacyProductUsecase) *PharmacyProductHandler {
	return &PharmacyProductHandler{pharmacyProductUsecase: u}
}

func (h *PharmacyProductHandler) GetAllPharmacy(c *gin.Context) {
	var request dto.ListPharmacyProductQueryParam
	var requestUri dto.RequestPharmacyUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	query.Condition("pharmacy", valueobject.Equal, requestUri.Id)
	pagedResult, err := h.pharmacyProductUsecase.FindAllPharmacyProduct(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	tempData := []*dto.ProductPharmacyRes{}
	for _, product := range pagedResult.Data.([]*entity.PharmacyProduct) {
		tempProduct := dto.NewProductPhamarcyRes(product)
		tempData = append(tempData, tempProduct)
	}
	c.JSON(http.StatusOK, dto.Response{
		Data:        tempData,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *PharmacyProductHandler) GetPharmacyProductDetail(c *gin.Context) {
	var requestUri dto.PharmacyProductUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}

	pharmacyProduct, err := h.pharmacyProductUsecase.FindOnePharmacyPeoduct(c.Request.Context(), &entity.PharmacyProduct{ProductId: requestUri.ProductId, PharmacyId: requestUri.PharmacyId})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{
		Data: dto.NewProductPhamarcyRes(pharmacyProduct),
	})
}

func (h *PharmacyProductHandler) PostPharmacyProduct(c *gin.Context) {
	var pharmacyProductReq dto.PharmacyProductReq
	var requestUri dto.RequestPharmacyUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := c.ShouldBindJSON(&pharmacyProductReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacyProduct, err := pharmacyProductReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacyProduct.PharmacyId = requestUri.Id
	_, err = h.pharmacyProductUsecase.CreatePharmacyProduct(c.Request.Context(), pharmacyProduct)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "created success"})
}

func (h *PharmacyProductHandler) PutPharmacyProduct(c *gin.Context) {
	var pharmacyProductUpdateReq dto.PharmacyProductUpdateReq
	var requestProductUri dto.PharmacyProductUri
	if err := c.ShouldBindUri(&requestProductUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := c.ShouldBindJSON(&pharmacyProductUpdateReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacyProduct, err := pharmacyProductUpdateReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacyProduct.PharmacyId = requestProductUri.PharmacyId
	pharmacyProduct.Id = requestProductUri.ProductId
	_, err = h.pharmacyProductUsecase.UpdatePharmacyProduct(c.Request.Context(), pharmacyProduct)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "updated success"})
}
