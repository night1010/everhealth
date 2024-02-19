package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
)

type StockRecordUsecase interface {
	FindAllStockRecord(ctx context.Context, querry *valueobject.Query) (*valueobject.PagedResult, error)
	CreateStockRecord(ctx context.Context, StockRecord *entity.StockRecord) (*entity.StockRecord, error)
	MonthlyReport(ctx context.Context, querry *valueobject.Query) (*valueobject.PagedResult, error)
}

type stockRecordUsecase struct {
	stockRecordRepository     repository.StockRecordRepository
	pharmacyProductRepository repository.PharmacyProductRepository
	manager                   transactor.Manager
}

func NewStockRecordUsecase(rp repository.StockRecordRepository, cr repository.PharmacyProductRepository, m transactor.Manager) StockRecordUsecase {
	return &stockRecordUsecase{stockRecordRepository: rp, pharmacyProductRepository: cr, manager: m}
}

func (u *stockRecordUsecase) FindAllStockRecord(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.stockRecordRepository.FindAllStockRecord(ctx, query)
}

func (u *stockRecordUsecase) CreateStockRecord(ctx context.Context, stockRecord *entity.StockRecord) (*entity.StockRecord, error) {
	var newStockRecord *entity.StockRecord
	stockRecord.ChangeAt = time.Now()
	pharmacyProduct, err := u.pharmacyProductRepository.FindOne(ctx, valueobject.NewQuery().Condition("\"pharmacy_products\".id", valueobject.Equal, stockRecord.PharmacyProductId).WithJoin("Pharmacy"))
	if err != nil {
		return nil, err
	}
	if pharmacyProduct == nil {
		return nil, apperror.NewClientError(fmt.Errorf("product with id %v not found", stockRecord.PharmacyProductId))
	}
	if pharmacyProduct.Pharmacy.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("cannot have access to change stock")
	}
	err = u.manager.Run(ctx, func(c context.Context) error {
		_, err = u.pharmacyProductRepository.FindOne(c, valueobject.NewQuery().Condition("\"pharmacy_products\".id", valueobject.Equal, stockRecord.PharmacyProductId).Lock())
		if err != nil {
			return err
		}

		newStockRecord, err = u.stockRecordRepository.Create(c, stockRecord)
		if err != nil {
			return err
		}
		number := int(pharmacyProduct.Stock)
		if stockRecord.IsReduction {
			number -= int(stockRecord.Quantity)
		} else {
			number += int(stockRecord.Quantity)
		}
		if number < 0 {
			return apperror.NewClientError(errors.New("product's stock cannot below zero"))
		}
		pharmacyProduct.Stock = number
		_, err = u.pharmacyProductRepository.Update(c, pharmacyProduct)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newStockRecord, nil
}

func (u *stockRecordUsecase) MonthlyReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.stockRecordRepository.MonthlyReport(ctx, query)
}
