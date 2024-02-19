package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
)

type StockMutationUsecase interface {
	FindAllStockMutation(ctx context.Context, querry *valueobject.Query) (*valueobject.PagedResult, error)
	GetStockMutationDetail(ctx context.Context, stockMutation *entity.StockMutation) (*entity.StockMutation, error)
	CreateStockMutation(ctx context.Context, stockMutation *dto.StockMutationReq) (*entity.StockMutation, error)
	UpdateStockMutation(ctx context.Context, stockMutation *entity.StockMutation) (*entity.StockMutation, error)
	FindPharmacyAvailable(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) ([]*entity.Pharmacy, error)
}
type stockMutationUsecase struct {
	stockMutationRepository   repository.StockMutationRepository
	stockRecordRepository     repository.StockRecordRepository
	pharmacyProductRepository repository.PharmacyProductRepository
	manager                   transactor.Manager
}

func NewStockMutationUsecase(rp repository.StockMutationRepository, cr repository.PharmacyProductRepository, sr repository.StockRecordRepository, m transactor.Manager) StockMutationUsecase {
	return &stockMutationUsecase{stockMutationRepository: rp, pharmacyProductRepository: cr, manager: m, stockRecordRepository: sr}
}

func (u *stockMutationUsecase) FindAllStockMutation(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {

	return u.stockMutationRepository.FindAllStockMutation(ctx, query)
}

func (u *stockMutationUsecase) GetStockMutationDetail(ctx context.Context, stockMutation *entity.StockMutation) (*entity.StockMutation, error) {
	query := valueobject.NewQuery().Condition("\"stock_mutations\".id", valueobject.Equal, stockMutation.Id).WithJoin("ToPharmacyProduct.Pharmacy").WithJoin("ToPharmacyProduct.Product").WithJoin("FromPharmacyProduct.Pharmacy").WithJoin("FromPharmacyProduct.Product")
	selectStockMutation, err := u.stockMutationRepository.FindOne(ctx, query)
	if err != nil {
		return nil, err
	}
	if selectStockMutation == nil {
		return nil, apperror.NewResourceNotFoundError("stock mutation", "id", stockMutation.Id)
	}
	adminId := ctx.Value("user_id").(uint)
	if selectStockMutation.ToPharmacyProduct.Pharmacy.AdminId == adminId || selectStockMutation.FromPharmacyProduct.Pharmacy.AdminId == adminId {
		return selectStockMutation, nil
	}
	return nil, apperror.NewForbiddenActionError("cannot access this stock mutation")
}

func (u *stockMutationUsecase) CreateStockMutation(ctx context.Context, stockMutation *dto.StockMutationReq) (*entity.StockMutation, error) {
	newStockMutation := &entity.StockMutation{}
	newStockMutation.MutatedAt = time.Now()
	newStockMutation.Status = entity.Pending
	newStockMutation.ToPharmacyProductId = stockMutation.ToPharmacyProductId
	newStockMutation.Quantity = stockMutation.Quantity
	toProduct, err := u.pharmacyProductRepository.FindOne(ctx, valueobject.NewQuery().Condition("\"pharmacy_products\".id", valueobject.Equal, newStockMutation.ToPharmacyProductId).WithJoin("Pharmacy"))
	if err != nil {
		return nil, err
	}
	if toProduct == nil {
		return nil, apperror.NewClientError(fmt.Errorf("pharmacy product with id %v not found", newStockMutation.ToPharmacyProductId))
	}
	if toProduct.Pharmacy.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError(fmt.Sprintf("u dont have access to stock mutation this pharmacy product %v", toProduct.Id))
	}
	pharmacyProduct, err := u.pharmacyProductRepository.FindOne(ctx, valueobject.NewQuery().Condition("pharmacy_id", valueobject.Equal, stockMutation.FromPharmacy).Condition("product_id", valueobject.Equal, toProduct.ProductId))
	if err != nil {
		return nil, err
	}
	if pharmacyProduct == nil {
		return nil, apperror.NewClientError(fmt.Errorf("pharmacy with id %v dont have selected product", stockMutation.FromPharmacy))
	}
	if pharmacyProduct.Id == stockMutation.ToPharmacyProductId {
		return nil, apperror.NewClientError(fmt.Errorf("cannot stock mutation to self"))
	}
	if pharmacyProduct.Stock < int(stockMutation.Quantity) {
		return nil, apperror.NewClientError(fmt.Errorf("insufficient stock request product"))
	}
	newStockMutation.FromPharmacyProductId = pharmacyProduct.Id
	newStockMutation, err = u.stockMutationRepository.Create(ctx, newStockMutation)
	if err != nil {
		return nil, err
	}

	return newStockMutation, nil
}

func (u *stockMutationUsecase) UpdateStockMutation(ctx context.Context, stockMutation *entity.StockMutation) (*entity.StockMutation, error) {
	updateStockMutation, err := u.stockMutationRepository.FindById(ctx, stockMutation.Id)
	if err != nil {
		return nil, err
	}
	if updateStockMutation.Status != entity.Pending {
		return nil, apperror.NewClientError(errors.New("status already change"))
	}
	isAccept := ctx.Value("is_accept").(*bool)
	if !*isAccept {
		updateStockMutation.Status = entity.Decline
		updateStockMutation, err = u.stockMutationRepository.Update(ctx, updateStockMutation)
		if err != nil {
			return nil, err
		}
		return updateStockMutation, nil
	}
	err = u.manager.Run(ctx, func(c context.Context) error {
		updateStockMutation.Status = entity.Accept
		toProduct, err := u.pharmacyProductRepository.FindOne(c, valueobject.NewQuery().Condition("\"pharmacy_products\".id", valueobject.Equal, updateStockMutation.ToPharmacyProductId).Lock())
		if err != nil {
			return err
		}
		fromProduct, err := u.pharmacyProductRepository.FindOne(c, valueobject.NewQuery().Condition("\"pharmacy_products\".id", valueobject.Equal, updateStockMutation.FromPharmacyProductId).WithJoin("Pharmacy"))
		if err != nil {
			return err
		}
		if fromProduct.Pharmacy.AdminId != ctx.Value("user_id").(uint) {
			return apperror.NewForbiddenActionError("cannot change status")
		}
		if fromProduct.Stock < int(updateStockMutation.Quantity) {
			updateStockMutation.Status = entity.Decline
			_, err = u.stockMutationRepository.Update(c, updateStockMutation)
			if err != nil {
				return err
			}
			return apperror.NewClientError(fmt.Errorf("insufficient stock request product"))
		}
		toStockRecord := entity.StockRecord{PharmacyProductId: toProduct.Id, Quantity: updateStockMutation.Quantity, IsReduction: false, ChangeAt: time.Now()}
		fromStockRecord := entity.StockRecord{PharmacyProductId: fromProduct.Id, Quantity: updateStockMutation.Quantity, IsReduction: true, ChangeAt: time.Now()}
		_, err = u.stockRecordRepository.Create(c, &toStockRecord)
		if err != nil {
			return err
		}
		_, err = u.stockRecordRepository.Create(c, &fromStockRecord)
		if err != nil {
			return err
		}
		toProduct.Stock = toProduct.Stock + int(updateStockMutation.Quantity)
		fromProduct.Stock = fromProduct.Stock - int(updateStockMutation.Quantity)
		_, err = u.pharmacyProductRepository.Update(c, toProduct)
		if err != nil {
			return err
		}
		_, err = u.pharmacyProductRepository.Update(c, fromProduct)
		if err != nil {
			return err
		}
		updateStockMutation, err = u.stockMutationRepository.Update(c, updateStockMutation)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updateStockMutation, nil

}

func (u *stockMutationUsecase) FindPharmacyAvailable(ctx context.Context, pharmacyProduct *entity.PharmacyProduct) ([]*entity.Pharmacy, error) {
	checkProduct, err := u.pharmacyProductRepository.FindById(ctx, pharmacyProduct.Id)
	if err != nil {
		return nil, apperror.NewClientError(fmt.Errorf("pharmacy product with id %v not found", pharmacyProduct.Id))
	}
	pharmacyProduct.ProductId = checkProduct.ProductId
	pharmacyProduct.PharmacyId = checkProduct.PharmacyId
	return u.pharmacyProductRepository.FindAllPharmacyAvailableProductId(ctx, pharmacyProduct)
}
