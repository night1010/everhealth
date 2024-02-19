package usecase

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
)

type OrderStatusUsecase interface {
	FindAllOrderStatus(ctx context.Context) ([]*entity.OrderStatus, error)
}

type orderStatusUsecase struct {
	orderStatusRepository repository.OrderStatusRepository
}

func NewOrderStatusUsecase(rp repository.OrderStatusRepository) OrderStatusUsecase {
	return &orderStatusUsecase{orderStatusRepository: rp}
}

func (u *orderStatusUsecase) FindAllOrderStatus(ctx context.Context) ([]*entity.OrderStatus, error) {
	return u.orderStatusRepository.FindAllOrderStatus(ctx)
}
