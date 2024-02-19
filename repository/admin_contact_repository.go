package repository

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type AdminContactRepository interface {
	BaseRepository[entity.AdminContact]
	HardDelete(ctx context.Context, contact *entity.AdminContact) error
}

type adminContactRepository struct {
	*baseRepository[entity.AdminContact]
	db *gorm.DB
}

func NewAdminContactRepository(db *gorm.DB) AdminContactRepository {
	return &adminContactRepository{
		db:             db,
		baseRepository: &baseRepository[entity.AdminContact]{db: db},
	}
}

func (r *adminContactRepository) HardDelete(ctx context.Context, contact *entity.AdminContact) error {
	result := r.conn(ctx).Unscoped().Delete(&contact)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
