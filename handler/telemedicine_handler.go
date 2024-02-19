package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type TelemedicineHandler struct {
	telemedicineUsecase usecase.TelemedicineUsecase
}

func NewTelemedicineHadnler(u usecase.TelemedicineUsecase) *TelemedicineHandler {
	return &TelemedicineHandler{telemedicineUsecase: u}
}

func (h *TelemedicineHandler) PostTelemedicine(c *gin.Context) {
	var request dto.TelemedicineCreateReq
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	telemedicine, err := h.telemedicineUsecase.CreateTelemedicine(c.Request.Context(), &entity.Telemedicine{DoctorId: request.DoctorId})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewTelemedicinRes(telemedicine), Message: "created success"})
}

func (h *TelemedicineHandler) PostPaymentTelemedicine(c *gin.Context) {
	var request dto.RequestUri
	if err := c.ShouldBindUri(&request); err != nil {
		_ = c.Error(err)
		return
	}
	telemedicine, err := h.telemedicineUsecase.PaymentProofTelemedicine(c.Request.Context(), &entity.Telemedicine{Id: uint(request.Id)})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewTelemedicinRes(telemedicine), Message: "payment success"})
}

func (h *TelemedicineHandler) PostSickLeave(c *gin.Context) {
	var request dto.RequestUri
	var sickLeaveReq dto.SickLeaveReq
	if err := c.ShouldBindUri(&request); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&sickLeaveReq); err != nil {
		_ = c.Error(err)
		return
	}
	sickLeave, err := sickLeaveReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}
	_, err = h.telemedicineUsecase.SickLeaveTelemedicine(c.Request.Context(), &entity.Telemedicine{Id: uint(request.Id)}, sickLeave)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "created sick leave document success"})
}

func (h *TelemedicineHandler) PostPrescription(c *gin.Context) {
	var request dto.RequestUri
	var prescriptionReq dto.Prescription
	if err := c.ShouldBindUri(&request); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&prescriptionReq); err != nil {
		_ = c.Error(err)
		return
	}
	_, err := h.telemedicineUsecase.PrescriptionTelemedicine(c.Request.Context(), &entity.Telemedicine{Id: uint(request.Id)}, &prescriptionReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "prescription document success"})
}

func (h *TelemedicineHandler) GetAllTelemedicine(c *gin.Context) {
	var prescriptionReq dto.TelemedicineParams
	if err := c.ShouldBindQuery(&prescriptionReq); err != nil {
		_ = c.Error(err)
		return
	}
	query := prescriptionReq.ToQuery()
	pageResult, err := h.telemedicineUsecase.FindAllTelemedicine(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	telemedicinesRes := []*dto.TelemedicineRes{}
	for _, telemedecine := range pageResult.Data.([]*entity.Telemedicine) {
		temp := dto.NewTelemedicinRes(telemedecine)
		telemedicinesRes = append(telemedicinesRes, temp)
	}
	c.JSON(http.StatusOK, dto.Response{Data: telemedicinesRes,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}

func (h *TelemedicineHandler) GetTelemedicineDetail(c *gin.Context) {
	var telemedicineUri dto.TelemedicineUri
	if err := c.ShouldBindUri(&telemedicineUri); err != nil {
		_ = c.Error(err)
		return
	}
	telemedicine, err := h.telemedicineUsecase.FindTelemedicine(c.Request.Context(), telemedicineUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewTelemedicinRes(telemedicine)})
}

func (h *TelemedicineHandler) EndChat(c *gin.Context) {
	var telemedicineUri dto.RequestUri
	if err := c.ShouldBindUri(&telemedicineUri); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.telemedicineUsecase.EndTelemedicine(c.Request.Context(), uint(telemedicineUri.Id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Message: "ended telemedicine successfully",
	})
}
