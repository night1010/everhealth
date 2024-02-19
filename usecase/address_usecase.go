package usecase

import (
	"context"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
)

type AddressUsecase interface {
	CreateAddress(context.Context, *entity.Address) (*entity.Address, error)
	GetAddress(context.Context) ([]*entity.Address, error)
	UpdateAddress(context.Context, *entity.Address) error
	DeleteAddress(context.Context, uint) error
	ChangeDefaultAddress(context.Context, uint) error
	ValidateCoordinate(ctx context.Context, cityId uint, location *valueobject.Coordinate) (bool, error)
}

type addressUsecase struct {
	manager      transactor.Manager
	addressRepo  repository.AddressRepository
	shippingRepo repository.ShippingMethodRepository
}

func NewAddressUsecase(
	addressRepo repository.AddressRepository,
	manager transactor.Manager,
	shippingRepo repository.ShippingMethodRepository,
) AddressUsecase {
	return &addressUsecase{
		addressRepo:  addressRepo,
		manager:      manager,
		shippingRepo: shippingRepo,
	}
}

func (u *addressUsecase) CreateAddress(ctx context.Context, address *entity.Address) (*entity.Address, error) {
	userId := ctx.Value("user_id").(uint)
	addressQuery := valueobject.NewQuery().Condition("profile_id", valueobject.Equal, userId)
	fetchedAddress, err := u.addressRepo.FindOne(ctx, addressQuery)
	if err != nil {
		return nil, err
	}
	if fetchedAddress == nil {
		address.IsDefault = true
	} else {
		address.IsDefault = false
	}
	address.ProfileId = userId

	createdAddrss, err := u.addressRepo.Create(ctx, address)
	if err != nil {
		return nil, err
	}
	return createdAddrss, nil
}

func (u *addressUsecase) GetAddress(ctx context.Context) ([]*entity.Address, error) {
	userId := ctx.Value("user_id").(uint)
	addressQuery := valueobject.NewQuery().Condition("profile_id", valueobject.Equal, userId).
		WithJoin("Province").
		WithJoin("City").
		WithSortBy("is_default").
		WithOrder("desc")
	fetchedAddresses, err := u.addressRepo.Find(ctx, addressQuery)
	if err != nil {
		return nil, err
	}
	return fetchedAddresses, nil

}

func (u *addressUsecase) UpdateAddress(ctx context.Context, address *entity.Address) error {
	userId := ctx.Value("user_id").(uint)
	addressQuery := valueobject.NewQuery().
		Condition("id", valueobject.Equal, address.Id).
		Condition("profile_id", valueobject.Equal, userId)
	fetchedAddress, err := u.addressRepo.FindOne(ctx, addressQuery)
	if err != nil {
		return err
	}
	if fetchedAddress == nil {
		return apperror.NewResourceNotFoundError("address", "id", address.Id)
	}

	fetchedAddress.Name = address.Name
	fetchedAddress.StreetName = address.StreetName
	fetchedAddress.PostalCode = address.PostalCode
	fetchedAddress.Phone = address.Phone
	fetchedAddress.Detail = address.Detail
	fetchedAddress.ProvinceId = address.ProvinceId
	fetchedAddress.CityId = address.CityId
	if address.Location != nil {
		fetchedAddress.Location = address.Location
	}
	_, err = u.addressRepo.Update(ctx, fetchedAddress)
	if err != nil {
		return err
	}
	return nil
}

func (u *addressUsecase) DeleteAddress(ctx context.Context, addressId uint) error {
	userId := ctx.Value("user_id").(uint)
	addressQuery := valueobject.NewQuery().
		Condition("id", valueobject.Equal, addressId).
		Condition("profile_id", valueobject.Equal, userId).Lock()

	err := u.manager.Run(ctx, func(c context.Context) error {
		address, err := u.addressRepo.FindOne(c, addressQuery)
		if err != nil {
			return err
		}
		if address == nil {
			return apperror.NewResourceNotFoundError("address", "id", userId)
		}
		err = u.addressRepo.Delete(c, address)
		if err != nil {
			return err
		}

		if address.IsDefault {
			addrressQueryDefault := valueobject.NewQuery().Condition("profile_id", valueobject.Equal, userId).Lock()
			fetchedAddress, err := u.addressRepo.FindOne(c, addrressQueryDefault)
			if err != nil {
				return err
			}
			if fetchedAddress == nil {
				return nil
			}
			fetchedAddress.IsDefault = true
			_, err = u.addressRepo.Update(c, fetchedAddress)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (u *addressUsecase) ChangeDefaultAddress(ctx context.Context, addressId uint) error {
	userId := ctx.Value("user_id").(uint)
	defaultAddressQuery := valueobject.NewQuery().
		Condition("is_default", valueobject.Equal, true).
		Condition("profile_id", valueobject.Equal, userId).Lock()
	addressQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, addressId).Lock()

	err := u.manager.Run(ctx, func(c context.Context) error {
		fetchedDefaultAddress, err := u.addressRepo.FindOne(c, defaultAddressQuery)
		if err != nil {
			return err
		}
		if fetchedDefaultAddress == nil {
			return apperror.NewResourceNotFoundError("address", "id", addressId)
		}

		fetchedAddress, err := u.addressRepo.FindOne(c, addressQuery)
		if err != nil {
			return err
		}
		if fetchedDefaultAddress == nil {
			return apperror.NewResourceNotFoundError("address", "id", addressId)
		}

		fetchedDefaultAddress.IsDefault = false
		fetchedAddress.IsDefault = true
		_, err = u.addressRepo.Update(c, fetchedDefaultAddress)
		if err != nil {
			return err
		}
		_, err = u.addressRepo.Update(c, fetchedAddress)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (u *addressUsecase) ValidateCoordinate(ctx context.Context, cityId uint, location *valueobject.Coordinate) (bool, error) {
	return u.addressRepo.ValidateCoordinate(ctx, cityId, location)
}
