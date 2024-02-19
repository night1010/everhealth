package repository

import (
	"context"
	"strings"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type DoctorProfileRepository interface {
	BaseRepository[entity.DoctorProfile]
	FindAllDoctors(context.Context, *valueobject.Query) (*valueobject.PagedResult, error)
}

type doctorProfileRepository struct {
	*baseRepository[entity.DoctorProfile]
	db *gorm.DB
}

func NewDoctorProfileRepository(db *gorm.DB) DoctorProfileRepository {
	return &doctorProfileRepository{
		db:             db,
		baseRepository: &baseRepository[entity.DoctorProfile]{db: db},
	}
}

func (r *doctorProfileRepository) FindAllDoctors(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "name":
			query.WithSortBy("Status ASC,\"Profile\".name")
		case "fee":
			query.WithSortBy("Status ASC,fee")
		}

		specialization := query.GetConditionValue("specialization")
		name := query.GetConditionValue("name")
		db.Joins("Profile").Joins("Specialist")

		if specialization != nil {
			db.Where("specialist_id", specialization)
		}

		if name != nil {
			db.Where("\"Profile\".name ILIKE ?", name)
		}
		return db
	})
}
