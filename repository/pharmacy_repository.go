package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type PharmacyRepository interface {
	BaseRepository[entity.Pharmacy]
	FindAllPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindNearestPharmacyFromAddress(ctx context.Context, addressId uint) ([]*entity.Pharmacy, error)
	FindAllPharmacySuperAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type pharmacyRepository struct {
	*baseRepository[entity.Pharmacy]
	db *gorm.DB
}

func NewPharmacyRepository(db *gorm.DB) PharmacyRepository {
	return &pharmacyRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Pharmacy]{db: db},
	}
}

func (r *pharmacyRepository) FindAllPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "name":
			query.WithSortBy("\"pharmacies\".name")
		case "id":
			query.WithSortBy("\"pharmacies\".id ")
		}
		db.Where("\"pharmacies\".admin_id = ?", ctx.Value("user_id").(uint))
		name := query.GetConditionValue("name")
		db.Joins("City").Joins("Province")

		if name != nil {
			db.Where("\"pharmacies\".name ILIKE ?", name)
		}
		return db
	})
}

func (r *pharmacyRepository) FindAllPharmacySuperAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "pharmacy_name":
			query.WithSortBy("\"pharmacies\".name")
		case "admin_name":
			query.WithSortBy("\"Admin__AdminContact\".name")

		case "id":
			query.WithSortBy("\"pharmacies\".id ")
		}
		name := query.GetConditionValue("name")
		db.Joins("City").Joins("Province").Joins("Admin.AdminContact")
		province := query.GetConditionValue("province")
		if name != nil {
			db.Where("\"pharmacies\".name ILIKE ?", name)
		}
		if province != nil {
			db.Where("\"pharmacies\".province_id = ?", province)
		}
		return db
	})
}

func (r *pharmacyRepository) FindNearestPharmacyFromAddress(ctx context.Context, addressId uint) ([]*entity.Pharmacy, error) {
	var pharmacy []*entity.Pharmacy
	err := r.db.
		Joins("City").
		Joins("JOIN addresses a ON st_dwithin(\"pharmacies\".location, a.location, 25000)").
		Where("a.id=?", addressId).
		Order("\"pharmacies\".location <-> a.location").
		Find(&pharmacy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return pharmacy, nil
}
