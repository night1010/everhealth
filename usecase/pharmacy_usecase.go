package usecase

import (
	"context"
	"fmt"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/valueobject"
)

type PharmacyUsecase interface {
	FindAllPharmacy(ctx context.Context, querry *valueobject.Query) (*valueobject.PagedResult, error)
	FindAllPharmacySuperAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindOnePharmacyDetail(ctx context.Context, pharmacy *entity.Pharmacy) (*entity.Pharmacy, error)
	CreatePharmacy(ctx context.Context, pharmacy *entity.Pharmacy) (*entity.Pharmacy, error)
	UpdatePharmacy(ctx context.Context, pharmacy *entity.Pharmacy) (*entity.Pharmacy, error)
	DeletePharmacy(ctx context.Context, pharmacy entity.Pharmacy) error
}

type pharmacyUsecase struct {
	pharmacyRepository repository.PharmacyRepository
	provinceRepository repository.ProvinceRepository
	cityRepository     repository.CityRepository
}

func NewPharmacyUsecase(rp repository.PharmacyRepository, pr repository.ProvinceRepository, cr repository.CityRepository) PharmacyUsecase {
	return &pharmacyUsecase{pharmacyRepository: rp, provinceRepository: pr, cityRepository: cr}
}

func (u *pharmacyUsecase) FindAllPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.pharmacyRepository.FindAllPharmacy(ctx, query)
}

func (u *pharmacyUsecase) FindAllPharmacySuperAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.pharmacyRepository.FindAllPharmacySuperAdmin(ctx, query)
}

func (u *pharmacyUsecase) FindOnePharmacyDetail(ctx context.Context, pharmacy *entity.Pharmacy) (*entity.Pharmacy, error) {
	selectPharmacy, err := u.pharmacyRepository.FindOne(ctx, valueobject.NewQuery().Condition("\"pharmacies\".id", valueobject.Equal, pharmacy.Id).WithJoin("Province").WithJoin("City"))
	if err != nil {
		return nil, err
	}
	if selectPharmacy == nil {
		return nil, apperror.NewResourceNotFoundError("pharmacy", "id", pharmacy.Id)
	}
	return selectPharmacy, nil
}

func (u *pharmacyUsecase) CreatePharmacy(ctx context.Context, pharmacy *entity.Pharmacy) (*entity.Pharmacy, error) {
	province, err := u.provinceRepository.FindById(ctx, pharmacy.ProvinceId)
	if err != nil {
		return nil, err
	}
	if province == nil {
		return nil, apperror.NewClientError(fmt.Errorf("province with id %v not found", pharmacy.ProvinceId))
	}
	city, err := u.cityRepository.FindById(ctx, pharmacy.CityId)
	if err != nil {
		return nil, err
	}

	if city == nil {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v not found", pharmacy.CityId))
	}
	if province.Id != city.ProvinceId {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v doesn't belong to province id %v", city.Id, pharmacy.ProvinceId))
	}
	pharmacy.AdminId = ctx.Value("user_id").(uint)
	newPharmacy, err := u.pharmacyRepository.Create(ctx, pharmacy)
	if err != nil {
		return nil, err
	}
	return newPharmacy, nil
}

func (u *pharmacyUsecase) UpdatePharmacy(ctx context.Context, pharmacy *entity.Pharmacy) (*entity.Pharmacy, error) {
	checkPharmacy, err := u.pharmacyRepository.FindById(ctx, pharmacy.Id)
	if err != nil {
		return nil, err
	}

	if checkPharmacy == nil {
		return nil, apperror.NewResourceNotFoundError("pharmacy", "id", pharmacy.Id)
	}
	province, err := u.provinceRepository.FindById(ctx, pharmacy.ProvinceId)
	if err != nil {
		return nil, err
	}

	if province == nil {
		return nil, apperror.NewClientError(fmt.Errorf("province with id %v not found", pharmacy.ProvinceId))
	}

	city, err := u.cityRepository.FindById(ctx, pharmacy.CityId)
	if err != nil {
		return nil, err
	}

	if city == nil {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v not found", pharmacy.CityId))
	}

	if province.Id != city.ProvinceId {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v doesn't belong to province id %v", city.Id, pharmacy.ProvinceId))
	}

	if checkPharmacy.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("cannot have access to update this pharmacy")
	}

	pharmacy.AdminId = checkPharmacy.AdminId
	updatedPharmacy, err := u.pharmacyRepository.Update(ctx, pharmacy)
	if err != nil {
		return nil, err
	}
	return updatedPharmacy, nil
}

func (u *pharmacyUsecase) DeletePharmacy(ctx context.Context, pharmacy entity.Pharmacy) error {
	checkPharmacy, err := u.pharmacyRepository.FindById(ctx, pharmacy.Id)
	if err != nil {
		return err
	}
	if checkPharmacy == nil {
		return apperror.NewResourceNotFoundError("product category", "id", pharmacy.Id)
	}
	if checkPharmacy.AdminId != ctx.Value("user_id").(uint) {
		return apperror.NewForbiddenActionError("cannot have access to delete this pharmacy")
	}
	return u.pharmacyRepository.Delete(ctx, &pharmacy)
}
