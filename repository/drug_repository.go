package repository

import (
	"context"
	"errors"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type DrugRepository interface {
	BaseRepository[entity.Drug]
	IsDrugAlreadyExist(ctx context.Context, name string, genericName string, manufacture string, content string, productId *uint) (bool, error)
}

type drugRepository struct {
	*baseRepository[entity.Drug]
	db *gorm.DB
}

func NewDrugRepository(db *gorm.DB) DrugRepository {
	return &drugRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Drug]{db: db},
	}
}

func (r *drugRepository) IsDrugAlreadyExist(ctx context.Context, name string, genericName string, manufacture string, content string, productId *uint) (bool, error) {
	var drug *entity.Drug

	conn := r.conn(ctx).
		Joins("Product").
		Where("name = ?", name).
		Where("generic_name = ?", genericName).
		Where("manufacture = ?", manufacture).
		Where("content = ?", content)
	if productId != nil {
		conn.Not("id = ?", *productId)
	}
	err := conn.First(&drug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, err
}
