package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/imagehelper"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
)

type ProductCategoryUsecase interface {
	GetProductCategories(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	GetProductCategoriesDetail(ctx context.Context, productCategory *entity.ProductCategory) (*entity.ProductCategory, error)
	CreateProductCategory(ctx context.Context, productCategory entity.ProductCategory) (*entity.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, productCategory entity.ProductCategory) (*entity.ProductCategory, error)
	DeleteProductCategories(ctx context.Context, productCategory entity.ProductCategory) error
}

type productCategoryUsecase struct {
	productCategoryRepo repository.ProductCategoryRepository
	imageHelper         imagehelper.ImageHelper
	manager             transactor.Manager
}

func NewProductCategoryUsecase(pcr repository.ProductCategoryRepository, img imagehelper.ImageHelper, manager transactor.Manager) ProductCategoryUsecase {

	return &productCategoryUsecase{
		productCategoryRepo: pcr,
		imageHelper:         img,
		manager:             manager,
	}
}

func (u *productCategoryUsecase) GetProductCategories(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	query.WithSortBy("name")
	return u.productCategoryRepo.FindProductCategories(ctx, query)
}

func (u *productCategoryUsecase) GetProductCategoriesDetail(ctx context.Context, productCategory *entity.ProductCategory) (*entity.ProductCategory, error) {
	selectProductCategory, err := u.productCategoryRepo.FindById(ctx, productCategory.Id)
	if err != nil {
		return nil, err
	}
	if selectProductCategory == nil {
		return nil, apperror.NewResourceNotFoundError("product category", "id", productCategory.Id)
	}
	return selectProductCategory, nil
}
func (u *productCategoryUsecase) CreateProductCategory(ctx context.Context, productCategory entity.ProductCategory) (*entity.ProductCategory, error) {
	checkProduct, err := u.productCategoryRepo.FindOne(ctx, valueobject.NewQuery().Condition("lower(name)", valueobject.Equal, strings.ToLower(productCategory.Name)))
	if err != nil {
		return nil, err
	}
	if checkProduct != nil {
		return nil, apperror.NewClientError(errors.New("cannot add duplicate product category"))
	}
	image := ctx.Value("image")
	imageKey := entity.ProductCategoryKeyPrefix + generateRandomString(10)
	imgUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.ProductCategoryFolder, imageKey)
	if err != nil {
		return nil, err
	}
	productCategory.Image = imgUrl
	productCategory.ImageKey = imageKey
	createdProductCategory, err := u.productCategoryRepo.Create(ctx, &productCategory)
	if err != nil {
		err2 := u.imageHelper.Destroy(ctx, entity.ProductCategoryFolder, imageKey)
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	return createdProductCategory, nil
}

func (u *productCategoryUsecase) UpdateProductCategory(ctx context.Context, productCategory entity.ProductCategory) (*entity.ProductCategory, error) {
	image := ctx.Value("image")
	var updatedProductCategory *entity.ProductCategory
	checkProductCategory, err := u.productCategoryRepo.FindById(ctx, productCategory.Id)
	if err != nil {
		return nil, err
	}
	if checkProductCategory == nil {
		return nil, apperror.NewResourceNotFoundError("product category", "id", productCategory.Id)
	}
	checkProduct, err := u.productCategoryRepo.FindOne(ctx, valueobject.NewQuery().Condition("lower(name)", valueobject.Equal, strings.ToLower(productCategory.Name)).Condition("id", valueobject.NotEqual, productCategory.Id))
	if err != nil {
		return nil, err
	}
	if checkProduct != nil {
		return nil, apperror.NewClientError(errors.New("cannot add duplicate product category"))
	}
	checkProductCategory.Name = productCategory.Name
	checkProductCategory.IsDrug = productCategory.IsDrug
	err = u.manager.Run(ctx, func(c context.Context) error {
		if image != nil {
			imgUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.ProductCategoryFolder, entity.ProductCategoryKeyPrefix+generateRandomString(10))
			if err != nil {
				return err
			}
			checkProductCategory.Image = imgUrl
		}

		updatedProductCategory, err = u.productCategoryRepo.Update(c, checkProductCategory)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedProductCategory, nil
}

func (u *productCategoryUsecase) DeleteProductCategories(ctx context.Context, productCategory entity.ProductCategory) error {
	checkProductCategory, err := u.productCategoryRepo.FindById(ctx, productCategory.Id)
	if err != nil {
		return err
	}
	if checkProductCategory == nil {
		return apperror.NewResourceNotFoundError("product category", "id", productCategory.Id)
	}
	return u.productCategoryRepo.Delete(ctx, &productCategory)
}
