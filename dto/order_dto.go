package dto

import (
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	AddressId     uint   `binding:"required" json:"address_id"`
	ShippingName  string `binding:"required" json:"shipping_name"`
	ShippingCost  string `binding:"required" json:"shipping_cost"`
	ShippingEta   string `binding:"required" json:"shipping_eta"`
	PaymentMethod string `binding:"required" json:"payment_method"`
}

type CreateOrderResponse struct {
	OrderId uint `json:"order_id"`
}

type OrderHistoryResponse struct {
	Id            string               `json:"id"`
	OrderItem     []*OrderItemResponse `json:"order_items"`
	ItemOrder     int                  `json:"item_order"`
	OrderDate     string               `json:"order_date"`
	ShippingName  string               `json:"shipping_name"`
	ShippingPrice string               `json:"shipping_price"`
	ShippingEta   string               `json:"shipping_eta"`
	TotalPrice    string               `json:"total_price"`
	OrderStatus   string               `json:"order_status"`
	Name          string               `json:"name,omitempty"`
}

type OrderItemResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	SubTotal string `json:"sub_total"`
	Image    string `json:"image"`
}

type OrderHistoryParam struct {
	Name   *string `form:"name"`
	Status *int    `form:"status" binding:"omitempty,numeric,min=1"`
	SortBy *string `form:"sort_by" binding:"omitempty,oneof=order_date price"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

type OrderUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

type OrderDetailResponse struct {
	Id              string               `json:"id"`
	OrderItem       []*OrderItemResponse `json:"order_items"`
	OrderedAt       string               `json:"ordered_at"`
	ExpiredAt       string               `json:"expired_at"`
	OrderStatus     string               `json:"order_status"`
	ProductPrice    string               `json:"product_price"`
	ShippingName    string               `json:"shipping_name"`
	ShippingPrice   string               `json:"shipping_price"`
	ShippingEta     string               `json:"shipping_eta"`
	TotalPrice      string               `json:"total_price"`
	PaymentProof    string               `json:"payment_proof"`
	Name            string               `json:"name"`
	StreetName      string               `json:"street"`
	PostalCode      string               `json:"postal_code"`
	Phone           string               `json:"phone"`
	Detail          string               `json:"detail"`
	Province        string               `json:"province"`
	City            string               `json:"city"`
	PharmacyContact string               `json:"pharmacy_contact,omitempty"`
	PharmacyEmail   string               `json:"pharmacy_email,omitempty"`
}

type CheckoutItemResponse struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	UnitInPack string `json:"unit_in_pack"`
	Price      string `json:"price"`
	Image      string `json:"image"`
	Quantity   int    `json:"qty"`
	SubTotal   string `json:"sub_total"`
}

type CheckoutResponse struct {
	OrderItems []CheckoutItemResponse `json:"order_item"`
	Total      string                 `json:"total_amount"`
	TotalItem  int                    `json:"total_item"`
}

type UserUpdateOrderStatusRequest struct {
	Status uint `json:"status" binding:"required,oneof=5 6"`
}

type AdminUpdateOrderStatusRequest struct {
	Status uint `json:"status" binding:"required,oneof=3 4 6"`
}

type OrderAddressUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

func (r *UserUpdateOrderStatusRequest) ToOrder(orderId uint) *entity.ProductOrder {
	return &entity.ProductOrder{
		Id:            orderId,
		OrderStatusId: r.Status,
	}
}

func (r *AdminUpdateOrderStatusRequest) ToOrder(orderId uint) *entity.ProductOrder {
	return &entity.ProductOrder{
		Id:            orderId,
		OrderStatusId: r.Status,
	}
}

func (r *CreateOrderRequest) ToOrder() (*entity.ProductOrder, error) {
	price, err := decimal.NewFromString(r.ShippingCost)
	if err != nil {
		return nil, err
	}
	return &entity.ProductOrder{
		AddressId:     r.AddressId,
		ShippingName:  r.ShippingName,
		ShippingPrice: price,
		ShippingEta:   r.ShippingEta,
		PaymentMethod: r.PaymentMethod,
	}, nil
}

func (op *OrderHistoryParam) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()

	if op.Page != nil {
		query.WithPage(*op.Page)
	}
	if op.Limit != nil {
		query.WithLimit(*op.Limit)
	}

	if op.Order != nil {
		query.WithOrder(valueobject.Order(*op.Order))
	}

	if op.SortBy != nil {
		query.WithSortBy(*op.SortBy)
	}

	if op.Name != nil {
		query.Condition("name", valueobject.ILike, *op.Name)
	}

	if op.Status != nil {
		query.Condition("order_status", valueobject.Equal, *op.Status)
	}

	return query, nil
}
