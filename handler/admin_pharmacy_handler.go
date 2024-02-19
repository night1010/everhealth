package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type AdminPharmacyHandler struct {
	adminPharmacyUsecase usecase.AdminPharmacyUsecase
}

func NewAdminPharmacyHandler(u usecase.AdminPharmacyUsecase) *AdminPharmacyHandler {
	return &AdminPharmacyHandler{adminPharmacyUsecase: u}
}

func (h *AdminPharmacyHandler) GetDetailAdminPharmacy(c *gin.Context) {
	var requestUri dto.RequestUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	adminPharmacy, err := h.adminPharmacyUsecase.FindOneAdminPharmacy(c.Request.Context(), &entity.User{Id: uint(requestUri.Id)})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewAdminPharmacyRes(adminPharmacy)})
}

func (h *AdminPharmacyHandler) GetAllAdminPharmacy(c *gin.Context) {
	var request dto.AdminPharmacyQueryReq
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.adminPharmacyUsecase.FindAllAdminPharmacy(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	adminPharmacies := pageResult.Data.([]*entity.User)
	adminPharmaciesRes := []*dto.AdminPharmacyRes{}
	for _, adminPharmacy := range adminPharmacies {
		adminPharmacyres := dto.NewAdminPharmacyRes(adminPharmacy)
		adminPharmaciesRes = append(adminPharmaciesRes, adminPharmacyres)
	}
	c.JSON(http.StatusOK, dto.Response{Data: adminPharmaciesRes,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}

func (h *AdminPharmacyHandler) PostAdminPharmacy(c *gin.Context) {
	var adminPharmacyReq dto.AdminPharmacyReq
	err := c.ShouldBindJSON(&adminPharmacyReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	adminPharmacy, adminContact := adminPharmacyReq.ToModel()

	_, err = h.adminPharmacyUsecase.CreateAdminPharmacy(c.Request.Context(), adminPharmacy, adminContact)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "create success"})

}

func (h *AdminPharmacyHandler) UpdateAdminPharmacy(c *gin.Context) {
	var adminPharmacyReq dto.AdminPharmacyUpdateReq
	var requestUri dto.RequestUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := c.ShouldBindJSON(&adminPharmacyReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	adminPharmacy, adminContact := adminPharmacyReq.ToModel()
	adminPharmacy.Id = uint(requestUri.Id)
	_, err = h.adminPharmacyUsecase.UpdateAdminPharmacy(c.Request.Context(), adminPharmacy, adminContact)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "Update success"})

}

func (h *AdminPharmacyHandler) DeleteAdminPharmacy(c *gin.Context) {
	var adminPharmacyUri dto.RequestUri
	err := c.ShouldBindUri(&adminPharmacyUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	adminPharmacy := &entity.User{Id: uint(adminPharmacyUri.Id)}

	err = h.adminPharmacyUsecase.DeleteAdminPharmacy(c.Request.Context(), adminPharmacy)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "delete success"})

}
