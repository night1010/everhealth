package dto

import "github.com/night1010/everhealth/entity"

type DrugResponse struct {
	GenericName          string `json:"generic_name"`
	DrugFormId           string   `json:"drug_form"`
	DrugClassificationId string   `json:"drug_classification"`
	Content              string `json:"content"`
}

func NewFromDrug(d *entity.Drug) *DrugResponse {
	return &DrugResponse{
		GenericName:          d.GenericName,
		DrugFormId:           d.DrugForm.Name,
		DrugClassificationId: d.DrugClassification.Name,
		Content:              d.Content,
	}
}
