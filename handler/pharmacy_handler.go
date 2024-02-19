package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type PharmacyHandler struct {
	pharmacyUsecase usecase.PharmacyUsecase
}

func NewPharmacyHandler(u usecase.PharmacyUsecase) *PharmacyHandler {
	return &PharmacyHandler{pharmacyUsecase: u}
}

func (h *PharmacyHandler) GetAllPharmacy(c *gin.Context) {
	var request dto.ListPharmacyQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.pharmacyUsecase.FindAllPharmacy(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacies := pageResult.Data.([]*entity.Pharmacy)
	pharmaciesRes := []*dto.PharmacyRes{}
	for _, pharmacy := range pharmacies {
		pharmacyres := dto.NewPharmacyRes(pharmacy)
		pharmaciesRes = append(pharmaciesRes, pharmacyres)
	}
	c.JSON(http.StatusOK, dto.Response{Data: pharmaciesRes,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}

func (h *PharmacyHandler) GetAllPharmacySuperAdmin(c *gin.Context) {
	var request dto.ListPharmacySuperAdminQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.pharmacyUsecase.FindAllPharmacySuperAdmin(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacies := pageResult.Data.([]*entity.Pharmacy)
	pharmaciesRes := []*dto.PharmacySuperAdminRes{}
	for _, pharmacy := range pharmacies {
		pharmacyres := dto.NewPharmacySuperAdminRes(pharmacy)
		pharmaciesRes = append(pharmaciesRes, pharmacyres)
	}
	c.JSON(http.StatusOK, dto.Response{Data: pharmaciesRes,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}

func (h *PharmacyHandler) GetPharmacyDetail(c *gin.Context) {
	var requestUri dto.RequestPharmacyUri
	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacy, err := h.pharmacyUsecase.FindOnePharmacyDetail(c.Request.Context(), &entity.Pharmacy{Id: uint(requestUri.Id)})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewPharmacyRes(pharmacy)})
}

func (h *PharmacyHandler) PostPharmacy(c *gin.Context) {
	var pharmacyReq dto.PharmacyReq
	err := c.ShouldBindJSON(&pharmacyReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacy, err := pharmacyReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}
	_, err = h.pharmacyUsecase.CreatePharmacy(c.Request.Context(), pharmacy)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "created success"})
}

func (h *PharmacyHandler) PutPharmacy(c *gin.Context) {
	var pharmacyReq dto.PharmacyReq
	var requestUri dto.RequestPharmacyUri

	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&pharmacyReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	pharmacy, err := pharmacyReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}

	pharmacy.Id = uint(requestUri.Id)
	_, err = h.pharmacyUsecase.UpdatePharmacy(c.Request.Context(), pharmacy)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}

func (h *PharmacyHandler) DeletePharmacy(c *gin.Context) {
	var requestUri dto.RequestPharmacyUri
	var pharmacy entity.Pharmacy
	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}

	pharmacy.Id = uint(requestUri.Id)
	err = h.pharmacyUsecase.DeletePharmacy(c.Request.Context(), pharmacy)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "delete success"})
}
