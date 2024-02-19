package dto

import "github.com/night1010/everhealth/entity"

type DrugFormRes struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func NewDrugFormres(d *entity.DrugForm) DrugFormRes {
	return DrugFormRes{Id: d.Id, Name: d.Name}
}
