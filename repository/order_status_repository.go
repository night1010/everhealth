package repository

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type OrderStatusRepository interface {
	BaseRepository[entity.OrderStatus]
	FindAllOrderStatus(ctx context.Context) ([]*entity.OrderStatus, error)
}

type orderStatusRepository struct {
	*baseRepository[entity.OrderStatus]
	db *gorm.DB
}

func NewOrderStatusRepository(db *gorm.DB) OrderStatusRepository {
	return &orderStatusRepository{
		db:             db,
		baseRepository: &baseRepository[entity.OrderStatus]{db: db},
	}
}

func (r *orderStatusRepository) FindAllOrderStatus(ctx context.Context) ([]*entity.OrderStatus, error) {
	OrderStatuses := make([]*entity.OrderStatus, 0)
	err := r.conn(ctx).Find(&OrderStatuses).Error
	if err != nil {
		return OrderStatuses, err
	}
	return OrderStatuses, nil

}
