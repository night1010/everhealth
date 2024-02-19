package dto

import (
	"github.com/night1010/everhealth/entity"
)

type DrugClassificationRes struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func NewDrugClassificationres(d *entity.DrugClassification) DrugClassificationRes {
	return DrugClassificationRes{Id: d.Id, Name: d.Name}
}
