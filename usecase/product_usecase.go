package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/imagehelper"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type ProductUsecase interface {
	ListAllProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, []*entity.Product, map[uint]string, map[uint]string, error)
	ListAllProductAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	ListNearbyProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, []*entity.Product, map[uint]string, map[uint]string, error)
	AddProduct(ctx context.Context, product *entity.Product, drug *entity.Drug) (*entity.Product, error)
	GetProductDetail(ctx context.Context, productId uint) (*entity.Product, decimal.Decimal, decimal.Decimal, error)
	UpdateProduct(ctx context.Context, product *entity.Product, drug *entity.Drug) (*entity.Product, error)
	GetProductDetailAdmin(ctx context.Context, productId uint) (*entity.Product, error)
}

type productUsecase struct {
	manager                transactor.Manager
	imageHelper            imagehelper.ImageHelper
	productRepo            repository.ProductRepository
	categoryRepo           repository.ProductCategoryRepository
	drugRepo               repository.DrugRepository
	drugFormRepo           repository.DrugFormRepository
	drugClassificationRepo repository.DrugClassificationRepository
	pharmacyProductRepo    repository.PharmacyProductRepository
}

func NewProductUsecase(
	manager transactor.Manager,
	imageHelper imagehelper.ImageHelper,
	productRepo repository.ProductRepository,
	categoryRepo repository.ProductCategoryRepository,
	drugRepo repository.DrugRepository,
	drugFormRepo repository.DrugFormRepository,
	drugClassificationRepo repository.DrugClassificationRepository,
	pharmacyProductRepo repository.PharmacyProductRepository,
) ProductUsecase {
	return &productUsecase{
		manager:                manager,
		imageHelper:            imageHelper,
		productRepo:            productRepo,
		categoryRepo:           categoryRepo,
		drugRepo:               drugRepo,
		drugFormRepo:           drugFormRepo,
		drugClassificationRepo: drugClassificationRepo,
		pharmacyProductRepo:    pharmacyProductRepo,
	}
}

func (u *productUsecase) ListAllProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, []*entity.Product, map[uint]string, map[uint]string, error) {
	pagedResult, err := u.productRepo.FindAllProducts(ctx, query)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	products := pagedResult.Data.([]*entity.Product)
	var listOfProduct []uint
	for _, product := range products {
		listOfProduct = append(listOfProduct, product.Id)
	}
	topPrice, err := u.pharmacyProductRepo.FindRangePrice(ctx, listOfProduct, true)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	floorPrice, err := u.pharmacyProductRepo.FindRangePrice(ctx, listOfProduct, false)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return pagedResult, products, topPrice, floorPrice, nil
}

func (u *productUsecase) ListAllProductAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	query.WithJoin("ProductCategory")
	pagedResult, err := u.productRepo.FindAllProducts(ctx, query)
	if err != nil {
		return nil, err
	}

	return pagedResult, nil
}

func (u *productUsecase) ListNearbyProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, []*entity.Product, map[uint]string, map[uint]string, error) {
	userId := ctx.Value("user_id").(uint)
	distanceInMeter := 25_000
	pagedResult, err := u.productRepo.FindNearbyProducts(ctx, query, userId, distanceInMeter)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	products := pagedResult.Data.([]*entity.Product)
	var listOfProduct []uint
	for _, product := range products {
		listOfProduct = append(listOfProduct, product.Id)
	}
	topPrice, err := u.pharmacyProductRepo.FindRangePrice(ctx, listOfProduct, true)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	floorPrice, err := u.pharmacyProductRepo.FindRangePrice(ctx, listOfProduct, false)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return pagedResult, products, topPrice, floorPrice, nil
}

func (u *productUsecase) AddProduct(ctx context.Context, product *entity.Product, drug *entity.Drug) (*entity.Product, error) {
	var createdProduct *entity.Product

	fetchedProductCategory, err := u.categoryRepo.FindById(ctx, product.ProductCategoryId)
	if err != nil {
		return nil, err
	}
	if fetchedProductCategory == nil {
		return nil, apperror.NewResourceNotFoundError("product category", "id", product.ProductCategoryId)
	}

	if fetchedProductCategory.IsDrug {
		if drug == nil {
			return nil, apperror.NewClientError(errors.New("drug should include drug data"))
		}

		fetchedDrugForm, err := u.drugFormRepo.FindById(ctx, drug.DrugFormId)
		if err != nil {
			return nil, err
		}
		if fetchedDrugForm == nil {
			return nil, apperror.NewResourceNotFoundError("drug form", "id", drug.DrugFormId)
		}

		fetchedDrugClassification, err := u.drugClassificationRepo.FindById(ctx, drug.DrugClassificationId)
		if err != nil {
			return nil, err
		}
		if fetchedDrugClassification == nil {
			return nil, apperror.NewResourceNotFoundError("drug classification", "id", drug.DrugClassificationId)
		}

		isDrugAlreadyExist, err := u.drugRepo.IsDrugAlreadyExist(ctx, product.Name, drug.GenericName, product.Manufacture, drug.Content, nil)
		if err != nil {
			return nil, err
		}
		if isDrugAlreadyExist {
			return nil, apperror.NewClientError(errors.New("drug with same name, generic name, manufacture, and content already exist"))
		}

	}

	image := ctx.Value("image")
	imageKey := entity.ProductKeyPrefix + generateRandomString(10)
	imgUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.ProductFolder, imageKey)
	if err != nil {
		return nil, err
	}

	product.Image = imgUrl
	product.ImageKey = imageKey

	err = u.manager.Run(ctx, func(c context.Context) error {
		createdProduct, err = u.productRepo.Create(c, product)
		if err != nil {
			return err
		}

		if fetchedProductCategory.IsDrug {
			drug.Product = *createdProduct
			createdDrug, err := u.drugRepo.Create(c, drug)
			if err != nil {
				return err
			}

			product.Drug = createdDrug
		}
		return nil
	})
	if err != nil {
		err2 := u.imageHelper.Destroy(ctx, entity.ProductFolder, imageKey)
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}

	return createdProduct, nil
}

func (u *productUsecase) UpdateProduct(ctx context.Context, product *entity.Product, drug *entity.Drug) (*entity.Product, error) {
	var updatedProduct *entity.Product
	fetchedProduct, err := u.productRepo.FindById(ctx, product.Id)
	if err != nil {
		return nil, err
	}
	if fetchedProduct == nil {
		return nil, apperror.NewResourceNotFoundError("product", "id", product.Id)
	}
	fetchedProductCategory, err := u.categoryRepo.FindById(ctx, product.ProductCategoryId)
	if err != nil {
		return nil, err
	}
	if fetchedProductCategory == nil {
		return nil, apperror.NewClientError(fmt.Errorf("product category id :%v not found", product.ProductCategoryId))
	}

	if fetchedProductCategory.IsDrug {
		if drug == nil {
			return nil, apperror.NewClientError(errors.New("drug should include drug data"))
		}

		fetchedDrugForm, err := u.drugFormRepo.FindById(ctx, drug.DrugFormId)
		if err != nil {
			return nil, err
		}
		if fetchedDrugForm == nil {
			return nil, apperror.NewClientError(fmt.Errorf("drug form id :%v not found", drug.DrugFormId))
		}

		fetchedDrugClassification, err := u.drugClassificationRepo.FindById(ctx, drug.DrugClassificationId)
		if err != nil {
			return nil, err
		}
		if fetchedDrugClassification == nil {
			return nil, apperror.NewClientError(fmt.Errorf("drug classification id :%v not found", drug.DrugClassificationId))
		}

		isDrugAlreadyExist, err := u.drugRepo.IsDrugAlreadyExist(ctx, product.Name, drug.GenericName, product.Manufacture, drug.Content, &product.Id)
		if err != nil {
			return nil, err
		}
		if isDrugAlreadyExist {
			return nil, apperror.NewClientError(errors.New("drug with same name, generic name, manufacture, and content already exist"))
		}

	}
	product.Image = fetchedProduct.Image
	product.ImageKey = fetchedProduct.ImageKey
	err = u.manager.Run(ctx, func(c context.Context) error {
		image := ctx.Value("image")
		if image != nil {
			imgUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.ProductFolder, entity.ProductCategoryKeyPrefix+generateRandomString(10))
			if err != nil {
				return err
			}
			product.Image = imgUrl
		}
		updatedProduct, err = u.productRepo.Update(c, product)
		if err != nil {
			return err
		}

		if fetchedProductCategory.IsDrug {
			drug.Product = *updatedProduct
			createdDrug, err := u.drugRepo.Update(c, drug)
			if err != nil {
				return err
			}

			product.Drug = createdDrug
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (u *productUsecase) GetProductDetail(ctx context.Context, productId uint) (*entity.Product, decimal.Decimal, decimal.Decimal, error) {
	query := valueobject.NewQuery().
		Condition("\"products\".id", valueobject.Equal, productId).WithPreload("Drug")

	fetchedProduct, err := u.productRepo.FindOne(ctx, query)
	if err != nil {
		return nil, decimal.Zero, decimal.Zero, err
	}
	if fetchedProduct == nil {
		return nil, decimal.Zero, decimal.Zero, apperror.NewResourceNotFoundError("product", "id", productId)
	}
	topPrice, err := u.pharmacyProductRepo.FindTopPrice(ctx, fetchedProduct.Id, true)
	if err != nil {
		return nil, decimal.Zero, decimal.Zero, err
	}
	floorPrice, err := u.pharmacyProductRepo.FindTopPrice(ctx, fetchedProduct.Id, false)
	if err != nil {
		return nil, decimal.Zero, decimal.Zero, err
	}
	return fetchedProduct, topPrice, floorPrice, nil
}

func (u *productUsecase) GetProductDetailAdmin(ctx context.Context, productId uint) (*entity.Product, error) {
	query := valueobject.NewQuery().
		Condition("\"products\".id", valueobject.Equal, productId).WithPreload("Drug.DrugForm").WithPreload("Drug.DrugClassification")

	fetchedProduct, err := u.productRepo.FindOne(ctx, query)
	if err != nil {
		return nil, err
	}
	if fetchedProduct == nil {
		return nil, apperror.NewResourceNotFoundError("product", "id", productId)
	}
	return fetchedProduct, nil
}
