package usecase

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
)

type DrugClassificationUsecase interface {
	FindAllDrugClassification(ctx context.Context) ([]*entity.DrugClassification, error)
}

type drugClassificationsUsecase struct {
	drugClassificationsRepo repository.DrugClassificationRepository
}

func NewDrugClassificationUsecase(r repository.DrugClassificationRepository) DrugClassificationUsecase {
	return &drugClassificationsUsecase{drugClassificationsRepo: r}
}

func (u *drugClassificationsUsecase) FindAllDrugClassification(ctx context.Context) ([]*entity.DrugClassification, error) {
	return u.drugClassificationsRepo.FindAllDrugClassification(ctx)
}
