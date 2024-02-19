package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type ShippingMethodHandler struct {
	usecase usecase.ShippingMethodUsecase
}

func NewShippingMethodHandler(shippingMethodUsecase usecase.ShippingMethodUsecase) *ShippingMethodHandler {
	return &ShippingMethodHandler{usecase: shippingMethodUsecase}
}

func (h *ShippingMethodHandler) GetShippingMethod(c *gin.Context) {
	var addressUri dto.AddressUri
	if err := c.ShouldBindUri(&addressUri); err != nil {
		_ = c.Error(err)
		return
	}

	shippingMethods, err := h.usecase.GetShippingMethod(c.Request.Context(), addressUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var response []dto.ShippingMethodResponse
	for _, method := range shippingMethods {
		response = append(response, dto.ShippingMethodResponse{
			Name:     method.Name,
			Duration: method.EstimatedDuration,
			Cost:     method.Cost,
		})
	}
	c.JSON(http.StatusOK, dto.Response{Data: response})
}
