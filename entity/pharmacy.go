package entity

import (
	"time"

	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type Pharmacy struct {
	Id                      uint   `gorm:"primaryKey;autoIncrement"`
	Name                    string `gorm:"not null"`
	AdminId                 uint   `gorm:"not null"`
	Admin                   User   `gorm:"foreignKey:AdminId;references:Id"`
	Address                 string `gorm:"not null"`
	CityId                  uint   `gorm:"not null"`
	City                    City
	ProvinceId              uint `gorm:"not null"`
	Province                Province
	Location                *valueobject.Coordinate `gorm:"not null;type:geography(POINT,4326)"`
	StartTime               time.Time               `gorm:"not null"`
	EndTime                 time.Time               `gorm:"not null"`
	OperationalDay          string                  `gorm:"not null"`
	PharmacistLicenseNumber string                  `gorm:"not null"`
	PharmacistPhoneNumber   string                  `gorm:"not null"`
	Products                []PharmacyProduct       `gorm:"many2many:pharmacy_products"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt
}
