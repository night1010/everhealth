package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository[entity.Product]
	FindAllProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindNearbyProducts(ctx context.Context, query *valueobject.Query, userId uint, distanceInMeter int) (*valueobject.PagedResult, error)
}

type productRepository struct {
	*baseRepository[entity.Product]
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Product]{db: db},
	}
}

func (r *productRepository) FindAllProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "price":
			query.WithSortBy("price")
		}

		category := query.GetConditionValue("category")
		name := query.GetConditionValue("name")
		isHidden := query.GetConditionValue("is_hidden")
		db.Joins("ProductCategory")

		if category != nil {
			db.Where("product_category_id", category)
		}

		if name != nil {
			db.Where("products.name ILIKE ?", name)
		}

		if isHidden != nil {
			db.Where("products.is_hidden = ?", isHidden)
		}
		return db
	})
}

func (r *productRepository) FindNearbyProducts(ctx context.Context, query *valueobject.Query, userId uint, distanceInMeter int) (*valueobject.PagedResult, error) {
	switch strings.Split(query.GetOrder(), " ")[0] {
	case "name":
		query.WithSortBy("\"products\".name")
	case "price":
		query.WithSortBy("\"products\".price")
	}
	pagedResult, err := r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		db.
			Select("\"products\".name as name", "\"products\".id as id", "\"products\".image as image", "\"products\".selling_unit as selling_unit").
			Joins("JOIN pharmacy_products pp on \"products\".id=pp.product_id").
			Joins("JOIN pharmacies ph ON ph.id=pp.pharmacy_id").
			Joins(fmt.Sprintf("JOIN addresses a on st_dwithin(ph.location, a.location, %d)", distanceInMeter)).
			Joins("JOIN users u on u.id=a.profile_id").
			Where("u.id = ?", userId).
			Where("a.is_default = ?", true)

		category := query.GetConditionValue("category")
		name := query.GetConditionValue("name")

		if category != nil {
			db.Where("product_category_id", category)
		}

		if name != nil {
			db.Where("\"products\".name ILIKE ?", name)
		}

		db.Group("\"products\".id")

		return db
	})
	if err != nil {
		return nil, err
	}

	return pagedResult, nil
}
