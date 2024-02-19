package repository

import (
	"context"
	"errors"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type ProvinceRepository interface {
	BaseRepository[entity.Province]
	FindAllProvince(ctx context.Context) ([]*entity.Province, error)
	FindProvinceDetail(ctx context.Context, id uint) (*entity.Province, error)
}

type provinceRepository struct {
	*baseRepository[entity.Province]
	db *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) ProvinceRepository {
	return &provinceRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Province]{db: db},
	}
}

func (r *provinceRepository) FindAllProvince(ctx context.Context) ([]*entity.Province, error) {
	provinces := make([]*entity.Province, 0)
	err := r.conn(ctx).Find(&provinces).Error
	if err != nil {
		return provinces, err
	}
	return provinces, nil

}

func (r *provinceRepository) FindProvinceDetail(ctx context.Context, id uint) (*entity.Province, error) {
	var province *entity.Province
	err := r.conn(ctx).
		Preload("Cities", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN cities_gadm cg ON \"cities\".city_gid = cg.gid")
		}).
		Where("\"provinces\".\"id\" = ?", id).
		First(&province).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return province, nil
}
