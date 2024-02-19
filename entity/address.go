package entity

import (
	"time"

	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type Address struct {
	Id         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"not null"`
	StreetName string `gorm:"not null"`
	PostalCode string `gorm:"not null"`
	Phone      string `gorm:"not null"`
	Detail     string
	Location   *valueobject.Coordinate `gorm:"not null;type:geography(POINT, 4326)"`
	IsDefault  bool                    `gorm:"not null"`
	ProfileId  uint                    `gorm:"not null"`
	Profile    Profile                 `gorm:"foreignKey:ProfileId;references:UserId"`
	ProvinceId uint                    `gorm:"not null"`
	Province   Province                `gorm:"foreignKey:ProvinceId;references:Id"`
	CityId     uint                    `gorm:"not null"`
	City       City                    `gorm:"foreignKey:CityId;references:Id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
