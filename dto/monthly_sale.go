package dto

import "github.com/night1010/everhealth/valueobject"

type MonthlySalesReport struct {
	Pharmacy          *string `json:"pharmacy,omitempty"`
	Product           *string `json:"product,omitempty"`
	ProductCategoryId *uint   `json:"product_category_id,omitempty"`
	TotalItem         uint    `json:"total_item"`
	TotalSales        uint    `json:"total_sales"`
	Month             string  `json:"month"`
}

type DataGraphReport struct {
	PharmacyGraph        []*MonthlySalesReport `json:"pharmacy_graph"`
	ProductCategoryGraph []*MonthlySalesReport `json:"product_category_graph"`
	ProductGraph         []*MonthlySalesReport `json:"product_graph"`
}

type MonthlySalesReportParams struct {
	Pharmacy        *uint `form:"pharmacy" binding:"numeric,min=1"`
	Product         *uint `form:"product" binding:"numeric,min=1"`
	ProductCategory *uint `form:"product_category" binding:"numeric,min=1"`
}

func (qp *MonthlySalesReportParams) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()
	if qp.Pharmacy != nil {
		query.Condition("pharmacy", valueobject.Equal, *qp.Pharmacy)
	}
	if qp.Product != nil {
		query.Condition("product", valueobject.Equal, *qp.Product)
	}
	if qp.ProductCategory != nil {
		query.Condition("product_category", valueobject.Equal, *qp.ProductCategory)
	}

	return query, nil
}
