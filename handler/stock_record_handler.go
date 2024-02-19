package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type StockRecordHandler struct {
	stockRecordUsecase usecase.StockRecordUsecase
}

func NewStockRecordHandler(u usecase.StockRecordUsecase) *StockRecordHandler {
	return &StockRecordHandler{stockRecordUsecase: u}
}

func (h *StockRecordHandler) GetAllStockRecord(c *gin.Context) {
	var request dto.StockRecordParams
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.stockRecordUsecase.FindAllStockRecord(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	stockRecords := pageResult.Data.([]*entity.StockRecord)
	stokcRecordsRes := []*dto.StockRecordRes{}
	for _, StockRecord := range stockRecords {
		StockRecordres := dto.NewStockRecordRes(StockRecord)
		stokcRecordsRes = append(stokcRecordsRes, StockRecordres)
	}
	c.JSON(http.StatusOK, dto.Response{Data: stokcRecordsRes,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}

func (h *StockRecordHandler) PostStockRecord(c *gin.Context) {
	var StockRecordReq dto.StockRecordReq
	err := c.ShouldBindJSON(&StockRecordReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	StockRecord := StockRecordReq.ToModel()
	_, err = h.stockRecordUsecase.CreateStockRecord(c.Request.Context(), StockRecord)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "created success"})
}

func (h *StockRecordHandler) GetStockMonthlyReport(c *gin.Context) {
	var request dto.StockMonthlyReportParams
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.stockRecordUsecase.MonthlyReport(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.Response{Data: pageResult.Data,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}
