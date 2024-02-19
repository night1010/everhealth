package usecase

import (
	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

type CartUsecase interface {
	GetCart(context.Context) (*entity.Cart, []*entity.CartItem, error)
	AddItem(context.Context, *entity.CartItem) error
	UpdateQty(context.Context, *entity.CartItem) error
	DeleteItem(context.Context, uint) error
	CheckItem(context.Context, uint, bool) error
	CheckAllItem(context.Context, bool) error
}

type cartUsecase struct {
	manager             transactor.Manager
	cartRepo            repository.CartRepository
	cartItemRepo        repository.CartItemRepository
	productRepo         repository.ProductRepository
	pharmacyProductRepo repository.PharmacyProductRepository
}

func NewCartUsecase(
	manager transactor.Manager,
	cartRepo repository.CartRepository,
	cartItemRepo repository.CartItemRepository,
	productRepo repository.ProductRepository,
	pharmacyProductRepo repository.PharmacyProductRepository,
) CartUsecase {
	return &cartUsecase{
		manager:             manager,
		cartRepo:            cartRepo,
		cartItemRepo:        cartItemRepo,
		productRepo:         productRepo,
		pharmacyProductRepo: pharmacyProductRepo,
	}
}

func (u *cartUsecase) GetCart(ctx context.Context) (*entity.Cart, []*entity.CartItem, error) {
	userId := ctx.Value("user_id").(uint)
	cartQuery := valueobject.NewQuery().Condition("user_id", valueobject.Equal, userId)
	cartItemQuery := valueobject.NewQuery().Condition("cart_id", valueobject.Equal, userId).WithJoin("Product").WithSortBy("id")
	fetchedCart, err := u.cartRepo.FindOne(ctx, cartQuery)
	if err != nil {
		return nil, nil, err
	}

	fetchedCartItem, err := u.cartItemRepo.Find(ctx, cartItemQuery)
	if err != nil {
		return nil, nil, err
	}
	fetchedCart.TotalAmount = decimal.Zero
	for _, item := range fetchedCartItem {
		if item.IsChecked {
			fetchedCart.TotalAmount = fetchedCart.TotalAmount.Add(item.SubAmount)
		}
	}
	cart, err := u.cartRepo.Update(ctx, fetchedCart)
	if err != nil {
		return nil, nil, err
	}
	return cart, fetchedCartItem, nil
}

func (u *cartUsecase) AddItem(ctx context.Context, item *entity.CartItem) error {
	userId := ctx.Value("user_id").(uint)
	item.CartId = userId
	cartItemQuery := valueobject.NewQuery().
		Condition("cart_id", valueobject.Equal, userId).
		Condition("product_id", valueobject.Equal, item.ProductId)
	fetchedCartItem, err := u.cartItemRepo.FindOne(ctx, cartItemQuery)
	if err != nil {
		return err
	}
	if fetchedCartItem != nil {
		totalQty := fetchedCartItem.Quantity + item.Quantity
		qty1 := decimal.NewFromInt(int64(fetchedCartItem.Quantity))
		qty2 := decimal.NewFromInt(int64(totalQty))
		fetchedCartItem.SubAmount = fetchedCartItem.SubAmount.Div(qty1).Mul(qty2)
		fetchedCartItem.Quantity = totalQty
		_, err = u.cartItemRepo.Update(ctx, fetchedCartItem)
		if err != nil {
			return err
		}
		return nil
	}
	topPrice, err := u.pharmacyProductRepo.FindTopPrice(ctx, item.ProductId, true)
	if err != nil {
		return err
	}
	if topPrice == decimal.Zero{
		return apperror.NewClientError(apperror.NewResourceStateError("Product not available"))
	}
	item.SubAmount = topPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
	_, err = u.cartItemRepo.Create(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

func (u *cartUsecase) UpdateQty(ctx context.Context, item *entity.CartItem) error {
	userId := ctx.Value("user_id").(uint)
	cartItemQuery := valueobject.NewQuery().
		Condition("id", valueobject.Equal, item.Id).
		Condition("cart_id", valueobject.Equal, userId)
	fetchedCartItem, err := u.cartItemRepo.FindOne(ctx, cartItemQuery)
	if err != nil {
		return err
	}
	if fetchedCartItem == nil {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("cartItem", "id", item.Id))
	}
	if fetchedCartItem.Quantity == 0 {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("cartItem", "id", item.Id))
	}
	qty1 := decimal.NewFromInt(int64(fetchedCartItem.Quantity))
	qty2 := decimal.NewFromInt(int64(item.Quantity))
	fetchedCartItem.SubAmount = fetchedCartItem.SubAmount.Div(qty1).Mul(qty2)
	fetchedCartItem.Quantity = item.Quantity
	_, err = u.cartItemRepo.Update(ctx, fetchedCartItem)
	if err != nil {
		return err
	}
	return nil
}

func (u *cartUsecase) DeleteItem(ctx context.Context, itemId uint) error {
	userId := ctx.Value("user_id").(uint)
	cartItemQuery := valueobject.NewQuery().
		Condition("id", valueobject.Equal, itemId).
		Condition("cart_id", valueobject.Equal, userId)
	fetchedCartItem, err := u.cartItemRepo.FindOne(ctx, cartItemQuery)
	if err != nil {
		return err
	}
	if fetchedCartItem == nil {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("cartItam", "id", itemId))
	}

	err = u.cartItemRepo.Delete(ctx, fetchedCartItem)
	if err != nil {
		return err
	}
	return nil
}

func (u *cartUsecase) CheckItem(ctx context.Context, itemId uint, check bool) error {
	userId := ctx.Value("user_id").(uint)
	cartItemQuery := valueobject.NewQuery().
		Condition("id", valueobject.Equal, itemId).
		Condition("cart_id", valueobject.Equal, userId)
	fetchedCartItem, err := u.cartItemRepo.FindOne(ctx, cartItemQuery)
	if err != nil {
		return err
	}
	if fetchedCartItem == nil {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("cartItam", "id", itemId))
	}
	if fetchedCartItem.IsChecked == check {
		return apperror.NewClientError(apperror.NewResourceStateError("already in that state"))
	}
	fetchedCartItem.IsChecked = check
	_, err = u.cartItemRepo.Update(ctx, fetchedCartItem)
	if err != nil {
		return err
	}
	return nil
}

func (u *cartUsecase) CheckAllItem(ctx context.Context, check bool) error {
	return u.cartItemRepo.CheckAllItem(ctx, check)
}
