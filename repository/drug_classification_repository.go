package repository

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type DrugClassificationRepository interface {
	BaseRepository[entity.DrugClassification]
	FindAllDrugClassification(ctx context.Context) ([]*entity.DrugClassification, error)
}

type drugClassificationRepository struct {
	*baseRepository[entity.DrugClassification]
	db *gorm.DB
}

func NewDrugClassificationRepository(db *gorm.DB) DrugClassificationRepository {
	return &drugClassificationRepository{
		db:             db,
		baseRepository: &baseRepository[entity.DrugClassification]{db: db},
	}
}

func (r *drugClassificationRepository) FindAllDrugClassification(ctx context.Context) ([]*entity.DrugClassification, error) {
	drugClassifications := make([]*entity.DrugClassification, 0)
	err := r.conn(ctx).Find(&drugClassifications).Error
	if err != nil {
		return drugClassifications, err
	}
	return drugClassifications, nil
}
