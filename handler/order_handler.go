package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(u usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		usecase: u,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	order, err := request.ToOrder()
	if err != nil {
		_ = c.Error(err)
		return
	}
	orderId, err := h.usecase.CreateOrder(c.Request.Context(), order)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.CreateOrderResponse{OrderId: orderId}})
}

func (h *OrderHandler) OrderHistory(c *gin.Context) {
	var request dto.OrderHistoryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pagedResult, err := h.usecase.ListAllOrders(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var listOrders []*dto.OrderHistoryResponse
	data := pagedResult.Data.([]*entity.ProductOrder)
	var orders *dto.OrderHistoryResponse
	for i, item := range data {
		orders = &dto.OrderHistoryResponse{
			Id:            fmt.Sprintf("%s%04d", "7", item.Id),
			OrderDate:     item.OrderedAt.Format(time.RFC3339),
			ShippingName:  item.ShippingName,
			ShippingPrice: item.ShippingPrice.String(),
			ShippingEta:   item.ShippingEta,
			TotalPrice:    item.TotalPayment.String(),
			OrderStatus:   item.OrderStatus.Name,
			ItemOrder:     item.ItemOrderQty,
		}
		if item.Profile.Name != "" {
			orders.Name = item.Profile.Name
		}
		for _, order := range data[i].OrderItems {
			orderItem := &dto.OrderItemResponse{
				Id:       order.PharmacyProduct.ProductId,
				Name:     order.PharmacyProduct.Product.Name,
				Quantity: order.Quantity,
				SubTotal: order.SubTotal.String(),
				Image:    order.PharmacyProduct.Product.Image,
			}
			orders.OrderItem = append(orders.OrderItem, orderItem)
		}
		listOrders = append(listOrders, orders)
	}
	c.JSON(200, dto.Response{
		Data:        listOrders,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *OrderHandler) UploadPaymentProof(c *gin.Context) {
	var requestUri dto.OrderUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.usecase.UploadPaymentProof(c.Request.Context(), requestUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "upload success"})
}

func (h *OrderHandler) OrderDetail(c *gin.Context) {
	var requestUri dto.OrderUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	order, itemOrder, address, err := h.usecase.OrderDetail(c.Request.Context(), requestUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var orderItems []*dto.OrderItemResponse
	for _, item := range itemOrder {
		orderItems = append(orderItems, &dto.OrderItemResponse{
			Id:       item.Id,
			Name:     item.PharmacyProduct.Product.Name,
			Quantity: item.Quantity,
			SubTotal: item.SubTotal.String(),
			Image:    item.PharmacyProduct.Product.Image,
		})
	}
	roleId := c.Request.Context().Value("role_id").(entity.RoleId)
	orderDetailRes := dto.OrderDetailResponse{
		Id:            fmt.Sprintf("%s%04d", "7", order.Id),
		OrderItem:     orderItems,
		OrderedAt:     order.OrderedAt.Format(time.RFC3339),
		ExpiredAt:     order.ExpiredAt.Format(time.RFC3339),
		OrderStatus:   order.OrderStatus.Name,
		PaymentProof:  order.PaymentProof,
		ProductPrice:  (order.TotalPayment.Sub(order.ShippingPrice)).String(),
		ShippingName:  order.ShippingName,
		ShippingPrice: order.ShippingPrice.String(),
		ShippingEta:   order.ShippingEta,
		TotalPrice:    order.TotalPayment.String(),
		Name:          address.Name,
		StreetName:    address.StreetName,
		PostalCode:    address.PostalCode,
		Phone:         address.Phone,
		Detail:        address.Detail,
		Province:      address.Province.Name,
		City:          address.City.Name,
	}
	if roleId == entity.RoleSuperAdmin {
		orderDetailRes.PharmacyContact = order.OrderItems[0].PharmacyProduct.Pharmacy.Admin.AdminContact.Phone
		orderDetailRes.PharmacyEmail = order.OrderItems[0].PharmacyProduct.Pharmacy.Admin.AdminContact.User.Email
	}
	c.JSON(http.StatusOK, dto.Response{Data: orderDetailRes})
}

func (h *OrderHandler) UserUpdateOrderStatus(c *gin.Context) {
	var requestUri dto.OrderUri
	var request dto.UserUpdateOrderStatusRequest
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	order := request.ToOrder(requestUri.Id)
	err := h.usecase.UserUpdateOrderStatus(c.Request.Context(), order)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}

func (h *OrderHandler) AdminUpdateOrderStatus(c *gin.Context) {
	var requestUri dto.OrderUri
	var request dto.AdminUpdateOrderStatusRequest
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	order := request.ToOrder(requestUri.Id)
	err := h.usecase.AdminUpdateOrderStatus(c.Request.Context(), order)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}

func (h *OrderHandler) GetAvailableProduct(c *gin.Context) {
	var orderAddressUri dto.OrderAddressUri
	if err := c.ShouldBindUri(&orderAddressUri); err != nil {
		_ = c.Error(err)
		return
	}
	total, pharmacyProducts, orderItems, _, err := h.usecase.GetAvailableProduct(c.Request.Context(), orderAddressUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var orderDto dto.CheckoutResponse
	orderDto.Total = total.String()
	orderDto.TotalItem = len(pharmacyProducts)
	for _, item := range pharmacyProducts {
		itemDto := dto.CheckoutItemResponse{
			Id:         item.Id,
			Name:       item.Product.Name,
			UnitInPack: item.Product.UnitInPack,
			Price:      item.Price.String(),
			Image:      item.Product.Image,
		}
		for _, orderItem := range orderItems {
			if orderItem.PharmacyProductId == item.Id {
				itemDto.Quantity = orderItem.Quantity
				itemDto.SubTotal = orderItem.SubTotal.String()
			}
		}
		orderDto.OrderItems = append(orderDto.OrderItems, itemDto)
	}
	c.JSON(http.StatusOK, dto.Response{Data: orderDto})
}

func (h *OrderHandler) GetMonthlySaleReport(c *gin.Context) {
	var request dto.MonthlySalesReportParams
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.usecase.MonthlyReport(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.Response{Data: pageResult.Data})
}
