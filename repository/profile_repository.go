package repository

import (
	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	BaseRepository[entity.Profile]
}

type profileRepository struct {
	*baseRepository[entity.Profile]
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Profile]{db: db},
	}
}
