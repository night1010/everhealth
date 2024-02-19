package dto

import "github.com/night1010/everhealth/entity"

type OrderStatusRes struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func NewOrderStatusRes(c *entity.OrderStatus) OrderStatusRes {
	return OrderStatusRes{Id: c.Id, Name: c.Name}
}
