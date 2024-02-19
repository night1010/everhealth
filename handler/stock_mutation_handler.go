package handler

import (
	"context"
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type StockMutationHandler struct {
	stockMutationUsecase usecase.StockMutationUsecase
}

func NewStockMutationHandler(u usecase.StockMutationUsecase) *StockMutationHandler {
	return &StockMutationHandler{stockMutationUsecase: u}
}

func (h *StockMutationHandler) GetAllStockMutation(c *gin.Context) {
	var request dto.StockMutationParams
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.stockMutationUsecase.FindAllStockMutation(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmacies := pageResult.Data.([]*entity.StockMutation)
	pharmaciesRes := []*dto.StockMutationRes{}
	for _, StockMutation := range pharmacies {
		StockMutationres := dto.NewStockMutationRes(StockMutation, c.Request.Context().Value("user_id").(uint))
		pharmaciesRes = append(pharmaciesRes, StockMutationres)
	}
	c.JSON(http.StatusOK, dto.Response{Data: pharmaciesRes,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}

func (h *StockMutationHandler) PostStockMutation(c *gin.Context) {
	var StockMutationReq dto.StockMutationReq
	err := c.ShouldBindJSON(&StockMutationReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	_, err = h.stockMutationUsecase.CreateStockMutation(c.Request.Context(), &StockMutationReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "created success"})
}

func (h *StockMutationHandler) ChangeStatusStockMutation(c *gin.Context) {
	var StockMutationReq dto.StockMutationAccept
	var requestUri dto.RequestUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := c.ShouldBindJSON(&StockMutationReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	ctx := context.WithValue(c.Request.Context(), "is_accept", StockMutationReq.IsAccept)
	c.Request = c.Request.WithContext(ctx)
	_, err = h.stockMutationUsecase.UpdateStockMutation(c.Request.Context(), &entity.StockMutation{Id: uint(requestUri.Id)})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "change status success"})
}

func (h *StockMutationHandler) GetStockMutationDetail(c *gin.Context) {
	var requestUri dto.RequestUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}

	stockMutation, err := h.stockMutationUsecase.GetStockMutationDetail(c.Request.Context(), &entity.StockMutation{Id: uint(requestUri.Id)})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewStockMutationRes(stockMutation, c.Request.Context().Value("user_id").(uint))})
}

func (h *StockMutationHandler) GetAllAvailablePharmacyStockMutation(c *gin.Context) {
	var request dto.StockMutationPharmacyReq
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	pharmacies, err := h.stockMutationUsecase.FindPharmacyAvailable(c.Request.Context(), &entity.PharmacyProduct{Id: request.ToPharmacyProductId})
	if err != nil {
		_ = c.Error(err)
		return
	}
	pharmaciesRes := []*dto.PharmacyStockMutationRes{}
	for _, p := range pharmacies {
		temp := dto.NewPharmacyStockMutationRes(p)
		pharmaciesRes = append(pharmaciesRes, temp)
	}
	c.JSON(http.StatusOK, dto.Response{Data: pharmaciesRes})
}
