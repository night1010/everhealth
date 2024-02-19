package dto

import (
	"time"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
)

type StockMutationReq struct {
	ToPharmacyProductId uint `json:"to_pharmacy_product_id" binding:"required,min=1"`
	FromPharmacy        uint `json:"from_pharmacy" binding:"required,min=1"`
	Quantity            int  `json:"quantity" binding:"required,min=1"`
}

type StockMutationPharmacyReq struct {
	ToPharmacyProductId uint `form:"pharmacy_product_id" binding:"required,min=1"`
}

type StockMutationAccept struct {
	IsAccept *bool `json:"is_accept" binding:"required"`
}

type ProductStockMutationRes struct {
	Id       uint                      `json:"id"`
	Name     string                    `json:"name"`
	Image    string                    `json:"image"`
	Pharmacy *PharmacyStockMutationRes `json:"pharmacy,omitempty"`
}

type PharmacyStockMutationRes struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type StockMutationRes struct {
	Id                  uint                       `json:"id"`
	ToPharmacyProduct   *ProductStockMutationRes   `json:"to_pharmacy_product,omitempty"`
	FormPharmacyProduct *ProductStockMutationRes   `json:"form_pharmacy_product,omitempty"`
	Quantity            int                        `json:"quantity"`
	Status              entity.StockMutationStatus `json:"status"`
	MutatedAt           time.Time                  `json:"mutated_at"`
	IsRequest           *bool                      `json:"is_request"`
}

func NewPharmacyStockMutationRes(p *entity.Pharmacy) *PharmacyStockMutationRes {
	return &PharmacyStockMutationRes{Id: p.Id, Name: p.Name}
}

func NewProductStockMutationRes(p *entity.PharmacyProduct) *ProductStockMutationRes {
	var pharmacy *PharmacyStockMutationRes
	if p.Pharmacy != nil {
		pharmacy = NewPharmacyStockMutationRes(p.Pharmacy)
	}
	return &ProductStockMutationRes{Id: p.Id, Name: p.Product.Name, Image: p.Product.Image, Pharmacy: pharmacy}
}
func NewStockMutationRes(u *entity.StockMutation, adminId uint) *StockMutationRes {
	var toProduct *ProductStockMutationRes
	var fromProduct *ProductStockMutationRes
	if u.ToPharmacyProduct != nil {
		toProduct = NewProductStockMutationRes(u.ToPharmacyProduct)
	}
	if u.FromPharmacyProduct != nil {
		fromProduct = NewProductStockMutationRes(u.FromPharmacyProduct)
	}
	isRequest := false
	if adminId == u.FromPharmacyProduct.Pharmacy.AdminId {
		isRequest = true
	}
	return &StockMutationRes{Id: u.Id, ToPharmacyProduct: toProduct, FormPharmacyProduct: fromProduct, Quantity: u.Quantity, Status: u.Status, MutatedAt: u.MutatedAt, IsRequest: &isRequest}
}

type StockMutationParams struct {
	PharmacyName *string `form:"pharmacy_name"`
	Status       *string `form:"status" binding:"omitempty,oneof=1 2 3"`
	SortBy       *string `form:"sort_by" binding:"omitempty,oneof=quantity mutated"`
	Order        *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit        *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page         *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (qp *StockMutationParams) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()
	if qp.PharmacyName != nil {
		query.Condition("pharmacy_name", valueobject.ILike, *qp.PharmacyName)
	}
	if qp.Status != nil {
		var status entity.StockMutationStatus
		switch *qp.Status {

		case "1":
			status = entity.Accept
		case "2":
			status = entity.Decline
		default:
			status = entity.Pending
		}
		query.Condition("status", valueobject.Equal, status)
	}
	if qp.Page != nil {
		query.WithPage(*qp.Page)
	}
	if qp.Limit != nil {
		query.WithLimit(*qp.Limit)
	}
	if qp.Order != nil {
		query.WithOrder(valueobject.Order(*qp.Order))
	}
	if qp.SortBy != nil {
		query.WithSortBy(*qp.SortBy)
	} else {
		query.WithSortBy("id")
	}

	return query, nil
}
