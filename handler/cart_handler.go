package handler

import (
	"net/http"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type CartHandler struct {
	usecase usecase.CartUsecase
}

func NewCartHandler(u usecase.CartUsecase) *CartHandler {
	return &CartHandler{
		usecase: u,
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	cart, cartItem, err := h.usecase.GetCart(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	var cartDto dto.CartResponse
	cartDto.Total = cart.TotalAmount.String()
	cartDto.TotalItem = len(cartItem)
	for _, item := range cartItem {
		cartItemRes := dto.CartItemResponse{
			Id:           item.Id,
			Name:         item.Product.Name,
			UnitInPack:   item.Product.UnitInPack,
			PricePerPack: (item.SubAmount.Div(decimal.NewFromInt(int64(item.Quantity)))).String(),
			Image:        item.Product.Image,
			Quantity:     item.Quantity,
			SubTotal:     item.SubAmount.String(),
			IsChecked:    item.IsChecked,
		}
		cartDto.CartItem = append(cartDto.CartItem, cartItemRes)
	}
	c.JSON(http.StatusOK, dto.Response{Data: cartDto})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var request dto.AddItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	item := request.ToItem()
	err := h.usecase.AddItem(c.Request.Context(), item)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "cart item created"})
}

func (h *CartHandler) ChangeQty(c *gin.Context) {
	var itemUri dto.CartItemUri
	var request dto.ChangeQtyRequest
	if err := c.ShouldBindUri(&itemUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	item := request.ToItem(itemUri.Id)
	err := h.usecase.UpdateQty(c.Request.Context(), item)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "item qty updated"})
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	var itemUri dto.CartItemUri
	if err := c.ShouldBindUri(&itemUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.usecase.DeleteItem(c.Request.Context(), itemUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "delete success"})
}

func (h *CartHandler) CheckItem(c *gin.Context) {
	var cartUri dto.CartItemUri
	var request dto.CartCheckRequest
	if err := c.ShouldBindUri(&cartUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.usecase.CheckItem(c.Request.Context(), cartUri.Id, request.IsCheck)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}

func (h *CartHandler) CheckAllItem(c *gin.Context) {
	var request dto.CartCheckRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.usecase.CheckAllItem(c.Request.Context(), request.IsCheck)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}
