package usecase

import (
	"context"
	"fmt"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type ShippingMethodUsecase interface {
	GetShippingMethod(ctx context.Context, addressId uint) ([]*entity.CalculatedShippingMethod, error)
}

type shippingMethodUsecase struct {
	addressRepo  repository.AddressRepository
	shippingRepo repository.ShippingMethodRepository
	pharmacyRepo repository.PharmacyRepository
	orderUsecase OrderUsecase
}

func NewShippingMethodUsecase(
	addressRepo repository.AddressRepository,
	shippingRepo repository.ShippingMethodRepository,
	pharmacyRepo repository.PharmacyRepository,
	orderUsecase OrderUsecase,
) ShippingMethodUsecase {
	return &shippingMethodUsecase{
		addressRepo:  addressRepo,
		shippingRepo: shippingRepo,
		pharmacyRepo: pharmacyRepo,
		orderUsecase: orderUsecase,
	}
}

func (u *shippingMethodUsecase) GetShippingMethod(ctx context.Context, addressId uint) ([]*entity.CalculatedShippingMethod, error) {
	userId := ctx.Value("user_id").(uint)

	_, products, _, _, err := u.orderUsecase.GetAvailableProduct(ctx, addressId)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, apperror.NewClientError(fmt.Errorf("no pharmacy available"))
	}

	addressQuery := valueobject.NewQuery().
		WithJoin("City").
		Condition("\"addresses\".id", valueobject.Equal, addressId)
	fetchedAddress, err := u.addressRepo.FindOne(ctx, addressQuery)
	if err != nil {
		return nil, err
	}
	if fetchedAddress == nil {
		return nil, apperror.NewResourceNotFoundError("address", "id", addressId)
	}

	if fetchedAddress.ProfileId != userId {
		return nil, apperror.NewForbiddenActionError("not the address of current logged in user")
	}

	fetchedPharmacy, err := u.pharmacyRepo.FindNearestPharmacyFromAddress(ctx, addressId)
	if err != nil {
		return nil, err
	}
	if fetchedPharmacy == nil || len(fetchedPharmacy) < 1 {
		return nil, apperror.NewClientError(fmt.Errorf("there's no pharmacy available near this address"))
	}

	distanceInKM, err := u.shippingRepo.FindDistanceBetween(ctx, fetchedPharmacy[0].Location, fetchedAddress.Location)
	if err != nil {
		return nil, err
	}
	distanceThresholdInKM := decimal.NewFromInt(25)
	calculatedShippingMethods := make([]*entity.CalculatedShippingMethod, 0)

	if distanceInKM.LessThan(distanceThresholdInKM) {
		fetchedShippingMethod, err := u.shippingRepo.Find(ctx, valueobject.NewQuery())
		if err != nil {
			return nil, err
		}
		distanceInKM = distanceInKM.DivRound(decimal.NewFromInt(1), 0)
		for _, fsm := range fetchedShippingMethod {
			calculatedShippingMethods = append(calculatedShippingMethods, &entity.CalculatedShippingMethod{
				Name:              fsm.Name,
				EstimatedDuration: fsm.Duration,
				Cost:              fsm.PricePerKM.Mul(distanceInKM).String(),
			})
		}
	}

	calculatedThirdPartyShippingMethods, _ := u.shippingRepo.GetThirdPartyShipping(ctx, fetchedPharmacy[0].City.Code, fetchedAddress.City.Code, "10")

	if len(calculatedThirdPartyShippingMethods) > 0 {
		calculatedShippingMethods = append(calculatedShippingMethods, calculatedThirdPartyShippingMethods...)
	}

	return calculatedShippingMethods, nil
}
