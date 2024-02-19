package dto

import (
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
)

type ProductCategoryParams struct {
	Name   *string `form:"name"`
	SortBy *string `form:"sort_by" binding:"omitempty,oneof=name"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	IsDrug *bool   `form:"is_drug"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
}

func (qp *ProductCategoryParams) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()
	if qp.IsDrug != nil {
		query.Condition("is_drug", valueobject.Equal, qp.IsDrug)
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

	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}

	return query, nil
}

type ProductCategoryUri struct {
	Id int64 `uri:"id" binding:"required,numeric"`
}

type ProductCategoryReq struct {
	Name   string `form:"name" binding:"required"`
	IsDrug *bool  `form:"is_drug" binding:"required"`
}

func (pcr *ProductCategoryReq) ToModel() entity.ProductCategory {
	return entity.ProductCategory{Name: pcr.Name, IsDrug: *pcr.IsDrug}
}

type ProductCategoryRes struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	IsDrug bool   `json:"is_drug"`
	Image  string `json:"image"`
}

func NewProductCategoryRes(pc *entity.ProductCategory) ProductCategoryRes {
	return ProductCategoryRes{Id: pc.Id, Name: pc.Name, IsDrug: pc.IsDrug, Image: pc.Image}
}
