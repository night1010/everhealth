package repository

import (
	"context"
	"strings"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type StockMutationRepository interface {
	BaseRepository[entity.StockMutation]
	FindAllStockMutation(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	BulkCreate(ctx context.Context, mutations []*entity.StockMutation) error
}

type stockMutationRepository struct {
	*baseRepository[entity.StockMutation]
	db *gorm.DB
}

func NewStockMutationRepository(db *gorm.DB) StockMutationRepository {
	return &stockMutationRepository{
		db:             db,
		baseRepository: &baseRepository[entity.StockMutation]{db: db},
	}
}

func (r *stockMutationRepository) FindAllStockMutation(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "id":
			query.WithSortBy("\"stock_mutations\".id ")
		case "mutated":
			query.WithSortBy("\"stock_mutations\".mutated_at")
		}
		db.Joins("ToPharmacyProduct.Pharmacy").Joins("ToPharmacyProduct.Product").Joins("FromPharmacyProduct.Pharmacy").Joins("FromPharmacyProduct.Product")

		status := query.GetConditionValue("status")
		name := query.GetConditionValue("pharmacy_name")
		adminId := ctx.Value("user_id").(uint)
		db.Where("\"ToPharmacyProduct__Pharmacy\".admin_id = ? OR \"FromPharmacyProduct__Pharmacy\".admin_id = ?", adminId, adminId)
		if status != nil {
			db.Where("\"stock_mutations\".status = ?", status)
		}
		if name != nil {
			db.Where("(\"ToPharmacyProduct__Pharmacy\".name ILIKE ? OR \"FromPharmacyProduct__Pharmacy\".name ILIKE ?)", name, name)
		}
		return db
	})
}

func (r *stockMutationRepository) BulkCreate(ctx context.Context, mutations []*entity.StockMutation) error {
	return r.conn(ctx).Model(&entity.StockMutation{}).Create(mutations).Error
}
