package usecase

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/imagehelper"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type OrderUsecase interface {
	CreateOrder(context.Context, *entity.ProductOrder) (uint, error)
	GetAvailableProduct(context.Context, uint) (decimal.Decimal, []*entity.PharmacyProduct, []*entity.OrderItem, []*entity.CartItem, error)
	ListAllOrders(context.Context, *valueobject.Query) (*valueobject.PagedResult, error)
	UploadPaymentProof(context.Context, uint) error
	OrderDetail(context.Context, uint) (*entity.ProductOrder, []*entity.OrderItem, *entity.Address, error)
	UserUpdateOrderStatus(context.Context, *entity.ProductOrder) error
	AdminUpdateOrderStatus(context.Context, *entity.ProductOrder) error
	MonthlyReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type orderUsecase struct {
	manager             transactor.Manager
	imageHelper         imagehelper.ImageHelper
	cartRepo            repository.CartRepository
	orderItemRepo       repository.OrderItemRepository
	productOrderRepo    repository.ProductOrderRepository
	cartItemRepo        repository.CartItemRepository
	addressRepo         repository.AddressRepository
	pharmacyRepo        repository.PharmacyRepository
	pharmacyProductRepo repository.PharmacyProductRepository
	stockMutationRepo   repository.StockMutationRepository
	stockRecordRepo     repository.StockRecordRepository
}

func NewOrderUsecase(
	manager transactor.Manager,
	imageHelper imagehelper.ImageHelper,
	cartRepo repository.CartRepository,
	orderItemRepo repository.OrderItemRepository,
	productOrderRepo repository.ProductOrderRepository,
	cartItemRepo repository.CartItemRepository,
	addressRepo repository.AddressRepository,
	pharmacyRepo repository.PharmacyRepository,
	pharmacyProductRepo repository.PharmacyProductRepository,
	stockMutationRepo repository.StockMutationRepository,
	stockRecordRepo repository.StockRecordRepository,
) OrderUsecase {
	return &orderUsecase{
		manager:             manager,
		imageHelper:         imageHelper,
		cartRepo:            cartRepo,
		orderItemRepo:       orderItemRepo,
		productOrderRepo:    productOrderRepo,
		cartItemRepo:        cartItemRepo,
		addressRepo:         addressRepo,
		pharmacyRepo:        pharmacyRepo,
		pharmacyProductRepo: pharmacyProductRepo,
		stockMutationRepo:   stockMutationRepo,
		stockRecordRepo:     stockRecordRepo,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, userOrder *entity.ProductOrder) (uint, error) {
	userId := ctx.Value("user_id").(uint)
	total, pharmacyProducts, orderItems, fetchedCartItem, err := u.GetAvailableProduct(ctx, userOrder.AddressId)
	if err != nil {
		return 0, err
	}
	var orderId uint
	err = u.manager.Run(ctx, func(c context.Context) error {
		var order entity.ProductOrder
		order.OrderedAt = time.Now()
		order.OrderStatusId = uint(entity.WaitingForPayment)
		order.ProfileId = userId
		order.ExpiredAt = time.Now().Add(time.Hour * 24).Truncate(time.Hour).Add(-time.Minute)
		order.ShippingName = userOrder.ShippingName
		order.ShippingPrice = userOrder.ShippingPrice
		order.ShippingEta = userOrder.ShippingEta
		order.TotalPayment = total.Add(userOrder.ShippingPrice)
		order.PaymentMethod = userOrder.PaymentMethod
		order.AddressId = userOrder.AddressId
		order.ItemOrderQty = len(pharmacyProducts)
		createdOrder, err := u.productOrderRepo.Create(c, &order)
		if err != nil {
			return err
		}
		for _, item := range orderItems {
			item.OrderId = createdOrder.Id
		}
		err = u.orderItemRepo.BulkCreate(c, orderItems)
		if err != nil {
			return err
		}
		err = u.cartItemRepo.BulkDelete(c, fetchedCartItem)
		if err != nil {
			return err
		}
		orderId = createdOrder.Id
		return nil
	})
	return orderId, err
}

func (u *orderUsecase) GetAvailableProduct(ctx context.Context, addressId uint) (decimal.Decimal, []*entity.PharmacyProduct, []*entity.OrderItem, []*entity.CartItem, error) {
	userId := ctx.Value("user_id").(uint)
	cartItemQuery := valueobject.NewQuery().
		Condition("cart_id", valueobject.Equal, userId).
		Condition("is_checked", valueobject.Equal, true)
	fetchedCartItem, err := u.cartItemRepo.Find(ctx, cartItemQuery)
	if err != nil {
		return decimal.Zero, nil, nil, nil, err
	}
	if len(fetchedCartItem) == 0 {
		return decimal.Zero, nil, nil, nil, apperror.NewClientError(apperror.NewResourceStateError("no item in cart"))
	}
	pharmacyProductM := make(map[uint]*entity.PharmacyProduct)
	var listOfProductId []uint
	for _, item := range fetchedCartItem {
		listOfProductId = append(listOfProductId, item.ProductId)
		pharmacyProductM[item.ProductId] = &entity.PharmacyProduct{
			Id:    0,
			Stock: item.Quantity,
		}
	}
	fetchedAddress, err := u.addressRepo.FindById(ctx, addressId)
	if err != nil {
		return decimal.Zero, nil, nil, nil, err
	}
	if fetchedAddress.ProfileId != userId {
		return decimal.Zero, nil, nil, nil, apperror.NewClientError(apperror.NewResourceNotFoundError("address", "id", addressId))
	}
	fetchedPP, err := u.pharmacyProductRepo.FindNearbyProductOrder(ctx, listOfProductId, fetchedAddress.Location)
	if err != nil {
		return decimal.Zero, nil, nil, nil, err
	}
	if len(fetchedPP) == 0 {
		return decimal.Zero, nil, nil, nil, apperror.NewResourceStateError("there are no pharmacies nearby that carry the product you are looking for.")
	}
	count := len(listOfProductId)
	var selectedPharmacy uint
	for key, value := range pharmacyProductM {
		if count == 0 {
			break
		}
		for _, pp := range fetchedPP {
			if selectedPharmacy != pp.PharmacyId {
				count = len(listOfProductId)
				selectedPharmacy = pp.PharmacyId
			}
			qty := value.Stock
			if (pp.ProductId == key) && pp.IsActive {
				value.Id = pp.Id
				value.ProductId = pp.ProductId
				value.PharmacyId = pp.PharmacyId
				value.Stock = qty
				value.Price = pp.Price
				count--
				break
			}
		}
		if value.Id == uint(0) {
			return decimal.Zero, nil, nil, nil, apperror.NewResourceStateError("No nearest pharmacy available with your product")
		}
	}
	if count != 0 {
		return decimal.Zero, nil, nil, nil, apperror.NewResourceStateError("No nearest pharmacy available with your product")
	}
	var totalPrice decimal.Decimal
	var orderItems []*entity.OrderItem
	var listOfProductPharmacyId []uint
	for _, value := range pharmacyProductM {
		orderItem := &entity.OrderItem{
			PharmacyProductId: value.Id,
			Quantity:          value.Stock,
			SubTotal:          value.Price.Mul(decimal.NewFromInt(int64(value.Stock))),
		}
		orderItems = append(orderItems, orderItem)
		totalPrice = totalPrice.Add(orderItem.SubTotal)
		listOfProductPharmacyId = append(listOfProductPharmacyId, value.Id)
	}
	productPharmacyQuery := valueobject.NewQuery().Condition("id", valueobject.In, listOfProductPharmacyId).WithPreload("Product")
	fetchedPPResult, err := u.pharmacyProductRepo.Find(ctx, productPharmacyQuery)
	if err != nil {
		return decimal.Zero, nil, nil, nil, err
	}
	return totalPrice, fetchedPPResult, orderItems, fetchedCartItem, nil
}

func (u *orderUsecase) ListAllOrders(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	userId := ctx.Value("user_id").(uint)
	roleId := ctx.Value("role_id").(entity.RoleId)
	return u.productOrderRepo.FindAllOrders(ctx, query, userId, roleId)
}

func (u *orderUsecase) UploadPaymentProof(ctx context.Context, orderId uint) error {
	userId := ctx.Value("user_id").(uint)
	orderQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, orderId)
	fetchedOrder, err := u.productOrderRepo.FindOne(ctx, orderQuery)
	if err != nil {
		return err
	}
	if fetchedOrder == nil {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", orderId))
	}
	if fetchedOrder.ProfileId != userId {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", orderId))
	}
	if fetchedOrder.ExpiredAt.Before(time.Now()) {
		return apperror.NewClientError(apperror.NewResourceStateError("order already expired"))
	}
	if fetchedOrder.OrderStatusId > uint(entity.WaitingForPaymentConfirmation) {
		return apperror.NewClientError(apperror.NewResourceStateError("payment already confirmed"))
	}
	image := ctx.Value("image")
	if image == nil {
		return apperror.NewClientError(apperror.NewResourceStateError("no image inputted"))
	}
	proofKey := fetchedOrder.ProofKey
	if fetchedOrder.ProofKey == "" {
		proofKey = entity.PaymentProofPrefix + generateRandomString(10)
		fetchedOrder.ProofKey = proofKey
	}
	proofUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.PaymentProofFolder, proofKey)
	if err != nil {
		return err
	}
	fetchedOrder.PaymentProof = proofUrl
	fetchedOrder.OrderStatusId = uint(entity.WaitingForPaymentConfirmation)
	_, err = u.productOrderRepo.Update(ctx, fetchedOrder)
	if err != nil {
		return err
	}
	return nil
}

func (u *orderUsecase) OrderDetail(ctx context.Context, orderId uint) (*entity.ProductOrder, []*entity.OrderItem, *entity.Address, error) {
	userId := ctx.Value("user_id").(uint)
	roleId := ctx.Value("role_id").(entity.RoleId)
	fetchedOrder, err := u.productOrderRepo.FindOrderDetail(ctx, orderId, userId, roleId)
	if err != nil {
		return nil, nil, nil, err
	}
	if fetchedOrder.Id == uint(0) {
		return nil, nil, nil, apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", orderId))
	}
	orderItemQuery := valueobject.NewQuery().Condition("order_id", valueobject.Equal, orderId).WithJoin("PharmacyProduct").WithJoin("PharmacyProduct.Product")
	fetchedItemOrder, err := u.orderItemRepo.Find(ctx, orderItemQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	addressQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, fetchedOrder.AddressId).WithPreload("Province").WithPreload("City")
	fetchedAddress, err := u.addressRepo.FindOne(ctx, addressQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	if fetchedAddress == nil {
		return nil, nil, nil, apperror.NewClientError(apperror.NewResourceNotFoundError("address", "id", fetchedOrder.AddressId))
	}
	return fetchedOrder, fetchedItemOrder, fetchedAddress, nil
}

func (u *orderUsecase) UserUpdateOrderStatus(ctx context.Context, order *entity.ProductOrder) error {
	userId := ctx.Value("user_id").(uint)
	orderQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, order.Id)
	fetchedOrder, err := u.productOrderRepo.FindOne(ctx, orderQuery)
	if err != nil {
		return err
	}
	if fetchedOrder == nil {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
	}
	if fetchedOrder.ProfileId != userId {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
	}
	if order.OrderStatusId == uint(entity.Canceled) {
		if fetchedOrder.OrderStatusId >= uint(entity.WaitingForPaymentConfirmation) {
			return apperror.NewClientError(apperror.NewResourceStateError("cant cancel order"))
		}
		fetchedOrder.OrderStatusId = uint(entity.Canceled)
	}
	if order.OrderStatusId == uint(entity.OrderConfirmed) {
		if fetchedOrder.OrderStatusId != uint(entity.Sent) {
			return apperror.NewClientError(apperror.NewResourceStateError("cant confirm order"))
		}
		fetchedOrder.OrderStatusId = uint(entity.OrderConfirmed)
	}
	_, err = u.productOrderRepo.Update(ctx, fetchedOrder)
	if err != nil {
		return err
	}
	return nil
}

func (u *orderUsecase) AdminUpdateOrderStatus(ctx context.Context, order *entity.ProductOrder) error {
	roleId := ctx.Value("role_id").(entity.RoleId)
	userId := ctx.Value("user_id").(uint)
	if roleId != entity.RoleAdmin {
		return apperror.NewForbiddenActionError("you're not admin")
	}
	err := u.manager.Run(ctx, func(c context.Context) error {
		orderQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, order.Id).Lock()
		fetchedOrder, err := u.productOrderRepo.FindOne(c, orderQuery)
		if err != nil {
			return err
		}
		if fetchedOrder == nil {
			return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
		}
		switch order.OrderStatusId {
		case uint(entity.Processed):
			if fetchedOrder.OrderStatusId != uint(entity.WaitingForPaymentConfirmation) {
				return apperror.NewClientError(apperror.NewResourceStateError("cant update order status to processed"))
			}
			fetchedOrderItem, err := u.orderItemRepo.ListOfOrderItem(c, order.Id, userId)
			if err != nil {
				return err
			}
			if len(fetchedOrderItem) == 0 {
				return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
			}
			pharmacyProductM := make(map[uint]*entity.PharmacyProduct)
			nearestPharmacyLoc := fetchedOrderItem[0].PharmacyProduct.Pharmacy.Location
			var listOfProductId []uint
			var stockMutations []*entity.StockMutation
			var stockRecords []*entity.StockRecord
			for _, item := range fetchedOrderItem {
				if item.PharmacyProduct.Stock >= item.Quantity {
					item.PharmacyProduct.Stock -= item.Quantity
					_, err := u.pharmacyProductRepo.Update(c, &item.PharmacyProduct)
					if err != nil {
						return err
					}
				} else {
					item.PharmacyProduct.Stock = item.Quantity - item.PharmacyProduct.Stock
					pharmacyProductM[item.PharmacyProduct.ProductId] = &item.PharmacyProduct
					listOfProductId = append(listOfProductId, item.PharmacyProduct.ProductId)
				}
			}
			if len(listOfProductId) != 0 {
				fetchedPP, err := u.pharmacyProductRepo.FindNearbyProductOrder(c, listOfProductId, nearestPharmacyLoc)
				if err != nil {
					return err
				}
				if len(fetchedPP) == 0 {
					fetchedOrder.OrderStatusId = uint(entity.Canceled)
					return apperror.NewResourceStateError("No nearest pharmacy found, order cancelled")
				}
				for key, value := range pharmacyProductM {
					for _, pp := range fetchedPP {
						if (pp.ProductId == key) && (pp.Stock >= value.Stock) && (value.PharmacyId != pp.PharmacyId) {
							sm, r1, r2 := createStockMutation(pp.Id, value.Id, value.Stock)
							sm.OrderId = fetchedOrder.Id
							stockMutations = append(stockMutations, sm)
							stockRecords = append(stockRecords, r1, r2)
							pp.Stock -= value.Stock
							value.Stock = 0
							_, err = u.pharmacyProductRepo.Update(c, pp)
							if err != nil {
								return err
							}
							_, err = u.pharmacyProductRepo.Update(c, value)
							if err != nil {
								return err
							}
							value = nil
							break
						}
					}
					if value != nil {
						fetchedOrder.OrderStatusId = uint(entity.Canceled)
						return apperror.NewResourceStateError("Cancel order")
					}
				}
				err = u.stockMutationRepo.BulkCreate(c, stockMutations)
				if err != nil {
					return err
				}
				err = u.stockRecordRepo.BulkCreate(c, stockRecords)
				if err != nil {
					return err
				}
			}
			fetchedOrder.OrderStatusId = uint(entity.Processed)
		case uint(entity.Sent):
			if fetchedOrder.OrderStatusId != uint(entity.Processed) {
				return apperror.NewClientError(apperror.NewResourceStateError("cant update order status to sent"))
			}
			fetchedOrder.OrderStatusId = uint(entity.Sent)
		case uint(entity.Canceled):
			if fetchedOrder.OrderStatusId >= uint(entity.Sent) {
				return apperror.NewClientError(apperror.NewResourceStateError("cant cancel order"))
			}
			if fetchedOrder.OrderStatusId == uint(entity.Processed) {
				stockMutationQuery := valueobject.NewQuery().Condition("order_id", valueobject.Equal, fetchedOrder.Id)
				fetchedSM, err := u.stockMutationRepo.Find(c, stockMutationQuery)
				if err != nil {
					return err
				}
				if len(fetchedSM) != 0 {
					var stockMutations []*entity.StockMutation
					var stockRecords []*entity.StockRecord
					for _, stockMutation := range fetchedSM {
						sm, r1, r2 := createStockMutation(stockMutation.ToPharmacyProductId, stockMutation.FromPharmacyProductId, stockMutation.Quantity)
						stockMutations = append(stockMutations, sm)
						stockRecords = append(stockRecords, r1, r2)
						ppId := [2]uint{stockMutation.FromPharmacyProductId, stockMutation.ToPharmacyProductId}
						ppQuery := valueobject.NewQuery().Condition("id", valueobject.In, ppId).Lock()
						fetchedPP, err := u.pharmacyProductRepo.Find(c, ppQuery)
						if err != nil {
							return err
						}
						itemQuery := valueobject.NewQuery().
							Condition("order_id", valueobject.Equal, fetchedOrder.Id).
							Condition("pharmacy_product_id", valueobject.Equal, stockMutation.ToPharmacyProductId)
						fetchedItem, err := u.orderItemRepo.FindOne(c, itemQuery)
						if err != nil {
							return err
						}
						for _, pp := range fetchedPP {
							if pp.Id == stockMutation.FromPharmacyProductId {
								pp.Stock += stockMutation.Quantity
							} else if pp.Id == stockMutation.ToPharmacyProductId {
								pp.Stock = pp.Stock + fetchedItem.Quantity - stockMutation.Quantity
							}
							_, err = u.pharmacyProductRepo.Update(c, pp)
							if err != nil {
								return err
							}
						}
					}
					err = u.stockMutationRepo.BulkCreate(c, stockMutations)
					if err != nil {
						return err
					}
					err = u.stockRecordRepo.BulkCreate(c, stockRecords)
					if err != nil {
						return err
					}
				}
			}
			fetchedOrder.OrderStatusId = uint(entity.Canceled)
		}
		_, err = u.productOrderRepo.Update(c, fetchedOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (u *orderUsecase) MonthlyReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	if ctx.Value("role_id").(entity.RoleId) == entity.RoleAdmin {
		return u.orderItemRepo.MonthlyReportAdminPharmacy(ctx, query)
	}

	return u.orderItemRepo.MonthlyReport(ctx, query)
}

func createStockMutation(from, to uint, qty int) (*entity.StockMutation, *entity.StockRecord, *entity.StockRecord) {
	stockMutation := &entity.StockMutation{
		ToPharmacyProductId:   to,
		FromPharmacyProductId: from,
		Quantity:              qty,
		Status:                entity.Accept,
		MutatedAt:             time.Now(),
	}
	record1 := &entity.StockRecord{
		PharmacyProductId: from,
		Quantity:          qty,
		IsReduction:       true,
		ChangeAt:          time.Now(),
	}
	record2 := &entity.StockRecord{
		PharmacyProductId: to,
		Quantity:          qty,
		IsReduction:       false,
		ChangeAt:          time.Now(),
	}
	return stockMutation, record1, record2
}
