package dto

import (
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type PharmacyProductUri struct {
	PharmacyId uint `uri:"pharmacy_id" binding:"required,numeric"`
	ProductId  uint `uri:"product_id" binding:"required,numeric"`
}

type PharmacyProductReq struct {
	ProductId uint   `json:"product_id" binding:"required"`
	Stock     *int   `json:"stock" binding:"required,min=0"`
	Price     string `json:"price" binding:"required,numeric,mind=1"`
	IsActive  *bool  `json:"is_active" binding:"required"`
}

type PharmacyProductUpdateReq struct {
	Price    string `json:"price" binding:"required,numeric,mind=1"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

func (p *PharmacyProductUpdateReq) ToModel() (*entity.PharmacyProduct, error) {
	price, err := decimal.NewFromString(p.Price)
	if err != nil {
		return nil, err
	}
	return &entity.PharmacyProduct{
		Price:    price,
		IsActive: *p.IsActive,
	}, nil
}
func (p *PharmacyProductReq) ToModel() (*entity.PharmacyProduct, error) {
	price, err := decimal.NewFromString(p.Price)
	if err != nil {
		return nil, err
	}
	return &entity.PharmacyProduct{
		ProductId: p.ProductId,
		Stock:     *p.Stock,
		Price:     price,
		IsActive:  *p.IsActive,
	}, nil
}

type ProductPharmacyRes struct {
	Id       uint             `json:"id"`
	Product  *ProductResponse `json:"product,omitempty"`
	Stock    int              `json:"stock"`
	Price    decimal.Decimal  `json:"price"`
	IsActive bool             `json:"is_active"`
}

func NewProductPhamarcyRes(p *entity.PharmacyProduct) *ProductPharmacyRes {
	var product *ProductResponse

	if p.Product != nil {
		product = NewFromProduct(p.Product)
	}

	return &ProductPharmacyRes{Id: p.Id,
		Product:  product,
		Stock:    p.Stock,
		Price:    p.Price,
		IsActive: p.IsActive}
}

type ListPharmacyProductQueryParam struct {
	Name     *string `form:"name"`
	Category *int    `form:"category" binding:"omitempty,numeric,min=1"`
	SortBy   *string `form:"sort_by" binding:"omitempty,oneof=name stock price"`
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit    *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page     *int    `form:"page" binding:"omitempty,numeric,min=1"`
	IsActive *bool   `form:"is_active" binding:"omitempty"`
}

func (qp *ListPharmacyProductQueryParam) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()

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

	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}

	if qp.Category != nil {
		query.Condition("category", valueobject.Equal, *qp.Category)
	}
	if qp.IsActive != nil {
		query.Condition("is_active", valueobject.Equal, *qp.IsActive)
	}

	return query, nil
}
