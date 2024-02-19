package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type StatusDoctor int

const (
	Online  StatusDoctor = 1
	Busy    StatusDoctor = 2
	Offline StatusDoctor = 3
)

var DoctorStatusMap = map[StatusDoctor]string{
	Online:  "online",
	Busy:    "busy",
	Offline: "offline",
}

type DoctorProfile struct {
	ProfileId        uint    `gorm:"primaryKey"`
	Profile          Profile `gorm:"foreignKey:ProfileId;references:UserId"`
	Certificate      string  `gorm:"not null"`
	CertificateKey   string
	SpecialistId     uint             `gorm:"not null"`
	Specialist       DoctorSpecialist `gorm:"foreignKey:SpecialistId;references:Id"`
	YearOfExperience uint             `gorm:"not null"`
	Status           StatusDoctor     `gorm:"not null"`
	Fee              decimal.Decimal  `gorm:"type:numeric"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

const (
	DoctorCertificateFolder = "doctor-certificate"
	DoctorCertificatePrefix = "doctor-certificate-"
)
