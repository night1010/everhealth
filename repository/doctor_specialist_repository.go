package repository

import (
	"context"

	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type DoctorSpecialistRepository interface {
	BaseRepository[entity.DoctorSpecialist]
	FindAllDoctorSpecialist(ctx context.Context) ([]*entity.DoctorSpecialist, error)
}

type doctorSpecialistRepository struct {
	*baseRepository[entity.DoctorSpecialist]
	db *gorm.DB
}

func NewDoctorSpecialistRepository(db *gorm.DB) DoctorSpecialistRepository {
	return &doctorSpecialistRepository{
		db:             db,
		baseRepository: &baseRepository[entity.DoctorSpecialist]{db: db},
	}
}

func (r *doctorSpecialistRepository) FindAllDoctorSpecialist(ctx context.Context) ([]*entity.DoctorSpecialist, error) {
	DoctorSpecialists := make([]*entity.DoctorSpecialist, 0)
	err := r.conn(ctx).Find(&DoctorSpecialists).Error
	if err != nil {
		return DoctorSpecialists, err
	}
	return DoctorSpecialists, nil

}
