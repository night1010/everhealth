package repository

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type DrugFormRepository interface {
	BaseRepository[entity.DrugForm]
	FindAllDrugForm(ctx context.Context) ([]*entity.DrugForm, error)
}

type drugFormRepository struct {
	*baseRepository[entity.DrugForm]
	db *gorm.DB
}

func NewDrugFormRepository(db *gorm.DB) DrugFormRepository {
	return &drugFormRepository{
		db:             db,
		baseRepository: &baseRepository[entity.DrugForm]{db: db},
	}
}

func (r *drugFormRepository) FindAllDrugForm(ctx context.Context) ([]*entity.DrugForm, error) {
	drugForms := make([]*entity.DrugForm, 0)
	err := r.conn(ctx).Find(&drugForms).Error
	if err != nil {
		return drugForms, err
	}
	return drugForms, nil
}
