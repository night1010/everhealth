package usecase

import (
	"context"
	"fmt"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/valueobject"
)

type PharmacyProductUsecase interface {
	FindAllPharmacyProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindOnePharmacyPeoduct(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) (*entity.PharmacyProduct, error)
	CreatePharmacyProduct(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) (*entity.PharmacyProduct, error)
	UpdatePharmacyProduct(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) (*entity.PharmacyProduct, error)
}

type pharmacyProductUsecase struct {
	pharmacyProductRepository repository.PharmacyProductRepository
	productRepository         repository.ProductRepository
	pharmacyRepository        repository.PharmacyRepository
}

func NewPharmacyProductUsecase(rp repository.PharmacyProductRepository, pr repository.PharmacyRepository, p repository.ProductRepository) PharmacyProductUsecase {
	return &pharmacyProductUsecase{pharmacyProductRepository: rp, productRepository: p, pharmacyRepository: pr}
}

func (u *pharmacyProductUsecase) FindAllPharmacyProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	idPharmacy := query.GetConditionValue("pharmacy").(uint)
	checkPharmacy, err := u.pharmacyRepository.FindById(ctx, idPharmacy)
	if err != nil {
		return nil, err
	}
	if checkPharmacy == nil {
		return nil, apperror.NewResourceNotFoundError("pharmacy", "id", idPharmacy)
	}
	if checkPharmacy.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("dont have access to this pharmacy")
	}
	return u.pharmacyProductRepository.FindAllPharmacyProducts(ctx, query)
}

func (u *pharmacyProductUsecase) FindOnePharmacyPeoduct(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	userId := ctx.Value("user_id").(uint)
	query := valueobject.NewQuery().
		Condition("\"pharmacy_products\".product_id", valueobject.Equal, pharmacyProduct.ProductId).
		Condition("\"pharmacy_products\".pharmacy_id", valueobject.Equal, pharmacyProduct.PharmacyId).
		WithJoin("Product.ProductCategory").WithPreload("Pharmacy")
	selectPharmacyProduct, err := u.pharmacyProductRepository.FindOne(ctx, query)
	if err != nil {
		return nil, err
	}
	if selectPharmacyProduct == nil {
		return nil, apperror.NewResourceNotFoundError("pharmacy product", "id", pharmacyProduct.ProductId)
	}
	if selectPharmacyProduct.Pharmacy.AdminId != userId {
		return nil, apperror.NewResourceNotFoundError("pharmacy product", "id", pharmacyProduct.ProductId)
	}
	return selectPharmacyProduct, nil
}

func (u *pharmacyProductUsecase) CreatePharmacyProduct(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	product, err := u.productRepository.FindById(ctx, pharmacyProduct.ProductId)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, apperror.NewClientError(fmt.Errorf("product with id %v not found", pharmacyProduct.ProductId))
	}

	pharmacy, err := u.pharmacyRepository.FindById(ctx, pharmacyProduct.PharmacyId)
	if err != nil {
		return nil, err
	}

	if pharmacy == nil {
		return nil, apperror.NewResourceNotFoundError("pharmacy", "id", pharmacyProduct.PharmacyId)
	}

	checkPharProduct, err := u.pharmacyProductRepository.FindOne(ctx, valueobject.NewQuery().Condition("pharmacy_id", valueobject.Equal, pharmacyProduct.PharmacyId).Condition("product_id", valueobject.Equal, pharmacyProduct.ProductId))
	if err != nil {
		return nil, err
	}

	if checkPharProduct != nil {
		return nil, apperror.NewClientError(fmt.Errorf("cannot add duplicate product on this pharmacy id %v", pharmacyProduct.PharmacyId))
	}

	if pharmacy.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("cannot have access to add product to this pharmacy")
	}

	newPharmacyProduct, err := u.pharmacyProductRepository.Create(ctx, pharmacyProduct)
	if err != nil {
		return nil, err
	}

	return newPharmacyProduct, nil
}

func (u *pharmacyProductUsecase) UpdatePharmacyProduct(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	pharmacy, err := u.pharmacyRepository.FindById(ctx, pharmacyProduct.PharmacyId)
	if err != nil {
		return nil, err
	}

	if pharmacy == nil {
		return nil, apperror.NewResourceNotFoundError("pharmacy", "id", pharmacyProduct.PharmacyId)
	}

	checkPharProduct, err := u.pharmacyProductRepository.FindById(ctx, pharmacyProduct.Id)
	if err != nil {
		return nil, err
	}

	if checkPharProduct == nil {
		return nil, apperror.NewResourceNotFoundError("pharmacy product", "id", pharmacyProduct.Id)
	}

	if pharmacy.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("cannot have access to add profuct to this pharmacy")
	}
	pharmacyProduct.ProductId = checkPharProduct.ProductId
	pharmacyProduct.Stock = checkPharProduct.Stock
	newPharmacyProduct, err := u.pharmacyProductRepository.Update(ctx, pharmacyProduct)
	if err != nil {
		return nil, err
	}

	return newPharmacyProduct, nil
}
