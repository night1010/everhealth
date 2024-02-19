package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type ProvinceHandler struct {
	provinceUsecase usecase.ProvinceUsecase
}

func NewProvinceHadnler(u usecase.ProvinceUsecase) *ProvinceHandler {
	return &ProvinceHandler{provinceUsecase: u}
}

func (h *ProvinceHandler) GetAllProvince(c *gin.Context) {
	provincies, err := h.provinceUsecase.FindAllProvince(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	provinciesRes := []*dto.ProvinceRes{}
	for _, province := range provincies {
		provinceres := dto.NewProvinceRes(province)
		provinciesRes = append(provinciesRes, provinceres)
	}
	c.JSON(http.StatusOK, provinciesRes)
}

func (h *ProvinceHandler) GetDetailProvince(c *gin.Context) {
	var provinceUri dto.ProvinceUri
	err := c.ShouldBindUri(&provinceUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	province, err := h.provinceUsecase.FindProvinceByIdWithCities(c.Request.Context(), uint(provinceUri.Id))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.NewProvinceRes(province))
}
