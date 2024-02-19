package dto

import "github.com/night1010/everhealth/entity"

type CartItemResponse struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	UnitInPack   string `json:"unit_in_pack"`
	PricePerPack string `json:"price"`
	Image        string `json:"image"`
	Quantity     int    `json:"qty"`
	SubTotal     string `json:"sub_total"`
	IsChecked    bool   `json:"is_checked"`
}

type CartResponse struct {
	CartItem  []CartItemResponse `json:"cart_item"`
	Total     string             `json:"total_amount"`
	TotalItem int                `json:"total_item"`
}

type AddItemRequest struct {
	ProductId uint `json:"product_id"`
	Quantity  int  `json:"qty"`
}

type ChangeQtyRequest struct {
	Quantity int `json:"qty" binding:"required,min=1"`
}

type CartItemUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

type CartCheckRequest struct {
	IsCheck bool `json:"is_check"`
}

func (r *AddItemRequest) ToItem() *entity.CartItem {
	return &entity.CartItem{
		ProductId: r.ProductId,
		Quantity:  r.Quantity,
		IsChecked: false,
	}
}

func (r *ChangeQtyRequest) ToItem(itemId uint) *entity.CartItem {
	return &entity.CartItem{
		Id:       itemId,
		Quantity: r.Quantity,
	}
}
