package repository

import (
	"context"
	"fmt"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type AddressRepository interface {
	BaseRepository[entity.Address]
	ValidateCoordinate(ctx context.Context, cityId uint, coordinate *valueobject.Coordinate) (bool, error)
}

type addressRepository struct {
	*baseRepository[entity.Address]
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Address]{db: db},
	}
}

func (r *addressRepository) ValidateCoordinate(ctx context.Context, cityId uint, coordinate *valueobject.Coordinate) (bool, error) {
	query := `SELECT COUNT(*)
FROM cities c
		 JOIN cities_gadm cg ON c.city_gid = cg.gid
WHERE c.id = ?
  AND st_intersects(cg.geom, st_geogfromtext(?))`

	var count int64
	err := r.db.WithContext(ctx).Raw(query, cityId, fmt.Sprintf("SRID=4326;POINT(%s %s)", coordinate.Longitude, coordinate.Latitude)).Scan(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
