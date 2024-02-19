package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PharmacyProductRepository interface {
	BaseRepository[entity.PharmacyProduct]
	FindAllPharmacyProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindNearbyProduct(ctx context.Context, query *valueobject.Query, userId uint, distanceInMeter int) (*valueobject.PagedResult, error)
	FindNearbyProductOrder(ctx context.Context, listProduct []uint, location *valueobject.Coordinate) ([]*entity.PharmacyProduct, error)
	BulkCreate(ctx context.Context, products []*entity.PharmacyProduct) ([]*entity.PharmacyProduct, error)
	FindTopPrice(ctx context.Context, productId uint, isTop bool) (decimal.Decimal, error)
	FindRangePrice(context.Context, []uint, bool) (map[uint]string, error)
	FindAllPharmacyAvailableProductId(ctx context.Context, pharmcyProduct *entity.PharmacyProduct) ([]*entity.Pharmacy, error)
}

type pharmacyProductRepository struct {
	*baseRepository[entity.PharmacyProduct]
	db *gorm.DB
}

func NewPharmacyProductRepository(db *gorm.DB) PharmacyProductRepository {
	return &pharmacyProductRepository{
		db:             db,
		baseRepository: &baseRepository[entity.PharmacyProduct]{db: db},
	}
}

func (r *pharmacyProductRepository) FindNearbyProduct(ctx context.Context, query *valueobject.Query, userId uint, distanceInMeter int) (*valueobject.PagedResult, error) {
	switch strings.Split(query.GetOrder(), " ")[0] {
	case "name":
		query.WithSortBy("\"Product\".name")
	case "price":
		query.WithSortBy("\"Product\".price")
	}
	pagedResult, err := r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		db.
			Joins("Product").
			Joins("Pharmacy").
			Joins(fmt.Sprintf("JOIN addresses a on st_dwithin(\"Pharmacy\".location, a.location, %d)", distanceInMeter)).
			Joins("JOIN users u on u.id=a.profile_id").
			Where("u.id = ?", userId).
			Where("a.is_default = ?", true)

		category := query.GetConditionValue("category")
		name := query.GetConditionValue("name")

		if category != nil {
			db.Where("product_category_id", category)
		}

		if name != nil {
			db.Where("\"Product\".name ILIKE ?", name)
		}

		return db
	})
	if err != nil {
		return nil, err
	}

	products := make([]*entity.Product, 0)
	fetchedPharmacyProducts := pagedResult.Data.([]*entity.PharmacyProduct)
	for _, pharmacyProduct := range fetchedPharmacyProducts {
		products = append(products, pharmacyProduct.Product)
	}
	pagedResult.Data = products
	return pagedResult, nil
}

func (r *pharmacyProductRepository) FindAllPharmacyProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "name":
			query.WithSortBy("\"Product\".name")
		case "id":
			query.WithSortBy("\"pharmacy_products\".id ")
		}

		category := query.GetConditionValue("category")
		name := query.GetConditionValue("name")
		pharmacy := query.GetConditionValue("pharmacy")
		isActive := query.GetConditionValue("is_active")
		db.Joins("Product").Joins("Product.ProductCategory").Joins("Pharmacy")
		if category != nil {
			db.Where("\"Product\".product_category_id", category)
		}

		if name != nil {
			db.Where("\"Product\".name ILIKE ?", name)
		}
		if isActive != nil {
			db.Where("is_active", isActive)
		}
		db.Where("Pharmacy.id", pharmacy)
		return db
	})
}

func (r *pharmacyProductRepository) FindNearbyProductOrder(ctx context.Context, listProduct []uint, location *valueobject.Coordinate) ([]*entity.PharmacyProduct, error) {
	var pharmacyProducts []*entity.PharmacyProduct
	longitude := location.Longitude.String()
	latitude := location.Latitude.String()
	err := r.conn(ctx).
		Model(&entity.PharmacyProduct{}).
		Joins("JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id ").
		Joins("JOIN addresses a ON st_dwithin(pharmacies.location,? , 25000)", location).
		Where("product_id IN ? ", listProduct).
		Order(fmt.Sprintf("pharmacies.location <-> ST_MakePoint(%s, %s) ", longitude, latitude)).
		Preload("Pharmacy").
		Find(&pharmacyProducts).Error
	if err != nil {
		return nil, err
	}
	return pharmacyProducts, nil
}

func (r *pharmacyProductRepository) BulkCreate(ctx context.Context, products []*entity.PharmacyProduct) ([]*entity.PharmacyProduct, error) {
	err := r.conn(ctx).Model(&entity.PharmacyProduct{}).Create(products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *pharmacyProductRepository) FindTopPrice(ctx context.Context, product uint, isTop bool) (decimal.Decimal, error) {
	var pharmacyProduct entity.PharmacyProduct
	var raw string
	if isTop {
		raw = "price DESC"
	} else {
		raw = "price ASC"
	}
	err := r.conn(ctx).
		Model(&entity.PharmacyProduct{}).
		Where("product_id = ? ", product).
		Order(raw).
		Limit(1).
		Find(&pharmacyProduct).Error
	if err != nil {
		return decimal.Zero, err
	}
	return pharmacyProduct.Price, nil
}

func (r *pharmacyProductRepository) FindRangePrice(ctx context.Context, listProduct []uint, isTop bool) (map[uint]string, error) {
	type Result struct {
		ProductId uint            `gorm:"column:product_id"`
		Price     decimal.Decimal `gorm:"column:range_price"`
	}
	var result []Result
	query := r.conn(ctx).Model(&entity.PharmacyProduct{})
	if isTop {
		query = query.Select("product_id, MAX(price) as range_price")
	} else {
		query = query.Select("product_id, MIN(price) as range_price")
	}
	err := query.Group("product_id").Scan(&result).Error
	if err != nil {
		return nil, err
	}
	listOfPrice := make(map[uint]string)
	for _, productUid := range listProduct {
		find := false
		for _, product := range result {
			if productUid == product.ProductId {
				listOfPrice[productUid] = product.Price.String()
				find = true
				break
			}
		}
		if !find {
			listOfPrice[productUid] = ""
		}
	}
	return listOfPrice, nil
}

func (r *pharmacyProductRepository) FindAllPharmacyAvailableProductId(ctx context.Context, pharmcyProduct *entity.PharmacyProduct) ([]*entity.Pharmacy, error) {
	listPharmacy := []*entity.Pharmacy{}
	err := r.conn(ctx).Raw(`SELECT p.* FROM pharmacy_products AS pp JOIN pharmacies AS p ON pp.pharmacy_id = p.id
	WHERE pp.product_id = ? AND NOT pp.pharmacy_id = ?`, pharmcyProduct.ProductId, pharmcyProduct.PharmacyId).Scan(&listPharmacy).Error
	if err != nil {
		return listPharmacy, err
	}
	return listPharmacy, nil
}
