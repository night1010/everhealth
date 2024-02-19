package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type OrderStatusHandler struct {
	OrderStatusUsecase usecase.OrderStatusUsecase
}

func NewOrderStatusHadnler(u usecase.OrderStatusUsecase) *OrderStatusHandler {
	return &OrderStatusHandler{OrderStatusUsecase: u}
}

func (h *OrderStatusHandler) GetAllOrderStatus(c *gin.Context) {
	provincies, err := h.OrderStatusUsecase.FindAllOrderStatus(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	provinciesRes := []*dto.OrderStatusRes{}
	for _, OrderStatus := range provincies {
		OrderStatusres := dto.NewOrderStatusRes(OrderStatus)
		provinciesRes = append(provinciesRes, &OrderStatusres)
	}
	c.JSON(http.StatusOK, provinciesRes)
}
