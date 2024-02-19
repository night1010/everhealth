package repository

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	BaseRepository[entity.CartItem]
	CheckAllItem(context.Context, bool) error
	BulkDelete(context.Context, []*entity.CartItem) error
}

type cartItemRepository struct {
	*baseRepository[entity.CartItem]
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{
		db:             db,
		baseRepository: &baseRepository[entity.CartItem]{db: db},
	}
}

func (r *cartItemRepository) CheckAllItem(ctx context.Context, check bool) error {
	userId := ctx.Value("user_id").(uint)
	err := r.conn(ctx).
		Model(&entity.CartItem{}).
		Where("cart_id = ?", userId).
		Where("is_checked = ?", !check).
		Update("is_checked", check).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *cartItemRepository) BulkDelete(ctx context.Context, items []*entity.CartItem) error {
	return r.conn(ctx).Model(&entity.CartItem{}).Delete(items).Error
}
