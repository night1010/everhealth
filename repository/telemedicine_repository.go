package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TelemedicineRepository interface {
	BaseRepository[entity.Telemedicine]
	FindAllTelemedicine(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindTelemedicineWhereStatusNotCancelAndEnd(ctx context.Context, telemedicine *entity.Telemedicine) (*entity.Telemedicine, error)
	FindByUserId(ctx context.Context, userId uint) (*entity.Telemedicine, error)
	FindByTelemedicineId(ctx context.Context, query *valueobject.Query) (*entity.Telemedicine, error)
}

type telemedicineRepository struct {
	*baseRepository[entity.Telemedicine]
	db *gorm.DB
}

func NewTelemedicineRepository(db *gorm.DB) TelemedicineRepository {
	return &telemedicineRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Telemedicine]{db: db},
	}
}

func (r *telemedicineRepository) FindTelemedicineWhereStatusNotCancelAndEnd(ctx context.Context, telemedicine *entity.Telemedicine) (*entity.Telemedicine, error) {
	var newTelemedicine *entity.Telemedicine
	err := r.conn(ctx).Where("profile_id =?", telemedicine.ProfileId).Where("status not in (?)", []entity.TelemedicineStatus{entity.Cancel, entity.End}).First(&newTelemedicine).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return newTelemedicine, nil
}

func (r *telemedicineRepository) FindByUserId(ctx context.Context, userId uint) (*entity.Telemedicine, error) {
	var fetchedTelemedicine *entity.Telemedicine
	err := r.conn(ctx).
		Where("profile_id=?", userId).
		Or("doctor_id=?", userId).
		First(&fetchedTelemedicine).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return fetchedTelemedicine, nil
}

func (r *telemedicineRepository) FindAllTelemedicine(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		query.WithSortBy("\"Doctor\".status DESC,ordered_at")
		db.Joins("Profile").
			Joins("Doctor.Profile").
			Preload("Chats", func(db *gorm.DB) *gorm.DB {
				return db.Order("chat_time DESC")
			})
		name := query.GetConditionValue("name")
		if ctx.Value("role_id").(entity.RoleId) == entity.RoleDoctor {
			if name != nil {
				db.Where("\"Profile\".name ILIKE ?", name)
			}
			db.Where("\"telemedicines\".Doctor_id =?", ctx.Value("user_id").(uint))
		} else {
			if name != nil {
				db.Where("\"Doctor__Profile\".name ILIKE ?", name)
			}
			db.Where("\"telemedicines\".Profile_id =?", ctx.Value("user_id").(uint))
		}
		status := query.GetConditionValue("status")
		if status != nil {
			db.Where("\"telemedicines\".status = ?", status)
		}
		db.Group("\"telemedicines\".id,\"Profile\".user_id,\"Doctor\".profile_id,\"Doctor__Profile\".user_id")
		return db
	})
}

func (r *telemedicineRepository) FindByTelemedicineId(ctx context.Context, q *valueobject.Query) (*entity.Telemedicine, error) {
	conditions := q.GetConditions()
	var t *entity.Telemedicine
	query := r.conn(ctx).Model(t)
	if q.IsLocked() {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	for _, s := range q.GetAssociations() {
		if s.Type == valueobject.AssociationTypeJoin {
			query.Joins(s.Entity)
		} else if s.Type == valueobject.AssociationTypePreload {
			query.Preload(s.Entity)
		}
	}

	query.Preload("Chats", func(db *gorm.DB) *gorm.DB {
		return db.Order("chat_time DESC")
	})

	for _, condition := range conditions {
		sql := fmt.Sprintf("%s %s ?", condition.Field, condition.Operation)
		query.Where(sql, condition.Value)
	}
	err := query.First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}
