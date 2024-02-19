package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type DrugClassificationHandler struct {
	drugClassificationUsecase usecase.DrugClassificationUsecase
}

func NewDrugClassificationHandler(u usecase.DrugClassificationUsecase) *DrugClassificationHandler {
	return &DrugClassificationHandler{drugClassificationUsecase: u}
}

func (h *DrugClassificationHandler) GetAllDrugClassification(c *gin.Context) {
	drugClassifications, err := h.drugClassificationUsecase.FindAllDrugClassification(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	drugClassificationsRes := []*dto.DrugClassificationRes{}
	for _, drugClassification := range drugClassifications {
		drugClassificationres := dto.NewDrugClassificationres(drugClassification)
		drugClassificationsRes = append(drugClassificationsRes, &drugClassificationres)
	}
	c.JSON(http.StatusOK, drugClassificationsRes)
}
