package repository

import (
	"context"
	"strings"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository[entity.User]
	FindAllAdminPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	HardDelete(ctx context.Context, user *entity.User) error
	FindAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type userRepository struct {
	*baseRepository[entity.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db:             db,
		baseRepository: &baseRepository[entity.User]{db: db},
	}
}

func (r *userRepository) FindAllAdminPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "email":
			query.WithSortBy("\"users\".email")
		case "id":
			query.WithSortBy("\"users\".id ")
		}
		db.Where("role_id", entity.RoleAdmin).Preload("AdminContact")
		name := query.GetConditionValue("email")
		if name != nil {
			db.Where("\"users\".email ILIKE ?", name)
		}
		return db
	})
}

func (r *userRepository) FindAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "email":
			query.WithSortBy("\"users\".email")
		case "id":
			query.WithSortBy("\"users\".id ")
		}
		roleId := query.GetConditionValue("role_id")
		name := query.GetConditionValue("email")
		isVerified := query.GetConditionValue("is_verified")
		db.Joins("Role")
		if name != nil {
			db.Where("\"users\".email ILIKE ?", name)
		}
		if roleId != nil {
			db.Where("\"users\".role_id = ?", roleId)

		} else {
			db.Where("\"users\".role_id IN ?", []entity.RoleId{entity.RoleAdmin, entity.RoleDoctor, entity.RoleUser})
		}
		db.Preload("AdminContact").Preload("Profile")
		if isVerified != nil {
			db.Where("\"users\".is_verified = ?", isVerified)
		}
		return db
	})
}

func (r *userRepository) HardDelete(ctx context.Context, user *entity.User) error {
	result := r.conn(ctx).Unscoped().Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
